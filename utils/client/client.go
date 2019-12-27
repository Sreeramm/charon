package client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/charon/errors"
)

// DoGET Do a get, this method will take care of all logging and all
func DoGET(url string, body map[string]string, headers map[string]string) ([]byte, int, errors.Error) {
	return doRead(http.MethodGet, url, body, headers, &http.Transport{})
}

// DoGETWithoutTLS Do a POST, this method will perform a post call without the TLS certification
func DoGETWithoutTLS(url string, body map[string]string, headers map[string]string) ([]byte, int, errors.Error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return doRead(http.MethodGet, url, body, headers, tr)
}

// DoPOST Do a POST, this method will take care of all logging and all
func DoPOST(url string, body []byte, headers map[string]string) ([]byte, int, errors.Error) {
	return doWrite(http.MethodPost, url, body, headers, &http.Transport{})
}

// DoPOSTWithoutTLS Do a POST, this method will perform a post call without the TLS certification
func DoPOSTWithoutTLS(url string, body []byte, headers map[string]string) ([]byte, int, errors.Error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return doWrite(http.MethodPost, url, body, headers, tr)
}

// DoPOSTWithCert Do a POST, this method will perform a post call including the specified certificates in the certpool
func DoPOSTWithCert(url string, body []byte, headers map[string]string, certFilePath string, backupWithInsecure bool) ([]byte, int, errors.Error) {
	// Get the SystemCertPool, continue with an empty pool on error.
	rootCAs := x509.NewCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	// Read in the cert file.
	certs, err := ioutil.ReadFile(certFilePath)
	if err != nil {
		return nil, 500, errors.InternalError{Err: err.Error()}
	}
	_ = rootCAs.AppendCertsFromPEM(certs)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			ClientCAs:          rootCAs,
			InsecureSkipVerify: backupWithInsecure,
		},
	}
	return doWrite(http.MethodPost, url, body, headers, tr)
}

// DoPUT Do a PUT, this method will take care of all logging and all
func DoPUT(url string, body []byte, headers map[string]string) ([]byte, int, errors.Error) {
	return doWrite(http.MethodPut, url, body, headers, &http.Transport{})
}

// DoPATCH Do a PATCH, this method will take care of all logging and all
func DoPATCH(url string, body []byte, headers map[string]string) ([]byte, int, errors.Error) {
	return doWrite(http.MethodPatch, url, body, headers, &http.Transport{})
}

// DoDELETE Do a DELETE, this method will take care of all logging and all
func DoDELETE() {

}

// internal method to do write operation
func doWrite(method string, url string, body []byte, headers map[string]string, tr *http.Transport) ([]byte, int, errors.Error) {
	//creating a new request
	request, err := http.NewRequest(method, url, bytes.NewBuffer(body))

	if err != nil {
		// fmt.Println("Error:  ", err.Error(), "  ", time.Now())
		return nil, 400, errors.InternalError{Err: err.Error()}
	}

	return doCall(request, method, url, headers, tr)
}

//internal methond to do read operation
func doRead(method string, url string, body map[string]string, headers map[string]string, tr *http.Transport) ([]byte, int, errors.Error) {

	//creating a new request
	request, err := http.NewRequest(method, url, nil)

	// request.SetBasicAuth("api_bo9pace", "bo9@2016Pace")
	if err != nil {
		// fmt.Println("Error while creating a request:  ", err.Error(), "  ", time.Now())
		return nil, 400, errors.InternalError{Err: err.Error()}
	}

	// adding GET params as query string for get request
	if body != nil {
		q := request.URL.Query()
		for key, val := range body {
			q.Add(key, val)
		}
		request.URL.RawQuery = q.Encode()
	}
	return doCall(request, method, url, headers, tr)
}

func doCall(request *http.Request, method string, url string, headers map[string]string, tr *http.Transport) ([]byte, int, errors.Error) {

	// fmt.Println("Attempting to do an external ", method, " call on url ", url, " at ", time.Now())

	timeout := time.Duration(15 * time.Second)
	//initialising client with a timeout of 5 secs
	client := http.Client{
		Timeout: timeout,
	}
	if tr.TLSClientConfig != nil {
		client.Transport = tr
	}

	//assigning headers to the request
	if headers != nil {
		for k, v := range headers {
			request.Header.Set(k, v)
		}
	}

	resp, err := client.Do(request)

	if err != nil {
		// fmt.Println("Error:  ", err.Error(), "  ", time.Now())
		return nil, 400, errors.InternalError{Err: err.Error()}
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// fmt.Println("Error:  ", err.Error(), "  ", time.Now())
		return nil, 400, errors.InternalError{Err: err.Error()}
	}

	return respBody, resp.StatusCode, nil
}

//GetAsURL forms a url with the
func GetAsURL(protocol string, host string, port string, path string) (string, errors.Error) {
	if protocol != "http" && protocol != "https" {
		return "", errors.InvalidInputError{Err: "Invalid protocol"}
	}
	if host == "" {
		return "", errors.InvalidInputError{Err: "No value for host provided"}
	}
	url := protocol + "://" + host
	if port != "" {
		url += ":" + port
	}
	if path != "" {
		url += path
	} else {
		url += "/"
	}
	return url, nil
}
