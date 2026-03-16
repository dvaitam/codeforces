package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const testcasesARaw = `2 2
3 3
4 4
5 5
6 6
7 7
8 8
9 9
10 10
11 11
12 12
13 13
14 14
15 15
16 16
17 17
18 18
19 19
20 20
21 21
22 22
23 23
24 24
25 25
26 2
27 2
28 3
29 4
30 5
31 6
32 7
33 8
34 9
35 10
36 11
37 12
38 13
39 14
40 15
41 16
42 17
43 18
44 19
45 20
46 21
47 22
48 23
49 24
50 25
51 26
52 2
53 2
54 3
55 4
56 5
57 6
58 7
59 8
60 9
61 10
62 11
63 12
64 13
65 14
66 15
67 16
68 17
69 18
70 19
71 20
72 21
73 22
74 23
75 24
76 25
77 26
78 2
79 2
80 3
81 4
82 5
83 6
84 7
85 8
86 9
87 10
88 11
89 12
90 13
91 14
92 15
93 16
94 17
95 18
96 19
97 20
98 21
99 22
100 23
100 26
`

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

func checkPassword(n, k int, s string) error {
	if len(s) != n {
		return fmt.Errorf("length mismatch: expected %d got %d", n, len(s))
	}
	set := make(map[rune]struct{})
	prev := rune(0)
	for i, r := range s {
		if r < 'a' || r > 'z' {
			return fmt.Errorf("invalid character %q", r)
		}
		if i > 0 && r == prev {
			return fmt.Errorf("consecutive characters equal at pos %d", i)
		}
		set[r] = struct{}{}
		prev = r
	}
	if len(set) != k {
		return fmt.Errorf("expected %d distinct letters got %d", k, len(set))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(testcasesARaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		var n, k int
		fmt.Sscan(line, &n, &k)
		input := fmt.Sprintf("%d %d\n", n, k)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if err := checkPassword(n, k, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
