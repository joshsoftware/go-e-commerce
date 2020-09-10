package service

import (
	"encoding/json"
	"net/http"

	logger "github.com/sirupsen/logrus"
)

type Footer []struct {
	Categories []struct {
		ID   int    `json:"id"`
		URL  string `json:"url"`
		Name string `json:"name"`
	} `json:"categories,omitempty"`
	Partners []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"partners,omitempty"`
	Contactus []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"contactus,omitempty"`
}

var (
	// Temporary hardcoded later can be fetched from DB
	defaultFooter = `
	[
		{"categories" : 
			[{"id": 1, "url":"https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_1280.png", "name": "About Us"},
			{"id": 2, "url":"https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_1280.png", "name": "Testomonials"},
			{"id": 3, "url":"https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_1280.png", "name": "Contact"},
			{"id": 4, "url":"https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_1280.png", "name": "Journal"},
			{"id": 5, "url":"https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_1280.png", "name": "Privacy Policy"}]},
		{"partners" : 
			[{"id": 6,  "name": "Support"},
			{"id": 7,  "name": "Shipping & Returns"},
			{"id": 8,  "name": "Size Guide"},
			{"id": 9,  "name": "Product Care"}]},
		{"contactus" : 
			[{"id": 11,  "name": "Josh"},
			{"id": 12,  "name": "Pune"},
			{"id": 13,  "name": "Bavdhan"},
			{"id": 14,  "name": "+974524379"}]}     

	]	
	`
)

// @Title getFooter
// @Description Return footer data
// @Router /footer [get]
// @Accept  json
// @Success 200 {object}
// @Failure 400 {object}
func getFooterHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		footer := Footer{}
		err := json.Unmarshal([]byte(defaultFooter), &footer)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while Unmarshalling request json")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		respBytes, err := json.Marshal(footer)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while marshalling footer data")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		rw.Header().Add("Content-Type", "application/json")
		rw.Write(respBytes)
	})
}
