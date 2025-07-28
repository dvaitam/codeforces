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
	"strings"
	"time"
)

func buildBinary(src, tag string) (string, error) {
	if strings.HasSuffix(src, ".go") {
		out := filepath.Join(os.TempDir(), tag)
		cmd := exec.Command("go", "build", "-o", out, src)
		if outb, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build %s: %v\n%s", src, err, string(outb))
		}
		return out, nil
	}
	return src, nil
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func solveCase(arr []int) string {
	n := len(arr)
	if n == 1 {
		if arr[0] == 1 {
			return "YES"
		}
		return "NO"
	}
	sort.Ints(arr)
	if arr[n-1]-arr[n-2] <= 1 {
		return "YES"
	}
	return "NO"
}

func generateCase(r *rand.Rand) (string, string) {
	n := r.Intn(10) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = r.Intn(10) + 1
	}
	var b strings.Builder
	fmt.Fprintf(&b, "1\n%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", arr[i])
	}
	b.WriteByte('\n')
	input := b.String()
	expect := fmt.Sprintf("%s\n", solveCase(arr))
	return input, expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	candSrc := os.Args[1]
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	refSrc := filepath.Join(dir, "1660B.go")

	cand, err := buildBinary(candSrc, "candB.bin")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	ref, err := buildBinary(refSrc, "refB.bin")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(ref, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
			os.Exit(1)
		}
		if err := runCase(cand, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
