package service

import (
	"encoding/json"
	"net/http"

	logger "github.com/sirupsen/logrus"
)

// CountryData - stores country with states and cities
type CountryData []struct {
	Name   string `json:"name"`
	States []struct {
		Name   string `json:"name"`
		Cities []struct {
			Name string `json:"name"`
		} `json:"cities"`
	} `json:"states"`
}

var (
	defaultdata = `
	[
   {
      "name":"India",
      "states":[
         {
            "name":"Maharashtra",
            "cities":[
               {
                  "name":"Mumbai"
               },
               {
                  "name":"Pune"
			   },
			   {
				"name":"Nagpur"
			 },
			 {
				"name":"Aurangabad"
			 },
			 {
				"name":"Nashik"
			 }
			   
            ]
         },
         {
            "name":"Madhya Pradesh",
            "cities":[
               {
                  "name":"Bhopal"
               },
               {
                  "name":"Indore"
			   },
			   {
				"name":"Gwalior"
			   },
			   {
				"name":"Ujjain"
			   },
			   {
				"name":"Jabalpur"
			   }
            ]
         }
      ]
   },
   {
      "name":"Pakistan",
      "states":[
         {
            "name":"Sindh",
            "cities":[
               {
                  "name":"Karachi"
               },
               {
                  "name":"Sukkur"
			   },
			   {
				"name":"Larkana"
			   },
			   {
				"name":"Thatta"
			   },
			   {
				"name":"Badin"
			   }
            ]
         },
         {
            "name":"Punjab",
            "cities":[
               {
                  "name":"Lahore"
               },
               {
                  "name":"Faisalabad"
			   },
			   {
				"name":"Multan"
			   },
			   {
				"name":"Rawalpindi"
			   },
			   {
				"name":"Bahawalpur"
			   }
            ]
         }
      ]
   }
]
	`
)

func countryDataHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		countryData := CountryData{}
		err := json.Unmarshal([]byte(defaultdata), &countryData)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while Unmarshalling request json")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		respBytes, err := json.Marshal(countryData)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while marshalling footer data")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		rw.Header().Add("Content-Type", "application/json")
		rw.Write(respBytes)
	})
}
