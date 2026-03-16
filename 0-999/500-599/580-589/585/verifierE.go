package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleE")
	cmd := exec.Command("go", "build", "-o", oracle, "585E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

const testcasesERaw = `6 9 17 6 5 19 20
10 16 12 7 4 12 7 10 13 5 5
7 12 12 11 12 2 8 11
3 7 19 7
8 9 4 18 10 5 3 14 4
7 16 17 6 6 3 5 6
2 14 5
5 4 20 20 5 8
5 14 7 18 16 10
6 7 14 10 7 7 4
4 18 13 16 4
6 18 15 11 16 3 6
6 7 17 9 9 4 6
8 15 5 13 9 5 12 6 8
4 19 12 17 12
10 17 18 10 11 2 7 10 8 18 19
2 11 16
3 5 3 16
9 18 19 20 19 14 15 12 12 9
6 17 6 18 2 17 6
1 4
9 2 6 5 3 12 11 8 4 3
4 7 13 14 18
10 13 15 8 14 3 9 20 2 16 2
7 6 8 4 20 6 12 10
5 16 12 18 15 17
6 10 18 11 16 19 12
7 12 14 10 2 8 16 11
7 15 6 7 17 8 7 8
2 10 18
9 7 20 7 17 5 17 9 15 7
2 3 18
1 4
2 20 6
3 20 8 13
6 5 5 14 17 6 12
10 18 8 9 20 16 11 5 6 15 5
4 6 7 10 10
6 7 16 12 19 5 4
7 2 3 8 12 15 9 12
7 20 5 12 4 17 17 11
4 14 3 13 12
3 4 19 10
8 2 7 6 19 20 12 16 4
7 10 16 11 2 14 2 16
2 6 11
6 11 6 17 16 2 12
3 15 2 11
3 20 7 4
7 18 16 19 6 15 12 6
9 13 20 7 4 4 9 5 8 17
4 2 14 12 10
6 20 10 16 18 3 11
8 15 9 17 2 18 14 12 8
10 13 11 7 7 10 8 5 13 14 6
5 10 7 14 2 4
3 14 15 4
2 20 9
9 4 8 12 16 19 5 7 7 20
10 3 16 10 8 19 3 17 4 2 3
10 8 7 17 8 18 11 20 20 4 6
6 19 2 12 11 19 9
10 11 2 20 15 13 4 12 10 13 5
5 20 6 13 3 2
10 14 2 15 9 15 20 16 11 2 7
10 16 19 4 2 3 4 6 14 14 7
4 19 5 2 3
7 14 11 7 14 2 12 15
9 8 18 13 9 3 15 17 18 3
4 12 16 7 20
4 8 3 8 5
5 12 12 11 6 7
10 15 10 3 14 5 20 17 7 14 9
7 14 6 19 16 13 3 18
4 8 18 14 17
4 20 4 19 5
3 7 11 2
6 5 6 18 7 19 15
5 2 13 6 3 18
4 16 10 15 3
4 20 10 8 12
2 10 17
6 11 3 10 5 13 11
4 12 9 9 8
4 14 5 19 10
3 13 14 2
9 3 19 15 5 10 8 18 11 4
2 18 15
9 13 7 6 5 11 11 7 17 11
4 19 20 9 4
8 8 13 15 10 12 5 6 16
6 18 6 4 16 14 17
6 11 17 6 11 3 15
2 16 20
3 4 5 14
4 11 5 15 20
5 19 2 20 3 14
7 18 13 9 9 7 5 5
9 6 11 13 17 18 9 15 6 7
6 7 2 12 6 5 18
`

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	scanner := bufio.NewScanner(strings.NewReader(testcasesERaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) < 1 {
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		if len(fields) != 1+n {
			fmt.Printf("test %d: expected %d numbers got %d\n", idx, n+1, len(fields))
			os.Exit(1)
		}
		input := strings.Join(fields, " ") + "\n"
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
