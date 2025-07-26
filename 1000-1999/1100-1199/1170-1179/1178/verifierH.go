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

const numTestsH = 100

func buildOracle() (string, error) {
	exe := filepath.Join(os.TempDir(), "oracleH")
	cmd := exec.Command("go", "build", "-o", exe, "1178H.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v: %s", err, out)
	}
	return exe, nil
}

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp := filepath.Join(os.TempDir(), "verifH_bin")
		cmd := exec.Command("go", "build", "-o", tmp, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, out)
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
	cmd.Stderr = &buf
	err := cmd.Run()
	return strings.TrimSpace(buf.String()), err
}

func genCaseH(rng *rand.Rand) int {
	return rng.Intn(3) + 1
}

func runCaseH(bin, oracle string, n int, a, b []int64) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < 2*n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", a[i], b[i]))
	}
	input := sb.String()
	exp, err := run(oracle, input)
	if err != nil {
		return fmt.Errorf("oracle error: %v", err)
	}
	got, err := run(bin, input)
	if err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		return
	}
	oracle, err := buildOracle()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer os.Remove(oracle)
	bin, cleanup, err := prepareBinary(os.Args[1])
	if err != nil {
		fmt.Println("compile error:", err)
		return
	}
	if cleanup != nil {
		defer cleanup()
	}
	rng := rand.New(rand.NewSource(9))
	for t := 0; t < numTestsH; t++ {
		n := genCaseH(rng)
		a := make([]int64, 2*n)
		b := make([]int64, 2*n)
		for i := 0; i < 2*n; i++ {
			a[i] = int64(rng.Intn(5))
			b[i] = int64(rng.Intn(5))
		}
		if err := runCaseH(bin, oracle, n, a, b); err != nil {
			fmt.Printf("case %d failed: %v\n", t+1, err)
			return
		}
	}
	fmt.Println("All tests passed")
}
