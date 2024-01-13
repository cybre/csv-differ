package handler

import (
	"bufio"
	"context"
	"net/http"
	"slices"
	"strings"

	"github.com/cybre/csv-differ/internal/differ"
	"github.com/cybre/csv-differ/internal/templates"
	"github.com/cybre/csv-differ/internal/utils"
)

func Diff() Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		err := r.ParseMultipartForm(64 << 20)
		if err != nil {
			return StatusError{http.StatusBadRequest, err}
		}

		file1, _, err := r.FormFile("file1")
		if err != nil {
			return StatusError{http.StatusBadRequest, err}
		}
		defer file1.Close()

		file2, _, err := r.FormFile("file2")
		if err != nil {
			return StatusError{http.StatusBadRequest, err}
		}
		defer file2.Close()

		file1Delimiter := r.FormValue("file1Delimiter")
		file2Delimiter := r.FormValue("file2Delimiter")

		scanner := bufio.NewScanner(file1)
		scanner.Scan()
		file1ColumnNames := utils.Map(strings.Split(scanner.Text(), file1Delimiter), func(item string) string {
			return strings.TrimFunc(item, func(r rune) bool {
				return r == '"'
			})
		})

		scanner = bufio.NewScanner(file2)
		scanner.Scan()
		file2ColumnNames := utils.Map(strings.Split(scanner.Text(), file2Delimiter), func(item string) string {
			return strings.TrimFunc(item, func(r rune) bool {
				return r == '"'
			})
		})

		// Extract the column names to compare for the diff.
		file1KeyColumnName := r.FormValue("file1Column")
		file2KeyColumnName := r.FormValue("file2Column")
		if file1KeyColumnName == "" || file2KeyColumnName == "" {
			if err := templates.PickColumns(file1ColumnNames, file2ColumnNames).Render(context.Background(), w); err != nil {
				return StatusError{http.StatusInternalServerError, err}
			}

			return nil
		}

		// Find the indexes of the columns to use for the diff.
		file1ColumnIndex := slices.Index(file1ColumnNames, file1KeyColumnName)
		file2ColumnIndex := slices.Index(file2ColumnNames, file2KeyColumnName)

		result, err := differ.Diff(differ.File{
			File:            file1,
			Delimiter:       file1Delimiter,
			DiffColumnIndex: file1ColumnIndex,
		}, differ.File{
			File:            file2,
			Delimiter:       file2Delimiter,
			DiffColumnIndex: file2ColumnIndex,
		})
		if err != nil {
			return StatusError{http.StatusInternalServerError, err}
		}

		if err := templates.Diff(result).Render(context.Background(), w); err != nil {
			return StatusError{http.StatusInternalServerError, err}
		}

		return nil
	}
}
