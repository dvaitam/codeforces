package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Embedded testcases (one n per line).
const embeddedTestcases = `17
39
36
10
25
40
32
42
39
6
40
2
32
18
37
16
14
47
32
36
37
32
27
42
11
16
42
11
35
26
49
2
44
6
12
50
39
4
21
3
19
32
40
48
26
47
29
27
48
38
30
10
25
8
4
10
33
15
18
45
29
42
21
28
34
26
38
24
36
39
28
39
16
23
45
3
19
40
44
46
12
46
22
36
38
38
8
47
43
15
42
38
19
20
9
6
32
42
32
7`

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

// Embedded reference solution (adapted from 66D.go).
func solve(n int) []string {
	if n == 2 {
		return []string{"-1"}
	}
	res := make([]string, 0, n)
	res = append(res, "10")
	res = append(res, "15")
	for i := 2; i < n; i++ {
		res = append(res, fmt.Sprint(i*6-6))
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(embeddedTestcases))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		input := line + "\n"
		n := 0
		fmt.Sscan(line, &n)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		outTokens := strings.Fields(got)
		expect := solve(n)
		if len(outTokens) != len(expect) {
			fmt.Printf("case %d failed: expected %d numbers got %d\n", idx, len(expect), len(outTokens))
			os.Exit(1)
		}
		ok := true
		for i := range expect {
			if outTokens[i] != expect[i] {
				ok = false
				break
			}
		}
		if !ok {
			fmt.Printf("case %d failed: expected %v got %v\n", idx, expect, outTokens)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
