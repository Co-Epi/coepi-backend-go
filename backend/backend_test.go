package backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

func getJSONSamples() (eas *ExposureAndSymptoms, check1 *ExposureCheck, check0 *ExposureCheck) {
	eas = new(ExposureAndSymptoms)
	eas.Contacts = []Contact{Contact{UUIDHash: "ax", DateStamp: "2020-03-04"}, Contact{UUIDHash: "by", DateStamp: "2020-03-15"}, Contact{UUIDHash: "cz", DateStamp: "2020-03-20"}}
	eas.Symptoms = []byte("JSONBLOB:severe fever,coughing")

	check1 = new(ExposureCheck)
	check1.Contacts = []Contact{Contact{UUIDHash: "by", DateStamp: "2020-03-04"}}

	check0 = new(ExposureCheck)
	check0.Contacts = []Contact{Contact{UUIDHash: "00", DateStamp: "2020-03-21"}}
	return eas, check1, check0
}

func TestJSONSample(t *testing.T) {
	eas, check1, _ := getJSONSamples()
	easJSON, err := json.Marshal(eas)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	fmt.Printf("ExposureAndSymptoms Sample: %s\n", easJSON)
	check1JSON, err := json.Marshal(check1)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	fmt.Printf("ExposureCheck Sample: %s\n", check1JSON)
}

func TestBackendSimple(t *testing.T) {
	backend, err := NewBackend(DefaultProject, DefaultInstance)
	if err != nil {
		t.Fatalf("%s", err)
	}

	eas, check1, check0 := getJSONSamples()
	err = backend.ProcessExposureAndSymptoms(eas)
	if err != nil {
		t.Fatalf("ProcessExposureAndSymptoms: %s", err)
	}

	response, err := backend.ProcessExposureCheck(check1)
	if err != nil {
		t.Fatalf("ProcessExposureCheck(check1): %s", err)
	}
	if len(response.Exposures) != 1 {
		t.Fatalf("ProcessExposureCheck(check1) Expected 1 response, got %d", len(response.Exposures))
	}

	exposure := response.Exposures[0]
	if !bytes.Equal(exposure.Symptoms, eas.Symptoms) {
		t.Fatalf("ProcessExposureCheck(check1) Symptoms Mismatch: expected %s, got %s", eas.Symptoms, exposure.Symptoms)
	}
	fmt.Printf("ProcessExposureCheck(check1) SUCCESS: [%s]\n", exposure.Symptoms)

	response, err = backend.ProcessExposureCheck(check0)
	if err != nil {
		t.Fatalf("ProcessExposureCheck(check0): %s", err)
	}
	if len(response.Exposures) != 0 {
		t.Fatalf("processExposureCheck(check0) Expected 0 responses, but got %d", len(response.Exposures))
	}
	fmt.Printf("processExposureCheck(check0) SUCCESS: []\n")
}

// TODO: setup N random users and have them generate Contact records between themselves in N goroutines, each doing a ExposureCheck every few seconds
//  Then have N/10 users post a ExposureAndSymptoms, which should have some of the go routines generating symptom responses
