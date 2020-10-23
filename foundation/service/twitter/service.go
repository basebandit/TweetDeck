package service

import (
	"context"
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
	TokenURL       string
}

//NewTwitter initializes a twitter service object and returns it.
func NewTwitter(key, secret, token, url string) *TwitterService {
	t := TwitterService{
		ConsumerKey:    key,
		ConsumerSecret: secret,
		AccessToken:    token,
		TokenURL:       url,
	}

	return &t
}

func (t *TwitterService) client() *twitter.Client {
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
func (t *TwitterService) Lookup(ctx context.Context, usernames []string) []twitter.User {
	ctx, span := global.Tracer("avatarlysis").Start(ctx, "foundation.service.twitter.lookup")
	defer span.End()

	client := t.client()

	userLookupParams := &twitter.UserLookupParams{
		ScreenName: usernames,
	}

	fmt.Println("Twitter lookup started...")

	users, tres, err := client.Users.Lookup(userLookupParams)
	if tres.StatusCode != 200 {
		fmt.Printf("oops something unexpected occured: %s response code [%v]\n", err, tres.StatusCode)
	}

	if err != nil {
		fmt.Printf("an error occured looking up twitter profiles: %v error: %s response code: [%v]\n", userLookupParams.ScreenName, err, tres.StatusCode)
	}

	fmt.Println("Twitter lookup finished successfully!")
	return users
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
