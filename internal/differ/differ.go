package differ

import (
	"bufio"
	"fmt"
	"io"
	"slices"
	"strings"
)

type File struct {
	File            io.Reader
	Delimiter       string
	DiffColumnIndex int
}

type Result struct {
	Additions []string
	Deletions []string
}

func Diff(file1, file2 File) (Result, error) {
	file1ColumnValues, err := extractColumnValues(file1)
	if err != nil {
		return Result{}, fmt.Errorf("error extracting column values from file 1: %v", err)
	}

	file2ColumnValues, err := extractColumnValues(file2)
	if err != nil {
		return Result{}, fmt.Errorf("error extracting column values from file 2: %v", err)
	}

	additions := make([]string, 0)
	deletions := make([]string, 0)
	for _, value := range file1ColumnValues {
		if !slices.Contains(file2ColumnValues, value) {
			deletions = append(deletions, value)
		}
	}

	for _, value := range file2ColumnValues {
		if !slices.Contains(file1ColumnValues, value) {
			additions = append(additions, value)
		}
	}

	return Result{
		Additions: additions,
		Deletions: deletions,
	}, nil
}

func extractColumnValues(file File) ([]string, error) {
	scanner := bufio.NewScanner(file.File)
	scanner.Scan() // Skip the header row.
	values := make([]string, 0)
	for scanner.Scan() {
		row := strings.Split(scanner.Text(), file.Delimiter)
		if (len(row) - 1) < file.DiffColumnIndex {
			break
		}

		values = append(values, strings.ToLower(row[file.DiffColumnIndex]))
	}

	return values, nil
}
