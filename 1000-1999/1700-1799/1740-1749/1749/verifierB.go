package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
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

func oracle(input string) string {
	r := strings.NewReader(input)
	var n int
	fmt.Fscan(r, &n)
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(r, &a[i])
	}
	b := make([]int64, n)
	var sumB, maxB int64
	for i := 0; i < n; i++ {
		fmt.Fscan(r, &b[i])
		sumB += b[i]
		if b[i] > maxB {
			maxB = b[i]
		}
	}
	var sumA int64
	for _, v := range a {
		sumA += v
	}
	ans := sumA + sumB - maxB
	return fmt.Sprint(ans)
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", rng.Intn(1000)+1)
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", rng.Intn(1000)+1)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	binPath := os.Args[1]
	bin, cleanup, err := buildExecutable(binPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to prepare binary: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	rng := rand.New(rand.NewSource(43))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		expect := oracle(tc)
		got, err := run(bin, "1\n"+tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
