package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesA = `50
98
54
6
34
66
63
52
39
62
46
75
28
65
18
37
18
97
13
80
33
69
91
78
19
40
13
94
10
88
43
61
72
13
46
56
41
79
82
27
71
62
57
67
34
8
71
2
12
93
52
91
86
81
1
79
64
43
32
94
42
91
9
25
73
29
31
19
70
58
12
11
41
66
63
14
39
71
38
91
16
71
43
70
27
78
71
76
37
57
12
77
50
41
74
31
38
24
25
24`

func solve(n int) string {
	if n%4 == 2 || n%4 == 3 {
		return "-1"
	}
	p := make([]int, n+1)
	half := n / 2
	for i := 1; i <= half; i += 2 {
		a := i
		b := i + 1
		c := n - a + 1
		d := n - b + 1
		p[a] = b
		p[b] = c
		p[c] = d
		p[d] = a
	}
	if n%2 == 1 {
		mid := (n + 1) / 2
		p[mid] = mid
	}
	var sb strings.Builder
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(p[i]))
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	scanner := bufio.NewScanner(strings.NewReader(testcasesA))
	scanner.Buffer(make([]byte, 0, 1024), 1<<20)
	if !scanner.Scan() {
		fmt.Fprintln(os.Stderr, "bad test data")
		os.Exit(1)
	}
	t, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "bad count: %v\n", err)
		os.Exit(1)
	}
	for i := 0; i < t; i++ {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "not enough testcases\n")
			os.Exit(1)
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			i--
			continue
		}
		n, err := strconv.Atoi(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d bad n: %v\n", i+1, err)
			os.Exit(1)
		}
		input := line + "\n"
		expected := solve(n)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", t)
}
