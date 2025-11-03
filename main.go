package main

import (
	"database/sql"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed templates/*.html static/*
var content embed.FS

//go:embed openapi.yaml
var openapiSpec string

var tpl *template.Template
var db *sql.DB
var sCookie *securecookie.SecureCookie

// SSE subscribers per room
var subsMu sync.Mutex
var subs = map[string]map[chan string]struct{}{}

func subscribe(room string) chan string {
	ch := make(chan string, 1)
	subsMu.Lock()
	defer subsMu.Unlock()
	if subs[room] == nil {
		subs[room] = map[chan string]struct{}{}
	}
	subs[room][ch] = struct{}{}
	return ch
}

func unsubscribe(room string, ch chan string) {
	subsMu.Lock()
	defer subsMu.Unlock()
	if subs[room] != nil {
		delete(subs[room], ch)
	}
	close(ch)
}

func publish(room, msg string) {
	subsMu.Lock()
	defer subsMu.Unlock()
	for ch := range subs[room] {
		select {
		case ch <- msg:
		default:
		}
	}
}

func main() {
	var err error
	db, err = sql.Open("sqlite3", "file:estimate.db?_foreign_keys=on")
	if err != nil {
		log.Fatal(err)
	}
	if err := migrate(db); err != nil {
		log.Fatal(err)
	}

	sCookie = securecookie.New(securecookie.GenerateRandomKey(32), securecookie.GenerateRandomKey(32))

	tpl = template.Must(template.ParseFS(content, "templates/*.html"))

	r := gin.Default()
	r.StaticFS("/static", http.FS(content))

	r.GET("/", indexHandler)
	r.POST("/rooms", createRoomHandler)
	r.GET("/login", loginHandler)
	r.POST("/login", doLogin)
	r.GET("/register", registerHandler)
	r.POST("/register", doRegister)
	r.GET("/logout", doLogout)
	// admin user management
	r.GET("/admin/users", adminListUsersHandler)
	r.POST("/admin/users/create", adminCreateUserHandler)
	r.POST("/admin/users/delete", adminDeleteUserHandler)
	r.POST("/admin/users/role", adminChangeRoleHandler)

	r.GET("/room/:id", roomHandler)
	r.GET("/room/:id/votes", votesPartialHandler)
	r.GET("/room/:id/stream", roomStreamHandler)
	r.POST("/room/:id/vote", voteHandler)
	r.POST("/room/:id/reveal", revealHandler)
	r.GET("/openapi.yaml", func(c *gin.Context) {
		c.Header("Content-Type", "application/yaml")
		c.String(200, openapiSpec)
	})
	r.GET("/docs", func(c *gin.Context) {
		tpl.ExecuteTemplate(c.Writer, "docs.html", nil)
	})

	fmt.Println("Listening on http://localhost:8080")
	r.Run()
}

func migrate(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, username TEXT UNIQUE, password TEXT, role TEXT);`,
		`CREATE TABLE IF NOT EXISTS rooms (id INTEGER PRIMARY KEY, name TEXT, created_at DATETIME);`,
		`CREATE TABLE IF NOT EXISTS votes (id INTEGER PRIMARY KEY, room_id INTEGER, user_id INTEGER, value TEXT, revealed INTEGER DEFAULT 0, FOREIGN KEY(room_id) REFERENCES rooms(id) ON DELETE CASCADE, FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE);`,
	}
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			return err
		}
	}
	// Ensure at least one admin exists
	var cnt int
	row := db.QueryRow(`SELECT COUNT(1) FROM users WHERE role='admin'`)
	_ = row.Scan(&cnt)
	if cnt == 0 {
		// hash default password
		hash, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
		_, err := db.Exec(`INSERT INTO users(username,password,role) VALUES(?,?,?)`, "admin", string(hash), "admin")
		if err != nil {
			return err
		}
		fmt.Println("Created default admin: username=admin password=admin")
	}
	return nil
}

// helpers: get current user from cookie (very small/simple auth for scaffold)
func getCurrentUser(c *gin.Context) (id int, username, role string) {
	if cookie, err := c.Request.Cookie("session"); err == nil {
		val := make(map[string]string)
		if err := sCookie.Decode("session", cookie.Value, &val); err == nil {
			// try to read user
			fmt.Sscanf(val["id"], "%d", &id)
			username = val["username"]
			role = val["role"]
			return
		}
	}
	return 0, "", ""
}

func indexHandler(c *gin.Context) {
	_, username, _ := getCurrentUser(c)
	rows, err := db.Query(`SELECT id, name FROM rooms ORDER BY created_at DESC`)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	defer rows.Close()
	type room struct {
		ID   int
		Name string
	}
	var rooms []room
	for rows.Next() {
		var r room
		rows.Scan(&r.ID, &r.Name)
		rooms = append(rooms, r)
	}
	tpl.ExecuteTemplate(c.Writer, "index.html", gin.H{"User": username, "Rooms": rooms})
}

func createRoomHandler(c *gin.Context) {
	_, username, role := getCurrentUser(c)
	if username == "" || role != "admin" {
		c.String(403, "admin only")
		return
	}
	name := c.PostForm("name")
	if name == "" {
		c.String(400, "name required")
		return
	}
	_, err := db.Exec(`INSERT INTO rooms(name, created_at) VALUES(?, datetime('now'))`, name)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	// if HTMX request, return updated rooms partial
	if c.GetHeader("HX-Request") == "true" {
		rows, _ := db.Query(`SELECT id, name FROM rooms ORDER BY created_at DESC`)
		type room struct {
			ID   int
			Name string
		}
		var rooms []room
		for rows.Next() {
			var r room
			rows.Scan(&r.ID, &r.Name)
			rooms = append(rooms, r)
		}
		tpl.ExecuteTemplate(c.Writer, "rooms_partial.html", gin.H{"Rooms": rooms})
		return
	}
	c.Redirect(302, "/")
}

// Admin: list users
func adminListUsersHandler(c *gin.Context) {
	_, username, role := getCurrentUser(c)
	if username == "" || role != "admin" {
		c.String(403, "admin only")
		return
	}
	rows, err := db.Query(`SELECT id, username, role FROM users ORDER BY id`)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	defer rows.Close()
	type u struct {
		ID             int
		Username, Role string
	}
	var users []u
	for rows.Next() {
		var uu u
		rows.Scan(&uu.ID, &uu.Username, &uu.Role)
		users = append(users, uu)
	}
	tpl.ExecuteTemplate(c.Writer, "admin_users.html", gin.H{"Users": users})
}

func adminCreateUserHandler(c *gin.Context) {
	_, username, role := getCurrentUser(c)
	if username == "" || role != "admin" {
		c.String(403, "admin only")
		return
	}
	uname := c.PostForm("username")
	pass := c.PostForm("password")
	r := c.PostForm("role")
	if uname == "" || pass == "" || r == "" {
		c.String(400, "missing fields")
		return
	}
	if r != "admin" && r != "voter" && r != "observer" {
		c.String(400, "invalid role")
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	_, err := db.Exec(`INSERT INTO users(username,password,role) VALUES(?,?,?)`, uname, string(hash), r)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.Redirect(302, "/admin/users")
}

func adminDeleteUserHandler(c *gin.Context) {
	_, username, role := getCurrentUser(c)
	if username == "" || role != "admin" {
		c.String(403, "admin only")
		return
	}
	id := c.PostForm("id")
	if id == "" {
		c.String(400, "id required")
		return
	}
	db.Exec(`DELETE FROM users WHERE id=?`, id)
	c.Redirect(302, "/admin/users")
}

func adminChangeRoleHandler(c *gin.Context) {
	_, username, role := getCurrentUser(c)
	if username == "" || role != "admin" {
		c.String(403, "admin only")
		return
	}
	id := c.PostForm("id")
	r := c.PostForm("role")
	if id == "" || r == "" {
		c.String(400, "missing fields")
		return
	}
	if r != "admin" && r != "voter" && r != "observer" {
		c.String(400, "invalid role")
		return
	}
	db.Exec(`UPDATE users SET role=? WHERE id=?`, r, id)
	c.Redirect(302, "/admin/users")
}

func votesPartialHandler(c *gin.Context) {
	id := c.Param("id")
	rows, _ := db.Query(`SELECT v.id, u.username, v.value, v.revealed FROM votes v JOIN users u ON u.id=v.user_id WHERE v.room_id=?`, id)
	type vrow struct {
		ID       int
		User     string
		Value    string
		Revealed int
	}
	var votes []vrow
	revealedAny := false
	for rows.Next() {
		var vr vrow
		rows.Scan(&vr.ID, &vr.User, &vr.Value, &vr.Revealed)
		votes = append(votes, vr)
		if vr.Revealed == 1 {
			revealedAny = true
		}
	}
	tpl.ExecuteTemplate(c.Writer, "votes_partial.html", gin.H{"Votes": votes, "Revealed": revealedAny})
}

func roomStreamHandler(c *gin.Context) {
	room := c.Param("id")
	ch := subscribe(room)
	clientGone := c.Writer.CloseNotify()
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.WriteHeader(200)
	// send initial ping
	fmt.Fprintf(c.Writer, "event: ping\ndata: ready\n\n")
	if f, ok := c.Writer.(http.Flusher); ok {
		f.Flush()
	}
	for {
		select {
		case <-clientGone:
			unsubscribe(room, ch)
			return
		case msg := <-ch:
			fmt.Fprintf(c.Writer, "data: %s\n\n", msg)
			if f, ok := c.Writer.(http.Flusher); ok {
				f.Flush()
			}
		}
	}
}

func loginHandler(c *gin.Context) {
	tpl.ExecuteTemplate(c.Writer, "login.html", nil)
}

func doLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	var id int
	var role string
	var hash string
	err := db.QueryRow(`SELECT id, role, password FROM users WHERE username=?`, username).Scan(&id, &role, &hash)
	if err != nil {
		tpl.ExecuteTemplate(c.Writer, "login.html", gin.H{"Error": "invalid credentials"})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		tpl.ExecuteTemplate(c.Writer, "login.html", gin.H{"Error": "invalid credentials"})
		return
	}
	value := map[string]string{"id": fmt.Sprintf("%d", id), "username": username, "role": role}
	if encoded, err := sCookie.Encode("session", value); err == nil {
		http.SetCookie(c.Writer, &http.Cookie{Name: "session", Value: encoded, Path: "/", Expires: time.Now().Add(24 * time.Hour)})
	}
	c.Redirect(302, "/")
}

func registerHandler(c *gin.Context) {
	tpl.ExecuteTemplate(c.Writer, "register.html", nil)
}

func doRegister(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	role := c.PostForm("role")
	if username == "" || password == "" || role == "" {
		tpl.ExecuteTemplate(c.Writer, "register.html", gin.H{"Error": "all fields required"})
		return
	}
	if role != "admin" && role != "voter" && role != "observer" {
		tpl.ExecuteTemplate(c.Writer, "register.html", gin.H{"Error": "invalid role"})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	_, err = db.Exec(`INSERT INTO users(username,password,role) VALUES(?,?,?)`, username, string(hash), role)
	if err != nil {
		tpl.ExecuteTemplate(c.Writer, "register.html", gin.H{"Error": "username taken or error"})
		return
	}
	c.Redirect(302, "/login")
}

func doLogout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{Name: "session", Value: "", Path: "/", Expires: time.Unix(0, 0)})
	c.Redirect(302, "/")
}

func roomHandler(c *gin.Context) {
	id := c.Param("id")
	var name string
	if err := db.QueryRow(`SELECT name FROM rooms WHERE id=?`, id).Scan(&name); err != nil {
		c.String(404, "room not found")
		return
	}
	// fetch votes but mask unless revealed
	rows, _ := db.Query(`SELECT v.id, u.username, v.value, v.revealed FROM votes v JOIN users u ON u.id=v.user_id WHERE v.room_id=?`, id)
	type vrow struct {
		ID       int
		User     string
		Value    string
		Revealed int
	}
	var votes []vrow
	revealedAny := false
	for rows.Next() {
		var vr vrow
		rows.Scan(&vr.ID, &vr.User, &vr.Value, &vr.Revealed)
		votes = append(votes, vr)
		if vr.Revealed == 1 {
			revealedAny = true
		}
	}
	_, username, role := getCurrentUser(c)
	tpl.ExecuteTemplate(c.Writer, "room.html", gin.H{"Room": name, "Votes": votes, "User": username, "Role": role, "Revealed": revealedAny})
}

func voteHandler(c *gin.Context) {
	id := c.Param("id")
	userID, _, role := getCurrentUser(c)
	if userID == 0 {
		c.String(403, "login to vote")
		return
	}
	if role == "observer" {
		c.String(403, "observers can't vote")
		return
	}
	value := c.PostForm("value")
	// find user id and upsert
	_, err := db.Exec(`INSERT OR REPLACE INTO votes(id, room_id, user_id, value, revealed) VALUES((SELECT id FROM votes WHERE room_id=? AND user_id=?), ?, ?, ?, 0)`, id, userID, id, userID, value)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	// notify subscribers
	publish(id, "vote")
	c.Redirect(302, "/room/"+id)
}

func revealHandler(c *gin.Context) {
	id := c.Param("id")
	_, _, role := getCurrentUser(c)
	if role != "admin" {
		c.String(403, "admin only")
		return
	}
	action := c.PostForm("action")
	if action == "reveal" {
		db.Exec(`UPDATE votes SET revealed=1 WHERE room_id=?`, id)
	} else if action == "reset" {
		db.Exec(`DELETE FROM votes WHERE room_id=?`, id)
	}
	publish(id, action)
	c.Redirect(302, "/room/"+id)
}
