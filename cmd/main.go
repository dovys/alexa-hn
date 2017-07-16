package main

import (
	"fmt"
	"net/http"
	"time"

	alexahn "github.com/dovys/alexa-hn"
	"github.com/dovys/alexa-hn/hn"
	"github.com/dovys/alexa-hn/stub"
	"github.com/kelseyhightower/envconfig"
	alexa "github.com/mikeflynn/go-alexa/skillserver"
)

type config struct {
	StubHN   bool          `envconfig:"stub_hn" default:"false"`
	CacheTTL time.Duration `envconfig:"cache_ttl" default:"1h"`
}

func main() {
	var c config
	envconfig.MustProcess("", &c)

	var hnc hn.Client

	if c.StubHN {
		hnc = stub.NewStubHNClient()
	} else {
		hnc = hn.NewClient(&http.Client{})
	}

	svc := alexahn.Cache(alexahn.NewService(hnc), c.CacheTTL)

	h := MakeEchoIntentHandler(svc)
	apps := map[string]interface{}{
		"/echo/hn": alexa.EchoApplication{
			AppID:    "amzn1.ask.skill.6897d6d0-1718-466a-833e-75e5efe261de",
			OnIntent: h,
			OnLaunch: h,
		},
		"/_ah/health": alexa.StdApplication{
			Methods: "GET",
			Handler: HealthCheckhandler,
		},
	}

	alexa.Run(apps, "8080")
}

func MakeEchoIntentHandler(svc alexahn.Service) func(*alexa.EchoRequest, *alexa.EchoResponse) {
	return func(echoReq *alexa.EchoRequest, echoResp *alexa.EchoResponse) {
		stories, err := svc.ReadTopStories()
		if err != nil {
			fmt.Println(err)
			echoResp.OutputSpeech("Please try again later")
			return
		}

		fmt.Printf(
			"%+v\n%+v\n%+v\n%+v\n",
			echoReq.GetIntentName(),
			echoReq.AllSlots(),
			echoReq.GetRequestType(),
			echoReq.GetSessionID(),
		)

		echoResp.OutputSpeech(stories.Speech)

		if stories.Card != nil {
			echoResp.Card(stories.Card.Title, stories.Card.Content)
		}
	}
}

func HealthCheckhandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}
