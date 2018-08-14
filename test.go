package main

import (
	"fmt"
	"net/http"
)

func GetTests(httpClient *http.Client, bucket *Bucket, testsChannel chan bool) {
	url := fmt.Sprintf(GetTestsUrlPattern, bucket.BucketKey)
	var tests Tests
	GetData(httpClient, url, &tests)
	// log.Printf("%+v", testsCollection)
	bucket.Tests = tests
	testsChannel <- true
}

func GetTestsMetrics(httpClient *http.Client, bucket *Bucket, testsChannel chan bool) {
	if len(bucket.Environments.Data) == 0 {
		// No envs dont care and no metrics required
		testsChannel <- true
		return
	}
	// Should have only one prod enviornment
	env := bucket.Environments.Data[0]

	noOfTests := len(bucket.Tests.Data)
	testChannel := make(chan bool)
	for i := range bucket.Tests.Data {
		go getTestMetrics(httpClient, &bucket.Tests.Data[i], env.Id, bucket.BucketKey, testChannel)
	}
	//Wait for all tests to complete
	for noOfTests > 0 {
		<-testChannel
		noOfTests--
	}
	// for _, test := range bucket.Tests.Data {
	// 	log.Printf("%+v", test.TestMetrics)
	// }
	testsChannel <- true
}

func getTestMetrics(httpClient *http.Client, test *Test, envID string, bucketID string, testChannel chan bool) {
	url := fmt.Sprintf(GetTestMetricsUrlPattern, bucketID, test.Id, MonthTimeFrame, envID)
	var testMetrics TestMetrics
	GetData(httpClient, url, &testMetrics)
	// log.Printf("For URL [ %+v ] Metrics are %+v", url, testMetrics)
	calcSuccessRateAndAvgRespTimeMs(&testMetrics)
	test.TestMetrics = testMetrics
	testChannel <- true
}

func calcSuccessRateAndAvgRespTimeMs(testMetrics *TestMetrics) {
	toatalNoResponseTimeToConsider := 0
	totalResponseTime := 0.0
	totalSuccessRatio := 0.0
	for _, responseTime := range testMetrics.ResponseTimes {
		// No Data available case
		if responseTime.SuccessRatio == 0 && responseTime.AvgRespTimeMs == 0 {
			continue
		}
		toatalNoResponseTimeToConsider++
		totalResponseTime = totalResponseTime + responseTime.AvgRespTimeMs
		totalSuccessRatio = totalSuccessRatio + responseTime.SuccessRatio
	}
	if toatalNoResponseTimeToConsider > 0 {
		testMetrics.AvgRespTimeMs = totalResponseTime / float64(toatalNoResponseTimeToConsider)
		testMetrics.SuccessRate = float64((totalSuccessRatio / float64(toatalNoResponseTimeToConsider)) * 100)
	}
}
