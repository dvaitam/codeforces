package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func compileRef() (string, error) {
	out := "ref_bin"
	cmd := exec.Command("g++", "solC.cpp", "-O2", "-std=c++17", "-o", out)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return out, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	ref, err := compileRef()
	if err != nil {
		fmt.Printf("failed to compile reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	file, err := os.Open("testcasesC.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		var x1, y1, r1, x2, y2, r2, x3, y3, r3 int
		fmt.Sscan(line, &x1, &y1, &r1, &x2, &y2, &r2, &x3, &y3, &r3)
		input := fmt.Sprintf("%d %d %d\n%d %d %d\n%d %d %d\n", x1, y1, r1, x2, y2, r2, x3, y3, r3)
		// run reference
		cmdRef := exec.Command("./" + ref)
		cmdRef.Stdin = bytes.NewBufferString(input)
		outRef, err := cmdRef.CombinedOutput()
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", idx, err)
			os.Exit(1)
		}
		expected := strings.TrimSpace(string(outRef))
		// run candidate
		cmd := exec.Command(binary)
		cmd.Stdin = bytes.NewBufferString(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != expected {
			fmt.Printf("Test %d failed:\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
