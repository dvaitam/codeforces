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
	out := "./oracleG.bin"
	cmd := exec.Command("go", "build", "-o", out, "1312G.go")
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
	if len(fields) < 1 {
		return ""
	}
	idx := 0
	n := fields[idx]
	idx++
	var sb strings.Builder
	sb.WriteString(n)
	sb.WriteByte('\n')
	for i := 0; i < toInt(n); i++ {
		p := fields[idx]
		c := fields[idx+1]
		idx += 2
		sb.WriteString(p)
		sb.WriteByte(' ')
		sb.WriteString(c)
		sb.WriteByte('\n')
	}
	k := fields[idx]
	idx++
	sb.WriteString(k)
	sb.WriteByte('\n')
	for i := 0; i < toInt(k); i++ {
		sb.WriteString(fields[idx+i])
		if i+1 < toInt(k) {
			sb.WriteByte(' ')
		}
	}
	sb.WriteByte('\n')
	return sb.String()
}

func toInt(s string) int {
	var x int
	fmt.Sscan(s, &x)
	return x
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	f, err := os.Open("testcasesG.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idxCase := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idxCase++
		input := lineToInput(line)
		exp, err := runProg(oracle, input)
		if err != nil {
			fmt.Printf("oracle error on test %d: %v\n", idxCase, err)
			fmt.Println(exp)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("candidate error on test %d: %v\n", idxCase, err)
			fmt.Println(got)
			os.Exit(1)
		}
		if exp != got {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%s\nGot:\n%s\n", idxCase, input, exp, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idxCase)
}
