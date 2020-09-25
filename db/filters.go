package db

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	logger "github.com/sirupsen/logrus"
)

var (
	filterProductCount = `SELECT
		COUNT(p.id)
		FROM   products p
		INNER JOIN category c
		ON p.cid = c.cid`

	filterProduct = `SELECT *
		FROM   products p
		INNER JOIN category c
		ON p.cid = c.cid`
)

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
func (s *pgStore) FilteredProducts(ctx context.Context, filter Filter, limitStr string, offsetStr string) (int, []Product, error) {

	var found bool
	var totalRecords int
	var products []Product

	// helper will be used in making query dynamic.
	// See how it's getting concatanation added in case a flag was Filter Flag is true
	sqlRegexp := ``
	isFiltered := `   `
	if filter.CategoryFlag == true {
		isFiltered += ` c.cid = ` + filter.CategoryId + ` AND`
		sqlRegexp += filter.CategoryId
		found = true
	}
	if filter.BrandFlag {
		isFiltered += ` LOWER(p.brand) = LOWER('` + filter.Brand + `') AND`
		sqlRegexp += filter.Brand
		found = true
	}
	if filter.SizeFlag {
		isFiltered += ` LOWER(p.size) = LOWER('` + filter.Size + `') AND`
		sqlRegexp += filter.Size
		found = true
	}
	if filter.ColorFlag {
		isFiltered += ` LOWER(p.color) =LOWER('` + filter.Color + `') AND`
		sqlRegexp += filter.Color
		found = true
	}
	if found {
		var validParameters = regexp.MustCompile(`^[\w ]+$`)
		if !validParameters.MatchString(sqlRegexp) {
			err := fmt.Errorf("Possible SQL Injection Attack.")
			logger.WithField("err", err.Error()).Error("Error In Parameters, special Characters are present.")
			return 0, []Product{}, err
		}
		isFiltered = ` WHERE ` + isFiltered[:len(isFiltered)-3]
	}

	getFilterProductCount := filterProductCount + isFiltered + `;`
	resultCount, err := s.db.Query(getFilterProductCount)

	if err != nil {
		logger.WithField("err", err.Error()).Error("Error getting Count of Filtered Products from database")
		return 0, []Product{}, err
	}

	for resultCount.Next() {
		err = resultCount.Scan(&totalRecords)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error fetching count of getFilterProductCount from database")
			return 0, []Product{}, err
		}
		break
	}

	offset, _ := strconv.Atoi(offsetStr)

	if totalRecords-1 < offset {
		err = fmt.Errorf("Page out of Range!")
		logger.WithField("err", err.Error()).Error("Error Offset is greater than total records")
		return 0, []Product{}, err

	}

	getFilterProduct := filterProduct + isFiltered

	if filter.PriceFlag {
		getFilterProduct += ` ORDER BY p.price ` + filter.Price + `, p.id LIMIT ` + limitStr + `  OFFSET  ` + offsetStr + `  ;`
	} else {
		getFilterProduct += ` ORDER BY p.id LIMIT ` + limitStr + `  OFFSET  ` + offsetStr + `  ;`
	}
	fmt.Println(getFilterProduct)

	err = s.db.Select(&products, getFilterProduct)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching Product Ids from database")
		return 0, []Product{}, err
	}
	if products == nil {
		err = fmt.Errorf("Desired page not found")
		logger.WithField("err", err.Error()).Error("Products don't exist by such filters!")
		return 0, []Product{}, err
	}

	return totalRecords, products, nil
}

// @Title SearchRecords
// @Description Get records that are searched as per request Parameter "text" along with count
// @Accept	request.Context, text as string, limitStr, pageStr
// @Success total= (count of search qualifying records), error=nil
// @Failure total=0, error= "Some Error"
func (s *pgStore) SearchProductsByText(ctx context.Context, text string, limitStr string, offsetStr string) (int, []Product, error) {

	var totalRecords int
	var products []Product
	var validParameters = regexp.MustCompile(`^[\w ]+$`)

	// if There are other chracters than word and space
	if validParameters.MatchString(text) == false {
		err := fmt.Errorf("Possible SQL Injection Attack.")
		logger.WithField("err", err.Error()).Error("Error In Parameters, special Characters are present.")
		return 0, []Product{}, err
	}

	// Split the text into slice of strings, max 10 first words will be considered
	textSlice := strings.SplitN(text, " ", 11)

	// If there are more than 10 words in search, ask user to be less verbose
	if len(textSlice) > 10 {
		err := fmt.Errorf("Unnecessary detailed text given.")
		logger.WithField("err", err.Error()).Error("Error In Parameters, very detailed!.")
		return 0, []Product{}, err
	}

	// Removing Duplicate words from textSlice
	textMap := make(map[string]bool, 10)
	for i := 0; i < len(textSlice); i++ {
		textMap[textSlice[i]] = true
	}

	// Query to help us get count of all such results
	getSearchCount := `SELECT COUNT(p.id) from products p
		INNER JOIN category c 
		ON p.cid = c.cid
		WHERE `

	isFiltered := `  `

	// iterate over all the textMap
	for key, _ := range textMap {
		isFiltered += ` 
		LOWER(p.name) LIKE LOWER('%` + key + `%') OR 
		LOWER(p.brand) LIKE LOWER('%` + key + `%') OR 
		LOWER(c.cname) LIKE LOWER('%` + key + `%') OR`
	}

	// remove that last OR
	isFiltered = isFiltered[:len(isFiltered)-2]

	getSearchCount += isFiltered + ` ;`
	fmt.Println("getsearchProduct---->", getSearchCount)

	countResult, err := s.db.Query(getSearchCount)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching count of getSearchCount from database")
		return 0, []Product{}, err
	}

	// countResult set should have only 1 record
	// It counts the number of records with the search results.
	for countResult.Next() {
		err = countResult.Scan(&totalRecords)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error fetching count of getSearchCount from database")
			return 0, []Product{}, err
		}
		break
	}

	offset, _ := strconv.Atoi(offsetStr)

	if totalRecords-1 < offset {
		err = fmt.Errorf("Page out of Range!")
		logger.WithField("err", err.Error()).Error("Error Offset is greater than total records")
		return 0, []Product{}, err

	}

	getSearchRecord := filterProduct + ` WHERE `
	getSearchRecord += isFiltered
	getSearchRecord += ` LIMIT ` + limitStr + ` OFFSET  ` + offsetStr + ` ;`

	fmt.Println("getsearchRecords---->", getSearchRecord)

	err = s.db.Select(&products, getSearchRecord)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching Product Results from database")
		return 0, []Product{}, err
	}

	if products == nil {
		err = fmt.Errorf("Desired page not found")
		logger.WithField("err", err.Error()).Error("page Out Of range")
		return 0, []Product{}, err
	}
	return totalRecords, products, nil
}
