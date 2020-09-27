package api

import (
	"log"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"pbz2/pkg/entities"
)

func (s *Server) initMuseumItemAPI(e *echo.Echo) {
	e.POST("/museumItem", s.createMuseumItem)
	e.GET("/museumItem", s.getMuseumItemPage)
	e.GET("/museumItem/:id", s.getMuseumItem)
	e.GET("/museumItems", s.getMuseumItems)
	e.GET("/deleteMuseumItem/:id", s.deleteMuseumItem)
	e.GET("/editMuseumItem/:id", s.editMuseumItem)
	e.POST("/editMuseumItem/:id", s.updateMuseumItem)

	e.GET("/museumItemSearch", s.getMuseumItemSearchPage)
	e.POST("/museumItemSearch", s.searchMuseumItems)
}

func (s *Server) createMuseumItem(c echo.Context) error {
	_, err := s.service.CreateMuseumItem(getItemFromForm(c))
	if err != nil {
		log.Print(err)
		return err
	}
	return c.Redirect(301, "/museumItems")
}

func (s *Server) getMuseumItem(c echo.Context) error {
	id := getIDFromURL(c)
	item, err := s.service.GetMuseumItemWithDetails(id)
	if err != nil {
		log.Printf("Failed to find item with details: %s", err)
		return err
	}

	_ = s.tmpl.ExecuteTemplate(c.Response(), "ShowMuseumItem", item)
	return nil
}

func (s *Server) editMuseumItem(c echo.Context) error {
	id := c.Param("id")
	parsedID, _ := strconv.ParseInt(id, 10, 64)
	item, err := s.service.GetMuseumItem(int(parsedID))
	if err != nil {
		log.Print(err)
		return err
	}
	_ = s.tmpl.ExecuteTemplate(c.Response(), "EditMuseumItem", item)
	return nil
}

func (s *Server) updateMuseumItem(c echo.Context) error {
	id := getIDFromURL(c)
	item := getItemFromForm(c)
	item.ID = id
	err := s.service.UpdateMuseumItem(item.MuseumItem)
	if err != nil {
		return err
	}
	return c.Redirect(301, "/museumItems")
}

func (s *Server) deleteMuseumItem(c echo.Context) error {
	id := getIDFromURL(c)
	err := s.service.DeleteMuseumItem(id)
	if err != nil {
		log.Print(err)
		return err
	}
	return c.Redirect(301, "/museumItems")
}

func (s *Server) getMuseumItemPage(c echo.Context) error {
	_ = s.tmpl.ExecuteTemplate(c.Response().Writer, "New museum item", nil)
	return nil
}

func (s *Server) getMuseumItemSearchPage(c echo.Context) error {
	_ = s.tmpl.ExecuteTemplate(c.Response().Writer, "Search museum items", nil)
	return nil
}

func (s *Server) searchMuseumItems(c echo.Context) error {
	args := getSearchArgsFromForm(c)
	items, err := s.service.FindMuseumItems(args)
	if err != nil {
		log.Printf("Failed to find museum item: %+v", err)
		return err
	}
	_ = s.tmpl.ExecuteTemplate(c.Response().Writer, "Museum items", items)
	return nil
}

func (s *Server) getMuseumItems(c echo.Context) error {
	items, err := s.service.FindMuseumItems(entities.SearchMuseumItemsArgs{})
	if err != nil {
		log.Printf("Failed to find museum item: %+v", err)
		return err
	}
	_ = s.tmpl.ExecuteTemplate(c.Response().Writer, "Museum items", items)
	return nil
}

func getIDFromURL(c echo.Context) int {
	id := c.Param("id")
	parsedID, _ := strconv.Atoi(id)
	return parsedID
}

func getItemFromForm(c echo.Context) entities.MuseumItemWithDetails {
	var item entities.MuseumItemWithDetails
	item.Name = c.FormValue("item_name")
	item.Annotation = c.FormValue("item_annotation")
	item.InventoryNumber = c.FormValue("inventory_number")

	creationDate, _ := time.Parse("2006-01-02", c.FormValue("creation_date"))
	item.CreationDate = entities.NewDate(creationDate)

	item.Keeper = getPersonFromForm(c)
	item.Fund = getFundFromForm(c)
	item.Set = getSetFromForm(c)
	return item
}

func getSearchArgsFromForm(c echo.Context) entities.SearchMuseumItemsArgs {
	var args entities.SearchMuseumItemsArgs
	set := getSetFromForm(c)
	if set.Name != "" {
		args.SetName = &set.Name
	}
	item := getItemFromForm(c)
	if item.Name != "" {
		args.ItemName = &item.Name
	}
	date, _ := time.Parse("2006-01-02", c.FormValue("date"))
	if !date.IsZero() {
		args.Date = &date
	}
	return args
}
func getPersonFromForm(c echo.Context) entities.Person {
	var person entities.Person
	person.FirstName = c.FormValue("first_name")
	person.LastName = c.FormValue("second_name")
	person.MiddleName = c.FormValue("middle_name")
	return person
}

func getFundFromForm(c echo.Context) entities.MuseumFund {
	var fund entities.MuseumFund
	fund.Name = c.FormValue("fund_name")
	return fund
}

func getSetFromForm(c echo.Context) entities.MuseumSet {
	var set entities.MuseumSet
	set.Name = c.FormValue("set_name")
	return set
}
