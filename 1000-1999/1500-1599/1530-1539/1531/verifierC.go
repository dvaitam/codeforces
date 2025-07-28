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
	exe := filepath.Join(dir, "oracleC")
	src := filepath.Join(dir, "1531C.go")
	cmd := exec.Command("go", "build", "-o", exe, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle: %v\n%s", err, out)
	}
	return exe, nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(100) + 1
	return fmt.Sprintf("%d\n", n)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		input := genCase(rng)
		candOut, cErr := runBinary(cand, input)
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
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", t+1, input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
