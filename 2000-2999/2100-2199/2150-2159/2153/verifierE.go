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

const testCount = 80

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "2153E.go")
	tmp, err := os.CreateTemp("", "oracle2153E")
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

func run(bin, input string) (string, error) {
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

func genCase(r *rand.Rand) string {
	t := r.Intn(4) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		n := r.Intn(5000) + 2
		maxM := n + r.Intn(5000) + 1
		if maxM > 1000000 {
			maxM = 1000000
		}
		m := r.Intn(maxM-n+1) + n
		fmt.Fprintf(&sb, "%d %d\n", n, m)
	}
	return sb.String()
}

func parseOutput(out string, t int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d outputs, got %d", t, len(fields))
	}
	res := make([]int64, t)
	for i := 0; i < t; i++ {
		val, err := strconv.ParseInt(fields[i], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer at test %d", i+1)
		}
		if val < 0 {
			return nil, fmt.Errorf("negative value at test %d", i+1)
		}
		res[i] = val
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
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
		input := genCase(r)
		expectStr, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\n", tcase+1, err)
			os.Exit(1)
		}
		gotStr, err := run(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", tcase+1, err)
			os.Exit(1)
		}
		reader := strings.NewReader(input)
		var t int
		fmt.Fscan(reader, &t)
		expectVals, err := parseOutput(expectStr, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid oracle output on test %d: %v\n", tcase+1, err)
			os.Exit(1)
		}
		gotVals, err := parseOutput(gotStr, t)
		if err != nil {
			fmt.Printf("test %d failed\ninput:\n%s\nerror: %v\n", tcase+1, input, err)
			os.Exit(1)
		}
		for i := 0; i < t; i++ {
			if expectVals[i] != gotVals[i] {
				fmt.Printf("test %d case %d failed\ninput:\n%s\nexpected: %d\ngot: %d\n", tcase+1, i+1, input, expectVals[i], gotVals[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", testCount)
}
