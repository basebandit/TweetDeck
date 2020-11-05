package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	"go.opentelemetry.io/otel/api/global"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

//TwitterService defines the parameters required by the twitter api
type TwitterService struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
	TokenURL       string
}

//NewTwitter initializes a twitter service object and returns it.
func NewTwitter(consumerKey, consumerSecret, accessSecret, accessToken, tokenURL string) *TwitterService {
	t := TwitterService{
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
		AccessToken:    accessToken,
		AccessSecret:   accessSecret,
		TokenURL:       tokenURL,
	}

	return &t
}

func (t *TwitterService) client() *twitter.Client {
	// //==== This is App Only Auth uses oauth 2
	config := &clientcredentials.Config{
		ClientID:     t.ConsumerKey,
		ClientSecret: t.ConsumerSecret,
		TokenURL:     t.TokenURL,
	}

	// http.Client will automatically authorize Requests
	httpClient := config.Client(oauth2.NoContext)
	//Twitter client
	client := twitter.NewClient(httpClient)

	return client
}

//Lookup fetches the profiles of the given username(s)
func (t *TwitterService) Lookup(ctx context.Context, usernames []string) ([]twitter.User, error) {
	ctx, span := global.Tracer("avatarlysis").Start(ctx, "foundation.service.twitter.lookup")
	defer span.End()

	client := t.client()

	userLookupParams := &twitter.UserLookupParams{
		ScreenName: usernames,
	}

	fmt.Println("Twitter lookup started now...")

	if client != nil {
		users, tres, err := client.Users.Lookup(userLookupParams)
		if tres != nil {
			if tres.StatusCode != 200 {
				fmt.Printf("oops something unexpected occured response code: [%v]\n", tres.StatusCode)
			}
		}

		if err != nil {
			fmt.Printf("an error occurred looking up twitter profiles error: %s\n", err)
			return nil, err
		}

		fmt.Println("Twitter lookup finished successfully!")

		return users, nil
	}

	return nil, errors.New("Twitter client object is uninitialized. it might be nil")
}

//LastTweetTime retrieves the given user id's last tweet time
func (t *TwitterService) LastTweetTime(id int64) (ltt string) {
	client := t.client()

	userTimelineParams := &twitter.UserTimelineParams{UserID: id, Count: 2}
	tweets, tres, err := client.Timelines.UserTimeline(userTimelineParams)

	if tres.StatusCode != 200 {
		fmt.Printf("error: %v response code: %v\n", err, tres.StatusCode)
		return
	}

	if len(tweets) > 0 {
		ltt = tweets[0].CreatedAt
		return
	}

	return
}
