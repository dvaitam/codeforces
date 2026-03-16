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
	cmd := exec.Command("go", "build", "-o", oracle, "909C.go")
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
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	const testcasesRaw = `1 f
2 s f
5 s f f f s
7 s s s f f s s
6 s s f f f f
1 f
6 f f s f s s
9 s s s s f s s f s
5 s s s s s
10 s s f s f s s s s s
5 f s s f s
1 f
2 f f
5 f f f s f
4 f s f f
6 s f f f f f
2 f f
1 s
5 f f f f s
10 f f f f f s f s s s
1 s
8 f s s f s f f s
2 f s
3 s s s
3 s s s
10 s f f f s f f f f f
8 f f f f s f s f
10 f s s s s f f f s s
3 f f f
2 f f
3 f f f
1 s
8 s s f f s s f f
7 f f s s f f s
5 s s f s f
2 s f
1 s
1 s
8 s f f f s f s s
2 f s
4 f s s s
3 s s f
5 f f f s s
4 f f s f
8 s s s f s s s s
7 s s s f f s s
9 s f f f s f f s f
2 f f
5 s f s s f
9 s f f s f s f s f
3 s f s
6 f s s f s f
7 f s s s s s s
6 s f f f s s
1 s
8 s f f f s f s f
2 s f
5 f s f f s
8 s f s f f s f s
9 f s s s s s f f s
3 s s f
1 f
8 s s s f f s s f
6 f s s f s s
3 s s f
5 s s s s s
9 f f s f f s f f s
9 f s f f f s s f s
6 f s s s f s
8 s s f s f s s f
2 f s
4 s f s f
3 f s s
7 s s s s s f s
5 s s f s s
8 s f f s s f f f
6 f s f s f s
4 s s s f
6 s f s f f f
6 s f s s s f
6 f f f f s s
10 f f f f s s s f s f
6 s f f s s s
10 s f s s s f s s s f
2 s f
2 f f
7 s s f s f s f
3 f f s
2 s s
1 s
4 f s f f
3 f f f
9 f s f f s s f f f
4 f s f f
10 f s s s f f f f f s
4 s s s f
7 f s s s f f s
5 s s s f s
3 f f f
4 f s f f`

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		input := "1\n" + strings.Join(strings.Fields(line), " ") + "\n"
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
			fmt.Printf("test %d failed. Expected %s got %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
