package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"test_accelbyte_muhammad_huzair/cmd"
)

func main() {
	// flags
	var (
		directory  string
		fileType   string
		startTime  string
		endTime    string
		outputType string
		outputName string
	)

	flag.StringVar(&directory, "directory", "", "(REQUIRED PARAMETER) directory path")
	flag.StringVar(&directory, "d", "", "(REQUIRED PARAMETER) directory path (shorthand)")
	flag.StringVar(&fileType, "type", "", "(REQUIRED PARAMETER) type of input files (json or csv)")
	flag.StringVar(&fileType, "t", "", "(REQUIRED PARAMETER) type of input files (shorthand)")
	flag.StringVar(&startTime, "startTime", "", "(REQUIRED PARAMETER) starting time in RFC3339 format (example format: 2021-12-12T00:00:00.00Z)")
	flag.StringVar(&endTime, "endTime", "", "(REQUIRED PARAMETER) ending time in RFC3339 format (example format: 2022-11-30T00:23:59.21Z)")
	flag.StringVar(&outputType, "outputFileType", "json", "(optional) output file type (json or yaml)")
	flag.StringVar(&outputName, "outputFileName", "", "(optional) output file type (shorthand)")

	//parse the flag
	flag.Parse()

	// validate required flags
	if directory == "" || fileType == "" || startTime == "" || endTime == "" {
		flag.Usage()
		os.Exit(1)
	}

	// parse start time
	startTimeParsed, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		fmt.Println("failed to parse start time:", err)
		os.Exit(1)
	}
	// parse end time
	endTimeParsed, err := time.Parse(time.RFC3339, endTime)
	if err != nil {
		fmt.Println("failed to parse end time:", err)
		os.Exit(1)
	}

	// read files from the directory
	files, err := os.ReadDir(directory)
	if err != nil {
		fmt.Println("failed to read directory:", err)
		os.Exit(1)
	}

	// Process the metrics data
	newProcessMetric := cmd.NewProcessMetric(files, &directory, &fileType, &outputType, &outputName, startTimeParsed, endTimeParsed)
	errProcess := newProcessMetric.ProcessMetricData()
	if errProcess != nil {
		fmt.Printf("failed process, got error : %v \n", errProcess)
		os.Exit(1)
	}

	fmt.Println("summary result job was done!")
}
