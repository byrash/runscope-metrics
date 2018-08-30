package main

import "os"

var (
	ProxyUrl             = os.Getenv(ProxyUrlEnvVarKey)
	RunScopeSecretKey    = os.Getenv(RunScopeSecretKeyEnvVarKey)
	CriticalAppsBucketId = os.Getenv(CriticalAppsBucketIdEnvVarKey)
	RestBucketId         = os.Getenv(RestBucketIdEnvVarKey)
)
