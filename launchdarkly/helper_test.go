package launchdarkly

import (
	"errors"
	"github.com/hashicorp/terraform/helper/schema"
	"reflect"
	"testing"
)

func TestParseCompositeID(t *testing.T) {
	expectedErr := errors.New("error: Import composite ID requires two parts separated by colon, eg x:y")

	testCases := []struct {
		name      string
		id        string
		wantedP1  string
		wantedP2  string
		wantedErr error
	}{
		{
			name:      "expected",
			id:        "id:test",
			wantedP1:  "id",
			wantedP2:  "test",
			wantedErr: nil,
		},
		{
			name:      "with more than one separator",
			id:        "id:test:id",
			wantedP1:  "id",
			wantedP2:  "test:id",
			wantedErr: nil,
		},
		{
			name:      "without separator",
			id:        "test",
			wantedP1:  "",
			wantedP2:  "",
			wantedErr: expectedErr,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			p1, p2, err := parseCompositeID(testCase.id)
			testParseCompositeIDVerify(t, p1, p2, err, testCase)
		})
	}
}

func testReadMethod(d *schema.ResourceData, m interface{}) error { return nil }

func TestResourceImport(t *testing.T) {
	resourceKey := "resource"
	projectKey := "project"
	dTest := new(schema.ResourceData)
	dTest.SetId(projectKey + ":" + resourceKey)

	dTestWithError := new(schema.ResourceData)
	dTestWithError.SetId(resourceKey)

	wantedResourceData := new(schema.ResourceData)
	wantedResourceData.SetId(resourceKey)
	wantedResourceData.Set("project_key", projectKey)
	wantedResourceData.Set("key", resourceKey)

	expectedErr := errors.New("error: Import composite ID requires two parts separated by colon, eg x:y")

	testCases := []struct {
		name               string
		readMethod         importFunc
		d                  *schema.ResourceData
		meta               interface{}
		wantedResourceData []*schema.ResourceData
		wantedErr          error
	}{
		{
			name:               "expected",
			readMethod:         testReadMethod,
			d:                  dTest,
			meta:               nil,
			wantedResourceData: []*schema.ResourceData{wantedResourceData},
			wantedErr:          nil,
		},
		{
			name:               "with bad resource/project key formatting",
			readMethod:         testReadMethod,
			d:                  dTestWithError,
			meta:               nil,
			wantedResourceData: nil,
			wantedErr:          expectedErr,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			resourceData, err := resourceImport(testCase.readMethod, testCase.d, testCase.meta)
			testResourceImportVerify(t, resourceData, err, testCase.wantedResourceData, testCase.wantedErr)
		})
	}
}

func testParseCompositeIDVerify(t *testing.T, p1 string, p2 string, err error, testCase struct {
	name      string
	id        string
	wantedP1  string
	wantedP2  string
	wantedErr error
}) {
	if p1 != testCase.wantedP1 {
		t.Errorf("got string (%s) but want (%s)", p1, testCase.wantedP1)
	}

	if p2 != testCase.wantedP2 {
		t.Errorf("got string (%s) but want (%s)", p2, testCase.wantedP2)
	}

	if testCase.wantedErr != nil {
		if err.Error() != testCase.wantedErr.Error() {
			t.Errorf("got error (%s) but want (%s)", err, testCase.wantedErr)
		}
	}
}

func testResourceImportVerify(t *testing.T, resourceData []*schema.ResourceData, err error, wantedResourceData []*schema.ResourceData, wantedErr error) {

	if !reflect.DeepEqual(resourceData, wantedResourceData) {
		t.Errorf("resourceData is not equal to wantedResourceData")
	}

	if wantedErr != nil {
		if err.Error() != wantedErr.Error() {
			t.Errorf("got error (%s) but want (%s)", err, wantedErr)
		}
	}
}
