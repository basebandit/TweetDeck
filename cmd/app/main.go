package main

import (
	"context"
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"ekraal.org/avatarlysis/api"
	"ekraal.org/avatarlysis/business/data/auth"
	"ekraal.org/avatarlysis/foundation/database"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

type Config struct {
	TwitterConsumerKey    string
	TwitterConsumerSecret string
	TwitterAccessToken    string
	TwitterAccessSecret   string
	TwitterTokenURL       string
	DatabaseUser          string
	DatabasePassword      string
	DatabaseName          string
	DatabaseHost          string
	DatabaseDisableTLS    bool
	PrivateKeyFile        string
	KeyID                 string
	Algorithm             string
}

func init() {
	//Lets confirm that all env vars are set
	if os.Getenv("CONSUMER_SECRET") == "" || os.Getenv("CONSUMER_KEY") == "" || os.Getenv("ACCESS_TOKEN") == "" || os.Getenv("ACCESS_SECRET") == "" || os.Getenv("TOKEN_URL") == "" || os.Getenv("AVATARLYSIS_DB_USER") == "" || os.Getenv("AVATARLYSIS_DB_PASSWORD") == "" || os.Getenv("AVATARLYSIS_DB_NAME") == "" || os.Getenv("AVATARLYSIS_DB_HOST") == "" || os.Getenv("AVATARLYSIS_DB_DISABLE_TLS") == "" || os.Getenv("AVATARLYSIS_PRIVATE_KEY") == "" || os.Getenv("AVATARLYSIS_KEY_ID") == "" || os.Getenv("AVATARLYSIS_ALGORITHM") == "" {
		log.Fatal("there is a missing config field")
	}
}

func main() {

	log := log.New(os.Stdout, "AVATARLYSIS : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	ctx := context.Background()

	cfg := env()

	db, err := database.Open(database.Config{
		User:       cfg.DatabaseUser,
		Password:   cfg.DatabasePassword,
		Host:       cfg.DatabaseHost,
		Name:       cfg.DatabaseName,
		DisableTLS: cfg.DatabaseDisableTLS,
	})

	if err != nil {
		log.Printf("main: %s", errors.Wrap(err, "connecting to db"))
		os.Exit(1)
	}

	if err := database.Ping(ctx, db); err != nil {
		log.Printf("main: %s", errors.Wrap(err, "pinging db"))
		os.Exit(1)
	}

	defer func() {
		log.Printf("main: Database Stopping : %s", cfg.DatabaseHost)
		db.Close()
	}()

	auth, err := authSetup(cfg)
	if err != nil {
		log.Printf("main: Auth setup failed : %s", err)
		os.Exit(1)
	}

	router := chi.NewRouter()

	api := api.NewServer(ctx, db, log, auth, router)

	srv := http.Server{
		Addr:    ":8880",
		Handler: api,
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()

	log.Printf("Server listening on %s\n", srv.Addr)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	sig := <-quit
	log.Printf("shutting down Avatarlysis API server...Reason: %s", sig)
	if err := srv.Shutdown(context.Background()); err != nil {
		panic(err)
	}

	log.Println("Avatarlysis API server stopped gracefully")
}

func authSetup(cfg *Config) (*auth.Auth, error) {

	privatePEM, err := ioutil.ReadFile(cfg.PrivateKeyFile)
	if err != nil {
		return nil, errors.Wrap(err, "reading auth private key")
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
	if err != nil {
		return nil, errors.Wrap(err, "parsing auth private key")
	}

	keyLookupFunc := func(publicKID string) (*rsa.PublicKey, error) {
		switch publicKID {
		case cfg.KeyID:
			return privateKey.Public().(*rsa.PublicKey), nil
		}
		return nil, fmt.Errorf("no public key found for the specified kid: %s", publicKID)
	}
	a, err := auth.New(privateKey, cfg.KeyID, cfg.Algorithm, keyLookupFunc)
	if err != nil {
		return nil, errors.Wrap(err, "constructing auth")
	}

	return a, nil
}

// func twitterLookUp() {
// 	fmt.Println("Hello")

// cfg := env()

// 	config := &clientcredentials.Config{
// 		ClientID:     cfg.TwitterConsumerKey,
// 		ClientSecret: cfg.TwitterConsumerSecret,
// 		TokenURL:     cfg.TwitterTokenURL,
// 	}

// 	httpClient := config.Client(oauth2.NoContext)

// 	client := twitter.NewClient(httpClient)

// 	userLookupParams := &twitter.UserLookupParams{ScreenName: []string{
// 		"OlandoWanda",
// 		"TPS_Ke",
// 		"EngKanyiri",
// 		"FelistusQ",
// 		"Otiisteve",
// 		"SashaShazlin",
// 		"Jean_Wangari",
// 		"JoyceKamande9",
// 		"Louiskandie",
// 		"DKJnr3",
// 	},
// 	}

// 	users, tres, err := client.Users.Lookup(userLookupParams)

// 	if err != nil {
// 		log.Fatalf("Your request failed with %v, error: %s", tres, err)
// 	}

// 	fmt.Printf("USER LOOKUP:\n%+v\n", users)
// }

func env() *Config {

	tls, err := strconv.ParseBool(os.Getenv("AVATARLYSIS_DB_DISABLE_TLS"))

	if err != nil {
		log.Printf("env: %s", errors.Wrap(err, "parsing DisableTLS environment config"))
		os.Exit(1)
	}

	cfg := Config{
		TwitterAccessSecret:   os.Getenv("CONSUMER_SECRET"),
		TwitterAccessToken:    os.Getenv("ACCESS_TOKEN"),
		TwitterConsumerKey:    os.Getenv("CONSUMER_KEY"),
		TwitterConsumerSecret: os.Getenv("CONSUMER_SECRET"),
		TwitterTokenURL:       os.Getenv("TOKEN_URL"),
		DatabaseUser:          os.Getenv("AVATARLYSIS_DB_USER"),
		DatabasePassword:      os.Getenv("AVATARLYSIS_DB_PASSWORD"),
		DatabaseName:          os.Getenv("AVATARLYSIS_DB_NAME"),
		DatabaseHost:          os.Getenv("AVATARLYSIS_DB_HOST"),
		DatabaseDisableTLS:    tls,
		Algorithm:             os.Getenv("AVATARLYSIS_ALGORITHM"),
		PrivateKeyFile:        os.Getenv("AVATARLYSIS_PRIVATE_KEY"),
		KeyID:                 os.Getenv("AVATARLYSIS_KEY_ID"),
	}

	return &cfg
}