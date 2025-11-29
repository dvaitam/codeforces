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

const testcases = `100
4
6
6
24
11
48
43
20
17
39
14
39
3
38
44
11
28
41
26
47
33
24
35
29
33
18
3
2
24
30
21
25
28
34
11
36
12
16
15
2
12
21
12
9
33
33
24
33
44
36
12
29
27
48
34
49
24
38
23
24
29
11
49
26
46
48
30
42
34
16
32
18
32
33
33
23
43
30
30
23
37
47
36
47
30
32
43
15
21
45
11
40
18
50
31
20
20
46
33
36`

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	cases, err := loadTestcases()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for idx, n := range cases {
		want := referenceSolve(n)
		input := fmt.Sprintf("%d\n", n)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if want != got {
			fmt.Printf("test %d failed\nexpected: %s\ngot: %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

func loadTestcases() ([]int, error) {
	scanner := bufio.NewScanner(strings.NewReader(testcases))
	if !scanner.Scan() {
		return nil, fmt.Errorf("invalid testcases")
	}
	t, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil {
		return nil, fmt.Errorf("failed to parse test count: %w", err)
	}
	nums := make([]int, 0, t)
	for i := 0; i < t; i++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("missing n for case %d", i+1)
		}
		n, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
		if err != nil {
			return nil, fmt.Errorf("failed to parse n on case %d: %w", i+1, err)
		}
		nums = append(nums, n)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}
	return nums, nil
}

func referenceSolve(n int) string {
	var b strings.Builder
	b.Grow(n * 6) // rough prealloc
	x := 0
	for i := 0; i < n; i++ {
		switch i % 3 {
		case 0, 1:
			x++
			fmt.Fprintf(&b, "%d %d\n", x, 0)
		case 2:
			fmt.Fprintf(&b, "%d %d\n", x, 3)
		}
	}
	return strings.TrimSpace(b.String())
}
