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

const testcasesRaw = `
3 19 1 3 18
1 20 1
2 5 4 5
4 4 1 2 3 4
4 11 1 4 7 11
1 20 13
2 19 10 14
4 5 1 2 3 4
2 20 3 6
5 20 7 9 10 13 19
1 2 1
1 10 3
2 18 4 5
3 3 1 2 3
3 5 1 3 5
3 17 8 11 16
5 7 1 2 3 4 5
5 20 3 8 14 17 18
2 14 9 14
3 7 1 2 6
1 13 11
5 8 2 3 5 6 7
3 18 11 13 15
5 16 1 8 9 12 13
4 15 1 6 9 13
2 20 3 16
1 12 7
5 12 1 2 9 10 12
1 19 12
1 12 3
3 11 7 8 11
1 8 8
3 8 1 6 7
4 20 6 10 13 16
2 2 1 2
5 19 5 9 11 14 19
4 16 5 6 7 8
2 8 2 5
5 16 4 6 10 14 15
2 12 3 4
3 20 4 6 18
4 19 6 7 8 12
1 18 16
2 11 5 7
5 11 1 6 7 8 11
2 8 3 7
5 19 3 9 10 12 17
4 17 1 11 13 14
1 2 2
2 10 9 10
4 13 6 7 8 9
2 4 1 3
1 17 1
5 15 2 10 12 13 14
2 2 1 2
3 5 3 4 5
2 20 17 20
2 10 3 5
2 20 5 8
2 9 5 8
3 7 3 5 6
2 4 2 3
4 9 1 3 4 5
4 8 1 5 6 8
3 10 1 3 4
5 5 1 2 3 4 5
5 15 3 4 6 8 11
4 11 2 4 5 11
3 18 9 10 12
2 11 7 9
1 6 4
1 19 4
5 9 3 4 5 8 9
2 7 1 6
3 17 9 11 13
1 12 7
4 8 1 3 6 8
5 19 4 7 14 15 17
4 6 1 4 5 6
5 19 2 4 6 7 8
1 20 13
4 6 1 2 5 6
4 6 1 2 3 6
2 19 17 19
5 18 9 11 13 15 16
5 5 1 2 3 4 5
4 10 1 2 4 9
2 16 2 16
5 13 5 6 7 11 12
1 20 4
4 15 2 10 11 15
2 6 2 4
1 16 10
2 18 5 12
1 17 11
3 5 1 2 3
3 19 13 14 18
1 5 4
1 10 3
5 11 3 5 6 8 10
`

var rawTestcases = func() []string {
	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scanner.Buffer(make([]byte, 0, 1024), 1024*1024)
	var cases []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			cases = append(cases, line)
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return cases
}()

func splitInts(s string) []int {
	var res []int
	num := 0
	neg := false
	started := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '-' {
			neg = true
		} else if c >= '0' && c <= '9' {
			started = true
			num = num*10 + int(c-'0')
		} else {
			if started {
				if neg {
					num = -num
				}
				res = append(res, num)
				num = 0
				neg = false
				started = false
			}
		}
	}
	if started {
		if neg {
			num = -num
		}
		res = append(res, num)
	}
	return res
}

func solveCase(n, m int, a []int) (string, error) {
	if len(a) != n {
		return "", fmt.Errorf("invalid a length")
	}

	inA := make([]bool, m+1)
	for _, v := range a {
		inA[v] = true
	}

	reachable := make([]bool, m+1)
	reachable[0] = true

	var P []int
	for _, v := range a {
		if !reachable[v] {
			P = append(P, v)
			// Update reachable with v (unbounded)
			// Iterate from v to m. 
			// If reachable[j-v] is true, then j is reachable.
			for j := v; j <= m; j++ {
				if reachable[j-v] {
					reachable[j] = true
				}
			}
		}
	}

	// Check condition 2: all reachable sums <= m must be in A
	for j := 1; j <= m; j++ {
		if reachable[j] && !inA[j] {
			return "NO", nil
		}
	}

	var sb strings.Builder
	sb.WriteString("YES\n")
	sb.WriteString(strconv.Itoa(len(P)))
	sb.WriteByte('\n')
	for i, v := range P {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return strings.TrimRight(sb.String(), "\n"), nil
}

func parseCase(line string) (int, int, []int, error) {
	nums := splitInts(line)
	if len(nums) < 2 {
		return 0, 0, nil, fmt.Errorf("invalid line")
	}
	n, m := nums[0], nums[1]
	if len(nums) != 2+n {
		return 0, 0, nil, fmt.Errorf("expected %d numbers got %d", 2+n, len(nums))
	}
	a := nums[2:]
	return n, m, a, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v, output: %s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for idx, line := range rawTestcases {
		n, m, a, err := parseCase(line)
		if err != nil {
			fmt.Printf("case %d invalid: %v\n", idx+1, err)
			os.Exit(1)
		}
		expected, err := solveCase(n, m, a)
		if err != nil {
			fmt.Printf("case %d solve error: %v\n", idx+1, err)
			os.Exit(1)
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i, v := range a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("case %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(rawTestcases))
}