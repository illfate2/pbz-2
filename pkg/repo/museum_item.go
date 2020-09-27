package repo

import (
	"context"
	"strconv"

	"pbz2/pkg/entities"
)

func (s *Repo) FindMuseumItem(id int) (entities.MuseumItem, error) {
	var item entities.MuseumItem
	err := s.conn.QueryRow(context.Background(),
		`SELECT
		id,name,creation_date,annotation,set_id,fund_id,
		keeper_id,inventory_number
		FROM museum_items WHERE id = $1`, id).
		Scan(
			&item.ID, &item.Name,
			&item.CreationDate.Time,
			&item.Annotation,
			&item.MuseumSetID, &item.MuseumFundID, &item.KeeperID, &item.InventoryNumber,
		)
	if err != nil {
		return entities.MuseumItem{}, err
	}
	return item, nil
}

func (s *Repo) FindMuseumItemByName(name string) (entities.MuseumItem, error) {
	var item entities.MuseumItem
	err := s.conn.QueryRow(context.Background(),
		`SELECT
		id,name,creation_date,annotation,set_id,fund_id,
		keeper_id,inventory_number
		FROM museum_items WHERE name = $1`, name).
		Scan(
			&item.ID, &item.Name,
			&item.CreationDate.Time,
			&item.Annotation,
			&item.MuseumSetID, &item.MuseumFundID, &item.KeeperID, &item.InventoryNumber,
		)
	if err != nil {
		return entities.MuseumItem{}, err
	}
	return item, nil
}

func (s *Repo) FindMuseumItemWithDetails(id int) (entities.MuseumItemWithDetails, error) {
	var item entities.MuseumItemWithDetails
	err := s.conn.QueryRow(context.Background(),
		`SELECT
		mi.id ,mi.name,mi.creation_date,mi.annotation,mi.set_id,mi.fund_id,
		mi.keeper_id,mi.inventory_number,p.first_name,p.second_name,p.middle_name,
		mis.name, mf.name
		FROM museum_items mi
		JOIN persons p ON mi.keeper_id=p.id
		JOIN museum_item_sets mis ON mi.set_id = mis.id
		JOIN museum_funds mf ON mi.fund_id = mf.id
		WHERE mi.id = $1`, id).
		Scan(
			&item.ID, &item.Name,
			&item.CreationDate.Time,
			&item.Annotation,
			&item.MuseumSetID, &item.MuseumFundID, &item.KeeperID, &item.InventoryNumber,
			&item.Keeper.FirstName,
			&item.Keeper.LastName,
			&item.Keeper.MiddleName,
			&item.Set.Name,
			&item.Fund.Name,
		)
	if err != nil {
		return entities.MuseumItemWithDetails{}, err
	}
	return item, nil
}

type queryBuilder struct {
	selects string
	from    string
	wheres  map[string][]interface{}
	joins   map[string]struct{}
	args    []interface{}
}

func newQueryBuilder() *queryBuilder {
	return &queryBuilder{
		wheres: make(map[string][]interface{}),
		joins:  make(map[string]struct{}),
	}
}

func (b *queryBuilder) withSelect(s string) *queryBuilder {
	b.selects = s
	return b
}

func (b *queryBuilder) withFrom(s string) *queryBuilder {
	b.from = s
	return b
}

func (b *queryBuilder) withWhere(column string, values ...interface{}) *queryBuilder {
	b.wheres[column] = append(b.wheres[column], values...)
	return b
}

func (b *queryBuilder) withJoin(s string) *queryBuilder {
	b.joins[s] = struct{}{}
	return b
}

func (b *queryBuilder) buildQuery() (string, []interface{}) {
	res := b.selects + "\n"
	var joins string
	for j := range b.joins {
		joins += j + "\n"
	}
	res += joins
	var wheres string
	if len(b.wheres) != 0 {
		wheres += "WHERE "
	}
	counter := 0
	for statement, v := range b.wheres {
		b.args = append(b.args, v...)
		if counter == 0 {
			wheres += statement + " = $" + strconv.Itoa(counter+1) + "\n"
			continue
		}
		wheres += "AND " + statement + " = $" + strconv.Itoa(counter+1) + "\n"
		counter++
	}
	res += wheres
	return res, b.args
}

func (s *Repo) FindMuseumItems(args entities.SearchMuseumItemsArgs) ([]entities.MuseumItem, error) {
	selects := `SELECT 
			mi.id,mi.name,mi.creation_date,mi.annotation,mi.set_id,mi.fund_id,
			mi.keeper_id,mi.inventory_number		
			FROM museum_items mi`
	q := newQueryBuilder()
	q.withSelect(selects)
	if args.SetName != nil {
		q.withJoin("JOIN museum_item_sets mis ON mis.id=mi.set_id").
			withWhere("mis.name", *args.SetName)
	}
	if args.ItemName != nil {
		q.withWhere("mi.name", *args.ItemName)
	}
	if args.Date != nil {
		q.withJoin("JOIN museum_item_movement mim ON mim.item_id = mi.id").
			withWhere("mim.exhibit_transfer_date >= ? AND mim.exhibit_return_date <= ?",
				*args.Date, *args.Date)
	}
	queryStr, queryArgs := q.buildQuery()
	rows, err := s.conn.Query(context.Background(), queryStr, queryArgs...)

	if err != nil {
		return nil, err
	}
	var items []entities.MuseumItem
	for rows.Next() {
		var item entities.MuseumItem
		err := rows.Scan(
			&item.ID, &item.Name,
			&item.CreationDate.Time,
			&item.Annotation,
			&item.MuseumSetID, &item.MuseumFundID, &item.KeeperID, &item.InventoryNumber,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *Repo) InsertMuseumItem(item entities.MuseumItem) (entities.MuseumItem, error) {
	err := s.conn.QueryRow(context.Background(),
		`INSERT INTO museum_items(name,creation_date,annotation,inventory_number,keeper_id,set_id,fund_id)
		VALUES($1,$2,$3,$4,$5,$6,$7) RETURNING id`,
		item.Name, item.CreationDate.Time, item.Annotation, item.InventoryNumber, item.KeeperID, item.MuseumSetID, item.MuseumFundID).
		Scan(&item.ID)
	if err != nil {
		return entities.MuseumItem{}, err
	}
	return item, nil
}

func (s *Repo) UpdateMuseumItem(item entities.MuseumItem) error {
	_, err := s.conn.Exec(context.Background(),
		`UPDATE museum_items 
			SET name = $1, creation_date = $2, annotation = $3
			WHERE id = $4`,
		item.Name, item.CreationDate.Time, item.Annotation, item.ID)
	return err
}

func (s *Repo) DeleteMuseumItem(id int) error {
	_, err := s.conn.Exec(context.Background(),
		`DELETE FROM museum_items WHERE id = $1`, id)
	return err
}
