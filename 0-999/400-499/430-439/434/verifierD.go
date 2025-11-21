package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleD")
	cmd := exec.Command("go", "build", "-o", oracle, "434D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v\n%s", err, string(out))
	}
	return oracle, nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func buildInput(n, m int, coeffs [][3]int, bounds [][2]int, cons [][3]int) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d %d %d\n", coeffs[i][0], coeffs[i][1], coeffs[i][2])
	}
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d %d\n", bounds[i][0], bounds[i][1])
	}
	for i := 0; i < m; i++ {
		fmt.Fprintf(&sb, "%d %d %d\n", cons[i][0], cons[i][1], cons[i][2])
	}
	return sb.String()
}

func parseOutput(out string) (int64, error) {
	fields := strings.Fields(out)
	if len(fields) != 1 {
		if len(fields) == 0 {
			return 0, fmt.Errorf("empty output")
		}
		return 0, fmt.Errorf("expected single integer, got %d", len(fields))
	}
	val, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q", fields[0])
	}
	return val, nil
}

func deterministicTests() []string {
	var tests []string
	// single generator, no constraints
	tests = append(tests, buildInput(1, 0,
		[][3]int{{1, 0, 0}},
		[][2]int{{0, 3}},
		[][3]int{},
	))
	// two generators, simple constraint
	tests = append(tests, buildInput(2, 1,
		[][3]int{{1, 2, 3}, {-1, 0, 0}},
		[][2]int{{-2, 2}, {0, 4}},
		[][3]int{{1, 2, 1}},
	))
	// chain constraints
	tests = append(tests, buildInput(3, 3,
		[][3]int{{0, 1, 0}, {0, 2, 0}, {0, 3, 0}},
		[][2]int{{0, 5}, {0, 5}, {0, 5}},
		[][3]int{{1, 2, 1}, {2, 3, 1}, {3, 1, 0}},
	))
	// negative ranges
	tests = append(tests, buildInput(2, 2,
		[][3]int{{-1, 0, 0}, {-1, 0, 0}},
		[][2]int{{-5, -1}, {-4, -2}},
		[][3]int{{1, 2, 3}, {2, 1, -1}},
	))
	// zero width ranges
	tests = append(tests, buildInput(2, 0,
		[][3]int{{2, 0, 0}, {3, 0, 0}},
		[][2]int{{1, 1}, {2, 2}},
		[][3]int{},
	))
	return tests
}

func randInt(rnd *rand.Rand, l, r int) int {
	return rnd.Intn(r-l+1) + l
}

func randomTests(count int) []string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]string, 0, count)
	for i := 0; i < count; i++ {
		n := randInt(rnd, 1, 50)
		m := randInt(rnd, 0, 100)
		coeffs := make([][3]int, n)
		bounds := make([][2]int, n)
		for j := 0; j < n; j++ {
			coeffs[j][0] = randInt(rnd, -10, 10)
			coeffs[j][1] = randInt(rnd, -1000, 1000)
			coeffs[j][2] = randInt(rnd, -1000, 1000)
			l := randInt(rnd, -100, 100)
			r := randInt(rnd, l, 100)
			bounds[j][0] = l
			bounds[j][1] = r
		}
		values := make([]int, n)
		for j := 0; j < n; j++ {
			values[j] = randInt(rnd, bounds[j][0], bounds[j][1])
		}
		constraints := make([][3]int, m)
		for j := 0; j < m; j++ {
			u := randInt(rnd, 1, n)
			v := randInt(rnd, 1, n)
			for v == u {
				v = randInt(rnd, 1, n)
			}
			diff := values[u-1] - values[v-1]
			maxAdd := 200 - diff
			add := 0
			if maxAdd > 0 {
				add = rnd.Intn(maxAdd + 1)
			}
			d := diff + add
			if d < -200 {
				d = -200
			}
			if d > 200 {
				d = 200
			}
			constraints[j] = [3]int{u, v, d}
		}
		tests = append(tests, buildInput(n, m, coeffs, bounds, constraints))
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	tests := deterministicTests()
	tests = append(tests, randomTests(400)...)

	for idx, input := range tests {
		expOut, err := runProgram(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle error: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotOut, err := runProgram(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		expVal, err := parseOutput(expOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle output invalid: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotVal, err := parseOutput(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if expVal != gotVal {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx+1, expVal, gotVal)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
