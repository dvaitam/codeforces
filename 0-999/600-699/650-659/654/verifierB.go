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
	oracle := filepath.Join(dir, "oracleB")
	cmd := exec.Command("go", "build", "-o", oracle, "654B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("oracle build failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func lineToInput(line string) (string, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return "", fmt.Errorf("empty line")
	}
	if _, err := strconv.Atoi(line); err != nil {
		return "", err
	}
	return line + "\n", nil
}

func run(bin string, input string) (string, string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	return out.String(), stderr.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	const testcasesBRaw = `4
18
27
25
24
2
8
3
15
24
14
15
20
12
25
6
3
15
0
28
26
12
13
19
24
24
0
22
14
8
23
25
7
18
30
3
28
10
0
0
0
20
17
0
30
28
12
21
6
13
23
0
16
7
24
14
30
15
17
7
11
7
21
7
24
14
30
9
29
0
13
26
29
17
29
20
3
5
20
23
27
9
3
23
10
28
23
22
16
29
30
13
16
26
29
21
6
9
9
18`

	scanner := bufio.NewScanner(strings.NewReader(testcasesBRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		input, err := lineToInput(line)
		if err != nil {
			fmt.Printf("test %d: parse error: %v\n", idx, err)
			os.Exit(1)
		}
		expOut, expErr, err := run(oracle, input)
		if err != nil {
			fmt.Printf("oracle run failed on test %d: %v\nstderr: %s\n", idx, err, expErr)
			os.Exit(1)
		}
		gotOut, gotErr, err := run(target, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, gotErr)
			os.Exit(1)
		}
		if strings.TrimSpace(gotOut) != strings.TrimSpace(expOut) {
			fmt.Printf("test %d failed\nexpected: %s\n got: %s\n", idx, strings.TrimSpace(expOut), strings.TrimSpace(gotOut))
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
