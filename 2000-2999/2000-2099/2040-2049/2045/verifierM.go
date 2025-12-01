package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const refSource = "./2045M.go"

type testCase struct {
	name string
	R, C int
	grid []string
}

type sideEntry struct {
	dir   int
	index int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierM.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	for idx, tc := range tests {
		input := buildInput(tc)

		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, refOut)
			os.Exit(1)
		}
		refK, refList, err := parseOutput(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, refOut)
			os.Exit(1)
		}
		expected := evaluate(tc)
		if err := compareAnswers(tc, refK, refList, expected); err != nil {
			fmt.Fprintf(os.Stderr, "reference answer mismatch on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, candOut)
			os.Exit(1)
		}
		candK, candList, err := parseOutput(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, candOut)
			os.Exit(1)
		}
		if err := compareAnswers(tc, candK, candList, expected); err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2045M-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2045M.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.R, tc.C))
	for _, row := range tc.grid {
		sb.WriteString(row)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(output string) (int, []string, error) {
	scanner := bufio.NewScanner(strings.NewReader(output))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return 0, nil, fmt.Errorf("failed to read k: %v", err)
		}
		return 0, nil, fmt.Errorf("missing k")
	}
	kVal, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return 0, nil, fmt.Errorf("invalid k value %q", scanner.Text())
	}
	list := make([]string, 0, kVal)
	for scanner.Scan() {
		list = append(list, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return 0, nil, fmt.Errorf("failed to read locations: %v", err)
	}
	if kVal != len(list) {
		return 0, nil, fmt.Errorf("k=%d but got %d location strings", kVal, len(list))
	}
	return kVal, list, nil
}

func compareAnswers(tc testCase, k int, list []string, expected map[string]bool) error {
	if k != len(expected) {
		return fmt.Errorf("reported k=%d but expected %d", k, len(expected))
	}
	seen := make(map[string]bool, len(list))
	for _, loc := range list {
		if seen[loc] {
			return fmt.Errorf("duplicate location %q", loc)
		}
		seen[loc] = true
		if !expected[loc] {
			return fmt.Errorf("location %q is invalid", loc)
		}
	}
	if len(seen) != len(expected) {
		return fmt.Errorf("provided %d unique locations but expected %d", len(seen), len(expected))
	}
	return nil
}

func evaluate(tc testCase) map[string]bool {
	totalMirrors := countMirrors(tc.grid)
	valid := make(map[string]bool)

	for c := 0; c < tc.C; c++ {
		if simulate(tc, -1, c, 2, totalMirrors) {
			valid[fmt.Sprintf("N%d", c+1)] = true
		}
	}
	for c := 0; c < tc.C; c++ {
		if simulate(tc, tc.R, c, 0, totalMirrors) {
			valid[fmt.Sprintf("S%d", c+1)] = true
		}
	}
	for r := 0; r < tc.R; r++ {
		if simulate(tc, r, tc.C, 3, totalMirrors) {
			valid[fmt.Sprintf("E%d", r+1)] = true
		}
	}
	for r := 0; r < tc.R; r++ {
		if simulate(tc, r, -1, 1, totalMirrors) {
			valid[fmt.Sprintf("W%d", r+1)] = true
		}
	}
	return valid
}

func countMirrors(grid []string) int {
	count := 0
	for _, row := range grid {
		for j := 0; j < len(row); j++ {
			if row[j] == '/' || row[j] == '\\' {
				count++
			}
		}
	}
	return count
}

var dr = []int{-1, 0, 1, 0}
var dc = []int{0, 1, 0, -1}

func reflect(dir int, ch byte) int {
	switch ch {
	case '/':
		switch dir {
		case 0:
			return 3
		case 1:
			return 2
		case 2:
			return 1
		case 3:
			return 0
		}
	case '\\':
		switch dir {
		case 0:
			return 1
		case 1:
			return 0
		case 2:
			return 3
		case 3:
			return 2
		}
	}
	return dir
}

func simulate(tc testCase, startR, startC, dir int, totalMirrors int) bool {
	visited := make([]bool, tc.R*tc.C)
	seen := make(map[[3]int]struct{})
	hit := 0

	r, c := startR, startC
	for {
		r += dr[dir]
		c += dc[dir]
		if r < 0 || r >= tc.R || c < 0 || c >= tc.C {
			break
		}
		state := [3]int{r, c, dir}
		if _, ok := seen[state]; ok {
			break
		}
		seen[state] = struct{}{}

		cell := tc.grid[r][c]
		if cell == '/' || cell == '\\' {
			idx := r*tc.C + c
			if !visited[idx] {
				visited[idx] = true
				hit++
			}
			dir = reflect(dir, cell)
		}
	}
	return hit == totalMirrors
}

func buildTests() []testCase {
	tests := []testCase{
		{
			name: "sample1",
			R:    4, C: 4,
			grid: []string{".//.", ".\\.", ".\\/.", "...."},
		},
		{
			name: "sample2",
			R:    4, C: 6,
			grid: []string{"./..\\.", ".\\...\\", "./../\\", "......"},
		},
		{
			name: "sample3",
			R:    4, C: 4,
			grid: []string{"....", "./\\.", ".\\/.", "...."},
		},
	}

	rng := rand.New(rand.NewSource(123456789))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng, i))
	}
	return tests
}

func randomTest(rng *rand.Rand, idx int) testCase {
	R := rng.Intn(5) + 1
	C := rng.Intn(5) + 1
	grid := make([]string, R)
	hasMirror := false
	for i := 0; i < R; i++ {
		row := make([]byte, C)
		for j := 0; j < C; j++ {
			val := rng.Intn(4)
			switch val {
			case 0:
				row[j] = '.'
			case 1:
				row[j] = '/'
				hasMirror = true
			case 2:
				row[j] = '\\'
				hasMirror = true
			default:
				row[j] = '.'
			}
		}
		grid[i] = string(row)
	}
	if !hasMirror {
		r := rng.Intn(R)
		c := rng.Intn(C)
		if rng.Intn(2) == 0 {
			grid[r] = replaceAt(grid[r], c, '/')
		} else {
			grid[r] = replaceAt(grid[r], c, '\\')
		}
	}
	return testCase{
		name: fmt.Sprintf("random_%d", idx+1),
		R:    R,
		C:    C,
		grid: grid,
	}
}

func replaceAt(s string, idx int, ch byte) string {
	b := []byte(s)
	b[idx] = ch
	return string(b)
}
