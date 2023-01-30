package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRun_FromTS(t *testing.T) {
	var b bytes.Buffer
	err := run("1674898343", "", &b)
	require.Nil(t, err, "unexpected error: %s", err)

	assert.Equal(t, "Sat, 28 Jan 2023 09:32:23 UTC\n", b.String())
}

func TestRun_InvalidTS(t *testing.T) {
	var b bytes.Buffer
	err := run("invalid", "", &b)
	require.NotNil(t, err, "expected error")
	assert.Equal(t, "invalid timestamp: invalid", err.Error())
}

func TestRun_FromDate(t *testing.T) {
	tt := []struct {
		name     string
		date     string
		expected string
	}{
		// Apparently, GMT+XX timezone parsing is broken: https://github.com/golang/go/issues/56392
		//{name: "date with time zone", date: "28.01.2023, 16:32:23 GMT+05", expected: "1674898343\n"},
		{name: "date with numeric time zone", date: "28.01.2023, 16:32:23 +0700", expected: "1674898343\n"},
		{name: "date without time zone", date: "28.01.2023, 09:32:23", expected: "1674898343\n"},
		{name: "only date", date: "28.01.2023", expected: "1674864000\n"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var b bytes.Buffer
			err := run("", tc.date, &b)
			require.Nil(t, err, "unexpected error: %s", err)

			assert.Equal(t, tc.expected, b.String())
		})
	}
}

func TestRun_FromDate_InvalidTZFallback(t *testing.T) {
	tt := []struct {
		name     string
		date     string
		expected string
	}{
		{name: "date with named timezone", date: "28.01.2023, 09:32:23 PST", expected: "1674898343\n"},  // fallback to UTC
		{name: "date with GMT timezone", date: "28.01.2023, 09:32:23 GMT+03", expected: "1674898343\n"}, // fallback to UTC
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var b bytes.Buffer
			err := run("", tc.date, &b)
			require.Nil(t, err, "unexpected error: %s", err)

			assert.Equal(t, tc.expected, b.String())
		})
	}
}

func TestRun_FromDate_InvalidDate(t *testing.T) {
	tt := []struct {
		name string
		date string
	}{
		{name: "invalid date", date: "2023.01.18"},
		{name: "invalid time", date: "28.01.2023, blah"},
		{name: "no time", date: "28.01.2023,"},
		{name: "invalid time format", date: "28.01.2023, 03:15PM"},
		{name: "invalid timezone", date: "28.01.2023, 12:45Z03:00"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var b bytes.Buffer
			err := run("", tc.date, &b)
			assert.ErrorContains(t, err, "invalid date")
		})
	}
}
