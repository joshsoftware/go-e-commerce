package service

import (
	"joshsoftware/go-e-commerce/config"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	// "github.com/dgrijalva/jwt-go"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
	"github.com/urfave/negroni"
)

func TestExampleTestSuite(t *testing.T) {
	config.Load()
	suite.Run(t, new(CartHandlerTestSuite))
}

func makeHTTPCall(method, path, requestURL, body string, handlerFunc http.HandlerFunc) (recorder *httptest.ResponseRecorder) {
	// create a http request using the given parameters
	req, _ := http.NewRequest(method, requestURL, strings.NewReader(body))
	// test recorder created for capturing api responses
	recorder = httptest.NewRecorder()
	// create a router to serve the handler in test with the prepared request
	router := mux.NewRouter()
	router.HandleFunc(path, handlerFunc).Methods(method)

	// serve the request and write the response to recorder
	router.ServeHTTP(recorder, req)
	return
}

func makeHTTPCallWithJWTMiddleware(method, path, requestURL, body string, handlerFunc http.HandlerFunc) (recorder *httptest.ResponseRecorder) {
	// create jwt token with userID
	JWTToken, _ := generateJwt(1)

	req, _ := http.NewRequest(method, requestURL, strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+JWTToken)
	// fmt.Println("kklklkl", req.Header)
	recorder = httptest.NewRecorder()
	// create a router to serve the handler in test with the prepared request
	router := mux.NewRouter()

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return config.JWTKey(), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
	// fmt.Printf("%T-=-Type-------------------------", jwtMiddleware) //output - *jwtmiddleware.JWTMiddleware
	router.Handle(path, negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(handlerFunc)),
	)).Methods(method)

	router.ServeHTTP(recorder, req)
	// fmt.Println("recorder", recorder)
	return
}

// func makeHTTPCallWithJWTMiddleware(method, path, requestURL, body string, handlerFunc http.HandlerFunc) (recorder *httptest.ResponseRecorder) {
// 	// create jwt token with userID
// 	JWTToken := ""
// 	signingMethod := jwt.SigningMethodHS256
// 	if body == "2" {
// 		// JWTToken = "00000"
// 		// JWTToken, _ = generateJwt(1)
// 		// JWTToken = JWTToken[:len(JWTToken)-3] + "abc"
// 		// body = ""
// 		signingMethod = nil
// 	} else {
// 		JWTToken, _ = generateJwt(1)
// 	}

// 	req, _ := http.NewRequest(method, requestURL, strings.NewReader(body))
// 	req.Header.Set("Authorization", "Bearer "+JWTToken)
// 	fmt.Println("kklklkl", req.Header)
// 	recorder = httptest.NewRecorder()
// 	// create a router to serve the handler in test with the prepared request
// 	router := mux.NewRouter()

// 	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
// 		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
// 			return config.JWTKey(), nil
// 		},
// 		SigningMethod: signingMethod,
// 	})
// 	fmt.Printf("%T-=-77777777777777777777777777", jwtMiddleware)
// 	router.Handle(path, negroni.New(
// 		negroni.HandlerFunc(JWT(signingMethod).HandlerWithNext),
// 		negroni.Wrap(http.HandlerFunc(handlerFunc)),
// 	)).Methods(method)

// 	router.ServeHTTP(recorder, req)
// 	fmt.Println("recorder", recorder)
// 	return
// }

// const defaultAuthorizationHeaderName = "Authorization"

// func makeAuthenticatedRequest(m string, p string, u string, c jwt.Claims, e jwt.SigningMethod) *httptest.ResponseRecorder {
// 	r, _ := http.NewRequest(m, u, nil)
// 	if c != nil {
// 		token := jwt.New(jwt.SigningMethodHS256)
// 		token.Claims = c
// 		// private key generated with http://kjur.github.io/jsjws/tool_jwt.html
// 		s, e := token.SignedString(config.JWTKey())
// 		if e != nil {
// 			panic(e)
// 		}
// 		r.Header.Set(defaultAuthorizationHeaderName, fmt.Sprintf("bearer %v", s))
// 	}
// 	w := httptest.NewRecorder()
// 	n := createNegroniMiddleware(m, p, e)
// 	n.ServeHTTP(w, r)
// 	return w
// }

// func createNegroniMiddleware(m string, p string, e jwt.SigningMethod) *negroni.Negroni {
// 	publicRouter := mux.NewRouter().StrictSlash(true)

// 	negProtected := negroni.New()
// 	//add the JWT negroni handler
// 	negProtected.Use(negroni.HandlerFunc(JWT(e).HandlerWithNext))
// 	negProtected.UseHandler(publicRouter)

// 	//Create the main router
// 	mainRouter := mux.NewRouter().StrictSlash(true)

// 	mainRouter.Handle(p, negProtected)
// 	// mainRouter.Handle("/protected", negProtected)
// 	//if routes match the handle prefix then I need to add this dummy matcher {_dummy:.*}
// 	// mainRouter.Handle("/protected/{_dummy:.*}", negProtected)

// 	n := negroni.Classic()
// 	// This are the "GLOBAL" middlewares that will be applied to every request
// 	// examples are listed below:
// 	//n.Use(gzip.Gzip(gzip.DefaultCompression))
// 	//n.Use(negroni.HandlerFunc(SecurityMiddleware().HandlerFuncWithNext))
// 	n.UseHandler(mainRouter)

// 	return n
// 	// return
// }

// func JWT(e jwt.SigningMethod) *jwtmiddleware.JWTMiddleware {
// 	var privateKey = config.JWTKey()

// 	return jwtmiddleware.New(jwtmiddleware.Options{
// 		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
// 			// if privateKey == nil {
// 			// 	var err error
// 			// 	privateKey, err = readPrivateKey()
// 			// 	if err != nil {
// 			// 		panic(err)
// 			// 	}
// 			// }
// 			return privateKey, nil
// 		},
// 		SigningMethod: e,
// 	})
// }
