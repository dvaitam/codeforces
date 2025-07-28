package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func buildOracle() (string, error) {
	oracle := "oracleH"
	cmd := exec.Command("go", "build", "-o", oracle, "1482H.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, string(out))
	}
	return oracle, nil
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("run error: %v\nstderr: %s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func randWord() string {
	l := rand.Intn(4) + 1
	b := make([]byte, l)
	for i := 0; i < l; i++ {
		b[i] = byte('a' + rand.Intn(26))
	}
	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rand.Seed(8)
	const T = 100
	for i := 0; i < T; i++ {
		n := rand.Intn(4) + 1
		var input bytes.Buffer
		fmt.Fprintln(&input, n)
		for j := 0; j < n; j++ {
			fmt.Fprintln(&input, randWord())
		}
		inp := input.String()
		expected, err := run("./"+oracle, inp)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		got, err := run(bin, inp)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if expected != got {
			fmt.Printf("test %d mismatch\nexpected: %s\ngot: %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
