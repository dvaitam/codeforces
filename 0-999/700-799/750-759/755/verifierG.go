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

const testcasesGRaw = `30 5
48 3
28 5
50 6
17 7
27 10
34 5
37 8
28 4
46 3
43 8
19 3
7 10
12 3
17 5
18 2
21 2
1 4
5 10
45 6
34 1
18 10
42 8
28 2
11 5
44 4
15 10
21 3
3 1
1 4
7 1
32 2
20 4
33 6
46 7
11 4
46 9
45 9
38 3
30 7
13 1
34 10
40 7
37 8
10 7
40 9
33 2
12 10
11 6
40 3
40 7
24 10
28 7
19 4
26 7
38 7
33 4
47 8
6 1
44 10
31 9
3 6
22 6
31 10
49 5
49 10
34 2
4 2
45 5
12 5
48 4
43 2
7 10
13 1
20 7
17 2
33 10
3 1
48 3
43 4
6 7
9 10
48 10
23 2
14 4
31 10
6 9
15 9
30 3
30 10
24 4
9 2
20 8
33 1
26 4
28 8
38 3
49 7
3 9
1 9
`

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleG")
	cmd := exec.Command("go", "build", "-o", oracle, "755G.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)


	scanner := bufio.NewScanner(strings.NewReader(testcasesGRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) != 2 {
			fmt.Printf("bad test %d\n", idx)
			os.Exit(1)
		}
		input := line + "\n"
		cmdO := exec.Command(oracle)
		cmdO.Stdin = strings.NewReader(input)
		var outO bytes.Buffer
		cmdO.Stdout = &outO
		if err := cmdO.Run(); err != nil {
			fmt.Printf("oracle error %v\n", err)
			os.Exit(1)
		}
		expected := strings.TrimSpace(outO.String())

		cmd := exec.Command(exe)
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
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
