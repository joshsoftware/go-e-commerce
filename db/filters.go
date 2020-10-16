package db

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	logger "github.com/sirupsen/logrus"
)

var (
	filterSearchProduct = `SELECT count(*) OVER() AS total,*
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

// TODO Add condition and sort capabilities for both these APIs
// These capabilities are suppossed to make filter and search APIs
// really really dynamic and much robust.
// eg for condition -> WHERE cid >= 5 AND tax <= 4
// eg for sort ->  category_id = desc, price asc

// @Title FilteredProducts
// @Description Get the products that are filtered as per request Parameters
// @Accept	request.Context, Filter struct's object
// @Success total= (count of filtered products), error=nil
// @Failure total=0, error= "Some Error"
func (s *pgStore) FilteredProducts(ctx context.Context, filter Filter, limitStr string, offsetStr string) (int, []Product, error) {

	var found bool
	var records Records
	var err error

	// helper will be used in making query dynamic.
	// See how it's getting concatanation added in case a flag was Filter Flag is true
	sqlRegexp := ``
	filterQuery := `   `
	if filter.CategoryFlag == true {
		filterQuery += ` c.cid = ` + filter.CategoryId + ` AND`
		sqlRegexp += filter.CategoryId
		found = true
	}
	if filter.BrandFlag {
		filterQuery += ` LOWER(p.brand) = LOWER('` + filter.Brand + `') AND`
		sqlRegexp += filter.Brand
		found = true
	}
	if filter.SizeFlag {
		filterQuery += ` LOWER(p.size) = LOWER('` + filter.Size + `') AND`
		sqlRegexp += filter.Size
		found = true
	}
	if filter.ColorFlag {
		filterQuery += ` LOWER(p.color) =LOWER('` + filter.Color + `') AND`
		sqlRegexp += filter.Color
		found = true
	}
	if found {
		var validParameters = regexp.MustCompile(`^[\w ]+$`)
		if validParameters.MatchString(sqlRegexp) == false {
			err = fmt.Errorf("Possible SQL Injection Attack.")
			logger.WithField("err", err.Error()).Error("Error In Parameters, special Characters are present.")
			return 0, []Product{}, err
		}
		filterQuery = ` WHERE ` + filterQuery[:len(filterQuery)-3]
	}

	filterQuery = filterSearchProduct + filterQuery

	if filter.PriceFlag {
		filterQuery += ` ORDER BY p.price ` + filter.Price + `, p.id LIMIT ` + limitStr + `  OFFSET  ` + offsetStr + `  ;`
	} else {
		filterQuery += ` ORDER BY p.id LIMIT ` + limitStr + `  OFFSET  ` + offsetStr + `  ;`
	}

	err = s.db.Select(&records, filterQuery)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching Products from database")
		return 0, []Product{}, err
	} else if len(records) == 0 {
		err = fmt.Errorf("Desired page not found, Offset was big")
		logger.WithField("err", err.Error()).Error("Products don't exist by such filters!")
		return 0, []Product{}, err
	}

	return records[0].TotalRecords, records.Products(), nil
}

// @Title SearchRecords
// @Description Get records that are searched as per request Parameter "text" along with count
// @Accept	request.Context, text as string, limitStr, pageStr
// @Success total= (count of search qualifying records), error=nil
// @Failure total=0, error= "Some Error"
func (s *pgStore) SearchProductsByText(ctx context.Context, text string, limitStr string, offsetStr string) (int, []Product, error) {

	var records Records
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

	searchQuery := ` WHERE `

	// iterate over all the textMap
	for key := range textMap {
		searchQuery += ` 
		LOWER(p.name) LIKE LOWER('%` + key + `%') OR 
		LOWER(p.brand) LIKE LOWER('%` + key + `%') OR 
		LOWER(c.cname) LIKE LOWER('%` + key + `%') OR`
	}

	// remove that last OR from searchQuery
	searchQuery = filterSearchProduct + searchQuery[:len(searchQuery)-2] +
		` LIMIT ` + limitStr + ` OFFSET  ` + offsetStr + ` ;`

	err := s.db.Select(&records, searchQuery)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching Products from database")
		return 0, []Product{}, err
	} else if len(records) == 0 {
		err = fmt.Errorf("Either Offset was big or No Records Present in database!")
		logger.WithField("err", err.Error()).Error("database Returned total record count as 0")
		return 0, []Product{}, err
	}

	return records[0].TotalRecords, records.Products(), nil
}
