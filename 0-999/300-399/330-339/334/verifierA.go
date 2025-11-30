package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesA.txt so the verifier is self-contained.
const testcasesA = `2
4
6
8
10
12
14
16
18
20
22
24
26
28
30
32
34
36
38
40
42
44
46
48
50
52
54
56
58
60
62
64
66
68
70
72
74
76
78
80
82
84
86
88
90
92
94
96
98
100
2
4
6
8
10
12
14
16
18
20
22
24
26
28
30
32
34
36
38
40
42
44
46
48
50
52
54
56
58
60
62
64
66
68
70
72
74
76
78
80
82
84
86
88
90
92
94
96
98
100`

// Embedded solution from 334A.go.
func buildExpected(n int) []int {
	res := make([]int, 0, n*n)
	for i := 1; i <= n; i++ {
		for k := 0; k < n; k++ {
			var val int
			if k < n/2 {
				val = k*n + i
			} else {
				val = k*n + n - (i - 1)
			}
			res = append(res, val)
		}
	}
	return res
}

func parseOutput(n int, output string) ([]int, error) {
	var nums []int
	r := strings.NewReader(output)
	for {
		var v int
		_, err := fmt.Fscan(r, &v)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		nums = append(nums, v)
	}
	if len(nums) != n*n {
		return nil, fmt.Errorf("expected %d numbers, got %d", n*n, len(nums))
	}
	return nums, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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
		fmt.Println("usage: verifierA /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	lines := strings.Fields(testcasesA)
	for idx, nStr := range lines {
		n, err := strconv.Atoi(strings.TrimSpace(nStr))
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad testcase at %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		input := fmt.Sprintf("%d\n", n)
		expect := buildExpected(n)
		gotStr, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := parseOutput(n, gotStr)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		for i, v := range expect {
			if got[i] != v {
				fmt.Printf("test %d failed at position %d: expected %d got %d\n", idx+1, i+1, v, got[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(lines))
}
