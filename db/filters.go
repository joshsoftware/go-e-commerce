package db

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
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

// @Title FilteredProducts
// @Description Get the products that are filtered as per request Parameters
// @Accept	request.Context, Filter struct's object
// @Success total= (count of filtered products), error=nil
// @Failure total=0, error= "Some Error"
func (s *pgStore) FilteredProducts(ctx context.Context, filter Filter, limit string, page string) (count int, products []Product, err error) {
	// We will be checking for SQL Injection as well in this Method only
	// found flag will help us find out if any of Filter flags were true
	var found bool
	// helper will be used in making query dynamic.
	// See how it's getting concatanation added in case a flag was Filter Flag is true
	injection := `  `
	helper := `  `
	if filter.CategoryFlag == true {
		helper += ` category_id = ` + filter.CategoryId + ` AND`
		injection += filter.CategoryId
		found = true
	}
	if filter.BrandFlag == true {
		// Since ' existed, we had to use ` instead of " , as compiler gave error otherwise
		helper += ` LOWER(brand) = LOWER('` + filter.Brand + `') AND`
		injection += filter.Brand
		found = true
	}
	if filter.SizeFlag == true {
		// Since ' existed, we had to use ` instead of " , as compiler gave error otherwise
		helper += ` LOWER(size) = LOWER('` + filter.Size + `') AND`
		injection += filter.Size
		found = true
	}
	if filter.ColorFlag == true {
		// Since ' existed, we had to use ` instead of " , as compiler gave error otherwise
		helper += ` LOWER(color) =LOWER('` + filter.Color + `') AND`
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
		helper = ` WHERE ` + helper[:len(helper)-3]
	}

	getFilterProductCount := `SELECT COUNT(id) FROM products ` + helper + `;`

	resultCount, err := s.db.Query(getFilterProductCount)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error getting Count of Filtered Products from database")
		return
	}

	// resultCount set should have only 1 record
	for resultCount.Next() {
		err = resultCount.Scan(&count)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error fetching count of getFilterProductCount from database")
			return
		}
		break
	}

	if count == 0 {
		err = fmt.Errorf("No records present")
		logger.WithField("err", err.Error()).Error("No records were in db for Products")
		return
	}

	// error already handled in filters_http
	ls, _ := strconv.Atoi(limit)
	ps, _ := strconv.Atoi(page)

	if (count - 1) < (int(ls) * (int(ps) - 1)) {
		err = fmt.Errorf("Desired Page not found")
		logger.WithField("err", err.Error()).Error("Page Out Of range")
		return
	}

	getFilterProduct := `SELECT id from Products` + helper

	if filter.PriceFlag == true {
		getFilterProduct += ` ORDER BY price ` + filter.Price
	}

	//fmt.Println(limit, page)
	getFilterProduct += ` LIMIT ` + limit + `  OFFSET  (` + page + ` -1) * ` + limit + ` ;`

	result, err := s.db.Query(getFilterProduct)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching Product Ids from database")
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

// @Title SearchRecords
// @Description Get records that are searched as per request Parameter "text" along with count
// @Accept	request.Context, text as string, limit, page
// @Success total= (count of search qualifying records), error=nil
// @Failure total=0, error= "Some Error"
func (s *pgStore) SearchRecords(ctx context.Context, text string, limit string, page string) (count int, products []Product, err error) {
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

	// If there are more than 10 words in search, ask user to be less verbose
	if len(textSlice) > 10 {
		err = fmt.Errorf("Unnecessary detailed text given.")
		logger.WithField("err", err.Error()).Error("Error In Parameters, very detailed!.")
		return
	}

	// Removing Duplicate words from textSlice
	textMap := make(map[string]bool, 10)
	for i := 0; i < len(textSlice); i++ {
		textMap[textSlice[i]] = true
	}

	// Query to help us get count of all such results
	getSearchCount := `SELECT COUNT(p.id) from products p
		INNER JOIN category c 
		ON p.category_id = c.id
		WHERE 
		`

	helper := `  `

	// iterate over all the textMap
	for key, _ := range textMap {
		helper += ` 
		LOWER(p.name) LIKE LOWER('%` + key + `%') OR 
		LOWER(p.brand) LIKE LOWER('%` + key + `%') OR 
		LOWER(c.name) LIKE LOWER('%` + key + `%') OR`
	}

	// remove that last OR
	helper = helper[:len(helper)-2]

	getSearchCount += helper + ` ;`
	countResult, err := s.db.Query(getSearchCount)

	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching count of getSearchCount from database")
		return
	}

	// countResult set should have only 1 record
	// It counts the number of records with the search results.
	for countResult.Next() {
		err = countResult.Scan(&count)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error fetching count of getSearchCount from database")
			return
		}
		break
	}

	fmt.Println(count)

	if count == 0 {
		err = fmt.Errorf("No records  present")
		logger.WithField("err", err.Error()).Error("No records were present for that search keyword")
		return
	}

	// error already handled in filters_http
	ls, _ := strconv.Atoi(limit)
	ps, _ := strconv.Atoi(page)

	if (count - 1) < (int(ls) * (int(ps) - 1)) {
		err = fmt.Errorf("Desired Page not found")
		logger.WithField("err", err.Error()).Error("Page Out Of range")
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

	getSearchRecordIds += helper

	getSearchRecordIds += ` LIMIT ` + limit + ` OFFSET  ( ` + page + ` -1) * ` + limit + ` ;`

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
