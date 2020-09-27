package api

import (
	"context"
	"html/template"
	"net/http"

	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"

	"pbz2/pkg/repo"
	"pbz2/pkg/service"
)

type Server struct {
	http.Handler
	service *service.Service
	tmpl    *template.Template
}

func NewServer() *Server {
	e := echo.New()
	conn, err := pgx.Connect(context.TODO(), "postgresql://pbz2:pbz2@localhost:5433/pbz2?sslmode=disable") // todo
	repo := repo.New(conn)
	if err != nil {
		panic(err)
	}
	s := &Server{
		Handler: e,
		service: service.NewService(repo),
		tmpl:    template.Must(template.ParseGlob("static/html/*")),
	}
	s.initMuseumItemAPI(e)
	s.initMuseumItemMovement(e)
	s.initMuseumSet(e)

	e.Static("/", "static")
	return s
}
