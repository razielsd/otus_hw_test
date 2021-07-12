// +build !bench

package hw10programoptimization

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDomainStat(t *testing.T) {
	data := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

	t.Run("find 'com'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("find 'gov'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
	})

	t.Run("find 'unknown'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "unknown")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})
}

func TestExtractEmail_ValidJson(t *testing.T) {
	js := `{"Id":2,"Name":"Brian Olson","Username":"non_quia_id","Email":"FrancesEllis@Quinu.edu","Phone":"237-75-34","Password":"cmEPhX8","Address":"Butterfield Junction 74"}`
	email, err := ExtractEmail(js)
	require.NoError(t, err)
	require.Equal(t, "\"FrancesEllis@Quinu.edu", email)
}

func TestExtractEmail_InvalidJson(t *testing.T) {
	js := `{"Id":2,"Name":"Brian Olson","Username":"non_quia_id","NoMail":"FrancesEllis@Quinu.edu","Phone":"237-75-34","Password":"cmEPhX8","Address":"Butterfield Junction 74"}`
	_, err := ExtractEmail(js)
	require.Error(t, err)
}

func TestGetHostExtractor_FoundEmail(t *testing.T) {
	tests := []struct {
		input    string
		domain   string
		expected string
	}{
		{input: "va.sya@test.gov", domain: "gov", expected: "test.gov"},
		{input: "PETRO.VICH@gmail.com", domain: "com", expected: "gmail.com"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			extractor := GetHostExtractor(tc.domain)
			host, found := extractor(tc.input)
			require.True(t, found)
			require.Equal(t, tc.expected, host)
		})
	}
}

func TestGetHostExtractor_NotFoundEmail(t *testing.T) {
	tests := []struct {
		input    string
		domain   string
		expected string
	}{
		{input: "va.sya@test.gov", domain: "com", expected: "test.gov"},
		{input: "gmail.com", domain: "com", expected: "gmail.com"}, // it is not email
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			extractor := GetHostExtractor(tc.domain)
			_, found := extractor(tc.input)
			require.False(t, found)
		})
	}
}
