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

const testCount = 200

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "883D.go")
	tmp, err := os.CreateTemp("", "oracle883D")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()
	os.Remove(path)
	cmd := exec.Command("go", "build", "-o", path, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v\n%s", err, out)
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
	n := 2 + r.Intn(40)
	bytes := make([]byte, n)
	hasP, hasStar := false, false
	for i := 0; i < n; i++ {
		switch r.Intn(3) {
		case 0:
			bytes[i] = '.'
		case 1:
			bytes[i] = 'P'
			hasP = true
		default:
			bytes[i] = '*'
			hasStar = true
		}
	}
	if !hasP {
		pos := r.Intn(n)
		bytes[pos] = 'P'
		hasP = true
	}
	if !hasStar {
		pos := r.Intn(n)
		bytes[pos] = '*'
		hasStar = true
	}
	return fmt.Sprintf("%d\n%s\n", n, string(bytes))
}

func parseOutput(out string) (int64, int64, error) {
	fields := strings.Fields(out)
	if len(fields) != 2 {
		return 0, 0, fmt.Errorf("expected 2 integers, got %d tokens", len(fields))
	}
	a, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	b, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	if a < 0 || b < 0 {
		return 0, 0, fmt.Errorf("negative output")
	}
	return a, b, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
	for t := 0; t < testCount; t++ {
		input := genCase(r)
		expectStr, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		gotStr, err := run(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		ea, eb, err := parseOutput(expectStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle output invalid on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		ga, gb, err := parseOutput(gotStr)
		if err != nil {
			fmt.Printf("test %d failed\ninput:%serror: %v\n", t+1, input, err)
			os.Exit(1)
		}
		if ea != ga || eb != gb {
			fmt.Printf("test %d failed\ninput:%sexpected: %d %d\ngot: %d %d\n", t+1, input, ea, eb, ga, gb)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", testCount)
}
