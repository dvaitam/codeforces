package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const testcasesERaw = `50 7
6 5
66 8
52 5
62 6
75 4
65 3
37 3
97 2
80 5
69 10
19 5
13 2
88 6
61 9
13 6
56 6
79 4
71 8
57 9
34 1
71 1
12 7
91 1
79 8
43 4
94 6
91 2
25 10
29 4
19 9
58 2
11 6
66 8
14 5
71 5
91 2
71 6
70 4
78 9
76 5
57 2
77 7
41 10
31 5
24 4
24 1
79 5
61 2
12 3
20 1
11 9
88 7
91 9
36 9
31 4
87 10
54 10
36 8
64 6
11 6
79 2
63 10
81 6
25 4
3 5
15 4
48 3
43 7
8 2
19 4
6 10
82 9
78 2
4 2
82 4
78 10
16 7
12 6
15 1
78 1
25 3
92 2
62 4
94 1
87 1
70 7
80 2
34 2
29 2
83 5
45 7
24 1
65 8
6 10
13 7
26 5
46 8
73 3
90 4
`

const mod = 1000000007

func modPow(x, y int) int {
	res := 1
	base := x % mod
	exp := y
	for exp > 0 {
		if exp&1 == 1 {
			res = (res * base) % mod
		}
		base = (base * base) % mod
		exp >>= 1
	}
	return res
}

func runCandidate(bin, input string) (string, error) {
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	scanner := bufio.NewScanner(strings.NewReader(testcasesERaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		var x, y int
		fmt.Sscan(line, &x, &y)
		expected := fmt.Sprintf("%d", modPow(x, y))
		got, err := runCandidate(bin, line+"\n")
		if err != nil {
			fmt.Printf("test %d: %v\n", idx, err)
			os.Exit(1)
		}
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
