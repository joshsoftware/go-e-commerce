package db

import (
	"context"
	"fmt"
	"regexp"

	logger "github.com/sirupsen/logrus"
)

// Filter struct is used to help us in Generating a dynamic Filter Query
type Filter struct {
	// Below fields are what we may receive as Parameters in request body
	CategoryId string
	Price      string
	Brand      string
	Size       string
	Color      string
	// These Flags will help us format our query, true means that field exists in Request Parameters
	CategoryFlag bool
	PriceFlag    bool
	BrandFlag    bool
	SizeFlag     bool
	ColorFlag    bool
}

// @Title FilteredRecordsCount
// @Description Get Count of records that are filtered as per request Parameters
// @Accept	request.Context, Filter struct's object
// @Success total= (count of filtered records), error=nil
// @Failure total=0, error= "Some Error"
func (s *pgStore) FilteredRecordsCount(ctx context.Context, filter Filter) (total int, err error) {
	// We will be checking for SQL Injection as well in this Method only
	// found flag will help us find out if any of Filter flags were true
	var found bool
	// helper will be used in making query dynamic.
	// See how it's getting concatanation added in case a flag was Filter Flag is true
	injection := " "
	helper := " "
	if filter.CategoryFlag == true {
		helper += " category_id = " + string(filter.CategoryId) + "' AND"
		injection += filter.CategoryId
		found = true
	}
	if filter.BrandFlag == true {
		// Since ' existed, we had to use ` instead of " , as compiler gave error otherwise
		helper += ` brand = '` + string(filter.Brand) + `' AND`
		injection += filter.Brand
		found = true
	}
	if filter.SizeFlag == true {
		// Since ' existed, we had to use ` instead of " , as compiler gave error otherwise
		helper += ` size ='` + string(filter.Size) + `' AND`
		injection += filter.Size
		found = true
	}
	if filter.ColorFlag == true {
		// Since ' existed, we had to use ` instead of " , as compiler gave error otherwise
		helper += ` color ='` + string(filter.Color) + `' AND`
		injection += filter.Color
		found = true
	}
	if found == true {
		// check for SQL Injection
		// Only allow words characters like [a-z0-9A-Z] and a space [ ]
		var validParameters = regexp.MustCompile(`^[\w ]+$`)
		// if There are other chracters than word and space
		if validParameters.MatchString(injection) == false {
			err = fmt.Errorf("Possible SQL Injection Attack.")
			logger.WithField("err", err.Error()).Error("Error In Parameters, special Characters are present.")
			return
		}
		// remove that last AND as it will make query invalid
		helper = " WHERE" + helper[:len(helper)-3]
	}
	// Ending the Query in a safe way
	helper += " ;"

	getFilterRecordCount := `SELECT COUNT(id) FROM products `
	getFilterRecordCount += string(helper)
	fmt.Println("getFilterRecordCount---->", getFilterRecordCount)

	result, err := s.db.Query(getFilterRecordCount)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching Product Ids from database")
		return
	}

	// result set should have only 1 record
	for result.Next() {
		err = result.Scan(&total)
		break
	}
	return
}

// @Title FilteredRecords
// @Description Get the records that are filtered as per request Parameters
// @Accept	request.Context, Filter struct's object
// @Success total= (count of filtered records), error=nil
// @Failure total=0, error= "Some Error"
func (s *pgStore) FilteredRecords(ctx context.Context, filter Filter, limit string, page string) (products []Product, err error) {
	// found flag will help us find out if any of Filter flags were true
	var found bool

	// helper will be used in making query dynamic.
	// See how it's getting concatanation added in case a flag was Filter Flag is true
	helper := " "
	if filter.CategoryFlag == true {
		helper += " category_id = " + string(filter.CategoryId) + " AND"
		found = true
	}
	if filter.BrandFlag == true {
		// Since ' existed, we had to use ` instead of " , as compiler gave error otherwise
		helper += " brand = '" + filter.Brand + "' AND"
		found = true
	}
	if filter.SizeFlag == true {
		// Since ' existed, we had to use ` instead of " , as compiler gave error otherwise
		helper += " size = '" + string(filter.Size) + "' AND"
		found = true
	}
	if filter.ColorFlag == true {
		// Since ' existed, we had to use ` instead of " , as compiler gave error otherwise
		helper += " color = '" + string(filter.Color) + "' AND"
		found = true
	}
	if found == true {
		// remove that last comma as it will make query invalid
		helper = " WHERE" + helper[:len(helper)-3]
	}

	getFilterRecord := "SELECT id from Products" + helper

	if filter.PriceFlag == true {
		getFilterRecord += " ORDER BY price " + string(filter.Price)
	}

	//fmt.Println(limit, page)
	getFilterRecord += " LIMIT " + string(limit) + "  OFFSET  (" + string(page) + " -1) * " + string(limit) + " ;"
	fmt.Println("getFilterRecord---->", getFilterRecord)

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

	// get All Filtered Products by their ids
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
