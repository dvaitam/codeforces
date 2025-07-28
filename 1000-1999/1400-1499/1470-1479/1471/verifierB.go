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
	"time"
)

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "1471B.go")
	bin := filepath.Join(os.TempDir(), "oracle1471B.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return bin, nil
}

func runProg(prog, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(prog, ".go") {
		cmd = exec.Command("go", "run", prog)
	} else {
		cmd = exec.Command(prog)
	}
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

func genCase(r *rand.Rand) string {
	t := r.Intn(4) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		n := r.Intn(20) + 1
		x := r.Int63n(5) + 2
		fmt.Fprintf(&sb, "%d %d\n", n, x)
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			val := r.Int63n(50) + 1
			sb.WriteString(strconv.FormatInt(val, 10))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func extraCases() []string {
	cases := []string{
		"1\n1 2\n1\n",
		fmt.Sprintf("1\n2 2\n4 8\n"),
		fmt.Sprintf("1\n3 3\n9 9 9\n"),
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
	tests := make([]string, 0, 105)
	for i := 0; i < 100; i++ {
		tests = append(tests, genCase(rng))
	}
	tests = append(tests, extraCases()...)

	for i, input := range tests {
		expect, err := runProg(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failure on case %d: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:\n%s\n got:\n%s\ninput:%s", i+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
