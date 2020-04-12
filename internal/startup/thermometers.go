package startup

import (
	"github.com/benfollis/calor/internal/config"
	"github.com/benfollis/calor/internal/pubsub"
	"github.com/benfollis/calor/internal/readings"
	"github.com/benfollis/calor/internal/utils"
	"sync"
	"time"
)

func StartThermometers(config config.BoundConfig, ps *pubsub.PubSub) sync.WaitGroup {
	var twg sync.WaitGroup
	for _, boundTherm := range config.Thermometers {
		thermometer := boundTherm.Thermometer
		ticker := time.NewTicker(time.Duration(boundTherm.UpdateInterval) * time.Second)
		twg.Add(1)
		go func (group *sync.WaitGroup) {
			defer group.Done()
			for {
				select {
				case <-ticker.C:
					reading, err := thermometer.Read()
					utils.CheckLog(err)
					if err == nil {
						ps.Publish(readings.Topic, reading)
					}
				}
			}
		}(&twg)
	}
	return twg
}
