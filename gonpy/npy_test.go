package gonpy

import (
    "bytes"
    "os"
    "path/filepath"
    "strings"
    "testing"
)

func TestMissingMagicString(t *testing.T) {
    // Given a string with a bad magic string
    var testData = []byte("\x93NUMP\x01\x02")
    reader := bytes.NewReader(testData)

    _, err := ParseData[float64](reader)

    if err == nil {
        t.Fatal("Expected an error to be returned.")
    }

    correctError := strings.Contains(err.Error(), "Incorrect file format")
    if !correctError {
        t.Logf("Different error returned than expected. Error: %v", err.Error())
        t.Failed()
    }
}

func TestWrongVersion(t *testing.T) {
    // Given a string with a bad version
    var testData = []byte("\x93NUMPY\x02\x02")
    reader := bytes.NewReader(testData)

    _, err := ParseData[float64](reader)

    if err == nil {
        t.Fatal("Expected an error to be returned.")
    }

    correctError := strings.Contains(err.Error(), "Unsupported version")
    if !correctError {
        t.Logf("Different error returned than expected. Error: %v", err.Error())
        t.Failed()
    }
}
