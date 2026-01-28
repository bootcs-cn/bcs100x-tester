package stages

import (
	"database/sql"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/bootcs-dev/tester-utils/test_case_harness"
	"github.com/bootcs-dev/tester-utils/tester_definition"
)

func moviesTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:     "movies",
		Timeout:  60 * time.Second,
		TestFunc: testMovies,
	}
}

func testMovies(harness *test_case_harness.TestCaseHarness) error {
	logger := harness.Logger
	workDir := harness.SubmissionDir

	// 1. 检查 SQL 文件存在 (1.sql - 13.sql)
	logger.Infof("Checking SQL files exist...")
	for i := 1; i <= 13; i++ {
		filename := fmt.Sprintf("%d.sql", i)
		if !harness.FileExists(filename) {
			return fmt.Errorf("%s does not exist", filename)
		}
	}
	logger.Successf("SQL files exist")

	// 2. 打开数据库
	dbPath := filepath.Join(workDir, "movies.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open movies.db: %v", err)
	}
	defer db.Close()

	// Test 1: 2008 年电影 (无序)
	logger.Infof("Testing 1.sql produces correct result...")
	if err := testMoviesSingleColUnordered(db, workDir, "1.sql", expectedMovies1); err != nil {
		return fmt.Errorf("1.sql: %v", err)
	}
	logger.Successf("✓ 1.sql produces correct result")

	// Test 2: Emma Stone 出生年份 (单值)
	logger.Infof("Testing 2.sql produces correct result...")
	if err := testMoviesSingleValue(db, workDir, "2.sql", "1988"); err != nil {
		return fmt.Errorf("2.sql: %v", err)
	}
	logger.Successf("✓ 2.sql produces correct result")

	// Test 3: 2018+ 电影按字母排序 (有序)
	logger.Infof("Testing 3.sql produces correct result...")
	if err := testMoviesSingleColOrdered(db, workDir, "3.sql", expectedMovies3); err != nil {
		return fmt.Errorf("3.sql: %v", err)
	}
	logger.Successf("✓ 3.sql produces correct result")

	// Test 4: 10.0 评分电影数量 (单值)
	logger.Infof("Testing 4.sql produces correct result...")
	if err := testMoviesSingleValue(db, workDir, "4.sql", "2"); err != nil {
		return fmt.Errorf("4.sql: %v", err)
	}
	logger.Successf("✓ 4.sql produces correct result")

	// Test 5: Harry Potter 电影 (双列有序)
	logger.Infof("Testing 5.sql produces correct result...")
	if err := testMoviesDoubleColOrdered(db, workDir, "5.sql", expectedMovies5); err != nil {
		return fmt.Errorf("5.sql: %v", err)
	}
	logger.Successf("✓ 5.sql produces correct result")

	// Test 6: 2012 年平均评分 (浮点数)
	logger.Infof("Testing 6.sql produces correct result...")
	if err := testMoviesFloatValue(db, workDir, "6.sql", 7.74, 0.01); err != nil {
		return fmt.Errorf("6.sql: %v", err)
	}
	logger.Successf("✓ 6.sql produces correct result")

	// Test 7: 2010 年电影及评分 (双列有序)
	logger.Infof("Testing 7.sql produces correct result...")
	if err := testMoviesDoubleColOrdered(db, workDir, "7.sql", expectedMovies7); err != nil {
		return fmt.Errorf("7.sql: %v", err)
	}
	logger.Successf("✓ 7.sql produces correct result")

	// Test 8: Toy Story 演员 (无序)
	logger.Infof("Testing 8.sql produces correct result...")
	if err := testMoviesSingleColUnordered(db, workDir, "8.sql", expectedMovies8); err != nil {
		return fmt.Errorf("8.sql: %v", err)
	}
	logger.Successf("✓ 8.sql produces correct result")

	// Test 9: 2004 年电影演员按出生年份排序 (有序)
	logger.Infof("Testing 9.sql produces correct result...")
	if err := testMoviesSingleColOrdered(db, workDir, "9.sql", expectedMovies9); err != nil {
		return fmt.Errorf("9.sql: %v", err)
	}
	logger.Successf("✓ 9.sql produces correct result")

	// Test 10: 9.0+ 评分电影导演 (无序)
	logger.Infof("Testing 10.sql produces correct result...")
	if err := testMoviesSingleColUnordered(db, workDir, "10.sql", expectedMovies10); err != nil {
		return fmt.Errorf("10.sql: %v", err)
	}
	logger.Successf("✓ 10.sql produces correct result")

	// Test 11: Chadwick Boseman 电影按评分排序 (有序)
	logger.Infof("Testing 11.sql produces correct result...")
	if err := testMoviesSingleColOrdered(db, workDir, "11.sql", expectedMovies11); err != nil {
		return fmt.Errorf("11.sql: %v", err)
	}
	logger.Successf("✓ 11.sql produces correct result")

	// Test 12: Johnny Depp & Helena Bonham Carter 共同电影 (无序，支持两种答案)
	logger.Infof("Testing 12.sql produces correct result...")
	if err := testMovies12(db, workDir); err != nil {
		return fmt.Errorf("12.sql: %v", err)
	}
	logger.Successf("✓ 12.sql produces correct result")

	// Test 13: Kevin Bacon 合作演员 (无序)
	logger.Infof("Testing 13.sql produces correct result...")
	if err := testMoviesSingleColUnordered(db, workDir, "13.sql", expectedMovies13); err != nil {
		return fmt.Errorf("13.sql: %v", err)
	}
	logger.Successf("✓ 13.sql produces correct result")

	logger.Successf("All tests passed!")
	return nil
}

// readMoviesSQLFile 读取 SQL 文件内容
func readMoviesSQLFile(workDir, filename string) (string, error) {
	content, err := os.ReadFile(filepath.Join(workDir, filename))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}

// executeMoviesQuery 执行 SQL 查询并返回单列结果
func executeMoviesQuery(db *sql.DB, query string) ([]string, error) {
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

// executeMoviesQueryDouble 执行 SQL 查询并返回双列结果
func executeMoviesQueryDouble(db *sql.DB, query string) ([][2]string, error) {
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

// testMoviesSingleColUnordered 测试无序单列结果
func testMoviesSingleColUnordered(db *sql.DB, workDir, filename string, expected []string) error {
	query, err := readMoviesSQLFile(workDir, filename)
	if err != nil {
		return err
	}

	actual, err := executeMoviesQuery(db, query)
	if err != nil {
		return err
	}

	if !moviesEqualSets(actual, expected) {
		return fmt.Errorf("result mismatch: expected %d rows, got %d rows", len(expected), len(actual))
	}
	return nil
}

// testMoviesSingleColOrdered 测试有序单列结果
func testMoviesSingleColOrdered(db *sql.DB, workDir, filename string, expected []string) error {
	query, err := readMoviesSQLFile(workDir, filename)
	if err != nil {
		return err
	}

	actual, err := executeMoviesQuery(db, query)
	if err != nil {
		return err
	}

	if !moviesEqualSlices(actual, expected) {
		return fmt.Errorf("result mismatch: expected %v, got %v", expected, actual)
	}
	return nil
}

// testMoviesSingleValue 测试单个值结果
func testMoviesSingleValue(db *sql.DB, workDir, filename, expected string) error {
	query, err := readMoviesSQLFile(workDir, filename)
	if err != nil {
		return err
	}

	actual, err := executeMoviesQuery(db, query)
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

// testMoviesFloatValue 测试浮点数结果
func testMoviesFloatValue(db *sql.DB, workDir, filename string, expected, tolerance float64) error {
	query, err := readMoviesSQLFile(workDir, filename)
	if err != nil {
		return err
	}

	var actual float64
	err = db.QueryRow(query).Scan(&actual)
	if err != nil {
		return fmt.Errorf("query error: %v", err)
	}

	if math.Abs(actual-expected) > tolerance {
		return fmt.Errorf("expected %.2f (±%.2f), got %.2f", expected, tolerance, actual)
	}
	return nil
}

// testMoviesDoubleColOrdered 测试有序双列结果
func testMoviesDoubleColOrdered(db *sql.DB, workDir, filename string, expected [][2]string) error {
	query, err := readMoviesSQLFile(workDir, filename)
	if err != nil {
		return err
	}

	actual, err := executeMoviesQueryDouble(db, query)
	if err != nil {
		return err
	}

	if len(actual) != len(expected) {
		return fmt.Errorf("expected %d rows, got %d", len(expected), len(actual))
	}

	for i := range actual {
		// 比较时忽略列顺序（作为集合比较），并处理浮点数格式差异
		if !rowsMatch(actual[i], expected[i]) {
			return fmt.Errorf("row %d mismatch: expected {%s, %s}, got {%s, %s}",
				i+1, expected[i][0], expected[i][1], actual[i][0], actual[i][1])
		}
	}
	return nil
}

// rowsMatch 比较两行是否匹配（处理浮点数格式差异）
func rowsMatch(actual, expected [2]string) bool {
	// 尝试直接匹配
	if (actual[0] == expected[0] && actual[1] == expected[1]) ||
		(actual[0] == expected[1] && actual[1] == expected[0]) {
		return true
	}

	// 尝试作为集合匹配（处理浮点数格式）
	actualSet := map[string]bool{}
	expectedSet := map[string]bool{}

	for _, v := range actual {
		actualSet[normalizeNumeric(v)] = true
	}
	for _, v := range expected {
		expectedSet[normalizeNumeric(v)] = true
	}

	return mapsEqual(actualSet, expectedSet)
}

// normalizeNumeric 标准化数字字符串（去除尾部零）
func normalizeNumeric(s string) string {
	// 尝试解析为浮点数
	var f float64
	if _, err := fmt.Sscanf(s, "%f", &f); err == nil {
		// 如果是整数，返回整数格式
		if f == float64(int64(f)) {
			return fmt.Sprintf("%d", int64(f))
		}
		// 否则返回浮点数格式（去除尾部零）
		return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", f), "0"), ".")
	}
	return s
}

// testMovies12 处理 test12 的两种可能答案
func testMovies12(db *sql.DB, workDir string) error {
	query, err := readMoviesSQLFile(workDir, "12.sql")
	if err != nil {
		return err
	}

	actual, err := executeMoviesQuery(db, query)
	if err != nil {
		return err
	}

	// 尝试第一种答案 (Johnny Depp & Helena Bonham Carter)
	if moviesEqualSets(actual, expectedMovies12a) {
		return nil
	}

	// 尝试第二种答案 (Bradley Cooper & Jennifer Lawrence)
	if moviesEqualSets(actual, expectedMovies12b) {
		return nil
	}

	return fmt.Errorf("result does not match either expected answer")
}

// moviesEqualSets 比较两个字符串切片是否包含相同元素（忽略顺序）
func moviesEqualSets(a, b []string) bool {
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

// moviesEqualSlices 比较两个字符串切片是否完全相等（有序）
func moviesEqualSlices(a, b []string) bool {
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

// mapsEqual 比较两个 map 是否相等
func mapsEqual(a, b map[string]bool) bool {
	if len(a) != len(b) {
		return false
	}
	for k := range a {
		if !b[k] {
			return false
		}
	}
	return true
}

// 预期结果数据 (对齐 CS50 check50)

// Test 1: 2008 年电影
var expectedMovies1 = []string{
	"Iron Man",
	"The Dark Knight",
	"Slumdog Millionaire",
	"Kung Fu Panda",
}

// Test 3: 2018+ 电影按字母排序
var expectedMovies3 = []string{
	"Avengers: Infinity War",
	"Black Panther",
	"Eighth Grade",
	"Gemini Man",
	"Happy Times",
	"Incredibles 2",
	"Kirklet",
	"Ma Rainey's Black Bottom",
	"Roma",
	"The Professor",
	"Toy Story 4",
}

// Test 5: Harry Potter 电影 (title, year)
var expectedMovies5 = [][2]string{
	{"Harry Potter and the Sorcerer's Stone", "2001"},
	{"Harry Potter and the Chamber of Secrets", "2002"},
	{"Harry Potter and the Prisoner of Azkaban", "2004"},
	{"Harry Potter and the Goblet of Fire", "2005"},
	{"Harry Potter and the Order of the Phoenix", "2007"},
	{"Harry Potter and the Half-Blood Prince", "2009"},
	{"Harry Potter and the Deathly Hallows: Part 1", "2010"},
	{"Harry Potter and the Deathly Hallows: Part 2", "2011"},
	{"Harry Potter: A History of Magic", "2017"},
}

// Test 7: 2010 年电影及评分 (title, rating)
var expectedMovies7 = [][2]string{
	{"Inception", "8.8"},
	{"Toy Story 3", "8.3"},
	{"How to Train Your Dragon", "8.1"},
	{"Shutter Island", "8.1"},
	{"The King's Speech", "8.0"},
	{"Harry Potter and the Deathly Hallows: Part 1", "7.7"},
	{"Iron Man 2", "7.0"},
	{"Alice in Wonderland", "6.4"},
}

// Test 8: Toy Story 演员
var expectedMovies8 = []string{
	"Don Rickles",
	"Jim Varney",
	"Tom Hanks",
	"Tim Allen",
}

// Test 9: 2004 年电影演员按出生年份排序
var expectedMovies9 = []string{
	"Craig T. Nelson",
	"Richard Griffifths",
	"Samuel L. Jackson",
	"Holly Hunter",
	"Jason Lee",
	"Rupert Grint",
	"Daniel Radcliffe",
	"Emma Watson",
}

// Test 10: 9.0+ 评分电影导演
var expectedMovies10 = []string{
	"Christopher Nolan",
	"Frank Darabont",
	"Yimou Zhang",
}

// Test 11: Chadwick Boseman 电影按评分排序
var expectedMovies11 = []string{
	"42",
	"Black Panther",
	"Marshall",
	"Ma Rainey's Black Bottom",
	"Get on Up",
	"Draft Day",
	"Message from the King",
}

// Test 12a: Johnny Depp & Helena Bonham Carter 共同电影
var expectedMovies12a = []string{
	"Corpse Bride",
	"Charlie and the Chocolate Factory",
	"Alice in Wonderland",
	"Alice Through the Looking Glass",
}

// Test 12b: Bradley Cooper & Jennifer Lawrence 共同电影 (备选答案)
var expectedMovies12b = []string{
	"Silver Linings Playbook",
	"Serena",
	"American Hustle",
	"Joy",
}

// Test 13: Kevin Bacon 合作演员
var expectedMovies13 = []string{
	"Bill Paxton",
	"Gary Sinise",
	"James McAvoy",
	"Jennifer Lawrence",
	"Tom Cruise",
	"Michael Fassbender",
	"Tom Hanks",
}
