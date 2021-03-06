// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package model

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewId(t *testing.T) {
	for i := 0; i < 1000; i++ {
		id := NewId()
		if len(id) > 26 {
			t.Fatal("ids shouldn't be longer than 26 chars")
		}
	}
}

func TestRandomString(t *testing.T) {
	for i := 0; i < 1000; i++ {
		r := NewRandomString(32)
		if len(r) != 32 {
			t.Fatal("should be 32 chars")
		}
	}
}

func TestAppError(t *testing.T) {
	err := NewAppError("TestAppError", "message", nil, "", http.StatusInternalServerError)
	json := err.ToJson()
	rerr := AppErrorFromJson(strings.NewReader(json))
	if err.Message != rerr.Message {
		t.Fatal()
	}

	t.Log(err.Error())
}

func TestAppErrorJunk(t *testing.T) {
	rerr := AppErrorFromJson(strings.NewReader("<html><body>This is a broken test</body></html>"))
	if "body: <html><body>This is a broken test</body></html>" != rerr.DetailedError {
		t.Fatal()
	}
}

func TestMapJson(t *testing.T) {

	m := make(map[string]string)
	m["id"] = "test_id"
	json := MapToJson(m)

	rm := MapFromJson(strings.NewReader(json))

	if rm["id"] != "test_id" {
		t.Fatal("map should be valid")
	}

	rm2 := MapFromJson(strings.NewReader(""))
	if len(rm2) > 0 {
		t.Fatal("make should be ivalid")
	}
}

func TestValidEmail(t *testing.T) {
	if !IsValidEmail("corey+test@hulen.com") {
		t.Error("email should be valid")
	}

	if IsValidEmail("@corey+test@hulen.com") {
		t.Error("should be invalid")
	}
}

func TestValidLower(t *testing.T) {
	if !IsLower("corey+test@hulen.com") {
		t.Error("should be valid")
	}

	if IsLower("Corey+test@hulen.com") {
		t.Error("should be invalid")
	}
}

func TestEtag(t *testing.T) {
	etag := Etag("hello", 24)
	if len(etag) <= 0 {
		t.Fatal()
	}
}

var hashtags = map[string]string{
	"#test":           "#test",
	"test":            "",
	"#test123":        "#test123",
	"#123test123":     "",
	"#test-test":      "#test-test",
	"#test?":          "#test",
	"hi #there":       "#there",
	"#bug #idea":      "#bug #idea",
	"#bug or #gif!":   "#bug #gif",
	"#hüllo":          "#hüllo",
	"#?test":          "",
	"#-test":          "",
	"#yo_yo":          "#yo_yo",
	"(#brakets)":      "#brakets",
	")#stekarb(":      "#stekarb",
	"<#less_than<":    "#less_than",
	">#greater_than>": "#greater_than",
	"-#minus-":        "#minus",
	"_#under_":        "#under",
	"+#plus+":         "#plus",
	"=#equals=":       "#equals",
	"%#pct%":          "#pct",
	"&#and&":          "#and",
	"^#hat^":          "#hat",
	"##brown#":        "#brown",
	"*#star*":         "#star",
	"|#pipe|":         "#pipe",
	":#colon:":        "#colon",
	";#semi;":         "#semi",
	"#Mötley;":        "#Mötley",
	".#period.":       "#period",
	"¿#upside¿":       "#upside",
	"\"#quote\"":      "#quote",
	"/#slash/":        "#slash",
	"\\#backslash\\":  "#backslash",
	"#a":              "",
	"#1":              "",
	"foo#bar":         "",
}

func TestParseHashtags(t *testing.T) {
	for input, output := range hashtags {
		if o, _ := ParseHashtags(input); o != output {
			t.Fatal("failed to parse hashtags from input=" + input + " expected=" + output + " actual=" + o)
		}
	}
}

func TestIsValidAlphaNum(t *testing.T) {
	cases := []struct {
		Input  string
		Result bool
	}{
		{
			Input:  "test",
			Result: true,
		},
		{
			Input:  "test-name",
			Result: true,
		},
		{
			Input:  "test--name",
			Result: true,
		},
		{
			Input:  "test__name",
			Result: true,
		},
		{
			Input:  "-",
			Result: false,
		},
		{
			Input:  "__",
			Result: false,
		},
		{
			Input:  "test-",
			Result: false,
		},
		{
			Input:  "test--",
			Result: false,
		},
		{
			Input:  "test__",
			Result: false,
		},
		{
			Input:  "test:name",
			Result: false,
		},
	}

	for _, tc := range cases {
		actual := IsValidAlphaNum(tc.Input)
		if actual != tc.Result {
			t.Fatalf("case: %v\tshould returned: %#v", tc, tc.Result)
		}
	}
}

func TestGetServerIpAddress(t *testing.T) {
	if len(GetServerIpAddress()) == 0 {
		t.Fatal("Should find local ip address")
	}
}

func TestIsValidAlphaNumHyphenUnderscore(t *testing.T) {
	casesWithFormat := []struct {
		Input  string
		Result bool
	}{
		{
			Input:  "test",
			Result: true,
		},
		{
			Input:  "test-name",
			Result: true,
		},
		{
			Input:  "test--name",
			Result: true,
		},
		{
			Input:  "test__name",
			Result: true,
		},
		{
			Input:  "test_name",
			Result: true,
		},
		{
			Input:  "test_-name",
			Result: true,
		},
		{
			Input:  "-",
			Result: false,
		},
		{
			Input:  "__",
			Result: false,
		},
		{
			Input:  "test-",
			Result: false,
		},
		{
			Input:  "test--",
			Result: false,
		},
		{
			Input:  "test__",
			Result: false,
		},
		{
			Input:  "test:name",
			Result: false,
		},
	}

	for _, tc := range casesWithFormat {
		actual := IsValidAlphaNumHyphenUnderscore(tc.Input, true)
		if actual != tc.Result {
			t.Fatalf("case: %v\tshould returned: %#v", tc, tc.Result)
		}
	}

	casesWithoutFormat := []struct {
		Input  string
		Result bool
	}{
		{
			Input:  "test",
			Result: true,
		},
		{
			Input:  "test-name",
			Result: true,
		},
		{
			Input:  "test--name",
			Result: true,
		},
		{
			Input:  "test__name",
			Result: true,
		},
		{
			Input:  "test_name",
			Result: true,
		},
		{
			Input:  "test_-name",
			Result: true,
		},
		{
			Input:  "-",
			Result: true,
		},
		{
			Input:  "_",
			Result: true,
		},
		{
			Input:  "test-",
			Result: true,
		},
		{
			Input:  "test--",
			Result: true,
		},
		{
			Input:  "test__",
			Result: true,
		},
		{
			Input:  ".",
			Result: false,
		},

		{
			Input:  "test,",
			Result: false,
		},
		{
			Input:  "test:name",
			Result: false,
		},
	}

	for _, tc := range casesWithoutFormat {
		actual := IsValidAlphaNumHyphenUnderscore(tc.Input, false)
		if actual != tc.Result {
			t.Fatalf("case: '%v'\tshould returned: %#v", tc.Input, tc.Result)
		}
	}
}

func TestIsValidId(t *testing.T) {
	cases := []struct {
		Input  string
		Result bool
	}{
		{
			Input:  NewId(),
			Result: true,
		},
		{
			Input:  "",
			Result: false,
		},
		{
			Input:  "junk",
			Result: false,
		},
		{
			Input:  "qwertyuiop1234567890asdfg{",
			Result: false,
		},
		{
			Input:  NewId() + "}",
			Result: false,
		},
	}

	for _, tc := range cases {
		actual := IsValidId(tc.Input)
		if actual != tc.Result {
			t.Fatalf("case: %v\tshould returned: %#v", tc, tc.Result)
		}
	}
}

func TestNowhereNil(t *testing.T) {
	t.Parallel()

	var nilStringPtr *string
	var nonNilStringPtr *string = new(string)
	var nilSlice []string
	var nilStruct *struct{}
	var nilMap map[bool]bool

	var nowhereNilStruct = struct {
		X *string
		Y *string
	}{
		nonNilStringPtr,
		nonNilStringPtr,
	}
	var somewhereNilStruct = struct {
		X *string
		Y *string
	}{
		nonNilStringPtr,
		nilStringPtr,
	}

	var privateSomewhereNilStruct = struct {
		X *string
		y *string
	}{
		nonNilStringPtr,
		nilStringPtr,
	}

	testCases := []struct {
		Description string
		Value       interface{}
		Expected    bool
	}{
		{
			"nil",
			nil,
			false,
		},
		{
			"empty string",
			"",
			true,
		},
		{
			"non-empty string",
			"not empty!",
			true,
		},
		{
			"nil string pointer",
			nilStringPtr,
			false,
		},
		{
			"non-nil string pointer",
			nonNilStringPtr,
			true,
		},
		{
			"0",
			0,
			true,
		},
		{
			"1",
			1,
			true,
		},
		{
			"0 (int64)",
			int64(0),
			true,
		},
		{
			"1 (int64)",
			int64(1),
			true,
		},
		{
			"true",
			true,
			true,
		},
		{
			"false",
			false,
			true,
		},
		{
			"nil slice",
			nilSlice,
			// A nil slice is observably the same as an empty slice, so allow it.
			true,
		},
		{
			"empty slice",
			[]string{},
			true,
		},
		{
			"slice containing nils",
			[]*string{nil, nil},
			true,
		},
		{
			"nil map",
			nilMap,
			false,
		},
		{
			"non-nil map",
			make(map[bool]bool),
			true,
		},
		{
			"non-nil map containing nil",
			map[bool]*string{true: nilStringPtr, false: nonNilStringPtr},
			// Map values are not checked
			true,
		},
		{
			"nil struct",
			nilStruct,
			false,
		},
		{
			"empty struct",
			struct{}{},
			true,
		},
		{
			"struct containing no nil",
			nowhereNilStruct,
			true,
		},
		{
			"struct containing nil",
			somewhereNilStruct,
			false,
		},
		{
			"struct pointer containing no nil",
			&nowhereNilStruct,
			true,
		},
		{
			"struct pointer containing nil",
			&somewhereNilStruct,
			false,
		},
		{
			"struct containing private nil",
			privateSomewhereNilStruct,
			true,
		},
		{
			"struct pointer containing private nil",
			&privateSomewhereNilStruct,
			true,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Description, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("panic: %v", r)
				}
			}()

			t.Parallel()
			require.Equal(t, testCase.Expected, checkNowhereNil(t, "value", testCase.Value))
		})
	}
}
