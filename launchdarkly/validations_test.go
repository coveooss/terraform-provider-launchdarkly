package launchdarkly

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestValidateKey(t *testing.T) {
	testCases := []struct {
		name      string
		v         interface{}
		k         string
		wantedErr []error
	}{
		{
			name:      "expected",
			v:         "k",
			k:         "a-key",
			wantedErr: nil,
		},
		{
			name:      "without character",
			v:         "",
			k:         "a-key",
			wantedErr: []error{errors.New(fmt.Sprintf("%s must be between 1 and 20 characters: %s", "a-key", ""))},
		},
		{
			name:      "with value more than 20 characters",
			v:         "a-very-long-value-that-exceeds-20-characters",
			k:         "a-key",
			wantedErr: []error{errors.New(fmt.Sprintf("%s must be between 1 and 20 characters: %s", "a-key", "a-very-long-value-that-exceeds-20-characters"))},
		},
		{
			name:      "with invalid match",
			v:         "(#*&$?(*@&$)",
			k:         "a-key",
			wantedErr: []error{errors.New(fmt.Sprintf("%s is not a valid key: %s", "a-key", "(#*&$?(*@&$)"))},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, errs := validateKey(testCase.v, testCase.k)
			testValidateVerifyGeneric(t, errs, testCase.wantedErr)
		})
	}
}

func TestValidateFeatureFlagKey(t *testing.T) {
	testCases := []struct {
		name      string
		v         interface{}
		k         string
		wantedErr []error
	}{
		{
			name:      "expected",
			v:         "k",
			k:         "a-key",
			wantedErr: nil,
		},
		{
			name:      "with invalid match",
			v:         "(#*&$?(*@&$)",
			k:         "a-key",
			wantedErr: []error{errors.New(fmt.Sprintf("%s is not a valid key: %s", "a-key", "(#*&$?(*@&$)"))},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, errs := validateFeatureFlagKey(testCase.v, testCase.k)
			testValidateVerifyGeneric(t, errs, testCase.wantedErr)
		})
	}
}

func TestValidateFeatureFlagVariationsType(t *testing.T) {
	testCases := []struct {
		name      string
		v         interface{}
		k         string
		wantedErr []error
	}{
		{
			name:      "expected",
			v:         "string",
			k:         "a-key",
			wantedErr: nil,
		},
		{
			name:      "with invalid type as value in string",
			v:         "long",
			k:         "a-key",
			wantedErr: []error{errors.New(fmt.Sprintf("expected %s to be one of %v, got %s", "a-key", []string{"number", "boolean", "string"}, "long"))},
		},
		{
			name:      "with invalid type as value",
			v:         1,
			k:         "a-key",
			wantedErr: []error{errors.New(fmt.Sprintf("expected %s to be a string", "a-key"))},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, errs := validateFeatureFlagVariationsType(testCase.v, testCase.k)
			testValidateVerifyGeneric(t, errs, testCase.wantedErr)
		})
	}
}

func TestValidateVariationValue(t *testing.T) {
	testCases := []struct {
		name      string
		v         interface{}
		k         string
		wantedErr []error
	}{
		{
			name:      "expected",
			v:         "string",
			k:         "a-key",
			wantedErr: nil,
		},
		{
			name:      "with empty string",
			v:         "",
			k:         "a-key",
			wantedErr: []error{errors.New(fmt.Sprintf("%s cannot be an empty string", "a-key"))},
		},
		{
			name:
			"with invalid type as value",
			v:         1,
			k:         "a-key",
			wantedErr: []error{errors.New(fmt.Sprintf("expected %s to be a string", "a-key"))},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, errs := validateVariationValue(testCase.v, testCase.k)
			testValidateVerifyGeneric(t, errs, testCase.wantedErr)
		})
	}
}

func TestValidateColor(t *testing.T) {
	testCases := []struct {
		name      string
		v         interface{}
		k         string
		wantedErr []error
	}{
		{
			name:      "expected",
			v:         "FF00FF",
			k:         "a-key",
			wantedErr: nil,
		},
		{
			name:      "with HEX sign before color code",
			v:         "#FF00FF",
			k:         "a-key",
			wantedErr: []error{errors.New(fmt.Sprintf("%s is not a valid HEX color code: %s", "a-key", "#FF00FF"))},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, errs := validateColor(testCase.v, testCase.k)
			testValidateVerifyGeneric(t, errs, testCase.wantedErr)
		})
	}
}

func testValidateVerifyGeneric(t *testing.T, errs []error, wantedErr []error) {
	if !reflect.DeepEqual(errs, wantedErr) {
		t.Errorf("got error (%s) but want (%s)", errs, wantedErr)
	}
}
