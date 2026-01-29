package helpers

import (
	"database/sql"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ReadSQLFile reads SQL file content from the working directory
func ReadSQLFile(workDir, filename string) (string, error) {
	content, err := os.ReadFile(filepath.Join(workDir, filename))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}

// ExecuteQuerySingleCol executes SQL query and returns single column results
func ExecuteQuerySingleCol(db *sql.DB, query string) ([]string, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var results []string
	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		results = append(results, value)
	}
	return results, nil
}

// ExecuteQueryFloat executes SQL query and returns a single float result
func ExecuteQueryFloat(db *sql.DB, query string) (float64, error) {
	var result float64
	err := db.QueryRow(query).Scan(&result)
	if err != nil {
		return 0, fmt.Errorf("query error: %v", err)
	}
	return result, nil
}

// ExecuteQueryDoubleCol executes SQL query and returns two column results
func ExecuteQueryDoubleCol(db *sql.DB, query string) ([][2]string, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var results [][2]string
	for rows.Next() {
		var col1, col2 string
		if err := rows.Scan(&col1, &col2); err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		results = append(results, [2]string{col1, col2})
	}
	return results, nil
}

// TestSQLSingleColUnordered tests unordered single column results
func TestSQLSingleColUnordered(db *sql.DB, workDir, filename string, expected []string) error {
	query, err := ReadSQLFile(workDir, filename)
	if err != nil {
		return err
	}

	actual, err := ExecuteQuerySingleCol(db, query)
	if err != nil {
		return err
	}

	if !EqualSets(actual, expected) {
		return fmt.Errorf("result mismatch: expected %d rows, got %d rows", len(expected), len(actual))
	}
	return nil
}

// TestSQLSingleColOrdered tests ordered single column results
func TestSQLSingleColOrdered(db *sql.DB, workDir, filename string, expected []string) error {
	query, err := ReadSQLFile(workDir, filename)
	if err != nil {
		return err
	}

	actual, err := ExecuteQuerySingleCol(db, query)
	if err != nil {
		return err
	}

	if !EqualSlices(actual, expected) {
		return fmt.Errorf("result mismatch: expected %v, got %v", expected, actual)
	}
	return nil
}

// TestSQLSingleValue tests single value results
func TestSQLSingleValue(db *sql.DB, workDir, filename, expected string) error {
	query, err := ReadSQLFile(workDir, filename)
	if err != nil {
		return err
	}

	actual, err := ExecuteQuerySingleCol(db, query)
	if err != nil {
		return err
	}

	if len(actual) != 1 {
		return fmt.Errorf("expected 1 row, got %d", len(actual))
	}

	if actual[0] != expected {
		return fmt.Errorf("expected %s, got %s", expected, actual[0])
	}
	return nil
}

// TestSQLFloat tests single float result with tolerance
func TestSQLFloat(db *sql.DB, workDir, filename string, expected, tolerance float64) error {
	query, err := ReadSQLFile(workDir, filename)
	if err != nil {
		return err
	}

	actual, err := ExecuteQueryFloat(db, query)
	if err != nil {
		return err
	}

	if math.Abs(actual-expected) > tolerance {
		return fmt.Errorf("result mismatch: expected %.5f (±%.2f), got %.5f", expected, tolerance, actual)
	}
	return nil
}

// TestSQLDoubleColOrdered tests ordered double column results
func TestSQLDoubleColOrdered(db *sql.DB, workDir, filename string, expected [][2]string) error {
	query, err := ReadSQLFile(workDir, filename)
	if err != nil {
		return err
	}

	actual, err := ExecuteQueryDoubleCol(db, query)
	if err != nil {
		return err
	}

	if len(actual) != len(expected) {
		return fmt.Errorf("expected %d rows, got %d", len(expected), len(actual))
	}

	for i := range actual {
		if !rowsMatch(actual[i], expected[i]) {
			return fmt.Errorf("row %d mismatch: expected {%s, %s}, got {%s, %s}",
				i+1, expected[i][0], expected[i][1], actual[i][0], actual[i][1])
		}
	}
	return nil
}

// EqualSets compares two string slices as sets (ignoring order)
func EqualSets(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	aCopy := make([]string, len(a))
	bCopy := make([]string, len(b))
	copy(aCopy, a)
	copy(bCopy, b)

	sort.Strings(aCopy)
	sort.Strings(bCopy)

	for i := range aCopy {
		if aCopy[i] != bCopy[i] {
			return false
		}
	}
	return true
}

// EqualSlices compares two string slices for exact equality (ordered)
func EqualSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// rowsMatch compares two rows for equality (handles column order and numeric formats)
func rowsMatch(actual, expected [2]string) bool {
	// Normalize both rows by converting numeric strings to consistent format
	normalizedActual := [2]string{normalizeValue(actual[0]), normalizeValue(actual[1])}
	normalizedExpected := [2]string{normalizeValue(expected[0]), normalizeValue(expected[1])}

	// Try direct match or reversed match on normalized values
	if (normalizedActual[0] == normalizedExpected[0] && normalizedActual[1] == normalizedExpected[1]) ||
		(normalizedActual[0] == normalizedExpected[1] && normalizedActual[1] == normalizedExpected[0]) {
		return true
	}

	return false
}

// normalizeValue normalizes a string value, converting numbers to consistent format
func normalizeValue(s string) string {
	// Try to parse as float
	var f float64
	if _, err := fmt.Sscanf(s, "%f", &f); err == nil {
		// It's a number - format consistently
		// If it's an integer, return as integer
		if f == float64(int64(f)) {
			return fmt.Sprintf("%d", int64(f))
		}
		// Otherwise return as float with consistent precision
		return fmt.Sprintf("%.10g", f)
	}
	// Not a number, return as-is
	return s
}

// tryFloatMatch tries to match rows as floating point numbers (deprecated, kept for compatibility)
func tryFloatMatch(actual, expected [2]string) bool {
	a0, err1 := parseFloat(actual[0])
	a1, err2 := parseFloat(actual[1])
	e0, err3 := parseFloat(expected[0])
	e1, err4 := parseFloat(expected[1])

	if err1 == nil && err2 == nil && err3 == nil && err4 == nil {
		tolerance := 0.01
		return (math.Abs(a0-e0) < tolerance && math.Abs(a1-e1) < tolerance) ||
			(math.Abs(a0-e1) < tolerance && math.Abs(a1-e0) < tolerance)
	}
	return false
}

// parseFloat parses a string as float64
func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}
