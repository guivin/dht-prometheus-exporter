package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	stdlibLog "log"
	"net/http"
)

var lg log.Logger

func main() {
	config := ReadConfig()
	lg = getLogger(config)
	w := lg.Writer()
	defer w.Close()
	sensor := newSensor(config)
	collector := newCollector(sensor)
	lg.Debug("Registering the prometheus collector")
	prometheus.MustRegister(collector)
	http.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{
		ErrorLog: stdlibLog.New(w, "", 0),
	}))
	lg.Info(fmt.Sprintf("Starting http server on TCP/%d port", config.listenPort))
	lg.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.listenPort), nil))
}
