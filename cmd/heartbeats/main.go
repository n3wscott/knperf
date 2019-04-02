package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/knative/eventing-sources/pkg/kncloudevents"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/kelseyhightower/envconfig"
)

type Heartbeat struct {
	Sequence int    `json:"id"`
	Label    string `json:"label"`
}

var (
	max       int
	label     string
	periodStr string
)

func init() {
	flag.IntVar(&max, "max", 10, "max count")
	flag.StringVar(&label, "label", "", "a special label")
	flag.StringVar(&periodStr, "period", "2", "the number of seconds between heartbeats")
}

type envConfig struct {
	// Target URL where to send heartbeat cloudevents
	Target string `envconfig:"TARGET" required:"true"`

	// Name of this pod.
	Name string `envconfig:"POD_NAME" required:"true"`

	// Namespace this pod exists in.
	Namespace string `envconfig:"POD_NAMESPACE" required:"true"`
}

func main() {
	flag.Parse()

	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		os.Exit(1)
	}

	c, err := kncloudevents.NewDefaultClient(env.Target)
	if err != nil {
		log.Fatalf("failed to create client: %s", err.Error())
	}

	var period time.Duration
	if p, err := strconv.Atoi(periodStr); err != nil {
		period = time.Duration(5) * time.Second
	} else {
		period = time.Duration(p) * time.Second
	}

	source := types.ParseURLRef(
		fmt.Sprintf("https://github.com/n3wscott/knperf/cmd/heartbeats/#%s/%s", env.Namespace, env.Name))
	log.Printf("Heartbeats Source: %s", source)

	if len(label) > 0 && label[0] == '"' {
		label, _ = strconv.Unquote(label)
	}
	hb := &Heartbeat{
		Sequence: 0,
		Label:    label,
	}
	ticker := time.NewTicker(period)
	for i := 0; i < max; i++ {
		hb.Sequence++

		event := cloudevents.Event{
			Context: cloudevents.EventContextV02{
				Type:   "dev.knative.eventing.samples.heartbeat",
				Source: *source,
				Extensions: map[string]interface{}{
					"the":   42,
					"heart": "yes",
					"beats": true,
				},
			},
			Data: hb,
		}

		if _, err := c.Send(context.Background(), event); err != nil {
			log.Printf("failed to send cloudevent: %s", err.Error())
		}
		// Wait for next tick
		<-ticker.C
	}

	log.Printf("quiting...")

	quitURL, _ := url.Parse("http://localhost:15000/quitquitquit")

	req := &http.Request{
		Method: http.MethodPost,
		URL:    quitURL,
	}

	if resp, err := http.DefaultClient.Do(req); err != nil {
		log.Printf("[ERROR] failed to call istio: %s", err)
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("istio quituquitquit: %s", string(body))
	}

	os.Exit(0)
}
