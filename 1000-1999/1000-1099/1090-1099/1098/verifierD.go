package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func buildRef() string {
	exe := "refD_bin"
	if _, err := os.Stat(exe); err == nil {
		return "./" + exe
	}
	cmd := exec.Command("g++", "-std=c++17", "solD.cpp", "-O2", "-o", exe)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	return "./" + exe
}

func run(path string, input []byte) ([]byte, error) {
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return out, err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref := buildRef()
	rand.Seed(4)

	for t := 1; t <= 100; t++ {
		q := rand.Intn(20) + 1
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "%d\n", q)
		weights := []int{}
		for i := 0; i < q; i++ {
			if len(weights) == 0 || rand.Intn(3) > 0 {
				// add
				x := rand.Intn(20) + 1
				weights = append(weights, x)
				fmt.Fprintf(&buf, "+ %d\n", x)
			} else {
				idx := rand.Intn(len(weights))
				x := weights[idx]
				weights = append(weights[:idx], weights[idx+1:]...)
				fmt.Fprintf(&buf, "- %d\n", x)
			}
		}
		input := buf.Bytes()
		exp, err := run(ref, input)
		if err != nil {
			fmt.Println("reference run error:", err)
			os.Exit(1)
		}
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", t, err)
			os.Exit(1)
		}
		if strings.TrimSpace(string(exp)) != strings.TrimSpace(string(out)) {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:\n%sGot:\n%s", t, string(input), string(exp), string(out))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
