package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)
const testcasesARaw = `
51
99
55
7
35
67
65
53
41
63
47
77
29
67
19
39
19
99
15
81
35
71
93
79
21
41
15
95
11
89
45
63
73
15
47
57
43
81
83
29
73
63
59
69
35
9
73
3
13
95
53
93
87
83
3
81
65
45
33
95
43
93
11
27
75
31
33
21
71
59
13
13
43
67
65
15
41
73
39
93
17
73
45
71
29
79
73
77
39
59
13
79
51
43
75
33
39
25
27
25
`


func expected(n int) string {
	var sb strings.Builder
	mid := n / 2
	for i := 0; i < n; i++ {
		d := i - mid
		if d < 0 {
			d = -d
		}
		stars := d
		ds := n - 2*d
		for j := 0; j < stars; j++ {
			sb.WriteByte('*')
		}
		for j := 0; j < ds; j++ {
			sb.WriteByte('D')
		}
		for j := 0; j < stars; j++ {
			sb.WriteByte('*')
		}
		if i+1 < n {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
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
		var n int
		fmt.Sscan(line, &n)
		input := fmt.Sprintf("%d\n", n)
		exp := expected(n)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("Test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if out != exp {
			fmt.Printf("Test %d wrong answer\nExpected:\n%s\nGot:\n%s\n", idx, exp, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
