package server

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/Co-Epi/coepi-backend-go/backend"
)

const (
	// adjust these below to your SSL Cert location
	sslBaseDir     = "/etc/pki/tls/certs/wildcard/wolk.com-new"
	sslKeyFileName = "www.wolk.com.key"
	caFileName     = "www.wolk.com.bundle"

	// DefaultPort is the port which the CEN HTTP server is listening in on
	DefaultPort = 8080

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
	go s.Start()
	return s, nil
}

func (s *Server) getConnection(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if strings.Contains(r.URL.Path, EndpointCENReport) {
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
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.HTTPPort),
		Handler:      s.Handler,
		ReadTimeout:  600 * time.Second,
		WriteTimeout: 600 * time.Second,
	}

	SSLKeyFile := path.Join(sslBaseDir, sslKeyFileName)
	CAFile := path.Join(sslBaseDir, caFileName)

	// Note: bringing the intermediate certs with CAFile into a cert pool and the tls.Config is *necessary*
	certpool := x509.NewCertPool() // https://stackoverflow.com/questions/26719970/issues-with-tls-connection-in-golang -- instead of x509.NewCertPool()
	pem, err := ioutil.ReadFile(CAFile)
	if err != nil {
		return fmt.Errorf("Failed to read client certificate authority: %v", err)
	}
	if !certpool.AppendCertsFromPEM(pem) {
		return fmt.Errorf("Can't parse client certificate authority")
	}

	config := tls.Config{
		ClientCAs:  certpool,
		ClientAuth: tls.NoClientCert, // tls.RequireAndVerifyClientCert,
	}
	config.BuildNameToCertificate()

	srv.TLSConfig = &config

	err = srv.ListenAndServeTLS(CAFile, SSLKeyFile)
	if err != nil {
		return err
	}
	return nil
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
	cenKey := ""
	pathpieces := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathpieces) >= 1 {
		cenKey = pathpieces[1]
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
	fmt.Printf("getCENReportHandler: %s\n", responsesJSON)
	w.Write(responsesJSON)
}

// /cenkeys/<timestamp>
func (s *Server) getCENKeysHandler(w http.ResponseWriter, r *http.Request) {
	ts := uint64(0)
	pathpieces := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathpieces) >= 1 {
		tsa, err := strconv.Atoi(pathpieces[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ts = uint64(tsa)
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
	fmt.Printf("genCENKeysHandler: %s\n", responsesJSON)
	w.Write(responsesJSON)
}

func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("CEN API Server v0.2"))
}
