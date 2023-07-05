package cmd

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type ProcessMetric struct {
	files          []os.DirEntry
	directory      *string
	fileType       *string
	outputFileType *string
	outputFileName *string
	startTime      time.Time
	endTime        time.Time
}

func NewProcessMetric(files []os.DirEntry, directory, fileType, outputFileType, outputFileName *string, startTime, endTime time.Time) ProcessMetric {
	return ProcessMetric{
		files:          files,
		directory:      directory,
		fileType:       fileType,
		outputFileType: outputFileType,
		outputFileName: outputFileName,
		startTime:      startTime,
		endTime:        endTime,
	}
}

type metric struct {
	Timestamp string `json:"timestamp"`
	LevelName string `json:"level_name"`
	Value     int    `json:"value"`
}

type summaryResult struct {
	LevelName  string `json:"level_name" yaml:"level_name"`
	TotalValue int    `json:"total_value" yaml:"total_value"`
}

const (
	DefaultPermission       = 0644
	DefaultOutputFolderName = "summary"
	DefaultOutputJson       = "out.json"
	DefaultOutputYaml       = "out.yaml"
	JsonType                = "json"
	YamlType                = "yaml"
	CsvType                 = "csv"
)

var (
	validExtensionFile    = []string{fmt.Sprintf(".%s", JsonType), fmt.Sprintf(".%s", CsvType)}
	errUnexpectedFileType = errors.New("unexpected file type")
)

// ProcessMetricData : Process data metric from parsed flags on the main function
func (pm *ProcessMetric) ProcessMetricData() error {
	// process the metrics data
	summary := make(map[string]int)

	//validate file type from the dir
	isValid := pm.validateFileTypeFromFile(pm.files, validExtensionFile)
	if !isValid {
		fmt.Println(fmt.Printf(" got error %v when reading metrics from file : %s ", errUnexpectedFileType, pm.files))
		return errUnexpectedFileType
	}

	for _, file := range pm.files {
		filename := filepath.Join(*pm.directory, file.Name())
		metrics, errReadMetrics := pm.readMetricsFromFile(filename, *pm.fileType)
		if errReadMetrics != nil {
			fmt.Printf(" got error %v when reading metrics from file : %s ", errReadMetrics, filename)
			return errReadMetrics
		}

		for _, metric := range metrics {
			timestamp, err := time.Parse(time.RFC3339, metric.Timestamp)
			if err != nil {
				fmt.Println("failed to parse timestamp:", metric.Timestamp)
				continue
			}

			if timestamp.After(pm.startTime) && timestamp.Before(pm.endTime) {
				summary[metric.LevelName] += metric.Value
			}
		}
	}

	// convert the summary to a slice of summaryResult objects
	var result []summaryResult
	for levelName, totalValue := range summary {
		result = append(result, summaryResult{LevelName: levelName, TotalValue: totalValue})
	}

	//display the summary metrics on file either json or yaml and display on stdout
	pm.summaryMetrics(pm.outputFileType, pm.outputFileName, result)

	return nil
}

// readMetricsFromFile : function to read metrics from files in a directory
func (pm *ProcessMetric) readMetricsFromFile(filename, fileType string) ([]metric, error) {
	var metrics []metric

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if fileType == JsonType {
		isFileTypeCorrect := pm.isJSONFile(filename)
		if !isFileTypeCorrect {
			return nil, errUnexpectedFileType
		}
		data, errReadAll := io.ReadAll(file)
		if errReadAll != nil {
			return nil, errReadAll
		}
		errUnmarshall := json.Unmarshal(data, &metrics)
		if errUnmarshall != nil {
			return nil, errUnmarshall
		}
	} else if fileType == CsvType {
		isFileTypeCorrect := pm.isCSVFile(filename)
		if !isFileTypeCorrect {
			return nil, errUnexpectedFileType
		}

		reader := csv.NewReader(file)
		records, errReaderReadAll := reader.ReadAll()
		if errReaderReadAll != nil {
			return nil, errReaderReadAll
		}

		isHeaderRow := true
		for _, record := range records {
			if isHeaderRow {
				isHeaderRow = false
				continue
			}

			if len(record) != 3 {
				return nil, fmt.Errorf("invalid CSV record: %v", record)
			}

			timestamp := record[0]
			levelName := record[1]
			value, errParseStrToInt := strconv.Atoi(strings.TrimSpace(record[2]))
			if errParseStrToInt != nil {
				// skip the record if the value cannot be parsed
				fmt.Printf("skipping record: Failed to parse value '%s': %v\n", record[2], errParseStrToInt)
				continue
			}

			metricSingle := metric{
				Timestamp: timestamp,
				LevelName: levelName,
				Value:     value,
			}
			metrics = append(metrics, metricSingle)
		}
	} else {
		return nil, fmt.Errorf("unsupported file type: %s", fileType)
	}

	return metrics, nil
}

// summaryMetrics : summarize the data metric into a file and display on stdout or console
func (pm *ProcessMetric) summaryMetrics(outputFileType, outputFileName *string, result []summaryResult) {
	// convert the result to the desired output format
	var output []byte
	var err error
	if *outputFileType == YamlType {
		output, err = yaml.Marshal(result)
		if err != nil {
			fmt.Println("failed to marshal YAML:", err)
			os.Exit(1)
		}
	} else {
		output, err = json.Marshal(result)
		if err != nil {
			fmt.Println("failed to marshal JSON:", err)
			os.Exit(1)
		}
	}

	// print the result in the console or stdout
	fmt.Println(string(output))

	// create the 'summary' folder if it doesn't exist
	err = os.MkdirAll(DefaultOutputFolderName, os.ModePerm)
	if err != nil {
		fmt.Println("failed to create 'summary' folder:", err)
		os.Exit(1)
	}

	// write the result to a file
	outputFileNameFinal := filepath.Join(DefaultOutputFolderName, pm.getOutputFileName(*outputFileName, *outputFileType))
	err = os.WriteFile(outputFileNameFinal, output, DefaultPermission)
	if err != nil {
		fmt.Println("failed to write output file:", err)
		os.Exit(1)
	}
}

// getOutputFileName : generate the output file name
func (pm *ProcessMetric) getOutputFileName(outputFileName, outputFileType string) string {
	if outputFileName != "" {
		return outputFileName
	}

	if outputFileType == YamlType {
		return DefaultOutputYaml
	}

	return DefaultOutputJson
}

func (pm *ProcessMetric) isJSONFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == fmt.Sprintf(".%s", JsonType)
}

func (pm *ProcessMetric) isCSVFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == fmt.Sprintf(".%s", CsvType)
}

func (pm *ProcessMetric) validateFileTypeFromFile(entries []os.DirEntry, validExtensions []string) bool {
	for _, entry := range entries {
		if entry.IsDir() {
			continue // Skip directories
		}
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		valid := false
		for _, validExt := range validExtensions {
			if ext == strings.ToLower(validExt) {
				valid = true
				break
			}
		}
		if !valid {
			return false
		}
	}
	return true
}
