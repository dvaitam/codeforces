package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const numTestsC = 100

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp := filepath.Join(os.TempDir(), fmt.Sprintf("candC_%d", time.Now().UnixNano()))
		cmd := exec.Command("go", "build", "-o", tmp, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("compile candidate: %v: %s", err, out)
		}
		return tmp, func() { os.Remove(tmp) }, nil
	}
	return path, func() {}, nil
}

func prepareOracle() (string, func(), error) {
	tmp := filepath.Join(os.TempDir(), fmt.Sprintf("oracleC_%d", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", tmp, "1675C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", nil, fmt.Errorf("compile oracle: %v: %s", err, out)
	}
	return tmp, func() { os.Remove(tmp) }, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(60) + 1
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		r := rng.Intn(3)
		switch r {
		case 0:
			b[i] = '0'
		case 1:
			b[i] = '1'
		default:
			b[i] = '?'
		}
	}
	return fmt.Sprintf("1\n%s\n", string(b))
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	cand, cleanCand, err := prepareBinary(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer cleanCand()

	oracle, cleanOracle, err := prepareOracle()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer cleanOracle()

	rng := rand.New(rand.NewSource(3))
	for i := 0; i < numTestsC; i++ {
		input := genCase(rng)
		want, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := run(cand, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
