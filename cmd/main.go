package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"ekraal.org/avatarlysis/api"
	"ekraal.org/avatarlysis/database"

	"github.com/go-chi/chi"
)

type Config struct {
	TwitterConsumerKey    string
	TwitterConsumerSecret string
	TwitterAccessToken    string
	TwitterAccessSecret   string
	TwitterTokenURL       string
}

func init() {
	//Lets confirm that all env vars are set
	if os.Getenv("CONSUMER_SECRET") == "" || os.Getenv("CONSUMER_KEY") == "" || os.Getenv("ACCESS_TOKEN") == "" || os.Getenv("ACCESS_SECRET") == "" || os.Getenv("TOKEN_URL") == "" {
		log.Fatal("there is a missing config field")
	}
}

func main() {
	ctx := context.Background()

	db, err := database.Connect(ctx, "localhost:27017", "avatar", "avatar", "avatars")

	if err != nil {
		log.Println("main: error:", err)
		os.Exit(1)
	}
	log := log.New(os.Stdout, "AVATARLYSIS : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	router := chi.NewRouter()

	api := api.NewServer(db, log, router)

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

// func twitterLookUp() {
// 	fmt.Println("Hello")

// 	cfg := env()

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

// func env() *Config {

// 	cfg := Config{
// 		TwitterAccessSecret:   os.Getenv("CONSUMER_SECRET"),
// 		TwitterAccessToken:    os.Getenv("ACCESS_TOKEN"),
// 		TwitterConsumerKey:    os.Getenv("CONSUMER_KEY"),
// 		TwitterConsumerSecret: os.Getenv("CONSUMER_SECRET"),
// 		TwitterTokenURL:       os.Getenv("TOKEN_URL"),
// 	}

// 	return &cfg
// }
