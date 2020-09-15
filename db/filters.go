package db

import (
	"context"
	"fmt"
	"regexp"
	"strings"

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
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error fetching count of getFilterRecordCount from database")
			return
		}
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

// @Title SearchRecords
// @Description Get records that are searched as per request Parameter "text" along with count
// @Accept	request.Context, text as string, limit, page
// @Success total= (count of search qualifying records), error=nil
// @Failure total=0, error= "Some Error"
func (s *pgStore) SearchRecords(ctx context.Context, text string, limit string, page string) (total int, products []Product, err error) {
	// check for SQL Injection
	// Only allow words characters like [a-z0-9A-Z] and a space [ ]
	var validParameters = regexp.MustCompile(`^[\w ]+$`)
	// if There are other chracters than word and space
	if validParameters.MatchString(text) == false {
		err = fmt.Errorf("Possible SQL Injection Attack.")
		logger.WithField("err", err.Error()).Error("Error In Parameters, special Characters are present.")
		return
	}

	// Split the text into slice of strings, max 10 first words will be considered
	textSlice := strings.SplitN(text, " ", 11)

	// TODO eliminate duplicate textSlice entries

	// If there are more than 10 words in search, ask user to be less verbose
	if len(textSlice) > 10 {
		err = fmt.Errorf("Unnecessary detailed text given.")
		logger.WithField("err", err.Error()).Error("Error In Parameters, very detailed!.")
		return
	}

	// Query to return Id's of Products where we may find a match in
	// product's name, description, brand, size, color or in
	// the category of that products category's name or description
	getSearchRecordIds := `SELECT p.id from products p
		INNER JOIN category c 
		ON p.category_id = c.id
		WHERE 
		`

	// Query to help us get count of all such results
	getSearchCount := `SELECT COUNT(p.id) from products p
		INNER JOIN category c 
		ON p.category_id = c.id
		WHERE 
		`

	helper := `  `

	// iterate over all the textSlice
	for i := 0; i < len(textSlice); i++ {
		helper += ` 
		LOWER(p.name) LIKE LOWER('%` + textSlice[i] + `%') OR 
		LOWER(p.description) LIKE LOWER('%` + textSlice[i] + `%') OR
		LOWER(p.brand) LIKE LOWER('%` + textSlice[i] + `%') OR 
		LOWER(p.color) LIKE LOWER('%` + textSlice[i] + `%') OR 
		LOWER(p.size) LIKE LOWER('%` + textSlice[i] + `%') OR 
		LOWER(c.name) LIKE LOWER('%` + textSlice[i] + `%') OR 
		LOWER(c.description) LIKE LOWER('%` + textSlice[i] + `%') OR`
	}

	// remove that last OR
	helper = helper[:len(helper)-2]

	getSearchRecordIds += helper
	getSearchCount += helper + ` ;`

	getSearchRecordIds += ` LIMIT ` + string(limit) + ` OFFSET  ( ` + string(page) + ` -1) * ` + string(limit) + ` ;`

	fmt.Println("getSearchRecordIds---->", getSearchRecordIds)

	countResult, err := s.db.Query(getSearchCount)

	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching count of getSearchCount from database")
		return
	}

	// countResult set should have only 1 record
	// It counts the number of records with the search results.
	for countResult.Next() {
		err = countResult.Scan(&total)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error fetching count of getSearchCount from database")
			return
		}
		break
	}

	fmt.Println(total)

	result, err := s.db.Query(getSearchRecordIds)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching Product Results from database")
		return
	}

	// idArr stores id's of all products
	var idArr []int

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
