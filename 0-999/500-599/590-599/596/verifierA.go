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
	oracle := filepath.Join(dir, "oracleA")
	cmd := exec.Command("go", "build", "-o", oracle, "596A.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

const testcasesARaw = `4 3 6 3 -2 -9 -2 -9 6
3 -6 -2 -6 9 -7 -2
4 7 1 -7 1 7 3 -7 3
1 -8 10
1 -4 -6
4 -7 7 -1 -1 -7 -1 -1 7
1 2 0
2 9 -2 -9 5
1 2 -2
4 10 -8 10 1 5 -8 5 1
2 -2 -7 -2 -3
1 -3 -9
1 2 1
2 5 -4 -7 -4
1 -3 10
4 -9 2 9 2 -9 -7 9 -7
1 -5 6
1 5 6
3 -6 4 7 -10 -6 -10
4 1 1 9 -1 9 1 1 -1
1 -4 0
4 2 3 -9 8 -9 3 2 8
4 9 -9 7 -10 9 -10 7 -9
1 -10 3
1 9 -7
3 -7 2 -5 2 -7 5
2 10 1 6 1
1 -2 0
4 3 1 10 1 3 7 10 7
3 0 0 1 -8 0 -8
3 7 -6 7 -1 2 -6
1 2 0
4 -7 5 -7 -7 0 -7 0 5
4 10 -8 -8 -4 10 -4 -8 -8
4 5 3 5 -4 8 -4 8 3
2 -7 -8 -7 -10
4 -2 10 -4 10 -4 -9 -2 -9
3 7 9 -6 -7 -6 9
3 5 -4 10 7 10 -4
2 -2 -6 9 -6
1 -8 -9
2 -6 4 10 1
1 3 -1
1 -5 -1
3 -1 -9 -1 2 3 -9
1 -9 -6
3 3 4 -9 4 3 9
1 -9 9
4 -7 -9 5 -9 -7 9 5 9
4 -7 6 -1 6 -1 -4 -7 -4
1 -9 -2
3 -6 8 -7 8 -7 3
4 -10 3 -6 4 -10 4 -6 3
3 -7 1 -7 -10 -8 -10
1 9 -4
3 1 -3 1 9 -10 -3
3 -10 1 -6 -4 -6 1
2 -8 7 8 7
2 0 -5 6 -5
1 2 -8
4 -1 -8 10 -8 10 1 -1 1
4 0 -4 4 1 4 -4 0 1
2 4 3 4 -5
2 4 -8 1 -8
4 4 -4 9 -10 4 -10 9 -4
4 7 -2 -3 -2 -3 2 7 2
4 -7 1 4 -2 4 1 -7 -2
3 0 6 7 -7 0 -7
2 6 10 8 -4
4 1 10 8 -6 1 -6 8 10
3 10 5 10 3 1 3
4 -6 5 -1 5 -1 -9 -6 -9
1 2 -9
3 -6 8 -6 -1 -3 8
4 -5 10 -5 3 0 10 0 3
2 4 1 -5 1
4 8 -9 -4 -9 -4 2 8 2
2 9 -8 -1 -8
4 4 -2 8 -2 8 5 4 5
2 -10 -1 1 10
2 9 -3 2 -3
1 7 4
3 -8 1 -4 6 -8 6
3 1 -3 1 -9 5 -3
4 -9 5 3 5 -9 4 3 4
4 -4 -8 4 9 4 -8 -4 9
2 -6 1 -6 0
3 6 9 6 4 4 9
2 3 -1 -3 -4
4 -2 8 -2 -7 10 8 10 -7
3 -3 1 0 1 0 10
3 8 -6 10 -6 8 -10
4 8 -6 8 -3 -4 -3 -4 -6
1 6 1
2 5 0 7 -8
3 2 -4 -9 -9 -9 -4
3 -5 -10 1 -10 -5 -1
4 10 -3 -9 -3 -9 -10 10 -10
4 -4 3 -4 4 -6 4 -6 3
4 3 7 3 -2 -10 7 -10 -2
`

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	scanner := bufio.NewScanner(strings.NewReader(testcasesARaw))
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
		if err := cmd.Run(); err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expected {
			fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
