package mazon

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/hiscaler/mazon-go/entity"
	"github.com/stretchr/testify/assert"
)

func TestScanFormService_Create(t *testing.T) {
	// This response struct is defined locally for testing purposes only.
	type mockCreateScanFormResponse struct {
		Code    int               `json:"code"`
		Message string            `json:"msg"`
		Result  []entity.ScanForm `json:"result"`
	}

	// Mock a server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/createScanForm" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var reqBody map[string]string
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		// Success case
		if reqBody["tracking_number"] == "TN123,TN456" {
			w.WriteHeader(http.StatusOK)
			response := mockCreateScanFormResponse{
				Code:    http.StatusOK,
				Message: "Success",
				Result:  []entity.ScanForm{{Url: "https://example.com/scanform.pdf"}},
			}
			_ = json.NewEncoder(w).Encode(response)
			return
		}

		// API error case
		w.WriteHeader(http.StatusOK)
		response := mockCreateScanFormResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request from API",
			Result:  nil,
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	// Create a service with a client pointing to the mock server
	httpClient := resty.New().SetBaseURL(mockServer.URL)
	// Since scanFormService is not exported, we need to create a client and get the service from it.
	c := &Client{httpClient: httpClient}
	s := scanFormService{httpClient: c.httpClient}

	// Define test cases
	tests := []struct {
		name            string
		trackingNumbers []string
		wantForms       []entity.ScanForm
		wantErr         bool
		wantErrMsg      string
	}{
		{
			name:            "Success",
			trackingNumbers: []string{"TN123", "  TN456  "},
			wantForms:       []entity.ScanForm{{Url: "https://example.com/scanform.pdf"}},
			wantErr:         false,
		},
		{
			name:            "Invalid Input - Empty Slice",
			trackingNumbers: []string{},
			wantForms:       nil,
			wantErr:         true,
			wantErrMsg:      ErrInvalidTrackingNumber.Error(),
		},
		{
			name:            "API Error",
			trackingNumbers: []string{"ANY_OTHER_NUMBER"},
			wantForms:       nil,
			wantErr:         true,
			wantErrMsg:      "400: Invalid request from API",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// The recheckError function needs to be accessible for the test to work.
			// We assume it's in the same package and handles the error conversion.
			forms, err := s.Create(context.Background(), tt.trackingNumbers...)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrMsg)
				assert.Nil(t, forms)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantForms, forms)
			}
		})
	}
}
