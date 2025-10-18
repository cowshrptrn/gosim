package main

import (
	"flag"
	"fmt"
	"io"
	"nrokkam/gosim/gonpy"
	"os"
)

func printMetdata[T gonpy.Numeric](file io.Reader) {
	fmt.Println("Starting parse")
	data, err := gonpy.ParseData[T](file)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Parsed %v elements with shape %v\n", len(data.Data()), data.Shape())
	}
}

func main() {
	dtypePtr := flag.String("type", "", "Numerical datatype.")
	filePathPtr := flag.String("file", "", "File to read for metadata.")

	flag.Parse()

	testFile, err := os.Open(*filePathPtr)
	if err != nil {
		fmt.Printf("Failed to open file with error %v\n", err)
		os.Exit(1)
	}

	switch *dtypePtr {
	case "i8":
		printMetdata[int8](testFile)
	case "i16":
		printMetdata[int16](testFile)
	case "i32":
		printMetdata[int32](testFile)
	case "i64":
		printMetdata[int64](testFile)
	case "u8":
		printMetdata[uint8](testFile)
	case "u16":
		printMetdata[uint16](testFile)
	case "u32":
		printMetdata[uint32](testFile)
	case "u64":
		printMetdata[uint64](testFile)
	case "f32":
		printMetdata[float32](testFile)
	case "f64":
		printMetdata[float64](testFile)
	default:
		fmt.Printf("Unrecognized datatype: %v\n", *dtypePtr)
		os.Exit(1)
	}

}
