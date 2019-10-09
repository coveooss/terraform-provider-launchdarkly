package launchdarkly

import (
	"errors"
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

func testParseCompositeIDVerify(t *testing.T, p1 string, p2 string, err error, testCase struct {
	name      string
	id        string
	wantedP1  string
	wantedP2  string
	wantedErr error
}) {
	if p1 != testCase.wantedP1 {
		t.Errorf("got string (%s) but want %s", p1, testCase.wantedP1)
	}

	if p2 != testCase.wantedP2 {
		t.Errorf("got string (%s) but want %s", p2, testCase.wantedP2)
	}

	if testCase.wantedErr != nil {
		if err.Error() != testCase.wantedErr.Error() {
			t.Errorf("got error (%s) but want (%s)", err, testCase.wantedErr)
		}
	}
}
