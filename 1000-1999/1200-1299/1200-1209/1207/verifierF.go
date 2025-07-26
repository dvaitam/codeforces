package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func buildRef() (string, func(), error) {
	tmp, err := os.CreateTemp("", "refF*")
	if err != nil {
		return "", nil, err
	}
	tmp.Close()
	if out, err := exec.Command("go", "build", "-o", tmp.Name(), "1207F.go").CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", nil, fmt.Errorf("ref build failed: %v\n%s", err, out)
	}
	return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
}

func run(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(bytes.TrimSpace(out)), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	cand, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()
	ref, rcleanup, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer rcleanup()

	file, err := os.Open("testcasesF.txt")
	if err != nil {
		fmt.Println("could not open testcasesF.txt:", err)
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
		idx++
		fields := strings.Fields(line)
		if len(fields) < 1 {
			fmt.Printf("test %d invalid\n", idx)
			os.Exit(1)
		}
		q := fields[0]
		var input strings.Builder
		input.WriteString(q)
		input.WriteByte('\n')
		rest := fields[1:]
		if len(rest)%3 != 0 {
			fmt.Printf("test %d bad query count\n", idx)
			os.Exit(1)
		}
		for i := 0; i < len(rest); i += 3 {
			input.WriteString(rest[i])
			input.WriteByte(' ')
			input.WriteString(rest[i+1])
			input.WriteByte(' ')
			input.WriteString(rest[i+2])
			input.WriteByte('\n')
		}
		expect, err := run(ref, []byte(input.String()))
		if err != nil {
			fmt.Printf("reference failed on test %d: %v\n", idx, err)
			os.Exit(1)
		}
		got, err := run(cand, []byte(input.String()))
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", idx, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:\n%s\nGot:\n%s\n", idx, input.String(), expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
