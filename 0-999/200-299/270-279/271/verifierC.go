package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases from testcasesC.txt.
const embeddedTestcasesC = `8 3
25 5
31 10
91 6
82 53
66 46
57 52
8 1
44 7
71 6
31 18
45 3
69 12
87 65
61 14
98 79
46 27
58 30
52 47
32 26
67 17
46 37
60 43
72 71
88 8
90 84
82 9
62 40
94 27
67 55
76 52
30 8
88 69
44 38
28 2
87 2
33 11
14 13
38 4
8 6
98 37
50 6
7 1
12 1
6 2
24 3
1 1
32 15
1 1
81 75
37 23
43 8
78 74
37 7
91 40
32 2
41 40
61 5
100 84
45 3
44 23
57 1
18 17
8 1
16 4
66 25
30 25
14 2
80 72
91 14
100 1
53 14
69 3
16 1
33 2
70 22
43 9
28 16
56 46
75 63
57 17
13 12
70 4
38 35
48 36
91 14
43 4
87 82
61 2
66 15
79 20
40 1
96 5
66 7
51 13
18 12
19 3
14 3
28 13
88 16`

func solve271C(n, k int) string {
	if n <= 3 || 3*k > n {
		return "-1"
	}

	ans := make([]int, n+5)
	for i := range ans {
		ans[i] = -1
	}
	i := 1
	heso := 4
	lastval := 0
	minNotAns := 0
	for id := 1; id <= k-1; id++ {
		lastval = i + 3
		if heso == 4 {
			minNotAns = i + 2
			ans[i] = id
			ans[i+1] = id
			ans[i+3] = id
			heso = 2
			i += 2
		} else {
			minNotAns = i + 4
			ans[i] = id
			ans[i+2] = id
			ans[i+3] = id
			heso = 4
			i += 4
		}
	}
	if lastval < minNotAns {
		if minNotAns+2 == n {
			ans[minNotAns] = k
			ans[minNotAns+2] = k
			ans[minNotAns+1] = k - 1
			ans[minNotAns-1] = k
		} else {
			for j := minNotAns; j <= n; j++ {
				ans[j] = k
			}
			if minNotAns+1 <= n {
				ans[minNotAns+1] = 1
			}
		}
	} else {
		for j := minNotAns; j <= n; j++ {
			if ans[j] == -1 {
				ans[j] = k
			}
		}
	}

	var sb strings.Builder
	for j := 1; j <= n; j++ {
		if j > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(ans[j]))
	}
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	lines := strings.Split(strings.TrimSpace(embeddedTestcasesC), "\n")
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, "test %d: invalid line\n", idx+1)
			os.Exit(1)
		}
		n, err1 := strconv.Atoi(parts[0])
		k, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			fmt.Fprintf(os.Stderr, "test %d: parse error\n", idx+1)
			os.Exit(1)
		}
		input := fmt.Sprintf("%d %d\n", n, k)
		want := solve271C(n, k)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %s got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(strings.Split(strings.TrimSpace(embeddedTestcasesC), "\n")))
}
