package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func buildExecutable(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "bin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		cmd := exec.Command("go", "build", "-o", tmp.Name(), path)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, stderr.String())
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func run(bin string, input string) (string, error) {
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

func oracle(exe string, test string) (string, error) {
	return run(exe, "1\n"+test)
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 2
	d := rng.Intn(100) + n + 1
	vals := rng.Perm(d)
	arr := vals[:n]
	for i := range arr {
		arr[i]++
	}
	sort.Ints(arr)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, d)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", arr[i])
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	binPath := os.Args[1]
	bin, cleanup, err := buildExecutable(binPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to prepare binary: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()
	oracleExe, oracleCleanup, err := buildExecutable("1650E.go")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build oracle: %v\n", err)
		os.Exit(1)
	}
	defer oracleCleanup()
	rng := rand.New(rand.NewSource(5))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		expected, err := oracle(oracleExe, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed: %v\n", err)
			os.Exit(1)
		}
		got, err := run(bin, "1\n"+tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
