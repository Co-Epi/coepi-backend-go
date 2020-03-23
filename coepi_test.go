package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/Co-Epi/coepi-backend-go/backend"
	"github.com/Co-Epi/coepi-backend-go/server"
)

// DefaultTransport contains all HTTP client operation parameters
var DefaultTransport http.RoundTripper = &http.Transport{
	Dial: (&net.Dialer{
		// limits the time spent establishing a TCP connection (if a new one is needed)
		Timeout:   120 * time.Second,
		KeepAlive: 120 * time.Second, // 60 * time.Second,
	}).Dial,
	//MaxIdleConns: 5,
	MaxIdleConnsPerHost: 25, // changed from 100 -> 25

	// limits the time spent reading the headers of the response.
	ResponseHeaderTimeout: 120 * time.Second,
	IdleConnTimeout:       120 * time.Second, // 90 * time.Second,

	// limits the time the client will wait between sending the request headers when including an Expect: 100-continue and receiving the go-ahead to send the body.
	ExpectContinueTimeout: 1 * time.Second,

	// limits the time spent performing the TLS handshake.
	TLSHandshakeTimeout: 5 * time.Second,
}

func httppost(url string, body []byte) (result []byte, err error) {

	httpclient := &http.Client{Timeout: time.Second * 120, Transport: DefaultTransport}
	bodyReader := bytes.NewReader(body)
	req, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		return result, fmt.Errorf("[coepi_test:httppost] %s", err)
	}

	resp, err := httpclient.Do(req)
	if err != nil {
		return result, fmt.Errorf("[coepi_test:httppost] %s", err)
	}

	result, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("[coepi_test:httppost] %s", err)
	}
	resp.Body.Close()

	return result, nil
}

func TestCoepiSimple(t *testing.T) {
	endpoint := fmt.Sprintf("coepi.wolk.com:%d", server.DefaultPort)

	// Post 
	eas := new(backend.ExposureAndSymptoms)
	eas.Contacts = []backend.Contact{backend.Contact{UUIDHash: "ax", DateStamp: "2020-03-04"}, backend.Contact{UUIDHash: "by", DateStamp: "2020-03-15"}, backend.Contact{UUIDHash: "cz", DateStamp: "2020-03-20"}}
	eas.Symptoms = []byte("JSONBLOB:severe fever,coughing")
	easJSON, err := json.Marshal(eas)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	result, err := httppost(fmt.Sprintf("https://%s/%s", endpoint, server.EndpointExposureAndSymptoms), easJSON)
	if err != nil {
		t.Fatalf("exposureandsymptoms: %s", err)
	}
	fmt.Printf("exposureandsymptoms[%s]\n", string(result))

	check1 := new(backend.ExposureCheck)
	check1.Contacts = []backend.Contact{backend.Contact{UUIDHash: "by", DateStamp: "2020-03-04"}}
	check1JSON, err := json.Marshal(check1)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	exposureCheckResponseRaw, err := httppost(fmt.Sprintf("https://%s/%s", endpoint, server.EndpointExposureCheck), check1JSON)
	if err != nil {
		t.Fatalf("exposurecheck: %s", err)
	}
	var ecr backend.ExposureCheckResponse
	fmt.Printf("exposureCheckResponseRaw: %s\n", exposureCheckResponseRaw)
	err = json.Unmarshal(exposureCheckResponseRaw, &ecr)
	if err != nil {
		t.Fatalf("exposurecheck(check1): %s", err)
	}
	if len(ecr.Exposures) != 1 {
		t.Fatalf("exposurecheck(check1) Expected 1 response, got %d", len(ecr.Exposures))
	}

	exposure := ecr.Exposures[0]
	if !bytes.Equal(eas.Symptoms, exposure.Symptoms) {
		t.Fatalf("exposurecheck(check1) Expected %s, got %s", eas.Symptoms, exposure.Symptoms)
	}
	fmt.Printf("exposurecheck(check1) SUCCESS: [%s]\n", exposure.Symptoms)

	check0 := new(backend.ExposureCheck)
	check0.Contacts = []backend.Contact{backend.Contact{UUIDHash: "00", DateStamp: "2020-03-21"}}
	check0JSON, err := json.Marshal(check0)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	exposureCheckResponseRaw, err = httppost(fmt.Sprintf("https://%s/%s", endpoint, server.EndpointExposureCheck), check0JSON)
	if err != nil {
		t.Fatalf("exposurecheck: %s", err)
	}
	var ecr0 backend.ExposureCheckResponse
	err = json.Unmarshal(exposureCheckResponseRaw, &ecr0)
	if err != nil {
		t.Fatalf("exposurecheck(check1): %s", err)
	}
	if len(ecr0.Exposures) != 0 {
		t.Fatalf("exposurecheck(check0) Expected 0 responses, but got %d [%s]", len(ecr0.Exposures), exposureCheckResponseRaw)
	}
	fmt.Printf("exposurecheck(check0) SUCCESS: []\n")
}
