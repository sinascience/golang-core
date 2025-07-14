package validator

import (
	"testing"
)

// TestValidateStruct tests our custom validator utility.
func TestValidateStruct(t *testing.T) {
	// Define a sample struct we want to validate, similar to our real payloads.
	type SamplePayload struct {
		Name  string `json:"name" validate:"required,min=2"`
		Email string `json:"email" validate:"required,email"`
	}

	// We use a "table-driven test" pattern, which is standard in Go.
	// It allows us to define a list of test cases easily.
	testCases := []struct {
		name          string      // The name of the test case
		payload       interface{} // The input data
		expectedError bool        // Whether we expect an error
		errorField    string      // Which field we expect the error on
	}{
		{
			name: "Success - Valid Payload",
			payload: &SamplePayload{
				Name:  "Venturo",
				Email: "test@venturo.dev",
			},
			expectedError: false,
			errorField:    "",
		},
		{
			name: "Failure - Missing Required Field",
			payload: &SamplePayload{
				Name:  "", // Name is empty, but required
				Email: "test@venturo.dev",
			},
			expectedError: true,
			errorField:    "name",
		},
		{
			name: "Failure - Invalid Email Format",
			payload: &SamplePayload{
				Name:  "Venturo",
				Email: "not-an-email", // Email format is invalid
			},
			expectedError: true,
			errorField:    "email",
		},
	}

	// Loop through all our defined test cases.
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Run the validator on the test case payload.
			errs := ValidateStruct(tc.payload)
			hasError := errs != nil

			// Check if the outcome matches what we expected.
			if hasError != tc.expectedError {
				t.Errorf("expected error: %v, but got: %v, errors: %v", tc.expectedError, hasError, errs)
			}

			// If we expected an error, check if the error is on the correct field.
			if tc.expectedError {
				if _, ok := errs[tc.errorField]; !ok {
					t.Errorf("expected error on field %s, but got none. errors: %v", tc.errorField, errs)
				}
			}
		})
	}
}
