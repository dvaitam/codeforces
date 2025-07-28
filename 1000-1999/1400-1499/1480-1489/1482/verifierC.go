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
	oracle := "oracleC"
	cmd := exec.Command("go", "build", "-o", oracle, "1482C.go")
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rand.Seed(3)
	const T = 100
	var input bytes.Buffer
	fmt.Fprintln(&input, T)
	for i := 0; i < T; i++ {
		n := rand.Intn(4) + 1
		m := rand.Intn(4) + 1
		fmt.Fprintf(&input, "%d %d\n", n, m)
		for j := 0; j < m; j++ {
			k := rand.Intn(n) + 1
			perm := rand.Perm(n)
			fmt.Fprintf(&input, "%d", k)
			for t := 0; t < k; t++ {
				fmt.Fprintf(&input, " %d", perm[t]+1)
			}
			input.WriteByte('\n')
		}
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
		fmt.Println("mismatch")
		fmt.Println("expected:\n" + expected)
		fmt.Println("got:\n" + got)
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
