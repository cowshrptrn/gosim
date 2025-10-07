package gonpy

import (
    "os"
    "path/filepath"
    "testing"
)

func runBaseCase[T numeric] (t *testing.T, filePath string, totalElms uint64, mod int) {
    testFile, err := os.Open(filePath)
    if err != nil {
        t.Fatalf("Failed to open file with error %v", err.Error())
    }
    defer testFile.Close()

    testResult, err := ParseData[T](testFile)

    if err != nil {
        t.Fatalf("Error in parsing: %v", err.Error())
    }

    if testResult.fortranOrder {
        t.Fatal("Incorrect fortran order.")
    }

    if len(testResult.shape) != 1 || testResult.shape[0]  != totalElms {
        t.Fatal("Incorrect array shape.")
    }

    for idx, val := range testResult.data {
        target := T(idx + 1)
        if mod > 0 {
            target = T((idx + 1) % mod)
        }
        if val != target {
            t.Errorf("Incorrect value. Expected %v, got %v", idx+1, val)
        }
    }
}

func TestBaseCase_int8(t *testing.T) {
    testFilePath := filepath.Join("testdata", "500_int8.npy")
    runBaseCase[int8](t, testFilePath, 500, 128)
}

func TestBaseCase_int16(t *testing.T) {
    testFilePath := filepath.Join("testdata", "500_int16.npy")
    runBaseCase[int16](t, testFilePath, 500, 0)
}

func TestBaseCase_int32(t *testing.T) {
    testFilePath := filepath.Join("testdata", "500_int32.npy")
    runBaseCase[int32](t, testFilePath, 500, 0)
}

func TestBaseCase_int64(t *testing.T) {
    testFilePath := filepath.Join("testdata", "500_int64.npy")
    runBaseCase[int64](t, testFilePath, 500, 0)
}


func TestBaseCase_uint8(t *testing.T) {
    testFilePath := filepath.Join("testdata", "500_uint8.npy")
    runBaseCase[uint8](t, testFilePath, 500, 256)
}

func TestBaseCase_uint16(t *testing.T) {
    testFilePath := filepath.Join("testdata", "500_uint16.npy")
    runBaseCase[uint16](t, testFilePath, 500, 0)
}

func TestBaseCase_uint32(t *testing.T) {
    testFilePath := filepath.Join("testdata", "500_uint32.npy")
    runBaseCase[uint32](t, testFilePath, 500, 0)
}

func TestBaseCase_uint64(t *testing.T) {
    testFilePath := filepath.Join("testdata", "500_uint64.npy")
    runBaseCase[uint64](t, testFilePath, 500, 0)
}

func TestBaseCase_f4(t *testing.T) {
    testFilePath := filepath.Join("testdata", "500_f4.npy")
    runBaseCase[float32](t, testFilePath, 500, 0)
}

func TestBaseCase_f8(t *testing.T) {
    testFilePath := filepath.Join("testdata", "500_f8.npy")
    runBaseCase[float64](t, testFilePath, 500, 0)
}
