package sparta

import (
	"Sparta/explore"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sirupsen/logrus"
)

func exploreTestHelloWorld(event *json.RawMessage,
	context *LambdaContext,
	w http.ResponseWriter,
	logger *logrus.Logger) {
	logger.Info("Hello World: ", string(*event))

	fmt.Fprint(w, string(*event))
}

func TestExplore(t *testing.T) {
	// Create the function to test
	var lambdaFunctions []*LambdaAWSInfo
	lambdaFn := NewLambda(IAMRoleDefinition{}, exploreTestHelloWorld, nil)
	lambdaFunctions = append(lambdaFunctions, lambdaFn)

	// Mock event specific data to send to the lambda function
	eventData := ArbitraryJSONObject{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3"}

	// Make the request and confirm
	logger, _ := NewLogger("warning")
	ts := httptest.NewServer(NewLambdaHTTPHandler(lambdaFunctions, logger))
	defer ts.Close()
	resp, err := explore.NewRequest(lambdaFn.lambdaFnName, eventData, ts.URL)
	if err != nil {
		t.Fatal(err.Error())
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	t.Log("Status: ", resp.Status)
	t.Log("Headers: ", resp.Header)
	t.Log("Body: ", string(body))
}