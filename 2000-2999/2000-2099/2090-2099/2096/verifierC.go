package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const refSource = "2000-2999/2000-2099/2090-2099/2096/2096C.go"

type testCase struct {
	n    int
	h    [][]int
	row  []int64
	col  []int64
	name string
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	candidate := args[0]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}
	refAns, err := parseOutput(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n%s", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	candAns, err := parseOutput(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\n%s", err, candOut)
		os.Exit(1)
	}

	for i, tc := range tests {
		if refAns[i] != candAns[i] {
			fmt.Fprintf(os.Stderr, "case %d (%s) mismatch: expected %d got %d\ninput:\n%s", i+1, tc.name, refAns[i], candAns[i], formatCase(tc))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	outPath := "./ref_2096C.bin"
	cmd := exec.Command("go", "build", "-o", outPath, refSource)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return outPath, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func buildTests() []testCase {
	var tests []testCase
	add := func(name string, h [][]int, row, col []int64) {
		n := len(h)
		cpyH := make([][]int, n)
		for i := 0; i < n; i++ {
			cpyH[i] = append([]int(nil), h[i]...)
		}
		cpyRow := append([]int64(nil), row...)
		cpyCol := append([]int64(nil), col...)
		tests = append(tests, testCase{n: n, h: cpyH, row: cpyRow, col: cpyCol, name: name})
	}

	// Small deterministic cases
	add("n1", [][]int{{0}}, []int64{5}, []int64{7})
	add("n2_simple", [][]int{{1, 1}, {1, 1}}, []int64{1, 1}, []int64{1, 1})

	// Force impossible by banning all transitions between rows and cols
	impH := [][]int{
		{0, 1, 0},
		{0, -1, 1},
		{1, 0, -1},
	}
	add("impossible", impH, []int64{1, 1, 1}, []int64{1, 1, 1})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	totalCells := 0
	const maxCells = 40000

	for len(tests) < 140 && totalCells < maxCells {
		n := rng.Intn(25) + 2 // 2..26
		if len(tests)%20 == 0 {
			n = rng.Intn(60) + 30 // occasional larger, still manageable
		}
		if totalCells+n*n > maxCells {
			break
		}
		h := make([][]int, n)
		for i := 0; i < n; i++ {
			h[i] = make([]int, n)
			for j := 0; j < n; j++ {
				// small heights to craft varied diffs
				h[i][j] = rng.Intn(7) - 3
			}
		}
		row := make([]int64, n)
		col := make([]int64, n)
		for i := 0; i < n; i++ {
			row[i] = int64(rng.Intn(50) + 1)
			col[i] = int64(rng.Intn(50) + 1)
		}
		add(fmt.Sprintf("random_%d", len(tests)), h, row, col)
		totalCells += n * n
	}

	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i := 0; i < tc.n; i++ {
			for j := 0; j < tc.n; j++ {
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(strconv.Itoa(tc.h[i][j]))
			}
			sb.WriteByte('\n')
		}
		for i := 0; i < tc.n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(tc.row[i], 10))
		}
		sb.WriteByte('\n')
		for i := 0; i < tc.n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(tc.col[i], 10))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string, expected int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d outputs, got %d", expected, len(fields))
	}
	ans := make([]int64, expected)
	for i, s := range fields {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", s)
		}
		ans[i] = v
	}
	return ans, nil
}

func formatCase(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i := 0; i < tc.n; i++ {
		for j := 0; j < tc.n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(tc.h[i][j]))
		}
		sb.WriteByte('\n')
	}
	for i := 0; i < tc.n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(tc.row[i], 10))
	}
	sb.WriteByte('\n')
	for i := 0; i < tc.n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(tc.col[i], 10))
	}
	sb.WriteByte('\n')
	return sb.String()
}
