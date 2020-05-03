package backend

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

const (
	// TableCENKeys stores the mapping between CENKeys and CENReports.
	TableCENKeys = "CENKeys"

	// TableCENKeys stores the mapping between CENKeys and CENReports.
	TableCENReport = "CENReport"

	// TCNReports is the only table for TCN
	TableTCNReport = "TCNReport"

	// Default Conn String 
	DefaultConnString = "root:CoEpi@/conn"
)

// Backend holds a client to connect  to the BigTable backend
type Backend struct {
	db *sql.DB
}

// CENReport payload is sent by client to /cenreport when user reports symptoms
type CENReport struct {
	ReportID        string `json:"reportID,omitempty"`
	Report          []byte `json:"report,omitempty"`  // this is expected to be a JSON blob but the server doesn't need to parse it
	CENKeys         string `json:"cenKeys,omitempty"` // comma separated list of hex AES-128 Keys
	ReportMimeType  string `json:"reportMimeType,omitempty"`
	ReportTimeStamp uint64 `json:"reportTimeStamp,omitempty"`
}

// TCNReport payload is sent by client to /tcnreport when user reports symptoms
// TCNReport is the original base64-encoded report

// raw bytes of appropriate sizes
// rvk32 tck32 j1 j2 TLV* sig64
// rvk || tck_{j1-1} || le_u16(j1) || le_u16(j2) || type || memo || sig
type TCNReport struct {
	Report          []byte				 // this is a blob ; parsing done only for signature validation on POST
}


// NewBackend sets up a client connection to BigTable to manage incoming payloads
func NewBackend(mysqlConnectionString string) (backend *Backend, err error) {
	backend = new(Backend)
	backend.db, err = sql.Open("mysql", mysqlConnectionString)
	if err != nil {
		return backend, err
	}

	return backend, nil
}

// ProcessTCNReport manages the API Endpoint to POST /v4/tcnreport
//  Input: TCNReport, rvk
//  Output: error
//  Behavior: write report bytes to "report" table
func (backend *Backend) ProcessTCNReport(tcnReport *TCNReport, tcnRVK []byte) (err error) {
	// WAS reportData, err := json.Marshal(tcnReport)

	// put the TCNReport in TCNReport table
	sReport := "insert into TCNReport (reportVK, report, reportTS) values ( ?, ?, ? ) on duplicate key update report = values(report)"
	stmtReport, err := backend.db.Prepare(sReport)
	if err != nil {
		return err
	}

	// store the tcnreportID in tcnReport table, one row per key

	// TimeStamp is epoch milliseconds
	currentTime := time.Now()
	var TimeStamp uint64
	TimeStamp = currentTime.Unix * 1000
	_, err = stmtReport.Exec(tcnRVK, tcnReport.Report, TimeStamp)
	if err != nil {
		panic(5)
		return err
	}

	return nil
}

// ProcessCENReport manages the API Endpoint to POST /cenreport
//  Input: CENReport
//  Output: error
//  Behavior: write report bytes to "report" table; write row for each CENKey with reportID
func (backend *Backend) ProcessCENReport(cenReport *CENReport) (err error) {
	reportData, err := json.Marshal(cenReport)
	if err != nil {
		return err
	}

	// put the CENReport in CENKeys table
	sKeys := "insert into CENKeys (cenKey, reportID, reportTS) values ( ?, ?, ? ) on duplicate key update reportTS = values(reportTS)"
	stmtKeys, err := backend.db.Prepare(sKeys)
	if err != nil {
		return err
	}

	// put the CENReport in CENReport table
	sReport := "insert into CENReport (reportID, report, reportMimeType, reportTS) values ( ?, ?, ?, ? ) on duplicate key update report = values(report)"
	stmtReport, err := backend.db.Prepare(sReport)
	if err != nil {
		return err
	}

	reportID := fmt.Sprintf("%x", Computehash(reportData))
	cenKeys := strings.Split(cenReport.CENKeys, ",")
	// store the cenreportID in cenkeys table, one row per key
	for _, cenKey := range cenKeys {
		cenKey := strings.Trim(cenKey, " \n")
		if len(cenKey) > 62 && len(cenKey) <= 64 {
			_, err = stmtKeys.Exec(cenKey, reportID, cenReport.ReportTimeStamp)

			if err != nil {
				return err
			}
		}
	}

	// store the cenreportID in cenReport table, one row per key
	_, err = stmtReport.Exec(reportID, cenReport.Report, cenReport.ReportMimeType, cenReport.ReportTimeStamp)
	if err != nil {
		panic(5)
		return err
	}

	return nil
}

// ProcessGetCENKeys manages the GET API endpoint /cenkeys
//  Input: timestamp
//  Output: array of CENKeys (in string form) for the last hour
func (backend *Backend) ProcessGetCENKeys(timestamp uint64) (cenKeys []string, err error) {
	cenKeys = make([]string, 0)

	s := "select cenKey From CENKeys where ReportTS >= 0" // TODO: ReportTS > ? and ReportTS <= ?"
	stmt, err := backend.db.Prepare(s)
	if err != nil {
		return cenKeys, err
	}
	rows, err := stmt.Query() // TODO: timestamp-3600, timestamp
	if err != nil {
		return cenKeys, err
	}
	for rows.Next() {
		var cenKey string
		err = rows.Scan(&cenKey)
		if err != nil {
			return cenKeys, err
		}
		cenKeys = append(cenKeys, cenKey)
	}
	return cenKeys, nil
}

// ProcessGetTCNReport manages the GET API endpoint /v4/tcnreport
//  Input: epochDay, intervalNumber, intervalLength
//  Output: array of TCNReports, already encoded as base64, in a list
func (backend *Backend) ProcessGetTCNReport(epochDay string, intervalNumber string, intervalLength string) (reports []*TCNReport, err error) {
	reports = make([]*TCNReport, 0)

// FIXME fix the "where" clause to use TS calculation from date, intervalNumber, intervalLength
// NB: IntervalNumber is relative to date
// NB: date is now "epochDay", days since start of epoch
	s := fmt.Sprintf("select TCNReport.report, TCNReport.reportTS From TCNReport where TCNReport.TS >= ((? * 86400) + (? * ?)) and TCNReport.TS <= ((? * 86400) + ((? + 1) * ?))")
	stmt, err := backend.db.Prepare(s)
	if err != nil {
		return reports, err
	}
	rows, err := stmt.Query(epochDay, intervalNumber, intervalLength, epochDay, intervalNumber, intervalLength)
	if err != nil {
		return reports, err
	}
	for rows.Next() {
		var r TCNReport
		err = rows.Scan(&(r.Report))
		if err != nil {
			return reports, err
		}
		reports = append(reports, &r)
	}
	return reports, nil
}

// ProcessGetCENReport manages the POST API endpoint /cenreport
//  Input: cenKey
//  Output: array of CENReports
func (backend *Backend) ProcessGetCENReport(cenKey string) (reports []*CENReport, err error) {
	reports = make([]*CENReport, 0)

	s := fmt.Sprintf("select CENKeys.reportID, report, reportMimeType, CENReport.reportTS From CENKeys, CENReport where CENKeys.CENKey = ? and CENKeys.reportID = CENReport.reportID")
	stmt, err := backend.db.Prepare(s)
	if err != nil {
		return reports, err
	}
	rows, err := stmt.Query(cenKey)
	if err != nil {
		return reports, err
	}
	for rows.Next() {
		var r CENReport
		err = rows.Scan(&(r.ReportID), &(r.Report), &(r.ReportMimeType), &(r.ReportTimeStamp))
		if err != nil {
			return reports, err
		}
		reports = append(reports, &r)
	}
	return reports, nil
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

func makeCENKeyString() string {
	key := make([]byte, 32)
	rand.Read(key)
	encoded := fmt.Sprintf("%x", key)
	return encoded
}

// GetSampleCENReportAndCENKeys generates a CENReport and an array of CENKeys (in string form)
func GetSampleCENReportAndCENKeys(nKeys int) (cenReport *CENReport, cenKeys []string) {
	cenKeys = make([]string, nKeys)
	for i := 0; i < nKeys; i++ {
		cenKeys[i] = makeCENKeyString()
	}
	CENKeys := fmt.Sprintf("%s,%s", cenKeys[0], cenKeys[1])
	curTS := uint64(time.Now().Unix())
	cenReport = new(CENReport)
	cenReport.ReportID = "1"
	cenReport.Report = []byte("severe fever,coughing,hard to breathe")
	cenReport.CENKeys = CENKeys
	cenReport.ReportTimeStamp = curTS
	return cenReport, cenKeys
}
