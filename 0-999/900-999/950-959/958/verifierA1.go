package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const refSource = "0-999/900-999/950-959/958/958A1.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA1.go /path/to/binary")
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
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		exp, err := parseAnswer(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		got, err := parseAnswer(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %s got %s\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				idx+1, tc.name, exp, got, tc.input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-958A1-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref958A1.bin")
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
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseAnswer(out string) (string, error) {
	ans := strings.TrimSpace(out)
	ans = strings.Title(strings.ToLower(ans))
	if ans == "Yes" || ans == "No" {
		return ans, nil
	}
	return "", fmt.Errorf("unexpected output %q", out)
}

func buildTests() []testCase {
	tests := []testCase{
		formatCase("n1_same", [][]byte{{'X'}}, [][]byte{{'X'}}),
		formatCase("n1_diff", [][]byte{{'X'}}, [][]byte{{'O'}}),
		formatCase("simple_yes", [][]byte{
			[]byte("XO"),
			[]byte("OX"),
		}, [][]byte{
			[]byte("OX"),
			[]byte("XO"),
		}),
		formatCase("simple_no", [][]byte{
			[]byte("XX"),
			[]byte("OO"),
		}, [][]byte{
			[]byte("XO"),
			[]byte("OX"),
		}),
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomCase(rng, i))
	}
	tests = append(tests, stressCase())
	return tests
}

func formatCase(name string, a, b [][]byte) testCase {
	n := len(a)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteString(string(a[i]))
		sb.WriteByte('\n')
	}
	for i := 0; i < n; i++ {
		sb.WriteString(string(b[i]))
		sb.WriteByte('\n')
	}
	return testCase{name: name, input: sb.String()}
}

func randomCase(rng *rand.Rand, idx int) testCase {
	n := rng.Intn(10) + 1
	a := randomGrid(rng, n)
	var b [][]byte
	if rng.Intn(2) == 0 {
		b = applyTransform(a, rng.Intn(8))
	} else {
		b = randomGrid(rng, n)
	}
	return formatCase(fmt.Sprintf("random_%d", idx+1), a, b)
}

func stressCase() testCase {
	n := 10
	a := randomGrid(rand.New(rand.NewSource(99)), n)
	b := applyTransform(a, 3) // rotate 270
	return formatCase("stress_max", a, b)
}

func randomGrid(rng *rand.Rand, n int) [][]byte {
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		row := make([]byte, n)
		for j := 0; j < n; j++ {
			if rng.Intn(2) == 0 {
				row[j] = 'X'
			} else {
				row[j] = 'O'
			}
		}
		grid[i] = row
	}
	return grid
}

func applyTransform(mat [][]byte, op int) [][]byte {
	n := len(mat)
	res := make([][]byte, n)
	for i := range res {
		res[i] = make([]byte, n)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			switch op {
			case 0:
				res[i][j] = mat[i][j]
			case 1:
				res[i][j] = mat[n-1-j][i]
			case 2:
				res[i][j] = mat[n-1-i][n-1-j]
			case 3:
				res[i][j] = mat[j][n-1-i]
			case 4:
				res[i][j] = mat[i][n-1-j]
			case 5:
				res[i][j] = mat[n-1-i][j]
			case 6:
				res[i][j] = mat[j][i]
			case 7:
				res[i][j] = mat[n-1-j][n-1-i]
			}
		}
	}
	return res
}
