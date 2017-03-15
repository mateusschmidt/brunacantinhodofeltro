package tweet

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/coreos/pkg/flagutil"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/gorilla/mux"
)

const consumerkey = "bBg6iwWB8a7sYM5G7cO92XCnu"
const consumersecret = "buPLEBoCNtXeTiPk02DdcQ7xWwXCOUqXVT1omdMftrMErgXRBV"
const accesstoken = "841445351559176196-TS9eqt1s51m0lIRU8MKhe67gbEHhvDR"
const accesssecret = "OtnANWNKLD62tnupDW9FPDnm0xc1QRz0eJZ86DywLy2bh"

type TweetMediaResume struct {
	URL  string
	Type string
}

type TweetResume struct {
	Text  string
	Media TweetMediaResume
}

func RegisterRouter(router *mux.Router) {
	router.HandleFunc("/api/tweet/posts", getPosts).Methods("GET")
}

func getPosts(w http.ResponseWriter, req *http.Request) {

	flags := flag.NewFlagSet("user-auth", flag.ExitOnError)
	consumerKey := flags.String("consumer-key", consumerkey, "Twitter Consumer Key")
	consumerSecret := flags.String("consumer-secret", consumersecret, "Twitter Consumer Secret")
	accessToken := flags.String("access-token", accesstoken, "Twitter Access Token")
	accessSecret := flags.String("access-secret", accesssecret, "Twitter Access Secret")
	flags.Parse(os.Args[1:])
	flagutil.SetFlagsFromEnv(flags, "TWITTER")

	if *consumerKey == "" || *consumerSecret == "" || *accessToken == "" || *accessSecret == "" {
		log.Fatal("Consumer key/secret and Access token/secret required")
	}

	config := oauth1.NewConfig(*consumerKey, *consumerSecret)
	token := oauth1.NewToken(*accessToken, *accessSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	// Home Timeline
	homeTimelineParams := &twitter.HomeTimelineParams{Count: 2}
	tweets, _, _ := client.Timelines.HomeTimeline(homeTimelineParams)
	log.Printf("User's HOME TIMELINE:\n%+v\n", tweets)

	resposta := make([]TweetResume, len(tweets))
	for index := 0; index < len(tweets); index++ {
		resposta[index].Text = tweets[index].Text
		if len(tweets[index].Entities.Media) > 0 {
			resposta[index].Media.URL = tweets[index].Entities.Media[0].DisplayURL
			resposta[index].Media.Type = tweets[index].Entities.Media[0].Type
		}
	}

	json.NewEncoder(w).Encode(resposta)
}
