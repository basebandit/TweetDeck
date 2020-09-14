package service

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/pkg/errors"
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

//Lookup fetches the profiles of the given username(s)
func (t *TwitterService) Lookup(usernames []string) ([]twitter.User, error) {
	client := t.client()

	userLookupParams := &twitter.UserLookupParams{
		ScreenName: usernames,
	}
	users, _, err := client.Users.Lookup(userLookupParams)
	if err != nil {
		return nil, errors.Wrap(err, "fetching twitter profiles")
	}

	return users, nil
}

//UserTimeline fetches the timeline of the user with the given id and returns a list of recent tweets  and retweets.
func (t *TwitterService) UserTimeline(id int64) ([]twitter.Tweet, error) {
	client := t.client()
	userTimelineParams := &twitter.UserTimelineParams{UserID: id}
	tweet, _, err := client.Timelines.UserTimeline(userTimelineParams)
	if err != nil {
		return nil, errors.Wrap(err, "fetching twitter user timeline")
	}
	return tweet, nil
}

func (t *TwitterService) client() *twitter.Client {
	config := &clientcredentials.Config{
		ClientID:     t.ConsumerKey,
		ClientSecret: t.ConsumerSecret,
		TokenURL:     t.TokenURL,
	}

	httpClient := config.Client(oauth2.NoContext)
	client := twitter.NewClient(httpClient)
	return client
}
