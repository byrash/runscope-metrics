package main

import (
	"fmt"
	"net/http"
	"strings"
)

func GetEnvironments(httpClient *http.Client, bucket *Bucket, envsChannel chan bool) {
	url := fmt.Sprintf(GetEnvsUrlPattern, bucket.BucketKey)
	var environments Environments
	GetData(httpClient, url, &environments)
	// log.Printf("%+v", environmentsCollection)
	bucket.Environments = environments
	envsChannel <- true
}

func ClearNonProdEnvs(bucket *Bucket) {
	prodEnvs := []Environment{}
	for _, env := range bucket.Environments.Data {
		if strings.Contains(strings.ToLower(env.Name), ProdEnvKey) {
			prodEnvs = append(prodEnvs, env)
		}
	}
	bucket.Environments.Data = prodEnvs
}
