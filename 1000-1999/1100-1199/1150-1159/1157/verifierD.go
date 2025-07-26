package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const numTestsD = 100

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: verifierD.go <binary>")
		os.Exit(1)
	}
	binPath, cleanup, err := prepareBinary(os.Args[1])
	if err != nil {
		fmt.Println("compile error:", err)
		os.Exit(1)
	}
	if cleanup != nil {
		defer cleanup()
	}
	r := rand.New(rand.NewSource(1))
	for t := 1; t <= numTestsD; t++ {
		k := int64(r.Intn(10) + 1)
		min := k * (k + 1) / 2
		n := min + int64(r.Intn(100))
		input := fmt.Sprintf("%d %d\n", n, k)
		expected := solveD(n, k)
		out, err := run(binPath, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", t, err)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if out != expected {
			fmt.Printf("test %d failed\ninput:%sexpected:%s got:%s\n", t, input, expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("OK")
}

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp := filepath.Join(os.TempDir(), "verify_binD")
		cmd := exec.Command("go", "build", "-o", tmp, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, string(out))
		}
		return tmp, func() { os.Remove(tmp) }, nil
	}
	return path, nil, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return buf.String(), err
}

func solveD(n, k int64) string {
	sig := k * (k + 1) / 2
	if k == 1 {
		return fmt.Sprintf("YES\n%d", n)
	}
	if n < sig {
		return "NO"
	}
	n -= sig
	q := n / k
	r := n % k
	if q > 0 || (q == 0 && r != k-1) {
		var sb strings.Builder
		sb.WriteString("YES\n")
		for i := int64(1); i < k; i++ {
			sb.WriteString(fmt.Sprintf("%d ", i+q))
		}
		sb.WriteString(fmt.Sprintf("%d", k+q+r))
		return sb.String()
	}
	if k >= 4 {
		var sb strings.Builder
		sb.WriteString("YES\n")
		for i := int64(1); i <= k-2; i++ {
			sb.WriteString(fmt.Sprintf("%d ", i+q))
		}
		sb.WriteString(fmt.Sprintf("%d %d", k+q, k+q+r-1))
		return sb.String()
	}
	return "NO"
}
