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

const testcasesDRaw = `19 2
33 2
31 6
17 3
42 13
25 12
23 16
26 21
31 17
26 11
38 23
33 31
26 17
47 23
30 13
49 2
18 13
24 13
32 25
48 11
16 9
17 4
44 17
9 5
37 34
15 11
16 13
36 7
35 12
32 7
12 5
38 5
49 8
11 4
14 5
46 27
28 5
11 8
36 11
26 23
38 15
43 16
43 38
33 17
11 4
31 28
11 5
50 13
16 13
22 13
48 19
15 8
50 29
39 7
10 7
5 2
17 7
32 9
26 17
42 11
26 5
35 31
36 17
20 11
50 7
27 17
15 2
40 13
34 3
44 15
14 9
48 37
34 31
49 6
29 10
34 11
30 7
34 29
12 5
24 11
24 17
14 11
33 2
43 22
14 11
24 7
41 12
10 7
37 30
39 35
13 8
26 7
49 43
28 25
15 2
10 3
12 5
35 32
17 2
29 12
`

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleD")
	cmd := exec.Command("go", "build", "-o", oracle, "755D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)


	scanner := bufio.NewScanner(strings.NewReader(testcasesDRaw))
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
		n := fields[0]
		k := fields[1]
		input := fmt.Sprintf("%s %s\n", n, k)
		// expected from oracle
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
