package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

func SetHeaders(httpReq *http.Request) {
	httpReq.Header.Set(ContentTypeKey, AppOrJSONContentType)
	httpReq.Header.Set(AuthKey, fmt.Sprintf(AuthStringPattern, os.Getenv(RunScopeSecretKeyEnvVarKey)))
}

func HandleError(err error, msg string) {
	if err != nil {
		panic(errors.Wrap(err, msg))
	}
}

func GetData(httpClient *http.Client, url string, requiredDataStruct interface{}) {
	httpReq, httpReqErr := http.NewRequest(GetMethod, url, nil)
	HandleError(httpReqErr, "Unable to create request object")
	SetHeaders(httpReq)
	httpResp, httpRespErr := httpClient.Do(httpReq)
	HandleError(httpRespErr, fmt.Sprintf("Unable to get from %+v", url))
	defer httpResp.Body.Close()
	response, readRespErr := ioutil.ReadAll(httpResp.Body)
	HandleError(readRespErr, "Unable to read responset")
	unMarshallErr := json.Unmarshal(response, requiredDataStruct)
	HandleError(unMarshallErr, "Unable to unmarshall JSON response")
}
