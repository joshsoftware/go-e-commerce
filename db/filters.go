package db

import (
	"context"
	"fmt"

	logger "github.com/sirupsen/logrus"
)

//"database/sql"

//const (
//getProductByFilterQuery=
//)

type Filter struct {
	CategoryId string // := req.URL.Query().Get("categoryid")
	Price      string //:= req.URL.Query().Get("price")
	Brand      string //:= req.URL.Query().Get("brand")
	Size       string //:= req.URL.Query().Get("size")
	Color      string //:= req.URL.Query().Get("color")

	CategoryFlag bool
	PriceFlag    bool
	BrandFlag    bool
	SizeFlag     bool
	ColorFlag    bool
}

func (s *pgStore) FilteredRecordsCount(ctx context.Context, filter Filter) (total int, err error) {
	var found bool

	helper := " "
	if filter.CategoryFlag == true {
		helper += " category_id = " + string(filter.CategoryId) + " ,"
		found = true
	}
	if filter.BrandFlag == true {
		helper += ` brand = '` + string(filter.Brand) + `',`
		found = true
	}
	if filter.SizeFlag == true {
		helper += ` size ='` + string(filter.Size) + `',`
		found = true
	}
	if filter.ColorFlag == true {
		helper += ` color ='` + string(filter.Color) + `',`
		found = true
	}
	if found == true {
		// remove that last comma
		helper = " WHERE" + helper[:len(helper)-1]
	}
	helper += " ;"

	getFilterRecord := `SELECT COUNT(id) FROM products `
	getFilterRecord += string(helper)

	fmt.Println(getFilterRecord)

	result, err := s.db.Query(getFilterRecord)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching Product Ids from database")
		return
	}
	var number int
	for result.Next() {
		err = result.Scan(&number)
	}
	fmt.Println(number)
	return number, err
}

func (s *pgStore) FilteredRecords(ctx context.Context, filter Filter, limit string, page string) (products []Product, err error) {

	var found bool

	helper := " "
	if filter.CategoryFlag == true {
		helper += " category_id = " + string(filter.CategoryId) + " ,"
		found = true
	}
	if filter.BrandFlag == true {
		helper += " brand = '" + filter.Brand + "' ,"
		found = true
	}
	if filter.SizeFlag == true {
		helper += " size = '" + string(filter.Size) + "' ,"
		found = true
	}
	if filter.ColorFlag == true {
		helper += " color = '" + string(filter.Color) + "' ,"
		found = true
	}
	if found == true {
		// remove that last comma
		helper = " WHERE" + helper[:len(helper)-1]
	}

	getFilterRecord := "SELECT id from Products" + helper

	if filter.PriceFlag == true {
		getFilterRecord += " ORDER BY price " + string(filter.Price)
	}

	//fmt.Println(limit, page)
	getFilterRecord += " LIMIT " + string(limit) + "  OFFSET  (" + string(page) + " -1) * " + string(limit) + " ;"
	fmt.Println(getFilterRecord)

	// idArr stores id's of all products
	var idArr []int

	result, err := s.db.Query(getFilterRecord)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching Product Ids from database")
		return
	}

	for result.Next() {
		var Id int
		err = result.Scan(&Id)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Couldn't Scan Resulted Product Ids into Id variable")
			return
		}
		idArr = append(idArr, Id)
	}

	for i := 0; i < len(idArr); i++ {
		var product Product
		product, err = s.GetProductByID(ctx, int(idArr[i]))
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error selecting Product from database by id " + string(idArr[i]))
			return
		}
		products = append(products, product)
	}

	return

}
