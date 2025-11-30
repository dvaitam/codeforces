package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	A int
	B int
	S string
}

// Embedded testcases from testcasesC.txt.
const testcaseData = `
3 6 ?1011?1?0
4 9 ?0?0?1110100?
2 9 1??0?11?000
7 5 1100?0??1000
6 5 10000?00?01
8 9 11?10??0001011100
3 4 1??101?
4 5 1?0000001
5 6 ?1010?00010
10 8 0?11001?100?0?0010
6 1 001?0?1
10 3 ?111?11??000?
9 2 11?????01?1
8 3 ??1?0?1???0
4 2 000010
0 10 0?1?10100?
0 7 ?01?10?
8 4 01101010?11?
2 4 10?101
3 8 ?011?1?100?
7 8 ??????01?0011?1
1 6 10?001?
8 5 111111????0?0
2 5 ?11?011
7 7 1?10?000?1?10?
9 5 1?1111?1??0001
10 3 ?0001??00??1?
1 3 10??
8 10 000?0?10?11?100?10
6 7 ?1?01?01100??
6 8 101????1011001
10 0 ?01??1????
4 8 1??1???1110?
7 9 0?01?01???01111?
10 0 00011111??
7 5 110110100011
2 9 111?1?1?111
7 3 ?11?100000
3 0 010
7 1 1??0001?
0 8 0?110?0?
8 10 1??01?101?00??010?
7 4 ??110?1001?
3 2 ?11??
4 7 ??0?1100??1
4 0 ??11
1 3 ?11?
3 6 00101???0
8 9 ?01111?111011?101
2 9 0???1?00?1?
2 10 01?111000?10
1 10 ??0?0?00?0?
9 10 111001?00?01??11???
6 0 011001
2 10 ?0100?01101?
0 5 00011
10 3 0?110?101??0?
4 6 0?????1?1?
6 4 0?1?00??00
3 3 11?01?
4 8 1101?0?10??1
8 8 ????0?11?000?0?0
7 5 11001110?011
10 10 ??0?0?01??0101??1111
6 9 10010???000?0??
5 8 1??11?0?00?10
1 0 ?
5 6 ?11??0?1111
4 5 1?0???001
10 1 ?0???11?0?0
7 3 1?01?00011
10 3 ?1?111???0111
8 1 ?11000110
0 10 0011??100?
1 4 011??
3 8 ?10?010?01?
3 1 ??00
6 9 0?0100??0??????
1 3 110?
10 0 10?0?10?00
0 4 11?1
2 9 0??1??1?1?0
6 6 011101?10??0
5 5 0111111?10
2 2 1000
9 8 ?0?1?0?0?11?00?00
10 0 0?111?0??0
10 8 ?001010??11?001?01
7 8 10?10?101??0?11
4 9 00100?001??0?
3 8 ?10?01??0??
4 0 ?00?
6 10 1??010?0????011?
6 6 110?00?01010
6 8 1000??01100??0
8 6 ?1001??01010?0
2 4 10??0?
6 1 1???111
5 7 0?1011??10??
8 4 0?111?0??00?
`

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	res := make([]testCase, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 3 {
			return nil, fmt.Errorf("case %d invalid format", i+1)
		}
		A, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("case %d bad A: %v", i+1, err)
		}
		B, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("case %d bad B: %v", i+1, err)
		}
		S := fields[2]
		res = append(res, testCase{A: A, B: B, S: S})
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("no test data")
	}
	return res, nil
}

// solve mirrors 1512C.go for one test case.
func solve(tc testCase) string {
	A, B := tc.A, tc.B
	arr := []byte(tc.S)
	n := len(arr)
	question := 0
	for i := 0; i < n; i++ {
		switch arr[i] {
		case '?':
			question++
		case '0':
			A--
		case '1':
			B--
		}
	}
	for i := 0; i < n; i++ {
		j := n - 1 - i
		if arr[i] != '?' && arr[j] == '?' {
			arr[j] = arr[i]
			if arr[i] == '0' {
				A--
			} else {
				B--
			}
			question--
		}
	}
	for i := n / 2; i < n; i++ {
		if arr[i] != '?' {
			continue
		}
		j := n - 1 - i
		if i == j {
			if A%2 == 1 {
				arr[i] = '0'
				A--
			} else if B%2 == 1 {
				arr[i] = '1'
				B--
			}
			question--
		} else {
			if A >= 2 {
				arr[i], arr[j] = '0', '0'
				A -= 2
				question -= 2
			} else if B >= 2 {
				arr[i], arr[j] = '1', '1'
				B -= 2
				question -= 2
			}
		}
	}
	ok := true
	for i := 0; i < n; i++ {
		if arr[i] != arr[n-1-i] {
			ok = false
			break
		}
	}
	if !ok || A != 0 || B != 0 || question != 0 {
		return "-1"
	}
	return string(arr)
}

func runCandidate(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.A, tc.B))
	sb.WriteString(tc.S)
	sb.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		expect := solve(tc)
		got, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
