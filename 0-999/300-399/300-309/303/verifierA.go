package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases (first line is count).
const embeddedTestcases = `100
50
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

func normalize(s string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	return strings.Join(lines, "\n")
}

func runBin(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return normalize(out.String()), nil
}

// Embedded solution logic from 303A.go.
func solveCase(n int) string {
	if n%2 == 0 {
		return "-1"
	}
	var b strings.Builder
	// line 1
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(strconv.Itoa(i))
	}
	b.WriteByte('\n')
	// line 2
	for i := 1; i < n; i++ {
		if i > 1 {
			b.WriteByte(' ')
		}
		b.WriteString(strconv.Itoa(i))
	}
	if n > 1 {
		b.WriteByte(' ')
	}
	b.WriteString("0\n")
	// line 3
	first := true
	for i := 1; i < n; i += 2 {
		if !first {
			b.WriteByte(' ')
		}
		b.WriteString(strconv.Itoa(i))
		first = false
	}
	for i := 0; i < n; i += 2 {
		if !first {
			b.WriteByte(' ')
		}
		b.WriteString(strconv.Itoa(i))
		first = false
	}
	return normalize(b.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	lines := strings.Split(strings.TrimSpace(embeddedTestcases), "\n")
	if len(lines) == 0 {
		fmt.Fprintln(os.Stderr, "no testcases embedded")
		os.Exit(1)
	}
	t, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil || t != len(lines)-1 {
		fmt.Fprintf(os.Stderr, "invalid testcase count header: %v\n", err)
		os.Exit(1)
	}
	for i := 1; i <= t; i++ {
		line := strings.TrimSpace(lines[i])
		n, err := strconv.Atoi(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: invalid n\n", i)
			os.Exit(1)
		}
		input := fmt.Sprintf("%d\n", n)
		expected := solveCase(n)
		got, err := runBin(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:%s\nexpected:%s\ngot:%s\n", i, line, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
