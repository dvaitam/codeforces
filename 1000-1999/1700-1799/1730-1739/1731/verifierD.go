package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func buildRef() (string, error) {
	out := filepath.Join(os.TempDir(), "refD_1731")
	cmd := exec.Command("go", "build", "-o", out, "1731D.go")
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return out, nil
}

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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	file, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to open testcasesD.txt:", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			fmt.Fprintf(os.Stderr, "bad testcase at line %d\n", idx+1)
			os.Exit(1)
		}
		n := 0
		m := 0
		fmt.Sscan(parts[0], &n)
		fmt.Sscan(parts[1], &m)
		if len(parts) != 2+n*m {
			fmt.Fprintf(os.Stderr, "bad matrix at line %d\n", idx+1)
			os.Exit(1)
		}
		input := fmt.Sprintf("1\n%d %d\n", n, m)
		for i := 0; i < n; i++ {
			row := strings.Join(parts[2+i*m:2+(i+1)*m], " ")
			input += row + "\n"
		}
		expect, err := runBinary(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := runBinary(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed:\ninput:\n%s\nexpected:%s\ngot:%s\n", idx+1, input, expect, got)
			os.Exit(1)
		}
		idx++
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
