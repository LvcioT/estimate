package http_gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func page1Handler(c *gin.Context) {
	c.HTML(http.StatusOK, "page.gohtml", gin.H{
		// this value goes to the "header.tmpl" file
		"title":  "Page 1 IS HERE",
		"number": 1,
	})
}

func page2Handler(c *gin.Context) {
	c.HTML(http.StatusOK, "page.gohtml", gin.H{
		// this value goes to the "header.tmpl" file
		"title":  "Page 2 IS HERE",
		"number": 2,
	})
}
