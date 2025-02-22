package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/andybalholm/brotli"
	"github.com/jsdelivr/globalping-cli/model"
)

var ApiUrl = "https://api.globalping.io/v1/measurements"

// Post measurement to Globalping API - boolean indicates whether to print CLI help on error
func PostAPI(measurement model.PostMeasurement) (model.PostResponse, bool, error) {
	// Format post data
	postData, err := json.Marshal(measurement)
	if err != nil {
		return model.PostResponse{}, false, errors.New("err: failed to marshal post data - please report this bug")
	}

	// Create a new request
	req, err := http.NewRequest("POST", ApiUrl, bytes.NewBuffer(postData))
	if err != nil {
		return model.PostResponse{}, false, errors.New("err: failed to create request - please report this bug")
	}
	req.Header.Set("User-Agent", userAgent())
	req.Header.Set("Accept-Encoding", "br")
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return model.PostResponse{}, false, errors.New("err: request failed - please try again later")
	}
	defer resp.Body.Close()

	// If an error is returned
	if resp.StatusCode != http.StatusAccepted {
		// Decode the response body as JSON
		var data model.PostError

		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return model.PostResponse{}, false, errors.New("err: invalid error format returned - please report this bug")
		}

		// 422 error
		if data.Error.Type == "no_probes_found" {
			return model.PostResponse{}, true, errors.New("no suitable probes found - please choose a different location")
		}

		// 400 error
		if data.Error.Type == "validation_error" {
			for _, v := range data.Error.Params {
				fmt.Printf("err: %s\n", v)
			}
			return model.PostResponse{}, true, errors.New("invalid parameters - please check the help for more information")
		}

		// 500 error
		if data.Error.Type == "api_error" {
			return model.PostResponse{}, false, errors.New("err: internal server error - please try again later")
		}

		// If the error type is unknown
		return model.PostResponse{}, false, fmt.Errorf("err: unknown error response: %s", data.Error.Type)
	}

	// Read the response body

	var bodyReader io.Reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "br" {
		bodyReader = brotli.NewReader(bodyReader)
	}

	var data model.PostResponse
	err = json.NewDecoder(bodyReader).Decode(&data)
	if err != nil {
		fmt.Println(err)
		return model.PostResponse{}, false, errors.New("err: invalid post measurement format returned - please report this bug")
	}

	return data, false, nil
}

func DecodeTimings(cmd string, timings json.RawMessage) (model.Timings, error) {
	var data model.Timings

	if cmd == "ping" {
		err := json.Unmarshal(timings, &data.Arr)
		if err != nil {
			return model.Timings{}, errors.New("invalid timings format returned (ping)")
		}
	} else {
		err := json.Unmarshal(timings, &data.Interface)
		if err != nil {
			return model.Timings{}, errors.New("invalid timings format returned (other)")
		}
	}

	return data, nil
}
