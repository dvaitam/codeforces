package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const refSource = "./736D.go"

type testCase struct {
	input string
	m     int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for i, tc := range tests {
		wantOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fail("reference runtime error on test %d: %v\ninput:\n%s", i+1, err, tc.input)
		}
		gotOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fail("candidate runtime error on test %d: %v\ninput:\n%s", i+1, err, tc.input)
		}
		want, err := parseAnswers(wantOut, tc.m)
		if err != nil {
			fail("reference output parse error on test %d: %v\noutput:\n%s", i+1, err, wantOut)
		}
		got, err := parseAnswers(gotOut, tc.m)
		if err != nil {
			fail("candidate output parse error on test %d: %v\noutput:\n%s", i+1, err, gotOut)
		}
		for idx := 0; idx < tc.m; idx++ {
			if want[idx] != got[idx] {
				fail("mismatch on test %d, pair %d: expected %s got %s\ninput:\n%s", i+1, idx+1, want[idx], got[idx], tc.input)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "736D-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("build reference failed: %v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", filepath.Clean(bin))
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

func parseAnswers(out string, m int) ([]string, error) {
	sc := bufio.NewScanner(strings.NewReader(out))
	sc.Split(bufio.ScanWords)
	sc.Buffer(make([]byte, 1024), 4<<20)
	var res []string
	for sc.Scan() {
		token := strings.ToUpper(sc.Text())
		if token == "" {
			continue
		}
		if token != "YES" && token != "NO" {
			return nil, fmt.Errorf("unexpected token %q", token)
		}
		res = append(res, token)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	if len(res) != m {
		return nil, fmt.Errorf("expected %d answers, got %d", m, len(res))
	}
	return res, nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	tests = append(tests, identityTest(1))
	tests = append(tests, manualSmall())
	tests = append(tests, randomMatrixTest(rng, 5, 0.5))
	tests = append(tests, randomMatrixTest(rng, 12, 0.35))
	tests = append(tests, randomMatrixTest(rng, 40, 0.4))
	tests = append(tests, randomMatrixTest(rng, 120, 0.25))
	tests = append(tests, randomMatrixTest(rng, 400, 0.18))
	tests = append(tests, randomMatrixTest(rng, 800, 0.12))
	tests = append(tests, randomMatrixTest(rng, 1200, 0.1))
	tests = append(tests, randomMatrixTest(rng, 2000, 0.08))
	tests = append(tests, randomMatrixTest(rng, 2000, 0.2))
	return tests
}

func identityTest(n int) testCase {
	mat := make([][]bool, n)
	for i := 0; i < n; i++ {
		row := make([]bool, n)
		row[i] = true
		mat[i] = row
	}
	return matrixToTest(mat)
}

func manualSmall() testCase {
	n := 3
	mat := make([][]bool, n)
	for i := 0; i < n; i++ {
		mat[i] = make([]bool, n)
		mat[i][i] = true
	}
	mat[0][1] = true
	mat[1][2] = true
	mat[2][0] = true
	return matrixToTest(mat)
}

func randomMatrixTest(rng *rand.Rand, n int, density float64) testCase {
	mat := buildInvertibleMatrix(rng, n, density)
	tc := matrixToTest(mat)
	if tc.m > 500000 {
		panic(fmt.Sprintf("generated test exceeds limit: n=%d m=%d", n, tc.m))
	}
	if tc.m < n {
		panic(fmt.Sprintf("invalid test: m < n (n=%d, m=%d)", n, tc.m))
	}
	return tc
}

func buildInvertibleMatrix(rng *rand.Rand, n int, density float64) [][]bool {
	base := make([][]bool, n)
	for i := 0; i < n; i++ {
		row := make([]bool, n)
		row[i] = true
		for j := i + 1; j < n; j++ {
			if rng.Float64() < density {
				row[j] = true
			}
		}
		base[i] = row
	}
	rows := rng.Perm(n)
	cols := rng.Perm(n)
	mat := make([][]bool, n)
	for i := 0; i < n; i++ {
		mat[i] = make([]bool, n)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if base[rows[i]][cols[j]] {
				mat[i][j] = true
			}
		}
	}
	return mat
}

func matrixToTest(mat [][]bool) testCase {
	n := len(mat)
	m := 0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if mat[i][j] {
				m++
			}
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if mat[i][j] {
				fmt.Fprintf(&sb, "%d %d\n", i+1, j+1)
			}
		}
	}
	return testCase{input: sb.String(), m: m}
}
