//go:build !bench
// +build !bench

package hw10programoptimization

import (
	"bytes"
	"strings"
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

func TestGetDomainStatAdditional(t *testing.T) {
	t.Run("empty input", func(t *testing.T) {
		result, err := GetDomainStat(strings.NewReader(""), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})

	t.Run("malformed JSON", func(t *testing.T) {
		data := `{"Id":1,"Name":"Test","Email":"test@example.com"`
		_, err := GetDomainStat(strings.NewReader(data), "com")
		require.Error(t, err)
	})

	t.Run("email without @", func(t *testing.T) {
		data := `{"Id":1,"Name":"Test","Username":"test","Email":"invalid-email","Phone":"123","Password":"pass","Address":"addr"}`
		result, err := GetDomainStat(strings.NewReader(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})

	t.Run("case sensitivity", func(t *testing.T) {
		data := `{"Id":1,"Name":"Test","Username":"test","Email":"Test@EXAMPLE.COM","Phone":"123","Password":"pass","Address":"addr"}`
		result, err := GetDomainStat(strings.NewReader(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"example.com": 1}, result)
	})

	t.Run("multiple empty lines", func(t *testing.T) {
		data := `{"Id":1,"Name":"Test","Username":"test","Email":"test@example.com","Phone":"123","Password":"pass","Address":"addr"}

{"Id":2,"Name":"Test2","Username":"test2","Email":"test2@example.org","Phone":"123","Password":"pass","Address":"addr"}`
		result, err := GetDomainStat(strings.NewReader(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"example.com": 1}, result)
	})
}
