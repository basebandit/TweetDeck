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
	"github.com/robfig/cron/v3"
)

//Config declares any external config that our service needs
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
	MaxQueue              int
	MaxWorkers            int
}

func init() {
	//Lets confirm that all env vars are set
	if os.Getenv("CONSUMER_SECRET") == "" || os.Getenv("CONSUMER_KEY") == "" || os.Getenv("ACCESS_TOKEN") == "" || os.Getenv("ACCESS_SECRET") == "" || os.Getenv("TOKEN_URL") == "" || os.Getenv("AVATARLYSIS_DB_USER") == "" || os.Getenv("AVATARLYSIS_DB_PASSWORD") == "" || os.Getenv("AVATARLYSIS_DB_NAME") == "" || os.Getenv("AVATARLYSIS_DB_HOST") == "" || os.Getenv("AVATARLYSIS_DB_DISABLE_TLS") == "" || os.Getenv("AVATARLYSIS_PRIVATE_KEY") == "" || os.Getenv("AVATARLYSIS_KEY_ID") == "" || os.Getenv("AVATARLYSIS_ALGORITHM") == "" {
		log.Fatal("there is a missing config field")
	}
}

func main() {

	logger := log.New(os.Stdout, "AVATARLYSIS : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

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
		logger.Printf("main: %s", errors.Wrap(err, "connecting to db"))
		os.Exit(1)
	}

	if err := database.Ping(ctx, db); err != nil {
		logger.Printf("main: %s", errors.Wrap(err, "pinging db"))
		os.Exit(1)
	}

	defer func() {
		logger.Printf("main: Database Stopping : %s", cfg.DatabaseHost)
		db.Close()
	}()

	auth, err := authSetup(cfg)
	if err != nil {
		logger.Printf("main: Auth setup failed : %s", err)
		os.Exit(1)
	}

	router := chi.NewRouter()

	apiCfg := api.ApiConfig{
		TwitterAccessSecret: cfg.TwitterAccessSecret, TwitterAccessToken: cfg.TwitterAccessToken, TwitterConsumerKey: cfg.TwitterConsumerKey, TwitterConsumerSecret: cfg.TwitterConsumerSecret,
		TwitterTokenURL: cfg.TwitterTokenURL,
	}

	api := api.NewServer(ctx, db, &apiCfg, logger, auth, router)

	srv := http.Server{
		Addr:    ":8880",
		Handler: api,
	}

	// pass in your specific zone name, using Kenya/Nairobi as example
	c := cron.New()
	c.AddFunc("CRON_TZ=Africa/Nairobi 00 09 * * *", func() {
		if err := api.TwitterLookup(); err != nil {
			logger.Println(err)
		}
	})
	c.Start()

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

// func chunk(users []string) [][]string {
// 	var divided [][]string
// 	// numCPU := 10
// 	chunkSize := 100

// 	for i := 0; i < len(users); i += chunkSize {
// 		end := i + chunkSize

// 		if end > len(users) {
// 			end = len(users)
// 		}

// 		divided = append(divided, users[i:end])
// 	}
// 	return divided
// }

// func twitterLookup(ctx context.Context, db *sqlx.DB, cfg *Config, log *log.Logger) error {
// 	ctx, span := global.Tracer("avatarlysis").Start(ctx, "main.twitterlookup")
// 	defer span.End()

// 	avs, err := avatar.GetUsernames(ctx, db)
// 	if err != nil {
// 		return err
// 	}

// 	dict := map[string]string{}
// 	for _, av := range avs {
// 		dict[av.Username] = av.ID
// 	}

// 	twitter := service.NewTwitter(cfg.TwitterConsumerKey, cfg.TwitterConsumerSecret, cfg.TwitterAccessToken, cfg.TwitterTokenURL)

// 	var unames []string

// 	for _, av := range avs {
// 		unames = append(unames, av.Username)
// 	}

// 	// fmt.Printf("%+v\n", unames)
// 	//Lets implement our queueing here
// 	chunks := chunk(unames)

// 	var count int
// 	for _, chunk := range chunks {
// 		fmt.Println("Chunk", len(chunk))
// 		users, err := twitter.Lookup(ctx, log, chunk)
// 		if err != nil {
// 			return err
// 		}

// 		if err := profiler(ctx, twitter, dict, users, db); err != nil {
// 			return err
// 		}

// 		count++
// 	}

// 	return nil
// }

// func profiler(ctx context.Context, twitter *service.TwitterService, dict map[string]string, users []twitter.User, db *sqlx.DB) error {
// 	var np profile.NewProfile
// 	var nps []profile.NewProfile
// 	now := time.Now()
// 	for _, user := range users {

// 		np.ID = uuid.New().String()
// 		np.CreatedAt = now
// 		np.UpdatedAt = now
// 		avatarID, ok := dict[user.ScreenName]
// 		if !ok {
// 			//lets skip to the next iteration
// 			continue
// 		}
// 		np.AvatarID = stringPointer(avatarID)
// 		np.Name = stringPointer(user.Name)
// 		np.Followers = intPointer(user.FollowersCount)
// 		np.Following = intPointer(user.FriendsCount)
// 		np.Likes = intPointer(user.FavouritesCount)
// 		np.Tweets = intPointer(user.StatusesCount)

// 		np.ProfileImageURL = stringPointer(strings.ReplaceAll(user.ProfileImageURLHttps, "_normal", ""))
// 		np.Bio = stringPointer(user.Description)
// 		np.TwitterID = stringPointer(user.IDStr)
// 		np.JoinDate = stringPointer(user.CreatedAt)
// 		ltt, err := twitter.UserTimeline(user.ID)
// 		if err != nil {
// 			log.Println(err)
// 			continue
// 		}
// 		if len(ltt) > 0 {
// 			np.LastTweetTime = stringPointer(ltt[0].CreatedAt)
// 		}

// 		nps = append(nps, np)
// 	}

// 	if err := profile.CreateMultiple(ctx, db, nps, time.Now()); err != nil {
// 		return err
// 	}

// 	return nil
// }

func stringPointer(val string) *string {
	str := val

	return &str
}

func intPointer(val int) *int {
	i := val

	return &i
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
