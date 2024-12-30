package nakama

import (
	"encoding/base64"
	"encoding/json"
	"net/url"
	"reflect"
)

// BuildFetchOptions constructs fetch options similar to the JavaScript version.
func BuildFetchOptions(method string, options map[string]interface{}, bodyJson string) (map[string]interface{}, error) {
	// Initialize fetchOptions with method and merge with provided options.
	fetchOptions := make(map[string]interface{})
	fetchOptions["method"] = method

	// Merge options into fetchOptions
	for key, value := range options {
		fetchOptions[key] = value
	}

	// Handle headers
	var headers map[string]string
	if optionsHeaders, ok := options["headers"]; ok {
		if h, ok := optionsHeaders.(map[string]string); ok {
			headers = h
		} else {
			return nil, &json.UnmarshalTypeError{Value: "headers", Type: reflect.TypeOf(map[string]string{})}
		}
	} else {
		headers = make(map[string]string)
	}
	fetchOptions["headers"] = headers

	// Set default headers if not already set
	if _, exists := headers["Accept"]; !exists {
		headers["Accept"] = "application/json"
	}
	if _, exists := headers["Content-Type"]; !exists {
		headers["Content-Type"] = "application/json"
	}

	// Remove headers with empty values
	for key, value := range headers {
		if value == "" {
			delete(headers, key)
		}
	}

	// Add body if provided
	if bodyJson != "" {
		fetchOptions["body"] = bodyJson
	}

	return fetchOptions, nil
}

// B64EncodeUnicode encodes a string in Base64, preserving Unicode characters.
func B64EncodeUnicode(str string) string {
	// Encode special characters into percent-encoded form, then Base64 encode.
	encoded := url.QueryEscape(str)
	return base64.StdEncoding.EncodeToString([]byte(encoded))
}

// B64DecodeUnicode decodes a Base64-encoded string with Unicode characters.
func B64DecodeUnicode(str string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}

	// Decode the percent-encoded string back into its original form.
	decoded, err := url.QueryUnescape(string(decodedBytes))
	if err != nil {
		return "", err
	}

	return decoded, nil
}
