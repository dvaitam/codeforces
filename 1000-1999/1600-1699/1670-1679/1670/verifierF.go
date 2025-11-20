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
)

func run(cmdPath string, input string) (string, error) {
	cmd := exec.Command(cmdPath)
	if strings.HasSuffix(cmdPath, ".go") {
		cmd = exec.Command("go", "run", cmdPath)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	path := filepath.Join(dir, "oracle1670F")
	cmd := exec.Command("go", "build", "-o", path, "1670F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	return path, nil
}

func genTests() []string {
	rand.Seed(6)
	tests := make([]string, 0, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(5) + 1
		l := rand.Int63n(50) + 1
		r := l + rand.Int63n(50)
		z := rand.Int63n(64)
		// Single test format: n l r z\n
		tests = append(tests, fmt.Sprintf("%d %d %d %d\n", n, l, r, z))
	}
	// Edge cases
	tests = append(tests, "1 1 1 0\n")
	tests = append(tests, "2 1 3 1\n")
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	tests := genTests()
	for i, tc := range tests {
		expStr, err := run(oracle, tc)
		if err != nil {
			fmt.Printf("oracle error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		expFields := strings.Fields(expStr)
		if len(expFields) == 0 {
			fmt.Printf("oracle produced empty output on test %d\ninput:\n%s", i+1, tc)
			os.Exit(1)
		}
		expVal, err := strconv.ParseInt(expFields[0], 10, 64)
		if err != nil {
			fmt.Printf("oracle output parse error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		gotStr, err := run(bin, tc)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		gotStr = strings.TrimSpace(gotStr)
		gotFields := strings.Fields(gotStr)
		if len(gotFields) == 0 {
			fmt.Printf("test %d failed\ninput:\n%sexpected:%d\ngot:%q\n", i+1, tc, expVal, gotStr)
			os.Exit(1)
		}
		gotVal, err := strconv.ParseInt(gotFields[0], 10, 64)
		if err != nil {
			fmt.Printf("test %d failed\ninput:\n%sexpected:%d\ngot:%q\n", i+1, tc, expVal, gotStr)
			os.Exit(1)
		}
		if gotVal != expVal {
			fmt.Printf("test %d failed\ninput:\n%sexpected:%d\ngot:%d\n", i+1, tc, expVal, gotVal)
			os.Exit(1)
		}
	}
	fmt.Println("ok")
}
