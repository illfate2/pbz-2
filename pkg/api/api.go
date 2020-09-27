package api

import (
	"context"
	"html/template"
	"net/http"

	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"

	"pbz2/pkg/service"
)

type Server struct {
	http.Handler
	service *service.Service
	tmpl    *template.Template
}

func NewServer() *Server {
	e := echo.New()
	conn, err := pgx.Connect(context.TODO(), "postgresql://pbz2:pbz2@localhost:5433/pbz2?sslmode=disable")
	if err != nil {
		panic(err)
	}
	s := &Server{
		Handler: e,
		service: service.NewService(conn),
		tmpl:    template.Must(template.ParseGlob("static/html/*")),
	}
	e.POST("/museumItem", s.createMuseumItem)
	e.GET("/museumItem", s.getMuseumItemPage)
	e.GET("/museumItem/:id", s.getMuseumItem)
	e.GET("/museumItems", s.getMuseumItems)
	e.GET("/deleteMuseumItem/:id", s.deleteMuseumItem)
	e.GET("/editMuseumItem/:id", s.editMuseumItem)
	e.POST("/editMuseumItem/:id", s.updateMuseumItem)

	e.GET("/museumItemSearch", s.getMuseumItemSearchPage)
	e.POST("/museumItemSearch", s.searchMuseumItems)

	e.POST("/museumItemMovement", s.createMuseumItemMovement)
	e.GET("/museumItemMovements", s.getMuseumItemMovements)
	e.GET("/museumItemMovement", s.getMuseumItemMovementPage)
	e.GET("/museumItemMovement/:id", s.getMuseumItemMovement)

	e.GET("/museumSets", s.getMuseumSets)
	e.GET("/museumSet/:id", s.getMuseumSet)


	e.Static("/", "static")
	return s
}
