# GOSIM Go Numeric Utilities

## About
This is a testing repo used to write numeric and numeric adjacent utilities in go.

## Getting Started
To build the local library and make a small executable that can be used to test the functionality run
```
go build
```

To test a specific portion of the repo, run "go test" For example, to test the gonpy functionality run

```
go test -v nrokkam/gosim/gonpy
```

To print metadata for a npy file containt float64 typed data use the build's target executable.
```
./gosim -type f64 -file gonpy/testdata/5_f8.npy
```

## Submodules
### gonpy
This is a package for deserializing a standard numpy file.

The numpy spec has a large number of variations for numeric and non-numeric data. Currently the library only supports rectangular tensors for the following dtypes serialized using both big and little endianness assuming it is properly specified in the npy file.

|np.dtype | golang type| C type (64-bit build) | Description |
|---------|------------|--------|-------------|
|int8 |int8 | signed char | 8-bit signed integer.  |
|int16 |int16 | short int / short | 16-bit signed integer.|
|int32 |int32 | int | 32-bit signed integer.|
|int64 |int64 | long int / long |64-bit signed integer.|
|uint8 |uint8 | char / byte | 8-bit unsigned integer.|
|uint16 |uint16 | unsigned short int /<br>unsigned short | 16-bit unsigned integer.|
|uint32 |uint32 | unsigned int | 32-bit unsigned integer.|
|uint64 |uint64 | unsigned long int /<br>unsigned long | 64-bit unsigned integer.|
|float32 |float32 | float | 32-bit floating point.|
|float64 |float64 | double | 32-bit floating point.|
