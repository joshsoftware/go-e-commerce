package service

import "joshsoftware/go-e-commerce/db"

//Dependencies Structure
type Dependencies struct {
	Store db.Storer
	// define other service dependencies
}
