package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestProcessMetricData(t *testing.T) {
	// Prepare test data
	directory := "../metrics/csv/"
	fileType := CsvType
	outputFileType := JsonType
	outputFileName := ""
	startTime := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2022, 11, 30, 0, 23, 51, 21000000, time.UTC)

	// Read the files in the test directory
	files, err := os.ReadDir(directory)
	if err != nil {
		t.Fatalf("Failed to read test directory: %v", err)
	}

	// Create a new ProcessMetric instance with test data
	pm := NewProcessMetric(files, &directory, &fileType, &outputFileType, &outputFileName, startTime, endTime)

	// Call the method being tested
	err = pm.ProcessMetricData()

	// Assert the expected behavior and results
	if err != nil {
		t.Errorf("ProcessMetricData returned an error: %v", err)
	}

	// Verify the output file exists and has the expected content
	outputFilePath := filepath.Join(DefaultOutputFolderName, pm.getOutputFileName(outputFileName, outputFileType))
	_, err = os.Stat(outputFilePath)
	if os.IsNotExist(err) {
		t.Errorf("Output file does not exist: %s", outputFilePath)
	}

	outputContent, err := ioutil.ReadFile(outputFilePath)
	if err != nil {
		t.Errorf("Failed to read output file: %v", err)
	}

	expectedOutput := `[{"level_name":"lobby_screen","total_value":505},{"level_name":"main_menu","total_value":80},{"level_name":"options","total_value":60},{"level_name":"help","total_value":50},{"level_name":"level1","total_value":520}]`

	if string(outputContent) != expectedOutput {
		t.Errorf("Output content does not match the expected value.\nExpected:\n%s\nActual:\n%s", expectedOutput, string(outputContent))
	}

	// Cleanup (delete the output file if it exists)
	err = os.Remove(outputFilePath)
	if err != nil {
		t.Errorf("Failed to remove output file: %v", err)
	}
}

func TestProcessMetricData_Negative_InvalidFileType(t *testing.T) {
	// Prepare test data
	directory := "../metrics/csv/"
	fileType := JsonType
	outputFileType := JsonType
	outputFileName := ""
	startTime := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2022, 11, 30, 0, 23, 51, 21000000, time.UTC)

	// Read the files in the test directory
	files, err := os.ReadDir(directory)
	if err != nil {
		t.Fatalf("Failed to read test directory: %v", err)
	}

	// Create a new ProcessMetric instance with test data
	pm := NewProcessMetric(files, &directory, &fileType, &outputFileType, &outputFileName, startTime, endTime)

	expectedError := errors.New("unexpected file type")

	// Call the method being tested
	err = pm.ProcessMetricData()

	// Assert the expected behavior and results
	fmt.Println("err ", err)
	if err != nil {
		if err.Error() != expectedError.Error() {
			t.Errorf("Output error does not match the expected error .\nExpected:\n%s\nActual:\n%s", expectedError.Error(), err.Error())
		}
	} else {
		t.Errorf("this should be error!")
	}
}
