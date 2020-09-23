package service

import (
	"context"
	"log"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/pkg/errors"
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

//Lookup fetches the profiles of the given username(s)
func (t *TwitterService) Lookup(ctx context.Context, l *log.Logger, usernames []string) ([]twitter.User, error) {
	ctx, span := global.Tracer("avatarlysis").Start(ctx, "foundation.service.twitter.lookup")
	defer span.End()

	client := t.client()

	userLookupParams := &twitter.UserLookupParams{
		ScreenName: usernames,
	}

	log.Println("twitter lookup started")

	users, tres, err := client.Users.Lookup(userLookupParams)
	if tres.StatusCode != 200 {
		log.Printf("fetching twitter profiles :%s res : [%v]\n", err, tres.StatusCode)
	}
	// if err != nil {
	// 	return nil, errors.Wrap(err, "fetching twitter profiles")
	// }

	log.Printf("twitter usernames lookup finished successful : res [%v]\n", tres.StatusCode)
	return users, nil
}

//UserTimeline fetches the timeline of the user with the given id and returns a list of recent tweets  and retweets.
func (t *TwitterService) UserTimeline(id int64) ([]twitter.Tweet, error) {
	client := t.client()
	userTimelineParams := &twitter.UserTimelineParams{UserID: id}
	tweet, tres, err := client.Timelines.UserTimeline(userTimelineParams)
	if err != nil {
		return nil, errors.Wrapf(err, "fetching twitter user timeline : res [%v]\n", tres)
	}

	log.Printf("twitter user timeline lookup finished successful : res [%v]\n", tres.StatusCode)

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
