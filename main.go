package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	dccHost      = flag.String("dccex-host", "0.0.0.0", "DCC-EX command station hostname or IP")
	dccPort      = flag.Int("dccex-port", 2560, "DCC-EX command station TCP port")
	port         = flag.Int("port", 9378, "Port to expose metrics on")
	pollInterval = time.Second
)

var (
	trackCurrent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dcc_ex_track_current_milliamps",
			Help: "Current draw of each DCC-EX track output in milliamps",
		},
		[]string{"track"},
	)
)

// Maintain a persistent TCP connection
func connect() net.Conn {
	for {
		addr := fmt.Sprintf("%s:%d", *dccHost, *dccPort)
		conn, err := net.Dial("tcp", addr)
		if err == nil {
			log.Printf("Connected to DCC-EX at %s\n", addr)
			return conn
		}
		log.Println("Reconnect in 2s:", err)
		time.Sleep(2 * time.Second)
	}
}

// Send a command and read until '>'
func sendCommand(conn net.Conn, cmd string) (string, error) {
	_, err := conn.Write([]byte(cmd + "\n"))
	if err != nil {
		return "", err
	}

	conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	reader := bufio.NewReader(conn)
	resp, err := reader.ReadString('>')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(resp), nil
}

// Parse <jI 399 210 ...>
func parseJI(resp string) ([]float64, bool) {
	resp = strings.Trim(resp, "<>")

	parts := strings.Fields(resp)
	if len(parts) < 2 || strings.ToLower(parts[0]) != "ji" {
		return nil, false
	}

	values := make([]float64, 0, len(parts)-1)

	for _, p := range parts[1:] {
		var ma float64
		if _, err := fmt.Sscanf(p, "%f", &ma); err == nil {
			values = append(values, ma)
		}
	}

	return values, true
}

func pollLoop() {
	conn := connect()

	for {
		resp, err := sendCommand(conn, "<J I>")
		if err != nil {
			conn = connect()
			continue
		}

		if values, ok := parseJI(resp); ok {
			for i, ma := range values {
				trackName := fmt.Sprintf("track%d", i+1)
				trackCurrent.WithLabelValues(trackName).Set(ma)
			}
		}

		time.Sleep(pollInterval)
	}
}

func main() {
	flag.Parse()

	if *dccHost == "0.0.0.0" {
		log.Fatal("Please specify the DCC-EX host with --host")
	}

	if *dccPort <= 0 || *dccPort > 65535 {
		log.Fatal("Please specify a valid port number with --port")
	}

	prometheus.MustRegister(trackCurrent)

	go pollLoop()

	http.Handle("/metrics", promhttp.Handler())
	log.Println("Exporter running on " + fmt.Sprintf(":%d", *port) + "/metrics")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
