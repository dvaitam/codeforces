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

func expected(n int) int {
	a := (n + 2) / 2
	b := (n + 1) / 2
	return a * b
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	const testcasesRaw = `18
73
98
9
33
16
64
98
58
61
84
49
27
13
63
4
50
56
78
98
99
1
90
58
35
93
30
76
14
41
4
3
4
84
70
2
49
88
28
55
93
4
68
29
98
57
64
71
30
45
30
87
29
98
59
38
3
54
72
83
13
24
81
93
38
16
96
43
93
92
65
55
65
86
25
39
37
76
64
65
51
76
5
62
32
96
52
54
86
23
47
71
90
100
87
95
48
12
57
85`

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		n, err := strconv.Atoi(line)
		if err != nil {
			fmt.Printf("test %d: invalid number\n", idx)
			os.Exit(1)
		}
		expect := expected(n)
		input := fmt.Sprintf("1\n%d\n", n)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(out.String())
		got, err := strconv.Atoi(gotStr)
		if err != nil {
			fmt.Printf("test %d: cannot parse output %q\n", idx, gotStr)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed: expected %d got %d\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
