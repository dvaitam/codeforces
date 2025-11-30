package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases from testcasesB.txt.
const embeddedTestcasesB = `1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
17
18
19
20
21
22
23
24
25
26
27
28
29
30
31
32
33
34
35
36
37
38
39
40
41
42
43
44
45
46
47
48
49
50
1000000000000000000
999999999999999999
999999999999999998
999999999999999997
999999999999999996
999999999999999995
999999999999999994
999999999999999993
999999999999999992
999999999999999991
999999999999999990
999999999999999989
999999999999999988
999999999999999987
999999999999999986
999999999999999985
999999999999999984
999999999999999983
999999999999999982
999999999999999981
999999999999999980
999999999999999979
999999999999999978
999999999999999977
999999999999999976
999999999999999975
999999999999999974
999999999999999973
999999999999999972
999999999999999971
999999999999999970
999999999999999969
999999999999999968
999999999999999967
999999999999999966
999999999999999965
999999999999999964
999999999999999963
999999999999999962
999999999999999961
999999999999999960
999999999999999959
999999999999999958
999999999999999957
999999999999999956
999999999999999955
999999999999999954
999999999999999953
999999999999999952
999999999999999951`

func digitSumStr(s string) int {
	sum := 0
	for i := 0; i < len(s); i++ {
		sum += int(s[i] - '0')
	}
	return sum
}

func solve770B(s string) int64 {
	x, _ := strconv.ParseInt(s, 10, 64)
	best := x
	bestSum := digitSumStr(s)

	for i := 0; i < len(s); i++ {
		if s[i] == '0' {
			continue
		}
		t := s[:i] + string(s[i]-1) + strings.Repeat("9", len(s)-i-1)
		num, _ := strconv.ParseInt(t, 10, 64)
		if num <= 0 {
			continue
		}
		sum := digitSumStr(t)
		if sum > bestSum || (sum == bestSum && num > best) {
			best = num
			bestSum = sum
		}
	}
	return best
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	lines := strings.Split(strings.TrimSpace(embeddedTestcasesB), "\n")
	for idx, line := range lines {
		want := strconv.FormatInt(solve770B(line), 10)
		input := line + "\n"
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(lines))
}
