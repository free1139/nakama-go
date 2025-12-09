package nakama

import (
	"encoding/base64"
	"encoding/json"
	"net/url"
	"reflect"
	"strconv"
	"time"
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

// Helper function to convert time.Time to *string
func timeToStringPointer(t *time.Time, layout string) *string {
	if t == nil {
		return nil
	}
	formattedTime := t.Format(layout)
	return &formattedTime
}

// Helper function to convert *string to *int
func stringPointerToIntPointer(s *string) *int {
	if s == nil {
		return nil
	}
	value, err := strconv.Atoi(*s)
	if err != nil {
		return nil
	}
	return &value
}

// Helper function to convert *int to *string
func intPointerToStringPointer(i *int) *string {
	if i == nil {
		return nil
	}
	value := strconv.Itoa(*i)
	return &value
}

// Helper function to convert *int64 to *int
func int64PointerToIntPointer(i *int64) *int {
	if i == nil {
		return nil
	}
	value := int(*i)
	return &value
}

// Helper function to convert a value to JSON string
func ToJSON(value interface{}) []byte {
	data, err := json.Marshal(value)
	if err != nil {
		panic(err) // Handle properly depending on your use case.
	}
	return data
}
