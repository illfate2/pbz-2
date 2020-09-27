package api

import (
	"log"

	"github.com/labstack/echo/v4"
)

func (s *Server) initMuseumSet(e *echo.Echo) {
	e.GET("/museumSets", s.getMuseumSets)
	e.GET("/museumSet/:id", s.getMuseumSet)
}

func (s *Server) getMuseumSets(c echo.Context) error {
	items, err := s.service.GetMuseumSets()
	if err != nil {
		log.Printf("Failed to find museum sets: %+v", err)
		return err
	}
	_ = s.tmpl.ExecuteTemplate(c.Response().Writer, "Museum sets", items)
	return nil
}

func (s *Server) getMuseumSet(c echo.Context) error {
	id := getIDFromURL(c)
	set, err := s.service.FindMuseumSet(id)
	if err != nil {
		log.Printf("Failed to find item with details: %s", err)
		return err
	}

	err = s.tmpl.ExecuteTemplate(c.Response(), "ShowMuseumSet", set)
	log.Print(err)
	return nil
}
