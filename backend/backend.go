package backend

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/bigtable"
)

const (
	// DefaultProject is your GC Project Name
	DefaultProject = "us-west1-wlk"
	// DefaultInstance is your GC BigTable Instance name
	DefaultInstance = "co-epi"
)

// Backend holds a client to connect to the BigTable backend
type Backend struct {
	client *bigtable.Client
}

// Contact represents a BLE pairing between 2 devices
type Contact struct {
	UUIDHash  string `json:"uuidHash,omitempty"`  // this is the SHA256 hash of a pair of BLE UUID ids
	DateStamp string `json:"dateStamp,omitempty"` // this is when the pair came into contact, used to set TTLs
}

// ExposureAndSymptoms payload is sent by client to /exposureandsymptoms when user reports symptoms
type ExposureAndSymptoms struct {
	ClientID string    `json:"clientID,omitempty"` // TBD
	Symptoms []byte    `json:"symptoms,omitempty"` // this is expected to be a JSON blob but the server doesn't need to parse it
	Contacts []Contact `json:"contacts,omitempty"` // these are the contacts
}

// ExposureCheck payload is sent by client to /exposurecheck to check for symptoms
type ExposureCheck struct {
	Contacts []Contact `json:"contacts,omitempty"`
}

// ExposureCheckResponse  is sent by server in response /exposurecheck
type ExposureCheckResponse struct {
	Exposures []Exposure `json:"exposures,omitempty"`
}

// Exposure  is sent by server in response /exposurecheck
type Exposure struct {
	Symptoms []byte    `json:"symptoms,omitempty"`
	Contacts []Contact `json:"contacts,omitempty"`
}

// TableContacts stores the mapping between UUIDs and symptomHash.
const TableContacts = "contacts"

// TableSymptoms stores the mapping between symptomHash and symptoms.   The content of the symptoms string is a JSON document that clients need to power the UI but the server does not need to process it
const TableSymptoms = "symptoms"

// NewBackend sets up a client connection to BigTable to manage incoming payloads
func NewBackend(project, instance string) (backend *Backend, err error) {
	backend = new(Backend)
	client, err := bigtable.NewClient(context.Background(), project, instance)
	if err != nil {
		return backend, err
	}
	backend.client = client
	return backend, nil
}

// ProcessExposureAndSymptoms manages the API Endpoint to POST /exposureandsymptoms
//  Input: ExposureAndSymptoms
//  Output: error
func (backend *Backend) ProcessExposureAndSymptoms(payload *ExposureAndSymptoms) (err error) {
	// store symptoms in the symptoms table
	symptoms := payload.Symptoms
	symptomsTable := backend.client.Open(TableSymptoms)
	symptomsHash := Computehash(symptoms)
	mut := bigtable.NewMutation()
	mut.Set("case", "symptoms", bigtable.Now(), []byte(symptoms))
	err = symptomsTable.Apply(context.Background(), fmt.Sprintf("%x", symptomsHash), mut)
	if err != nil {
		return err
	}

	contactsTable := backend.client.Open(TableContacts)
	// store the first 64 one cell per observation
	for _, contact := range payload.Contacts {
		mut := bigtable.NewMutation()
		mut.Set("symptoms", contact.DateStamp, bigtable.Now(), []byte(fmt.Sprintf("%x", symptomsHash)))
		err = contactsTable.Apply(context.Background(), contact.UUIDHash, mut)
		if err != nil {
			return err
		}
	}
	return nil
}

// ProcessExposureCheck manages the POST API endpoint /exposurecheck
//  Input: ExposureCheck
//  Output: array of ExposureCheckResponse
func (backend *Backend) ProcessExposureCheck(payload *ExposureCheck) (response ExposureCheckResponse, err error) {
	tableContacts := backend.client.Open(TableContacts)
	// store one cell per observation
	symptomsHashes := make(map[string][]Contact, 0)
	for _, contact := range payload.Contacts {
		rr := bigtable.PrefixRange(contact.UUIDHash)
		err := tableContacts.ReadRows(context.Background(), rr, func(r bigtable.Row) bool {
			for k, xv := range r {
				switch k {
				case "symptoms":
					for _, yv := range xv {
						dt := strings.Split(yv.Column, ":") // symptoms:2020-03-15 => dt[1] = "2020-03-15"
						if len(dt) == 2 {
							t, err := time.Parse("2006-01-02", dt[1])
							if err == nil && time.Since(t) < 24*7*time.Hour {
								// fmt.Printf("MATCH: %s|%s\n", date, string(yv.Value))
								// Question: how can we filter on the right number of days without a disease lookup demanding a peek into the symptoms data?
							}
							symptomsHash := string(yv.Value)
							if _, ok := symptomsHashes[symptomsHash]; !ok {
								symptomsHashes[symptomsHash] = make([]Contact, 0)
							}
							symptomsHashes[symptomsHash] = append(symptomsHashes[symptomsHash], contact)
						} else {
							fmt.Printf("Problem: %d\n", len(dt))
						}
					}
				}
			}
			return true // Keep going.
		}, bigtable.RowFilter(bigtable.LatestNFilter(1)))
		if err != nil {
			// TODO: handle err.
		}
	}

	// For all symptomHashes, get Symptoms byte blobs (containing reported symptoms + severity, dates, etc.)
	exposures := make([]Exposure, 0)
	for symptomsHash, contacts := range symptomsHashes {
		tableSymptoms := backend.client.Open(TableSymptoms)
		rr := bigtable.PrefixRange(symptomsHash)
		err := tableSymptoms.ReadRows(context.Background(), rr, func(r bigtable.Row) bool {
			for k, xv := range r {
				switch k {
				case "case":
					for _, yv := range xv {
						if yv.Column == "case:symptoms" {
							var exposure Exposure
							exposures = append(exposures, exposure)
							exposure.Contacts = contacts
							exposure.Symptoms = yv.Value
						}
					}
				}
			}
			return true // Keep going.
		}, bigtable.RowFilter(bigtable.LatestNFilter(1)))
		if err != nil {
			// TODO: handle err.
		}
	}
	response.Exposures = exposures

	return response, nil
}

// Computehash returns the hash of its inputs
func Computehash(data ...[]byte) []byte {
	hasher := sha256.New()
	for _, b := range data {
		_, err := hasher.Write(b)
		if err != nil {
			panic(1)
		}
	}
	return hasher.Sum(nil)
}
