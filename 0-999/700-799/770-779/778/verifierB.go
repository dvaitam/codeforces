package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	exe := filepath.Join(dir, "oracleB")
	src := filepath.Join(dir, "778B.go")
	cmd := exec.Command("go", "build", "-o", exe, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle: %v\n%s", err, out)
	}
	return exe, nil
}

func randName(rng *rand.Rand) string {
	letters := "abcdefghijklmnopqrstuvwxyz"
	l := rng.Intn(3) + 1
	b := make([]byte, l)
	for i := range b {
		b[i] = letters[rng.Intn(len(letters))]
	}
	return string(b)
}

func randBits(rng *rand.Rand, m int) string {
	b := make([]byte, m)
	for i := 0; i < m; i++ {
		if rng.Intn(2) == 0 {
			b[i] = '0'
		} else {
			b[i] = '1'
		}
	}
	return string(b)
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(4) + 1
	m := rng.Intn(5) + 1
	names := make([]string, n)
	for i := range names {
		names[i] = randName(rng)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	ops := []string{"AND", "OR", "XOR"}
	for i := 0; i < n; i++ {
		if i == 0 || rng.Intn(2) == 0 {
			sb.WriteString(fmt.Sprintf("%s := %s\n", names[i], randBits(rng, m)))
		} else {
			op1 := "?"
			if rng.Intn(2) == 0 {
				op1 = names[rng.Intn(i)]
			}
			op2 := "?"
			if rng.Intn(2) == 0 {
				op2 = names[rng.Intn(i)]
			}
			op := ops[rng.Intn(3)]
			sb.WriteString(fmt.Sprintf("%s := %s %s %s\n", names[i], op1, op, op2))
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		input := genCase(rng)
		candOut, cErr := runBinary(candidate, input)
		refOut, rErr := runBinary(oracle, input)
		if cErr != nil {
			fmt.Fprintf(os.Stderr, "case %d: candidate error: %v\ninput:\n%s", t+1, cErr, input)
			os.Exit(1)
		}
		if rErr != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle error: %v\ninput:\n%s", t+1, rErr, input)
			os.Exit(1)
		}
		if candOut != refOut {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%s\nexpected:%s\nactual:%s\n", t+1, input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
