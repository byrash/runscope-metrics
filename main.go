package main

import (
	"bufio"
	"crypto/tls"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {
	start := time.Now()
	defer handlePanic()
	proxyURL, _ := url.Parse(os.Getenv(ProxyUrlEnvVarKey))
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           http.ProxyURL(proxyURL),
		},
	}
	//Get All meta data about buckets
	// GetBuckets(httpClient)
	buckets := GetBuckets(httpClient)
	FillTestsMetricsForBuckets(httpClient, &buckets)
	//Get each buckets test metrics
	// for _, bucket := range buckets.Data {
	// 	if strings.EqualFold(bucket.ProjectName, "Sample Test Name") {
	// 		log.Printf("For Project --> %+v\n", bucket.ProjectName)
	// 		for _, test := range bucket.Tests.Data {
	// 			log.Printf("%+v --> Success Rate %+v and Avg Response Time %+v \n", test.Name, test.SuccessRate, test.AvgRespTimeMs)
	// 		}
	// 	}
	// }
	WriteStatsToExcel(&buckets)
	log.Println("Completed in", time.Since(start))
}

func handlePanic() {
	if err := recover(); err != nil {
		log.Printf("\n\n")
		log.Println(err)
		reader := bufio.NewReader(os.Stdin)
		log.Printf("\n\n")
		log.Println("Hit Enter to close this window")
		reader.ReadString('\n')
	}
}
