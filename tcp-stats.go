package main

import (
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

var version = ""
var hostname = ""
var uri = ""

func main() {
	fmt.Printf("Prometheus TCP Metrics Scraper, Written by Paul Schou  Source code: github.com/pschou/win-tcp-prom (version: %s)\n", version)
	if len(os.Args) < 2 {
		fmt.Printf("Please provide url to post to as argument.\n", version)
	} else {
		uri = os.Args[1]
	}
	hostname, _ = os.Hostname()
	fmt.Printf("Hostname = %s\n", hostname)
	fmt.Printf("URI = %s\n", uri)
	collect()
	for now := range time.Tick(5 * time.Minute) {
		fmt.Printf("  %q\n", now)
		collect()
	}
}

func collect() {
	fmt.Printf("Collecting TCP stats...\n")
	out, err := exec.Command("netstat", "-n", "-s", "-p", "tcp").Output()
	if err != nil {
		log.Fatal(err)
	}
	metrics := ""
	for _, line := range strings.Split(string(out), "\n") {
		parts := strings.SplitN(line, " = ", 2)
		if len(parts) == 2 {
			metric_name := strings.ReplaceAll(strings.TrimSpace(parts[0]), " ", "_")
			fmt.Printf("  %s %s\n", metric_name, parts[1])
			metrics += fmt.Sprintf("%s{hostname=%q} %s\n", metric_name, hostname, parts[1])
		}
	}

	fmt.Printf("Collecting Memory stats...\n")
	v, _ := mem.VirtualMemory()
	fmt.Printf("  Memory_Total %v\n  Memory_Free %v\n", v.Total, v.Free)
	metrics += fmt.Sprintf("Memory_Total{hostname=%q} %v\nMemory_Free{hostname=%q} %v\n", hostname, v.Total, hostname, v.Free)
	fmt.Printf("Collecting CPU stats...\n")
	c, _ := cpu.Percent(10*time.Second, true)
	for i, val := range c {
		fmt.Printf("  CPU_Load{core=\"%d\"} %f\n", i, val)
		metrics += fmt.Sprintf("CPU_Load{hostname=%q,core=\"%d\"} %f\n", hostname, i, val)
	}

	if uri != "" {
		SendPostRequest(metrics)
	}
}

func SendPostRequest(metrics string) {
	fmt.Printf("Sending metrics...\n")
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("data", "metrics.prom")
	if err != nil {
		fmt.Printf("Error when creating form file %g\n", err)
	}
	part.Write([]byte(metrics))

	//for key, val := range params {
	//	_ = writer.WriteField(key, val)
	//}
	err = writer.Close()
	if err != nil {
		fmt.Printf("Error when closing writer %g\n", err)
		return
	}

	request, err := http.NewRequest("POST", uri, body)
	if err != nil {
		fmt.Printf("Error when making http request %g\n", err)
		return
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Printf("Error in http do %g\n", err)
	}
	//fmt.Printf("Response: %g\n", resp)
	fmt.Printf("Sent!  Response: %d\n", resp.StatusCode)
}
