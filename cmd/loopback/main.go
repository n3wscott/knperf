package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	"github.com/knative/eventing-sources/pkg/kncloudevents"
	"github.com/montanaflynn/stats"
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
	flag.StringVar(&periodStr, "period", "5", "the number of seconds between heartbeats")
}

type config struct {
	// Target URL where to send heartbeat cloudevents
	Target string `envconfig:"TARGET" required:"true"`

	// Name of this pod.
	Name string `envconfig:"POD_NAME" required:"true"`

	// Namespace this pod exists in.
	Namespace string `envconfig:"POD_NAMESPACE" required:"true"`
}

type LoopBack struct {
	ce client.Client

	period time.Duration
	source types.URLRef
	count  int
	label  string

	sent    map[string]time.Time
	latency []time.Duration
}

func (lb *LoopBack) Start(ctx context.Context) {
	go lb.Send(ctx)
	if err := lb.ce.StartReceiver(ctx, lb.Receive); err != nil {
		log.Fatal(err)
	}
}

func (lb *LoopBack) Send(ctx context.Context) {
	hb := &Heartbeat{
		Sequence: 0,
		Label:    label,
	}
	ticker := time.NewTicker(lb.period)
	for ; lb.count > 0; lb.count-- {
		hb.Sequence++

		id := uuid.New().String()
		event := cloudevents.Event{
			Context: cloudevents.EventContextV02{
				Type:   "dev.knative.eventing.samples.heartbeat",
				ID:     id,
				Source: lb.source,
				Extensions: map[string]interface{}{
					"the":   42,
					"heart": "yes",
					"beats": true,
				},
			},
			Data: hb,
		}

		if _, ok := lb.sent[id]; ok {
			log.Printf("[WARN] resent %s", id)
		} else {
			lb.sent[id] = time.Now()
		}

		if _, err := lb.ce.Send(ctx, event); err != nil {
			log.Printf("failed to send cloudevent: %s", err.Error())
		}
		// Wait for next tick
		<-ticker.C
	}
}

func (lb *LoopBack) Receive(event cloudevents.Event) {
	id := event.Context.AsV02().ID
	if then, ok := lb.sent[id]; !ok {
		log.Printf("[WARN] unknown event %s", id)
	} else {
		latency := time.Since(then)
		log.Printf("latency %s", latency.String())
		lb.latency = append(lb.latency, latency)
		delete(lb.sent, id)
	}
}

func (lb *LoopBack) Validate() {
	for k, _ := range lb.sent {
		log.Printf("[ERR] undelivered event %s", k)
	}
	data := stats.LoadRawData(lb.latency)

	minimum, err := stats.Min(data)
	maximum, err := stats.Max(data)
	mean, err := stats.Mean(data)

	if err != nil {
		log.Printf("failed to calculate min, max and mean: %s", err)
		return
	}

	min := time.Duration(minimum)
	max := time.Duration(maximum)
	m := time.Duration(mean)
	log.Printf("min: %s; max: %s; mean: %s ", min, max, m)
}

func main() {
	flag.Parse()

	var env config
	envconfig.MustProcess("", &env)

	ce, err := kncloudevents.NewDefaultClient(env.Target)
	if err != nil {
		log.Fatalf("failed to create client: %s", err.Error())
	}

	var period time.Duration
	if p, err := strconv.Atoi(periodStr); err != nil {
		period = time.Duration(50) * time.Millisecond
	} else {
		period = time.Duration(p) * time.Millisecond
	}

	source := types.ParseURLRef(
		fmt.Sprintf("https://github.com/n3wscott/knperf/cmd/loopback/#%s/%s", env.Namespace, env.Name))
	log.Printf("Loopback Source: %s", source)

	if len(label) > 0 && label[0] == '"' {
		label, _ = strconv.Unquote(label)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(max)*period+time.Second*10)
	defer cancel()

	lb := &LoopBack{
		ce:      ce,
		source:  *source,
		count:   max,
		label:   label,
		period:  period,
		sent:    make(map[string]time.Time, max),
		latency: make([]time.Duration, 0),
	}
	lb.Start(ctx)

	lb.Validate()
}
