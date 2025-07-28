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

type testCase struct {
	input string
}

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleC")
	cmd := exec.Command("go", "build", "-o", oracle, "1478C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func genValidCase(rng *rand.Rand) testCase {
	n := rng.Intn(4) + 1
	xs := make([]int64, n)
	seen := make(map[int64]bool)
	for i := 0; i < n; i++ {
		for {
			v := rng.Int63n(1000) + 1
			if !seen[v] {
				xs[i] = v
				seen[v] = true
				break
			}
		}
	}
	m := 2 * n
	a := make([]int64, m)
	for i := 0; i < n; i++ {
		a[i] = xs[i]
		a[i+n] = -xs[i]
	}
	d := make([]int64, m)
	for i := 0; i < m; i++ {
		var sum int64
		for j := 0; j < m; j++ {
			sum += abs(a[i] - a[j])
		}
		d[i] = sum
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range d {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	return testCase{sb.String()}
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func genRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(4) + 1
	m := 2 * n
	d := make([]int64, m)
	for i := 0; i < m; i++ {
		d[i] = rng.Int63n(2000) + 1
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range d {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	return testCase{sb.String()}
}

func runExe(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{genValidCase(rng), genRandomCase(rng)}
	for len(cases) < 102 {
		if rng.Intn(2) == 0 {
			cases = append(cases, genValidCase(rng))
		} else {
			cases = append(cases, genRandomCase(rng))
		}
	}

	for i, tc := range cases {
		exp, err := runExe(oracle, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
