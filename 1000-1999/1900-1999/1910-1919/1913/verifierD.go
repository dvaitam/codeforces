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

const mod int64 = 998244353

// solution logic from 1913D.go
func countReachable(p []int) int64 {
	n := len(p)
	prefix := make([]bool, n)
	mn := int(^uint(0) >> 1)
	for i, x := range p {
		if x < mn {
			mn = x
			prefix[i] = true
		}
	}
	suffix := make([]int, 0)
	mn = int(^uint(0) >> 1)
	for i := n - 1; i >= 0; i-- {
		if p[i] < mn {
			mn = p[i]
			suffix = append(suffix, i)
		}
	}
	stack := make([]int, 0)
	ancSum := make([]int64, 0)
	dp := make([]int64, n)
	sub := make([]int64, n)

	for i, x := range p {
		base := int64(0)
		if prefix[i] {
			base = 1
		}
		poppedTotal := int64(0)
		for len(stack) > 0 && p[stack[len(stack)-1]] > x {
			idx := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			ancSum = ancSum[:len(ancSum)-1]
			poppedTotal = (poppedTotal + sub[idx]) % mod
			sub[i] = (sub[i] + sub[idx]) % mod
		}
		anc := int64(0)
		if len(ancSum) > 0 {
			anc = ancSum[len(ancSum)-1]
		}
		dp[i] = (base + poppedTotal + anc) % mod
		sub[i] = (sub[i] + dp[i]) % mod
		stack = append(stack, i)
		if len(ancSum) > 0 {
			ancSum = append(ancSum, (ancSum[len(ancSum)-1]+dp[i])%mod)
		} else {
			ancSum = append(ancSum, dp[i])
		}
	}

	res := int64(0)
	for _, idx := range suffix {
		res = (res + dp[idx]) % mod
	}
	return res
}

func solveCase(fields []string) (string, error) {
	if len(fields) < 2 {
		return "", fmt.Errorf("bad test line")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return "", err
	}
	if len(fields) != 1+n {
		return "", fmt.Errorf("expected %d numbers got %d", 1+n, len(fields))
	}
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i], _ = strconv.Atoi(fields[1+i])
	}
	return fmt.Sprintf("%d", countReachable(p)), nil
}

const testcasesData = `
5 27 11 9 19 24
5 45 10 42 46 32
3 24 41 42
3 48 37 9
6 29 42 16 22 20 14
3 19 24 8
5 17 41 16 47 13
4 31 14 34 1
6 24 31 25 30 3 26
1 42
2 25 21
1 49
1 7
5 17 33 44 19 8
3 30 48 23
3 32 43 31
2 48 8
1 19
1 33
2 36 31
5 26 16 46 41 7
1 16
3 21 43 39
4 11 32 41 12
6 18 21 17 49 40 19
6 38 1 40 37 39 5
2 26 27
1 43
1 47
1 36
5 37 17 49 33 31
6 24 1 15 4 32 46
5 15 13 16 5 24
6 34 18 14 42 36 28
4 25 22 10 9
3 4 15 29
3 32 28 19
1 41
3 38 40 8
4 3 5 1 14
1 29
4 47 14 12 24
4 1 20 23 2
5 24 35 5 21 36
4 22 7 36 2
2 14 38
5 40 14 3 22 41
2 24 25
1 33
4 18 30 14 27
1 18
2 17 27
5 32 33 18 2 41
4 43 31 13 39
3 10 19 4
1 10
6 14 10 35 28 3 45
3 49 32 17
6 5 19 4 2 39 29
2 35 26
6 25 49 18 13 11 38
1 24
1 16
2 19 37
1 38
1 23
1 5
4 34 4 41 36
3 25 48 38
6 18 7 4 21 22 27
1 47
5 7 5 23 12 9
6 40 20 47 15 36 24
4 24 4 23 25
6 42 44 16 12 24 33
2 7 35
6 43 42 37 32 14 6
5 19 14 39 27 35
4 19 23 5 21
6 49 47 4 8 15 42
3 24 26 47
1 43
4 44 42 3 40
6 39 35 44 42 14 16
6 45 8 3 31 37 11
1 19
6 49 13 19 12 43 4
3 42 45 37
2 22 14
2 33 19
3 9 29 14
4 24 47 14 43
3 10 29 1
3 22 38 18
6 37 33 10 49 42 3
4 39 29 30 47
1 38
5 3 43 14 48 2
5 7 24 45 3 39
4 11 37 36 38
`

func lineToInput(line string) (string, error) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return "", fmt.Errorf("bad test line")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return "", err
	}
	if len(fields) != 1+n {
		return "", fmt.Errorf("expected %d numbers got %d", 1+n, len(fields))
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fields[1+i])
	}
	sb.WriteByte('\n')
	return sb.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(testcasesData))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		input, err := lineToInput(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad test %d: %v\n", idx, err)
			os.Exit(1)
		}
		want, err := solveCase(strings.Fields(line))
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad solve on test %d: %v\n", idx, err)
			os.Exit(1)
		}

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		if err := cmd.Run(); err != nil {
			fmt.Printf("Test %d: runtime error: %v\nstderr: %s\n", idx, err, errBuf.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != want {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx, input, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
