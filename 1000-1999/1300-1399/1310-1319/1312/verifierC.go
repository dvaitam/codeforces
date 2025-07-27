package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func buildOracle() (string, error) {
	out := "./oracleC.bin"
	cmd := exec.Command("go", "build", "-o", out, "1312C.go")
	if b, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v: %s", err, string(b))
	}
	return out, nil
}

func runProg(prog, input string) (string, error) {
	cmd := exec.Command(prog)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func lineToInput(line string) string {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return ""
	}
	n := fields[0]
	k := fields[1]
	rest := strings.Join(fields[2:], " ")
	return fmt.Sprintf("1\n%s %s\n%s\n", n, k, rest)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	f, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		input := lineToInput(line)
		exp, err := runProg(oracle, input)
		if err != nil {
			fmt.Printf("oracle error on test %d: %v\n", idx, err)
			fmt.Println(exp)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("candidate error on test %d: %v\n", idx, err)
			fmt.Println(got)
			os.Exit(1)
		}
		if exp != got {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%s\nGot:\n%s\n", idx, input, exp, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
