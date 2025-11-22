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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
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

func genCase(n int, rng *rand.Rand) []byte {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, n)
	for row := 0; row < 2; row++ {
		for i := 0; i < n; i++ {
			if i > 0 {
				buf.WriteByte(' ')
			}
			fmt.Fprint(&buf, rng.Intn(200_000)+1)
		}
		buf.WriteByte('\n')
	}
	return buf.Bytes()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	target, cleanupTarget, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanupTarget()

	_, src, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(src)
	refPath := filepath.Join(baseDir, "2045E.go")
	refBin, cleanupRef, err := buildIfGo(refPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanupRef()

	rng := rand.New(rand.NewSource(2045))
	tests := make([][]byte, 0, 12)

	// Small and edge-shaped cases.
	tests = append(tests, []byte("1\n1\n1\n"))
	tests = append(tests, []byte("2\n1 1\n1 1\n"))
	tests = append(tests, []byte("3\n8 4 5\n5 4 8\n")) // sample

	for _, n := range []int{4, 5, 10, 50, 200, 1000, 5000, 15000} {
		tests = append(tests, genCase(n, rng))
	}

	for idx, tc := range tests {
		expOut, err := runProgram(refBin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotOut, err := runProgram(target, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		exp := strings.TrimSpace(expOut)
		got := strings.TrimSpace(gotOut)
		if exp != got {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected: %s\ngot     : %s\n", idx+1, exp, got)
			os.Exit(1)
		}
	}

	fmt.Println("all tests passed")
}
