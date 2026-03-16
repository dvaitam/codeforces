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

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleC")
	cmd := exec.Command("go", "build", "-o", oracle, "875C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
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

	const testcasesRaw = `2 6 3 6 3 2 3 6 4 3
3 4 1 3 3 1 2 4 2 2 1
3 4 4 2 2 2 2 4 3 2 4 1 5 1 3 4 4 3
2 3 5 3 1 3 2 3 4 3 2 3 1
4 5 5 4 2 3 3 4 1 2 1 3 2 3 4
3 5 3 2 3 1 3 5 5 1 2 3 1
5 7 1 1 2 1 1 2 6 3 1 7 1 3
5 5 2 4 4 2 3 3 2 3 4 4 1 4 3 5 5 4 1 5 1 4
5 5 1 5 2 5 5 2 1 3 2 2 2 1 2
3 3 4 3 1 3 3 4 3 1 3 3 1 2
4 7 1 6 1 4 1 3 3 6 2 4
3 7 3 2 6 7 4 3 3 7 4 4 1 3 5 3
5 5 5 5 5 3 1 4 4 1 4 3 4 1 1 3 1 3 5 3 2 5 5
4 7 3 2 1 5 3 7 2 5 5 6 3 2 7 2 3 7 6 7
2 6 1 5 2 3 3
4 6 3 4 4 5 4 2 1 2 5 1 4 2 3 1
5 7 3 6 5 2 1 5 4 3 7 2 5 2 1 6 2 5 1
5 6 3 3 3 1 4 6 3 3 5 5 5 3 3 2 6 4 2 1 5 2 5 4 3 4 1 5
5 7 5 7 7 7 2 1 3 5 2 6 2 6 3 4 1 6 6 2 5 2 3 2 7 4
2 6 4 2 6 6 4 5 6 1 2 2 4
2 5 4 3 3 4 3 5 3 1 3 1 4
2 4 2 4 4 4 4 1 4 3
3 7 3 3 7 1 4 4 5 2 1 1 5
4 6 1 1 2 6 1 5 4 1 3 1 3 4 2 6 3 4
3 7 2 4 7 3 5 1 2 2 2 4
2 7 5 1 7 5 6 6 4 2 3 7 5
2 2 5 1 2 1 2 2 3 2 2 1
5 7 2 1 6 5 2 6 5 1 7 3 7 1 2 4 5 2 7 4 5 2 6 3 6 2
5 7 1 5 1 7 5 3 1 7 7 1 1 4 5 3 5 4 3 7
5 7 3 7 7 6 1 2 3 1 2 6 2 1 3 5 2 1 4 6 1
4 7 1 5 2 3 7 1 4 1 4
4 6 3 5 4 4 2 1 4 3 2 4 1 3 1 1 1
4 4 2 4 1 3 4 4 4 1 1 5 4 3 1 1 1
4 6 1 4 1 2 5 6 3 1 1 4 3 2 5 6
3 4 2 2 2 5 3 1 4 4 1 2 2 3
4 5 2 3 2 2 4 4 3 5 2 4 5 1 2 5 1 4
3 4 1 1 3 1 1 1 1 2
3 6 1 4 4 5 6 3 6 2 2 3
3 7 4 1 2 4 5 3 1 4 4 2 4 4
3 6 2 6 4 4 5 1 3 3 5 3 5 1 5 5
5 5 3 1 5 3 4 1 5 1 2 4 1 4 4 4 1 4 5 4 2 2 3 4
5 5 5 3 5 1 3 3 2 3 4 5 1 2 3 2 5 3 3 4 1 3 5 4 1
5 6 4 6 2 3 2 2 1 5 3 2 1 4 5 4 4 1 6 6 1 4
3 4 2 2 2 2 1 2 4 3 3 3 1
5 6 2 3 6 5 3 5 6 2 4 5 3 2 4 3 4 4 5 3 4 4 1 6
3 4 1 1 2 1 3 5 2 4 1 1 1
2 4 4 3 1 4 3 5 3 4 2 2 1
3 3 2 3 2 5 3 1 2 3 3 4 2 1 3 1
5 5 4 5 1 4 1 5 4 3 2 5 4 1 4 5 1 3 5 4 5 4 4 3 2 2
4 5 3 5 2 4 1 1 4 4 5 4 5 1 3
5 5 1 3 2 3 1 2 1 2 5 5 2 1 1 4 5 5 4 1 4 5
4 6 1 6 1 2 2 4 3 3 2 1 4
4 6 1 1 3 5 4 2 3 4 2 2 3 2 6 4
4 6 4 5 2 1 4 3 5 4 2 3 6 3 4 3 3 4 6
4 4 4 3 1 4 2 3 2 2 3 4 2 1 4 4 4 4 2 4 3
5 5 3 5 1 1 4 3 4 2 4 1 5 1 4 1 3
3 5 2 1 2 3 1 2 1 1 4
2 3 1 3 5 1 1 3 2 1
2 6 5 4 1 5 2 2 2 4 4
5 5 4 2 4 1 5 3 1 5 3 3 3 4 2 5 5 1 3 1 1 3 1 2 5
3 6 2 5 6 3 2 4 5 5 1 5 1 6 4
2 6 4 4 2 4 5 3 2 4 3
5 5 5 5 3 2 4 2 3 4 4 1 4 1 5 4 4 2 4 2 3 4 4 2
3 6 3 3 6 4 5 2 5 6 1 2 3 5 6 6
2 2 2 2 1 3 1 2 2
3 6 4 2 2 4 1 3 6 5 2 3 2 5 5
5 5 5 4 5 2 5 4 1 2 4 5 3 5 2 1 1 1 2
3 6 5 2 6 6 1 5 4 6 6 3 6 3 4 6 2
3 5 5 3 4 4 4 3 4 1 2 2 1 2 4 4
3 5 1 1 4 2 1 4 5 5 4 5 2 1 2
5 6 4 5 5 1 1 5 6 4 1 2 4 5 3 3 1 6 1 2 5 1 5 5 3 2 4 3
2 6 4 3 2 4 4 5 3 2 5 3 1
2 3 3 3 3 3 4 1 3 2 2
4 7 1 1 4 5 6 2 7 2 1 1 4 2 1 6 3
4 7 1 3 3 3 5 7 2 6 7 5 5 1 3 5 5
2 3 4 1 3 2 2 2 2 1
5 5 3 5 1 5 1 5 2 4 1 5 2 3 1 5 2 2 2 4
3 7 4 3 2 3 6 3 5 7 2 2 4 7
4 5 3 5 1 1 3 3 3 1 1 5 1 3
2 3 3 1 3 3 3 2 1 2
5 5 1 5 4 3 2 3 2 4 3 1 1 1 3 1 3 3 5 3 1 2 5 2
4 4 4 4 3 2 3 3 4 4 1 4 3 2 1 2 4 4 1 2 1
4 4 2 2 4 1 3 5 4 1 2 4 3 3 4 4 4
2 3 5 2 1 1 1 1 3 1 1 3
4 4 2 2 4 2 2 2 5 2 2 3 1 4 2 1 4
3 6 5 4 2 4 5 5 4 1 2 1 5 1 1
2 7 2 6 4 4 6 3 4 3
5 6 1 2 5 2 6 2 1 1 5 1 6 4 2 4 2 5 6 5 3 1 5 2 3
4 6 3 5 6 6 3 2 1 4 2 5 4 3 4 3 4
3 6 3 6 3 3 5 4 4 5 4 6 5 3 6 5 2 2
5 5 3 5 2 5 2 5 5 5 5 4 3 3 4 3 4 5 1 2 3 4
2 2 2 1 2 1 1
3 7 2 6 3 5 7 3 1 6 6 1 1
3 5 2 4 4 3 3 1 5 4 1 1 2 1
4 6 2 4 2 3 2 2 4 2 5 5 5 4 5 6 2 1
2 5 5 1 2 5 1 4 3 3 5 2
2 4 2 3 3 4 4 3 3 1
5 5 1 5 3 5 4 5 4 4 2 5 4 2 1 5 1 2
5 6 1 6 4 3 3 1 6 3 5 5 1 1 1 5 1 2 1 4 4
5 5 3 2 1 4 5 2 1 2 4 1 1 5 4 2 5 3 4 4 5 1 4 5`

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		input := line + "\n"
		cmdO := exec.Command(oracle)
		cmdO.Stdin = strings.NewReader(input)
		var outO bytes.Buffer
		cmdO.Stdout = &outO
		if err := cmdO.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "oracle run error: %v\n", err)
			os.Exit(1)
		}
		expected := strings.TrimSpace(outO.String())

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expected {
			fmt.Printf("test %d failed\nexpected: %s\n got: %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
