package service

import "joshsoftware/go-e-commerce/db"

type Dependencies struct {
	Store db.Storer
	// define other service dependencies
}
