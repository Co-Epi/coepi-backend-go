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
		return result, fmt.Errorf("[cen_test:httppost] %s", err)
	}

	resp, err := httpclient.Do(req)
	if err != nil {
		return result, fmt.Errorf("[cen_test:httppost] %s", err)
	}

	result, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("[cen_test:httppost] %s", err)
	}
	resp.Body.Close()

	return result, nil
}

func httpget(url string) (result []byte, err error) {

	httpclient := &http.Client{Timeout: time.Second * 120, Transport: DefaultTransport}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return result, fmt.Errorf("[cen_test:httpget] %s", err)
	}

	resp, err := httpclient.Do(req)
	if err != nil {
		return result, fmt.Errorf("[cen_test:httpget] %s", err)
	}

	result, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("[cen_test:httpget] %s", err)
	}
	resp.Body.Close()

	return result, nil
}

func TestCENSimple(t *testing.T) {
	endpoint := fmt.Sprintf("coepi.wolk.com:%d", server.DefaultPort)

	// Post CENReport to /cenreport, along with CENKeys
	cenReport, cenReportKeys := backend.GetSampleCENReportAndCENKeys(2)
	cenReportJSON, err := json.Marshal(cenReport)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// POST CENReport
	result, err := httppost(fmt.Sprintf("https://%s/%s", endpoint, server.EndpointCENReport), cenReportJSON)
	if err != nil {
		t.Fatalf("EndpointCENReport: %s", err)
	}
	fmt.Printf("EndpointCENReport[%s]\n", string(result))

	// GET CENKeys
	curTS := uint64(time.Now().Unix())
	resp, err := httpget(fmt.Sprintf("https://%s/%s/%d", endpoint, server.EndpointCENKeys, curTS-10))
	if err != nil {
		t.Fatalf("EndpointCENKeys: %s", err)
	}
	fmt.Printf("EndpointCENKeys: %s\n", string(resp))

	var cenKeys []string
	err = json.Unmarshal(resp, &cenKeys)
	if err != nil {
		t.Fatalf("EndpointCENKeys(check1): [%s] [%s]", resp, err)
	}
	if len(cenKeys) < 2 {
		t.Fatalf("Incorrect response length [%d] -- should be at least 2", len(cenKeys))
	}
	found := make([]bool, len(cenKeys))
	for _, cenKey := range cenKeys {
		for j, reportKey := range cenReportKeys {
			if cenKey == reportKey {
				found[j] = true
			}
		}
	}

	// GET CENREport
	for i := 0; i < 2; i++ {
		if !found[i] {
			t.Fatalf("EndpointCENKey key %d in report [%s] not found", i, cenReportKeys[i])
		}

		cenKey := cenReportKeys[i]
		reportsRaw, err := httpget(fmt.Sprintf("https://%s/%s/%s", endpoint, server.EndpointCENReport, cenKey))
		if err != nil {
			t.Fatalf("EndpointCENReport: %s", err)
		}
		var reports []*backend.CENReport
		err = json.Unmarshal(reportsRaw, &reports)
		if err != nil {
			t.Fatalf("EndpointCENReport(check1): %s", err)
		}
		if len(reports) > 0 {
			report := reports[0]
			fmt.Printf("EndpointCENKeys SUCCESS: [%s]\n", report.Report)

			if !bytes.Equal(report.Report, cenReport.Report) {
				t.Fatalf("EndpointCENKeys(check1) Expected %s, got %s", cenReport.Report, report.Report)
			}
		} else {
			t.Fatalf("hmm, no report (%s)", cenKey)
		}
	}
}
