package api

import (
	"log"
	"time"

	"github.com/labstack/echo/v4"

	"pbz2/pkg/service"
)

func (s *Server) getMuseumItemMovements(c echo.Context) error {
	movements, err := s.service.FindMuseumItemMovements()
	if err != nil {
		log.Printf("Failed to find museum item movement: %+v", err)
		return err
	}
	_ = s.tmpl.ExecuteTemplate(c.Response().Writer, "Museum item Movements", movements)
	return nil
}

func (s *Server) getMuseumItemMovement(c echo.Context) error {
	id := getIDFromURL(c)
	item, err := s.service.FindMuseumItemMovement(id)
	if err != nil {
		log.Printf("Failed to find item with details: %s", err)
		return err
	}

	_ = s.tmpl.ExecuteTemplate(c.Response(), "ShowMuseumItemMovement", item)
	return nil
}

func (s *Server) createMuseumItemMovement(c echo.Context) error {
	_, err := s.service.CreateMuseumItemMovement(getMovementFromForm(c))
	if err != nil {
		log.Print(err)
		return err
	}
	return c.Redirect(301, "/museumItemMovements")
}

func (s *Server) getMuseumItemMovementPage(c echo.Context) error {
	_ = s.tmpl.ExecuteTemplate(c.Response().Writer, "New museum item movement", nil)
	return nil
}

func getMovementFromForm(c echo.Context) service.MuseumItemMovement {
	var movement service.MuseumItemMovement
	movement.AcceptDate = getParsedTime(c.FormValue("accept_date"))
	movement.ExhibitTransferDate = getParsedTime(c.FormValue("exhibit_transfer_date"))
	movement.ExhibitReturnDate = getParsedTime(c.FormValue("exhibit_return_date"))
	movement.ResponsiblePerson = getPersonFromForm(c)
	movement.Item.Name = c.FormValue("item_name")
	return movement
}

func getParsedTime(t string) *time.Time {
	if t == "" {
		return nil
	}
	t = t + ":00"
	parsed, err := time.Parse("2006-01-02T15:04:05", t)
	log.Print(err)
	return &parsed
}
