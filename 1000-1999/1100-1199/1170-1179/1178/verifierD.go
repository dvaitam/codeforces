package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const numTestsD = 100

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp := filepath.Join(os.TempDir(), "verifD_bin")
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

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	limit := int(math.Sqrt(float64(n))) + 1
	for i := 2; i < limit; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func nextPrime(n int) int {
	for !isPrime(n) {
		n++
	}
	return n
}

func solveD(n int) string {
	m := nextPrime(n)
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(m))
	sb.WriteByte('\n')
	for i := 1; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", i, i+1))
	}
	sb.WriteString(fmt.Sprintf("1 %d\n", n))
	extra := m - n
	half := n / 2
	for i := 1; i <= extra; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", i, i+half))
	}
	return strings.TrimSpace(sb.String())
}

func genCaseD(rng *rand.Rand) int {
	return rng.Intn(50) + 3
}

func runCaseD(bin string, n int) error {
	input := fmt.Sprintf("%d\n", n)
	expected := solveD(n)
	out, err := run(bin, input)
	if err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	if out != expected {
		return fmt.Errorf("expected:\n%s\n got:\n%s", expected, out)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin, cleanup, err := prepareBinary(os.Args[1])
	if err != nil {
		fmt.Println("compile error:", err)
		return
	}
	if cleanup != nil {
		defer cleanup()
	}
	rng := rand.New(rand.NewSource(4))
	for t := 0; t < numTestsD; t++ {
		n := genCaseD(rng)
		if err := runCaseD(bin, n); err != nil {
			fmt.Printf("case %d failed: %v\n", t+1, err)
			return
		}
	}
	fmt.Println("All tests passed")
}
