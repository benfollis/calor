package startup

import (
	"github.com/benfollis/calor/internal/config"
	"github.com/benfollis/calor/internal/pubsub"
	"github.com/benfollis/calor/internal/readings"
	"time"
)

func StartThermometers(config config.BoundConfig, ps *pubsub.PubSub){
	for _, boundTherm := range config.Thermometers {
		thermometer := boundTherm.Thermometer
		ticker := time.NewTicker(time.Duration(boundTherm.UpdateInterval) * time.Second)
		go func () {
			for {
				select {
				case <-ticker.C:
					reading := thermometer.Read()
					ps.Publish(readings.Topic, reading)
				}
			}
		}()
	}
}
