package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runProgram(bin string, input []byte) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout")
		}
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

func randBinary(k int, rng *rand.Rand) string {
	b := make([]byte, k)
	b[0] = '1'
	for i := 1; i < k; i++ {
		if rng.Intn(2) == 0 {
			b[i] = '0'
		} else {
			b[i] = '1'
		}
	}
	return string(b)
}

func buildTests() []byte {
	rng := rand.New(rand.NewSource(2056))
	type tc struct {
		k int
		m int64
		s string
	}

	var cases []tc

	cases = append(cases, tc{1, 1, "1"})
	cases = append(cases, tc{1, 5, "1"})
	cases = append(cases, tc{2, 7, "10"})
	cases = append(cases, tc{3, 6, "101"})

	cases = append(cases, tc{10, 123456789, randBinary(10, rng)})

	// Long strings to cover large k. Keep total length <= 2e5.
	long1 := strings.Repeat("10", 50000) // length 100000, starts with '1'
	cases = append(cases, tc{len(long1), 999_999_937, long1})

	long2 := "1" + strings.Repeat("0", 99999) // power of two length 100000
	cases = append(cases, tc{len(long2), 1_000_000_000, long2})

	var buf bytes.Buffer
	fmt.Fprintln(&buf, len(cases))
	for _, c := range cases {
		fmt.Fprintf(&buf, "%d %d\n", c.k, c.m)
		fmt.Fprintln(&buf, c.s)
	}
	return buf.Bytes()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF2.go /path/to/binary")
		return
	}

	target, cleanTarget, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanTarget()

	_, src, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(src)
	refPath := filepath.Join(baseDir, "2056F2.go")
	refBin, cleanRef, err := buildIfGo(refPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanRef()

	input := buildTests()

	expOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}
	gotOut, err := runProgram(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target failed: %v\n", err)
		os.Exit(1)
	}

	exp := strings.Fields(expOut)
	got := strings.Fields(gotOut)
	if len(exp) != len(got) {
		fmt.Fprintf(os.Stderr, "output length mismatch: expected %d lines, got %d\n", len(exp), len(got))
		os.Exit(1)
	}
	for i := range exp {
		if exp[i] != got[i] {
			fmt.Fprintf(os.Stderr, "mismatch at case %d: expected %s got %s\n", i+1, exp[i], got[i])
			os.Exit(1)
		}
	}

	fmt.Println("all tests passed")
}
