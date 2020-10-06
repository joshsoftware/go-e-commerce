package config

import (
	"fmt"
	"strconv"

	"github.com/spf13/viper"
)

var (
	appName                string
	appPort                int
	jwtKey                 string
	jwtExpiryDurationHours int
)

// Load - loads all the environment variables and/or params in application.yml
func Load() {
	viper.SetDefault("APP_NAME", "e-commerce")
	viper.SetDefault("APP_PORT", "8002")

	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
	viper.AddConfigPath("./..")
	viper.AddConfigPath("./../..")
	viper.ReadInConfig()
	viper.AutomaticEnv()
	// Check for the presence of JWT_KEY and JWT_EXPIRY_DURATION_HOURS
	JWTKey()
	JWTExpiryDurationHours()
}

// AppName - returns the app name
func AppName() string {
	if appName == "" {
		appName = ReadEnvString("APP_NAME")
	}
	return appName
}

// AppPort - returns application http port
func AppPort() int {
	if appPort == 0 {
		appPort = ReadEnvInt("APP_PORT")
	}
	return appPort
}

// JWTKey - returns the JSON Web Token key
func JWTKey() []byte {
	return []byte(ReadEnvString("JWT_SECRET"))
}

//MailerConfig - returns Configuration for mailer
func MailerConfig() (host string, port int, email string, username string, password string) {

	host = ReadEnvString("MAILER_HOST")
	port = ReadEnvInt("MAILER_PORT")
	email = ReadEnvString("MAILER_EMAIL")
	username = ReadEnvString("MAILER_USERNAME")
	password = ReadEnvString("MAILER_PASSWORD")

	return host, port, email, username, password
}

// JWTExpiryDurationHours - returns duration for jwt expiry in int
func JWTExpiryDurationHours() int {
	return int(ReadEnvInt("JWT_EXPIRY_DURATION_HOURS"))
}

// ReadEnvInt - reads an environment variable as an integer
func ReadEnvInt(key string) int {
	checkIfSet(key)
	v, err := strconv.Atoi(viper.GetString(key))
	if err != nil {
		panic(fmt.Sprintf("key %s is not a valid integer", key))
	}
	return v
}

// ReadEnvString - reads an environment variable as a string
func ReadEnvString(key string) string {
	checkIfSet(key)
	return viper.GetString(key)
}

// ReadEnvBool - reads environment variable as a boolean
func ReadEnvBool(key string) bool {
	checkIfSet(key)
	return viper.GetBool(key)
}

//CheckIfSet checks if all the necessary keys are set
func checkIfSet(key string) {
	if !viper.IsSet(key) {
		err := fmt.Errorf("Key %s is not set", key)
		panic(err)
		//logger.WithField("err", err.Error()).Error("Error Couldn't find db!")
	}
}
