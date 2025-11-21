package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const referenceSolutionRel = "0-999/900-999/900-909/907/907B.go"

var referenceSolutionPath string

func init() {
	referenceSolutionPath = referenceSolutionRel
	if _, file, _, ok := runtime.Caller(0); ok {
		dir := filepath.Dir(file)
		candidate := filepath.Join(dir, "907B.go")
		if _, err := os.Stat(candidate); err == nil {
			referenceSolutionPath = candidate
			return
		}
	}
	if abs, err := filepath.Abs(referenceSolutionRel); err == nil {
		if _, err := os.Stat(abs); err == nil {
			referenceSolutionPath = abs
		}
	}
}

type testCase struct {
	name  string
	input string
}

func formatInput(rows []string, x, y int) string {
	var sb strings.Builder
	for i := 0; i < 9; i++ {
		row := rows[i]
		if len(row) != 9 {
			panic(fmt.Sprintf("row %d length %d != 9", i, len(row)))
		}
		for j := 0; j < 9; j++ {
			sb.WriteByte(row[j])
			if j%3 == 2 && j != 8 {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')
		if i%3 == 2 && i != 8 {
			sb.WriteByte('\n')
		}
	}
	sb.WriteString(fmt.Sprintf("%d %d\n", x, y))
	return sb.String()
}

func deterministicTests() []testCase {
	cases := []struct {
		name string
		rows []string
		x, y int
	}{
		{
			name: "center_partial",
			rows: []string{
				".........",
				"..x.x....",
				".........",
				".........",
				"....x....",
				".........",
				"...o.....",
				".........",
				"......o..",
			},
			x: 5, y: 5,
		},
		{
			name: "target_full",
			rows: []string{
				"xxxooooxx",
				"oxoxoxoxo",
				"xxooxxoxo",
				".........",
				"...xx....",
				".........",
				"ooxoxoxox",
				"xoxoxoxoo",
				"oxoxoxoxo",
			},
			x: 2, y: 2,
		},
		{
			name: "edge_target",
			rows: []string{
				"..x......",
				"...o.....",
				"..x......",
				".........",
				".........",
				".........",
				"......x..",
				"...o..o..",
				"..x...x..",
			},
			x: 9, y: 8,
		},
	}
	var tests []testCase
	for _, c := range cases {
		tests = append(tests, testCase{
			name:  c.name,
			input: formatInput(c.rows, c.x, c.y),
		})
	}
	return tests
}

func randomBoard(rng *rand.Rand) ([]string, int, int) {
	grid := make([][]byte, 9)
	hasDot := false
	hasFilled := false
	for i := 0; i < 9; i++ {
		grid[i] = make([]byte, 9)
		for j := 0; j < 9; j++ {
			val := rng.Intn(3)
			switch val {
			case 0:
				grid[i][j] = 'x'
				hasFilled = true
			case 1:
				grid[i][j] = 'o'
				hasFilled = true
			default:
				grid[i][j] = '.'
				hasDot = true
			}
		}
	}
	if !hasDot {
		i := rng.Intn(9)
		j := rng.Intn(9)
		grid[i][j] = '.'
		hasDot = true
	}
	if !hasFilled {
		i := rng.Intn(9)
		j := rng.Intn(9)
		grid[i][j] = 'x'
		hasFilled = true
	}
	type coord struct{ x, y int }
	var filled []coord
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if grid[i][j] != '.' {
				filled = append(filled, coord{i, j})
			}
		}
	}
	if len(filled) == 0 {
		// Should not happen, but guard anyway.
		grid[0][0] = 'x'
		filled = append(filled, coord{0, 0})
	}
	last := filled[rng.Intn(len(filled))]
	rows := make([]string, 9)
	for i := 0; i < 9; i++ {
		rows[i] = string(grid[i])
	}
	return rows, last.x + 1, last.y + 1
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(20250305))
	var tests []testCase
	for i := 0; i < 200; i++ {
		rows, x, y := randomBoard(rng)
		tests = append(tests, testCase{
			name:  fmt.Sprintf("random_%d", i+1),
			input: formatInput(rows, x, y),
		})
	}
	return tests
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildReferenceBinary() (string, func(), error) {
	if referenceSolutionPath == "" {
		return "", nil, fmt.Errorf("reference solution path not set")
	}
	if _, err := os.Stat(referenceSolutionPath); err != nil {
		return "", nil, fmt.Errorf("reference solution not found: %v", err)
	}
	tmpDir, err := os.MkdirTemp("", "907B-ref")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref_907B")
	cmd := exec.Command("go", "build", "-o", binPath, referenceSolutionPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return binPath, cleanup, nil
}

func parseBoard(output string) ([]string, error) {
	lines := strings.Split(output, "\n")
	rows := make([]string, 0, 9)
	for _, line := range lines {
		clean := strings.ReplaceAll(line, " ", "")
		if clean == "" {
			continue
		}
		if len(rows) < 9 {
			if len(clean) != 9 {
				return nil, fmt.Errorf("row %d has length %d (want 9)", len(rows)+1, len(clean))
			}
			for i := 0; i < 9; i++ {
				ch := clean[i]
				if ch != '.' && ch != 'x' && ch != 'o' && ch != '!' {
					return nil, fmt.Errorf("invalid character %q in row %d", ch, len(rows)+1)
				}
			}
			rows = append(rows, clean)
			continue
		}
		// Already have 9 rows; any further non-space content is invalid.
		return nil, fmt.Errorf("unexpected extra output after board: %q", line)
	}
	if len(rows) != 9 {
		return nil, fmt.Errorf("expected 9 non-empty rows, got %d", len(rows))
	}
	return rows, nil
}

func compareBoards(a, b []string) (bool, string) {
	for i := 0; i < 9; i++ {
		if a[i] != b[i] {
			return false, fmt.Sprintf("row %d mismatch: expected %q, got %q", i+1, a[i], b[i])
		}
	}
	return true, ""
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	refBin, cleanup, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := append(deterministicTests(), randomTests()...)
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %s (%d): %v\ninput:\n%soutput:\n%s\n", tc.name, idx+1, err, tc.input, refOut)
			os.Exit(1)
		}
		refBoard, err := parseBoard(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference parse error on test %s (%d): %v\ninput:\n%soutput:\n%s\n", tc.name, idx+1, err, tc.input, refOut)
			os.Exit(1)
		}

		out, runErr := runProgram(bin, tc.input)
		if runErr != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %s (%d): %v\ninput:\n%soutput:\n%s\n", tc.name, idx+1, runErr, tc.input, out)
			os.Exit(1)
		}
		board, err := parseBoard(out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse error on test %s (%d): %v\ninput:\n%soutput:\n%s\n", tc.name, idx+1, err, tc.input, out)
			os.Exit(1)
		}
		if ok, msg := compareBoards(refBoard, board); !ok {
			fmt.Fprintf(os.Stderr, "wrong answer on test %s (%d): %s\ninput:\n%sreference output:\n%s\nparticipant output:\n%s\n", tc.name, idx+1, msg, tc.input, refOut, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
