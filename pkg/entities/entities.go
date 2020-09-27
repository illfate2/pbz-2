package entities

import "time"

type MuseumExhibition struct {
	ID      int
	Name    string
	Address string
	PhoneNo string
}

type MuseumFund struct {
	ID   int
	Name string
}

type MuseumItemMovement struct {
	ID                  int
	MuseumItemID        int
	Item                MuseumSet
	AcceptDate          *time.Time
	ExhibitTransferDate *time.Time
	ExhibitReturnDate   *time.Time
	MuseumExhibitionID  int
	Exhibition          MuseumExhibition
	ResponsiblePersonID int
	ResponsiblePerson   Person
}

type MuseumItem struct {
	ID              int
	InventoryNumber string
	Name            string
	CreationDate    Date
	KeeperID        int
	MuseumSetID     int
	MuseumFundID    int
	Annotation      string
}

type MuseumItemWithKeeper struct {
	MuseumItem
	Keeper Person
}

type SearchMuseumItemsArgs struct {
	ItemName *string
	SetName  *string
	Date     *time.Time
}

type MuseumItemWithDetails struct {
	MuseumItem
	Keeper Person
	Set    MuseumSet
	Fund   MuseumFund
}

type Person struct {
	ID         int
	FirstName  string
	LastName   string
	MiddleName string
}

type MuseumSet struct {
	ID   int
	Name string
}

type MuseumSetWithDetails struct {
	MuseumSet
	Items []MuseumItemWithKeeper
}

type Date struct {
	Time time.Time
}

func (d Date) String() string {
	return d.Time.Format("2006-01-02")
}

func NewDate(time time.Time) Date {
	return Date{
		Time: roundToDay(time),
	}
}

func roundToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
