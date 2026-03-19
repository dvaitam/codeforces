package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func buildReference() (string, error) {
	srcPath := os.Getenv("REFERENCE_SOURCE_PATH")
	if srcPath == "" {
		_, file, _, ok := runtime.Caller(0)
		if !ok {
			return "", fmt.Errorf("cannot determine verifier directory and REFERENCE_SOURCE_PATH not set")
		}
		srcPath = filepath.Join(filepath.Dir(file), "488B.go")
	}
	content, err := os.ReadFile(srcPath)
	if err != nil {
		return "", fmt.Errorf("read reference source: %v", err)
	}
	tmp, err := os.CreateTemp("", "488B-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	if strings.Contains(string(content), "#include") {
		cppPath := tmp.Name() + ".cpp"
		if err := os.WriteFile(cppPath, content, 0644); err != nil {
			os.Remove(tmp.Name())
			return "", fmt.Errorf("write cpp: %v", err)
		}
		defer os.Remove(cppPath)
		cmd := exec.Command("g++", "-O2", "-o", tmp.Name(), cppPath)
		if out, err := cmd.CombinedOutput(); err != nil {
			os.Remove(tmp.Name())
			return "", fmt.Errorf("%v\n%s", err, string(out))
		}
	} else {
		cmd := exec.Command("go", "build", "-o", tmp.Name(), srcPath)
		if out, err := cmd.CombinedOutput(); err != nil {
			os.Remove(tmp.Name())
			return "", fmt.Errorf("%v\n%s", err, string(out))
		}
	}
	return tmp.Name(), nil
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(5)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		v := rng.Intn(500) + 1
		sb.WriteString(fmt.Sprintf("%d\n", v))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)
	var cases []string
	// simple deterministic cases
	cases = append(cases, "0\n")
	cases = append(cases, "4\n1\n1\n3\n3\n")
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}
	for i, in := range cases {
		exp, err := run(refBin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if err := verify(in, out, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d wrong answer: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func chk(b []int) bool {
	dif := (b[3] - b[0]) * 4
	mid := (b[1] + b[2]) * 2
	if mid != dif {
		return false
	}
	sum := 0
	for _, v := range b {
		sum += v
	}
	return sum == dif
}

func verify(input, output, expected string) error {
	tokensIn := strings.Fields(input)
	if len(tokensIn) == 0 {
		return fmt.Errorf("empty input")
	}
	n, err := strconv.Atoi(tokensIn[0])
	if err != nil {
		return fmt.Errorf("bad input: %v", err)
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i], err = strconv.Atoi(tokensIn[i+1])
		if err != nil {
			return fmt.Errorf("bad input number: %v", err)
		}
	}

	tokensOut := strings.Fields(output)
	if len(tokensOut) == 0 {
		return fmt.Errorf("empty output")
	}

	tokensExp := strings.Fields(expected)
	if len(tokensExp) == 0 {
		return fmt.Errorf("empty reference output")
	}

	if tokensExp[0] == "NO" {
		if tokensOut[0] != "NO" {
			return fmt.Errorf("expected NO but got %s", tokensOut[0])
		}
		if len(tokensOut) > 1 {
			return fmt.Errorf("unexpected extra output after NO")
		}
		return nil
	}

	if tokensOut[0] != "YES" {
		return fmt.Errorf("expected YES but got %s", tokensOut[0])
	}

	need := 4 - n
	if len(tokensOut)-1 != need {
		return fmt.Errorf("expected %d numbers after YES, got %d", need, len(tokensOut)-1)
	}

	b := make([]int, 0, 4)
	b = append(b, a...)
	for i := 1; i < len(tokensOut); i++ {
		v, err := strconv.Atoi(tokensOut[i])
		if err != nil {
			return fmt.Errorf("invalid number %q", tokensOut[i])
		}
		b = append(b, v)
	}
	if len(b) != 4 {
		return fmt.Errorf("total numbers != 4")
	}
	sort.Ints(b)
	if !chk(b) {
		return fmt.Errorf("output numbers invalid: %v", b)
	}
	return nil
}
