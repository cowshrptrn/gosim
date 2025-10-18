package gonpy

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Numeric interface {
	int8 | int16 | int32 | int64 |
		uint8 | uint16 | uint32 | uint64 |
		float32 | float64
}
type NpyData[T Numeric] struct {
	fortranOrder bool
	shape        []uint64
	data         []T
}

func (nd *NpyData[T]) Data() []T {
	return nd.data
}

func (nd *NpyData[T]) Shape() []uint64 {
	return nd.shape
}

// For now, only supporting rectangular ND arrays
func findEndofSubstr(data *[]byte, target *[]byte) (uint64, error) {
	var data_idx uint64
	for data_idx = 0; data_idx < uint64(len(*data)); data_idx++ {
		var offset uint64
		for offset = 0; offset < uint64(len(*target)); offset++ {
			if data_idx+offset >= uint64(len(*data)) {
				return 0, errors.New("Failed to find substring.")
			}
			if (*data)[data_idx+offset] != (*target)[offset] {
				break
			}
		}
		if offset == uint64(len(*target)) {
			return data_idx + offset, nil
		}
	}

	return 0, errors.New("Failed to find substring.")
}

func findBoundedOffsets(data *[]byte, startOffset uint64, startChar byte, endChar byte) (uint64, uint64, error) {
	var start, end, idx uint64 // zero-initialized
	for idx = startOffset; idx < uint64(len(*data)); idx++ {
		if (*data)[idx] == startChar && start == 0 {
			start = idx + 1
		} else if (*data)[idx] == endChar && start > 0 {
			end = idx
			break
		}
	}

	if start > 0 && end > 0 {
		return start, end, nil
	} else {
		return start, end, errors.New("Failed to find ending character.")
	}
}

func extractHeaderData(data *[]byte, dtypeDescr *string, fortranOrder *bool, shape *[]uint64) error {
	// Don't bother validating the entire structure, just look for the keywords
	dtypeKey := []byte("'descr':")
	offset, err := findEndofSubstr(data, &dtypeKey)
	if err != nil {
		return err
	}

	var start, end uint64
	start, end, err = findBoundedOffsets(data, offset, '\'', '\'')

	if err != nil {
		return errors.New("Malformed descr entry")
	}

	*dtypeDescr = string((*data)[start:end])

	// Process fortran order
	orderKey := []byte("'fortran_order':")
	offset, err = findEndofSubstr(data, &orderKey)
	if err != nil {
		return err
	}

	trueStr := "True"
	falseStr := "False"
	for idx := offset; idx < uint64(len(*data)-5); idx++ {
		char := (*data)[idx]
		if (*data)[idx] != byte(' ') {
			if char == falseStr[0] && falseStr == string((*data)[idx:idx+uint64(len(falseStr))]) {
				*fortranOrder = false
				break
			} else if char == trueStr[0] && trueStr == string((*data)[idx:idx+uint64(len(trueStr))]) {
				*fortranOrder = true
				break
			} else {
				return fmt.Errorf("Malformed fortran_order entry in header. First 5 chars: %v ", string((*data)[offset:offset+5]))
			}
		}
	}

	// Process shape
	shapeKey := []byte("shape")
	offset, err = findEndofSubstr(data, &shapeKey)
	if err != nil {
		return err
	}

	start, end, err = findBoundedOffsets(data, offset, '(', ')')
	if err != nil {
		return err
	}

	sizeStr := string((*data)[start:end])
	sizes := strings.Split(sizeStr, ",")
	for _, size := range sizes {
		if len(size) == 0 {
			continue
		}
		val, err := strconv.Atoi(strings.TrimSpace(size))
		if err != nil {
			return err
		}
		if val <= 0 {
			return errors.New("Invalid dimension size.")
		}
		*shape = append(*shape, uint64(val))
	}

	return nil
}

type NumericType int

const (
	Float NumericType = iota
	SignedInteger
	UnsignedInteger
)

func isBigEndianToByteOrder(isBigEndian bool) binary.ByteOrder {
	if isBigEndian {
		return binary.BigEndian
	} else {
		return binary.LittleEndian
	}
}

func parseSignedInt8(value *int8, isBigEndian bool, reader io.Reader) error {
	return binary.Read(reader, isBigEndianToByteOrder(isBigEndian), value)
}

func parseSignedInt16(value *int16, isBigEndian bool, reader io.Reader) error {
	return binary.Read(reader, isBigEndianToByteOrder(isBigEndian), value)
}

func parseSignedInt32(value *int32, isBigEndian bool, reader io.Reader) error {
	return binary.Read(reader, isBigEndianToByteOrder(isBigEndian), value)
}

func parseSignedInt64(value *int64, isBigEndian bool, reader io.Reader) error {
	return binary.Read(reader, isBigEndianToByteOrder(isBigEndian), value)
}

func parseSignedIntFunc(size int, isBigEndian bool) (any, error) {
	if size == 1 {
		return func(reader io.Reader, value *int8) error {
			return parseSignedInt8(value, isBigEndian, reader)
		}, nil
	} else if size == 2 {
		return func(reader io.Reader, value *int16) error {
			return parseSignedInt16(value, isBigEndian, reader)
		}, nil
	} else if size == 4 {
		return func(reader io.Reader, value *int32) error {
			return parseSignedInt32(value, isBigEndian, reader)
		}, nil
	} else if size == 8 {
		return func(reader io.Reader, value *int64) error {
			return parseSignedInt64(value, isBigEndian, reader)
		}, nil
	}
	err := errors.New("Unsupported size.")
	return nil, err
}

func parseUnignedInt8(value *uint8, isBigEndian bool, reader io.Reader) error {
	return binary.Read(reader, isBigEndianToByteOrder(isBigEndian), value)
}

func parseUnignedInt16(value *uint16, isBigEndian bool, reader io.Reader) error {
	return binary.Read(reader, isBigEndianToByteOrder(isBigEndian), value)
}

func parseUnsignedInt32(value *uint32, isBigEndian bool, reader io.Reader) error {
	return binary.Read(reader, isBigEndianToByteOrder(isBigEndian), value)
}

func parseUnsignedInt64(value *uint64, isBigEndian bool, reader io.Reader) error {
	return binary.Read(reader, isBigEndianToByteOrder(isBigEndian), value)
}

func parseUnsignedIntFunc(size int, isBigEndian bool) (any, error) {
	if size == 1 {
		return func(reader io.Reader, value *uint8) error {
			return parseUnignedInt8(value, isBigEndian, reader)
		}, nil
	} else if size == 2 {
		return func(reader io.Reader, value *uint16) error {
			return parseUnignedInt16(value, isBigEndian, reader)
		}, nil
	} else if size == 4 {
		return func(reader io.Reader, value *uint32) error {
			return parseUnsignedInt32(value, isBigEndian, reader)
		}, nil
	} else if size == 8 {
		return func(reader io.Reader, value *uint64) error {
			return parseUnsignedInt64(value, isBigEndian, reader)
		}, nil
	}
	err := errors.New("Unsupported size.")
	return nil, err
}

func parseFloat32(value *float32, isBigEndian bool, reader io.Reader) error {
	return binary.Read(reader, isBigEndianToByteOrder(isBigEndian), value)
}

func parseFloat64(value *float64, isBigEndian bool, reader io.Reader) error {
	return binary.Read(reader, isBigEndianToByteOrder(isBigEndian), value)
}

func parseFloatFunc(size int, isBigEndian bool) (any, error) {
	if size == 4 {
		return func(reader io.Reader, value *float32) error {
			return parseFloat32(value, isBigEndian, reader)
		}, nil
	} else if size == 8 {
		return func(reader io.Reader, value *float64) error {
			return parseFloat64(value, isBigEndian, reader)
		}, nil
	}
	err := errors.New("Unsupported size.")
	return nil, err
}

func parseDtype(dtype string) (any, error) {
	// This needs to be extended, but for now we will only support the
	// real-values signed / unsigned integer and floating point types
	// Sources:
	// https://numpy.org/doc/stable/reference/arrays.dtypes.html#arrays-dtypes-constructing
	// https://numpy.org/doc/stable/reference/arrays.interface.html#arrays-interface

	currOffset := 0
	// Parse endianness
	// Options are >, <, =, or none.
	// Default and "=" are machine native (little endian) so only need to check for big endian.
	var isBigEndian = false
	if dtype[0] == '>' {
		isBigEndian = true
	}

	if strings.Contains("><=|", string(dtype[currOffset])) {
		currOffset += 1
	}

	// Unpack numeric type
	var numericType NumericType
	if dtype[currOffset] == 'f' {
		numericType = Float
	} else if dtype[currOffset] == 'i' {
		numericType = SignedInteger
	} else if dtype[currOffset] == 'u' {
		numericType = UnsignedInteger
	} else {
		return nil, errors.New("Unrecognized data type.")
	}

	currOffset += 1

	// float, int, and uint all default to 32bit
	// This number is the number of BYTES, not BITS
	// Support 8, 16, 32, and 64 bit numbers
	byteness := 4
	sizeStr := string(dtype[len(dtype)-1])
	if strings.Contains("1248", sizeStr) {
		var err error
		byteness, err = strconv.Atoi(sizeStr)
		if err != nil {
			return nil, err
		}
	} else {
		err := fmt.Errorf("Unrecognized size: %v", sizeStr)
		return nil, err
	}

	switch numericType {
	case Float:
		return parseFloatFunc(byteness, isBigEndian)
	case SignedInteger:
		return parseSignedIntFunc(byteness, isBigEndian)
	case UnsignedInteger:
		return parseUnsignedIntFunc(byteness, isBigEndian)
	default:
		err := errors.New("Unrecognized type signifier")
		return nil, err
	}

}

func ParseData[T Numeric](data io.Reader) (NpyData[T], error) {
	var RetVal NpyData[T]
	// Check initial magic string
	magicString := []byte("\x93NUMPY")
	prefixBuffer := make([]byte, len(magicString))

	numBytes, err := data.Read(prefixBuffer)
	if err != nil {
		return RetVal, err
	}

	if numBytes < len(magicString) || !bytes.Equal(magicString, prefixBuffer) {
		return RetVal, errors.New("Incorrect file format. Expected correct magic string.")
	}

	version := make([]byte, 2)
	numBytes, err = data.Read(version)
	if err != nil || numBytes < 2 {
		return RetVal, err
	}

	if version[0] != 1 {
		return RetVal, errors.New("Unsupported version.")
	}

	headerSize := make([]byte, 2)
	numBytes, err = data.Read(headerSize)
	if err != nil || numBytes < 2 {
		return RetVal, err
	}

	var headerSizeVal uint16
	headerSizeReader := bytes.NewReader(headerSize)
	err = binary.Read(headerSizeReader, binary.LittleEndian, &headerSizeVal)
	if err != nil {
		return RetVal, err
	}

	var header = make([]byte, headerSizeVal)
	numBytes, err = data.Read(header)
	if err != nil {
		return RetVal, err
	}

	if numBytes < int(headerSizeVal) {
		return RetVal, errors.New("Failed to read all header values.")
	}
	var dtypeDescr string

	err = extractHeaderData(&header, &dtypeDescr, &(RetVal.fortranOrder), &(RetVal.shape))

	if err != nil {
		return RetVal, err
	}

	// Parse dtype into a method that will take reader,
	// extract next N bytes, and convert them to the appropriate data type
	var parseFunAny any
	parseFunAny, err = parseDtype(dtypeDescr)
	if err != nil {
		return RetVal, err
	}

	var parseFun func(io.Reader, *T) error
	var ok bool
	if parseFun, ok = parseFunAny.(func(io.Reader, *T) error); !ok {
		castErr := errors.New("Recieved incorrect type of parser function")
		return RetVal, castErr
	}

	var totalSize uint64 = 1
	for _, val := range RetVal.shape {
		totalSize = totalSize * val
	}

	RetVal.data = make([]T, totalSize)
	for dataIdx := 0; dataIdx < len(RetVal.data); dataIdx++ {
		parseErr := parseFun(data, &(RetVal.data[dataIdx]))
		if parseErr != nil {
			return RetVal, parseErr
		}
	}

	return RetVal, nil
}
