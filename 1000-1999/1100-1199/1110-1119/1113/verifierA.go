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

const testcasesRaw = `100
51 98
55 6
35 66
64 52
40 62
47 75
29 65
19 37
19 97
14 80
34 69
92 78
20 40
14 94
11 88
44 61
73 13
47 56
42 79
83 27
72 62
58 67
35 8
72 2
13 93
53 91
87 81
2 79
65 43
33 94
43 91
10 25
74 29
32 19
71 58
13 11
42 66
64 14
40 71
39 91
17 71
44 70
28 78
72 76
38 57
13 77
51 41
75 31
39 24
26 24
6 79
86 34
62 9
13 87
98 17
21 5
12 90
71 88
52 91
69 36
68 31
29 87
77 54
76 36
59 64
86 83
91 46
12 42
80 15
64 76
82 43
26 32
4 94
36 15
92 29
49 22
44 55
9 13
20 90
30 6
75 82
70 78
89 10
5 16
83 25
79 74
17 51
13 48
16 5
79 3
26 24
93 16
63 27
95 8
88 3
71 55
81 13
35 9
30 10
84 39`

func expected(n, v int) int {
	if v >= n-1 {
		return n - 1
	}
	m := n - v
	cost := v + m*(m+1)/2 - 1
	return cost
}

func run(bin, input string) (string, error) {
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scanner.Text())
	for caseNum := 1; caseNum <= t; caseNum++ {
		if !scanner.Scan() {
			fmt.Printf("case %d: missing n\n", caseNum)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scanner.Text())
		if !scanner.Scan() {
			fmt.Printf("case %d: missing v\n", caseNum)
			os.Exit(1)
		}
		v, _ := strconv.Atoi(scanner.Text())
		input := fmt.Sprintf("%d %d\n", n, v)
		want := fmt.Sprintf("%d", expected(n, v))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", caseNum, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", caseNum, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
