package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

const testCount = 120

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "2005E1.go")
	tmp, err := os.CreateTemp("", "oracle2005E1")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()
	os.Remove(path)
	cmd := exec.Command("go", "build", "-o", path, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v\n%s", err, string(out))
	}
	return path, nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(r *rand.Rand) (string, int, int, int) {
	t := r.Intn(5) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	totalCells := 0
	totalLen := 0
	for i := 0; i < t; i++ {
		l := r.Intn(10) + 1
		if r.Intn(4) == 0 {
			l = r.Intn(20) + 1
		}
		n := r.Intn(15) + 1
		m := r.Intn(15) + 1
		totalCells += n * m
		totalLen += l
		fmt.Fprintf(&sb, "%d %d %d\n", l, n, m)
		for j := 0; j < l; j++ {
			val := r.Intn(min(7, n*m)) + 1
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", val)
		}
		sb.WriteByte('\n')
		for row := 0; row < n; row++ {
			for col := 0; col < m; col++ {
				val := r.Intn(min(7, n*m)) + 1
				if col > 0 {
					sb.WriteByte(' ')
				}
				fmt.Fprintf(&sb, "%d", val)
			}
			sb.WriteByte('\n')
		}
	}
	return sb.String(), t, totalCells, totalLen
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func parseOutput(out string, t int) ([]string, error) {
	lines := strings.Fields(out)
	if len(lines) != t {
		return nil, fmt.Errorf("expected %d outputs, got %d", t, len(lines))
	}
	for i, line := range lines {
		if line != "T" && line != "N" {
			return nil, fmt.Errorf("output %d invalid: %s", i+1, line)
		}
	}
	return lines, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE1.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	r := rand.New(rand.NewSource(1))
	for tcase := 0; tcase < testCount; tcase++ {
		input, t, _, _ := genCase(r)
		expectOut, err := runProgram(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\n", tcase+1, err)
			os.Exit(1)
		}
		gotOut, err := runProgram(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", tcase+1, err)
			os.Exit(1)
		}
		expectArr, err := parseOutput(expectOut, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle output invalid on test %d: %v\n", tcase+1, err)
			os.Exit(1)
		}
		gotArr, err := parseOutput(gotOut, t)
		if err != nil {
			fmt.Printf("test %d failed\ninput:\n%s\nerror: %v\n", tcase+1, input, err)
			os.Exit(1)
		}
		for i := 0; i < t; i++ {
			if expectArr[i] != gotArr[i] {
				fmt.Printf("test %d case %d failed\ninput:\n%s\nexpected: %s\ngot: %s\n", tcase+1, i+1, input, expectArr[i], gotArr[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", testCount)
}
