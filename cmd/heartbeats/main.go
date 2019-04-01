package main

import (
	"context"
	"flag"
	"fmt"
	"log"
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
	sink      string
	label     string
	periodStr string
)

func init() {
	flag.IntVar(&max, "max", 10, "max count")
	flag.StringVar(&label, "label", "", "a special label")
	flag.StringVar(&periodStr, "period", "5", "the number of seconds between heartbeats")
}

type envConfig struct {
	// Sink URL where to send heartbeat cloudevents
	Sink string `envconfig:"SINK"`

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

	if env.Sink != "" {
		sink = env.Sink
	}

	c, err := kncloudevents.NewDefaultClient(sink)
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
		fmt.Sprintf("https://github.com/knative/eventing-sources/cmd/heartbeats/#%s/%s", env.Namespace, env.Name))
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
}
