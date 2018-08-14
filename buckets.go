package main

import (
	"net/http"

	"github.com/pkg/errors"
)

func GetBuckets(httpClient *http.Client) Buckets {
	var buckets Buckets
	GetData(httpClient, GetBucketsUrl, &buckets)
	// log.Printf("%+v", buckets)
	getProdEnvsAndTestsInfo(httpClient, &buckets)
	// for _, bucket := range buckets.Data {
	// 	log.Printf("%+v\n\n", bucket)
	// }
	// log.Printf("Buckets after processing %v", len(bucketsCollection.Data))
	return buckets
}

func getProdEnvsAndTestsInfo(httpClient *http.Client, bucketsCollection *Buckets) {
	buckets := []Bucket{}
	totalNoOfBuckets := len(bucketsCollection.Data)
	// log.Printf("Totla No of Buckets %v \n", totalNoOfBuckets)
	bucketsChannel := make(chan bool)
	for _, bucket := range bucketsCollection.Data {

		go func(b Bucket) {
			// log.Println("Getting env & tests for " + bucket.ProjectName)
			envsChannel := make(chan bool)
			testsChannel := make(chan bool)
			go GetEnvironments(httpClient, &b, envsChannel)
			go GetTests(httpClient, &b, testsChannel)
			//Wait for both environments & test to be completed
			<-envsChannel
			//Once we get Env Channel go and clear non prod envs and wait
			//This is not an IO intensive and not blocking for naything so run sequentially
			ClearNonProdEnvs(&b)
			if len(b.Environments.Data) > 1 {
				panic(errors.New("More than one production environment exists for Project" + bucket.ProjectName))
			}
			<-testsChannel
			// log.Printf("%+v\n\n", bucket)
			buckets = append(buckets, b)
			bucketsChannel <- true
		}(bucket)

	}
	// for _, bucket := range buckets {
	// 	log.Printf("%+v\n\n", bucket)
	// }

	for totalNoOfBuckets > 0 {
		<-bucketsChannel
		totalNoOfBuckets--
	}
	bucketsCollection.Data = buckets
}

func FillTestsMetricsForBuckets(httpClient *http.Client, bucketsCollection *Buckets) {
	buckets := []Bucket{}
	totalNoOfBuckets := len(bucketsCollection.Data)
	bucketsChannel := make(chan bool)
	for _, bucket := range bucketsCollection.Data {

		go func(b Bucket) {
			testsChannel := make(chan bool)
			go GetTestsMetrics(httpClient, &b, testsChannel)
			<-testsChannel
			// log.Printf("%+v\n\n", bucket)
			buckets = append(buckets, b)
			bucketsChannel <- true
		}(bucket)

	}
	// for _, bucket := range buckets {
	// 	log.Printf("%+v\n\n", bucket)
	// }

	for totalNoOfBuckets > 0 {
		<-bucketsChannel
		totalNoOfBuckets--
	}
	bucketsCollection.Data = buckets
}
