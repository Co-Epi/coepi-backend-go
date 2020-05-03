package server

import (
//	"crypto/tls"
//	"crypto/x509"
	"crypto/ed25519"
	"encoding/json"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/Co-Epi/cen-server/backend"
)

const (
	// adjust these below to your SSL Cert location
	
	sslBaseDir     = "/etc/letsencrypt/live/v1.api.coepi.org"
	sslKeyFileName = "privkey.pem"
	caFileName     = "fullchain.pem"


	// DefaultPort is the port which the CEN HTTP server is listening in on
	DefaultPort = 443

	// DefaultAddr is the addr which the CEN HTTP server is listening in on
	DefaultAddr = "172.31.12.128"

	// EndpointTCNReport is the name of the HTTP endpoint for GET/POST of TCNReport for v4
	EndpointTCNReport = "v4/tcnreport"

	// EndpointCENReport is the name of the HTTP endpoint for GET/POST of CENReport
	EndpointCENReport = "cenreport"

	// EndpointCENKeys is the name of the HTTP endpoint for GET CenKeys
	EndpointCENKeys = "cenkeys"
)

// Server manages HTTP connections
type Server struct {
	backend  *backend.Backend
	Handler  http.Handler
	HTTPPort uint16
}

// NewServer returns an HTTP Server to handle simple-api-process-flow https://github.com/Co-Epi/data-models/blob/master/simple-api-process-flow.md
func NewServer(httpPort uint16, connString string) (s *Server, err error) {
	s = &Server{
		HTTPPort: httpPort,
	}

	backend, err := backend.NewBackend(connString)
	if err != nil {
		return s, err
	}
	s.backend = backend

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.getConnection)
	s.Handler = mux
	err =  s.Start()
	return s, nil
}

func (s *Server) getConnection(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if strings.Contains(r.URL.Path, EndpointTCNReport) {
		if r.Method == http.MethodPost {
			s.postTCNReportHandler(w, r)
		} else {
			s.getTCNReportHandler(w, r)
		}
	} else if strings.Contains(r.URL.Path, EndpointCENReport) {
		if r.Method == http.MethodPost {
			s.postCENReportHandler(w, r)
		} else {
			s.getCENReportHandler(w, r)
		}
	} else if strings.Contains(r.URL.Path, EndpointCENKeys) {
		s.getCENKeysHandler(w, r)
	} else {
		s.homeHandler(w, r)
	}
}

// Start kicks off the HTTP Server
func (s *Server) Start() (err error) {
	addrport := fmt.Sprintf("%s:%d", DefaultAddr, s.HTTPPort)
	srv := &http.Server{
		Addr:         addrport,
		Handler:      s.Handler,
		ReadTimeout:  600 * time.Second,
		WriteTimeout: 600 * time.Second,
	}

	SSLKeyFile := path.Join(sslBaseDir, sslKeyFileName)
	CAFile := path.Join(sslBaseDir, caFileName)

//	// Note: bringing the intermediate certs with CAFile into a cert pool and the tls.Config is *necessary*
//	certpool := x509.NewCertPool() // https://stackoverflow.com/questions/26719970/issues-with-tls-connection-in-golang -- instead of x509.NewCertPool()
//	pem, err := ioutil.ReadFile(CAFile)
//	if err != nil {
//		return fmt.Errorf("Failed to read client certificate authority: %v", err)
//	}
//	if !certpool.AppendCertsFromPEM(pem) {
//		return fmt.Errorf("Can't parse client certificate authority")
//	}
//
//	config := tls.Config{
//		ClientCAs:  certpool,
//		ClientAuth: tls.NoClientCert, // tls.RequireAndVerifyClientCert,
//	}
//	config.BuildNameToCertificate()
//
//	srv.TLSConfig = &config
//
	err = srv.ListenAndServeTLS(CAFile, SSLKeyFile)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) curtime() (t string) {
	currentTime := time.Now()
	return fmt.Sprintf("%s", currentTime.Format("2006-01-02 15:04:05"))
}

// POST /tcnreport/v0.4.0
func (s *Server) postTCNReportHandler(w http.ResponseWriter, r *http.Request) {
	// Read Post Body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	r.Body.Close()

	// fmt.Printf("%s: POST /tcnreport/v0.4.0: %s\n", s.curtime(), string(body))

	// Parse body as TCNReport
	var payload backend.TCNReport

	// keep the base64-encoded message, pass it to the handler as-is
	copy(payload.Report,body)

	// Need to decode whole message into a buffer
	decodedMessage, err := base64.StdEncoding.DecodeString( string(body) );
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lenDecoded := len( decodedMessage )

	// Use slices to reference rvk and sig, and use them to validate the sig of everything but the sig.
	// ((rvk)(everything-else))(sig)

	// ed25519.Verify(public, message, sig)
	if !(ed25519.Verify(decodedMessage[:32] , decodedMessage[:lenDecoded-64], decodedMessage[lenDecoded-64:])) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}


	// Process TCNReport payload and rvk (rvk used to create primary key only)
	err = s.backend.ProcessTCNReport(&payload, decodedMessage[:32])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write([]byte("OK"))
}

// GET /tcnreport?epochDay=<epochDay>&intervalNumber=<intervalNumber>&intervalLength=<intervalLength>
func (s *Server) getTCNReportHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Printf("%s: GET %s Request\n", s.curtime(), r.URL.Path)

	// tcnKey := ""
	// pathpieces := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	// if len(pathpieces) >= 1 {
	// 	tcnKey = pathpieces[1]
	// } else {
	// 	http.Error(w, "Usage: Usage: /tcnreport?epochDay=<n>&intervalNumber=<n>&intervalLength=<n>", http.StatusBadRequest)
	// 	return
	// }

	// Handle parameters
	q := r.URL.Query()
	epochDay := q.Get("epochDay")
	intervalNumber := q.Get("intervalNumber")
	intervalLength := q.Get("intervalLength")

	// pass parameters as arguments
	reports, err := s.backend.ProcessGetTCNReport(epochDay,intervalNumber,intervalLength)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// FIXME encode reports with base64 when returning them
	// some sort of encoder? iterate over list?
	responsesJSON, err := json.Marshal(reports)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// fmt.Printf("%s: GET %s Response: %s\n", s.curtime(), r.URL.Path, responsesJSON)
	// FIXME change what the Write has as an argument !!!!
	w.Write(responsesJSON)
}

// POST /cenreport
func (s *Server) postCENReportHandler(w http.ResponseWriter, r *http.Request) {
	// Read Post Body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	r.Body.Close()

	// fmt.Printf("%s: POST /cenreport: %s\n", s.curtime(), string(body))

	// Parse body as CENReport
	var payload backend.CENReport
	err = json.Unmarshal(body, &payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Process CENReport payload
	err = s.backend.ProcessCENReport(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write([]byte("OK"))
}

// GET /cenreport/<cenkey>
func (s *Server) getCENReportHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Printf("%s: GET %s Request\n", s.curtime(), r.URL.Path)

	cenKey := ""
	pathpieces := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathpieces) >= 1 {
		cenKey = pathpieces[1]
	} else {
		http.Error(w, "Usage: Usage: /cenreport/<cenkey>", http.StatusBadRequest)
		return
	}

	// Handle CenKey
	reports, err := s.backend.ProcessGetCENReport(cenKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	responsesJSON, err := json.Marshal(reports)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// fmt.Printf("%s: GET %s Response: %s\n", s.curtime(), r.URL.Path, responsesJSON)
	w.Write(responsesJSON)
}

// GET /cenkeys/<timestamp>
func (s *Server) getCENKeysHandler(w http.ResponseWriter, r *http.Request) {
	ts := uint64(0)
	// fmt.Printf("%s: GET %s Request: %s\n", s.curtime(), r.URL.Path)

	pathpieces := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathpieces) > 1 {
		tsa, err := strconv.Atoi(pathpieces[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ts = uint64(tsa)
	} else {
		ts = uint64(time.Now().Unix()) - 3600
	}

	cenKeys, err := s.backend.ProcessGetCENKeys(ts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	responsesJSON, err := json.Marshal(cenKeys)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// fmt.Printf("%s: GET %s Response: %s\n", s.curtime(), r.URL.Path, responsesJSON)
	w.Write(responsesJSON)
}

func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("CEN API Server v0.2"))
}
