# Sprint Names

API only project intended to handle Sprints' names.

The main goals of this project is to master Go language and Domain Driven Development.

## Setup

Install Wire and then generate the container

```shell
go install github.com/google/wire/cmd/wire@latest
cd container
wire
```

Build the project running

```shell
go build -o build/app cmd/main.go
```

If you want to directly run the binary, remember to put a valid `.env` file inside the `build` folder. 

## Todos

- [x] folder structure
- [x] module config
- [x] framework Gin
- [x] provider UUIDv7
- [x] DB sqlite in VCS (Gorm?) UUIDv7
- [ ] Domain Cartoons
  - [ ] repository
  - [ ] route handler
  - [ ] action
    - [ ] get cartoon
    - [ ] get all
    - [ ] save cartoon with characters
    - [ ] delete cartoon
    - [ ] edit cartoon
    - [ ] add character
    - [ ] delete character
    - [ ] edit character
- [ ] Domain Sprint
  - [x] repository
  - [ ] route handler
  - [ ] action
    - [ ] create sprint
    - [ ] get sprint
    - [x] get all
    - [ ] delete sprint
    - [ ] edit sprint
- [ ] Domain association
  - [ ] Service for access to other domains
  - [ ] repository
  - [ ] route handler
  - [ ] action
    - [ ] create association
    - [ ] edit association
    - [ ] delete association
- [ ] test
  - [ ] repository in memory
  - [ ] repository in CSV/JSON (textual DB in versioning)
- [ ] dependency injection
  - [ ] single container
  - [ ] multiple containers (sql, in memory)
- [ ] multilanguage
- [ ] OpenAPI
- [ ] OpenTelemetry
- [ ] Interface (HTMLX?)

## Credits

Property of [Progetti e Soluzioni Group](https://progettiesoluzioni.it)
