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
	src := filepath.Join(dir, "1357A2.go")
	tmp, err := os.CreateTemp("", "oracle1357A2")
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

func genInput(r *rand.Rand) string {
	idx := r.Intn(4)
	return fmt.Sprintf("%d\n", idx)
}

func parseOutput(out string) (int, error) {
	fields := strings.Fields(out)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected single number, got %d tokens", len(fields))
	}
	val, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, err
	}
	if val < 0 || val > 3 {
		return 0, fmt.Errorf("output %d out of range", val)
	}
	return val, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA2.go /path/to/binary")
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
		input := genInput(r)
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
		expectVal, err := parseOutput(expectStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid oracle output on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		gotVal, err := parseOutput(gotStr)
		if err != nil {
			fmt.Printf("test %d failed\ninput:%sexpected: %d\nerror: %v\n", t+1, input, expectVal, err)
			os.Exit(1)
		}
		if expectVal != gotVal {
			fmt.Printf("test %d failed\ninput:%sexpected: %d\ngot: %d\n", t+1, input, expectVal, gotVal)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", testCount)
}
