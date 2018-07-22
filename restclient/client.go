package restclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// RestClient is the interface through which Handlers can make http calls
// Each handler should iniitialize it's own RestClient
type RestClient struct {
	Client         *http.Client
	DefaultHeaders http.Header
}

// FetchResponse is a convenience struct for storing relevant http response info
type FetchResponse struct {
	StatusCode int
	Status     string
	Body       []byte
}

// ErrorResponse is a convenience struct for storing
type ErrorResponse struct {
	Type   string   `json:"type"`
	Errors []string `json:"errors"`
}

func (r *RestClient) mergeHeaders(src http.Header) http.Header {
	dest := r.DefaultHeaders

	for h, _ := range dest {
		if val, ok := src[h]; ok {
			dest[h] = val
		}
	}

	return dest
}

// New creates a RestClient
func New() *RestClient {
	return &RestClient{
		Client: &http.Client{},
		DefaultHeaders: http.Header{
			"Content-Type": []string{"application/json"},
		},
	}
}

// Get from a url
func (r *RestClient) Get(url string, request *http.Request, headers http.Header) (FetchResponse, error) {
	return r.fetch("GET", url, request, headers)
}

// Post to a url with a payload
func (r *RestClient) Post(url string, request *http.Request, headers http.Header) (FetchResponse, error) {
	return r.fetch("POST", url, request, headers)
}

// Put to a url with a body
func (r *RestClient) Put(url string, request *http.Request, headers http.Header) (FetchResponse, error) {
	return r.fetch("PUT", url, request, headers)
}

// Patch to a url with a body
func (r *RestClient) Patch(url string, request *http.Request, headers http.Header) (FetchResponse, error) {
	return r.fetch("PATCH", url, request, headers)
}

// Delete from a url with a body (optional)
func (r *RestClient) Delete(url string, request *http.Request, headers http.Header) (FetchResponse, error) {
	return r.fetch("DELETE", url, request, headers)
}

func CopyRequestBody(r *http.Request) (*bytes.Buffer, error) {
	var bodyBytes []byte
	bodyBytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Printf("copy body %s", fmt.Sprintf(err.Error()))
		return nil, err
	}

	return bytes.NewBuffer(bodyBytes), nil
}

// fetch makes an http call and returns a FetchResponse struct
func (r *RestClient) fetch(
	method string,
	url string,
	request *http.Request,
	headers http.Header,
) (FetchResponse, error) {

	body, err := CopyRequestBody(request)
	if err != nil {
		message := fmt.Sprintf("RestClient: error copying request body (%s)", method, url, err.Error())
		log.Printf(message)
		return FetchResponse{}, fmt.Errorf(message)
	}
	newRequest, err := http.NewRequest(method, url, body)
	newRequest.Header = r.mergeHeaders(headers)

	log.Printf("RestClient: %s %s", method, url)
	response, err := r.Client.Do(newRequest)
	defer response.Body.Close()

	if err != nil {
		message := fmt.Sprintf("RestClient: error doing %s %s. (%s)", method, url, err.Error())
		log.Printf(message)
		return FetchResponse{}, fmt.Errorf(message)
	}

	responseBody, err := ioutil.ReadAll(response.Body)

	if err != nil {
		message := fmt.Sprintf("RestClient: Error reading bytes. (%s)", err.Error())
		log.Printf(message)
		return FetchResponse{}, fmt.Errorf(message)
	}

	fetchResponse := FetchResponse{
		Status:     response.Status,
		StatusCode: response.StatusCode,
		Body:       responseBody,
	}

	return fetchResponse, nil
}

// newErrorResponse is a convenience method for create an ErrorResponse struct
func newErrorResponse(errorType string, errorMessage string) ErrorResponse {
	return ErrorResponse{
		Type:   errorType,
		Errors: []string{errorMessage},
	}
}

// WriteErrorResponse writes a JSON body error response to the client
func (r *RestClient) WriteErrorResponse(
	writer http.ResponseWriter,
	errorType string,
	errorMessage string,
) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusInternalServerError)

	errorResponse := newErrorResponse(errorType, errorMessage)
	if err := json.NewEncoder(writer).Encode(errorResponse); err != nil {
		http.Error(writer, "RestClient: Failed to encode error response", http.StatusInternalServerError)
		return
	}
}

// CheckForError returns a generic error response if the status code is above 400
func (r *RestClient) CheckForError(response FetchResponse) *ErrorResponse {
	if response.StatusCode < 400 {
		return nil
	} else if response.StatusCode == 404 {
		return &ErrorResponse{
			Type:   "not_found",
			Errors: []string{"Resource not found"},
		}
	} else if response.StatusCode < 500 {
		return &ErrorResponse{
			Type:   "bad_external_request",
			Errors: []string{"Bad external request"},
		}
	} else if response.StatusCode < 600 {
		return &ErrorResponse{
			Type:   "server_error",
			Errors: []string{"Server error occured"},
		}
	}

	return nil
}

// Response generates JSON response for endpoints. This ends the request
func (r *RestClient) WriteJSONResponse(
	writer http.ResponseWriter,
	response FetchResponse,
	responseObject interface{},
) {
	errorResponse := r.CheckForError(response)
	if errorResponse != nil {
		r.WriteErrorResponse(writer, errorResponse.Type, errorResponse.Errors[0])
		return
	}

	if err := json.Unmarshal(response.Body, &responseObject); err != nil {
		log.Printf("Couldn't unmarshall json (%v)", err.Error())
		r.WriteErrorResponse(writer, "server_error", "Unmarshaling JSON failed")
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(response.StatusCode)

	if err := json.NewEncoder(writer).Encode(responseObject); err != nil {
		http.Error(writer, "can't read body", http.StatusBadRequest)
		return
	}
}
