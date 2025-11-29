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

const testcasesRaw = `100
4 38
9 3 6 10
8 38
2 10 1 8 5 9 4 4
8 35
9 8 7 3 4 3 9 7
1 5
3
10 3
5 1 5 8 10 7 7 7 10 8
3 24
2 1 3
8 14
5 7 5 7 9 7 10 6
9 38
7 10 4 6 1 5 10 3 6
9 37
10 2 4 10 5 5 2 2 8
8 6
6 8 1 4 9 2 6 9
7 19
4 8 6 6 1 3 5
5 39
10 7 8 8 4
9 36
4 4 7 7 6 6 9 5 5
1 6
7
5 13
9 4 8 7 8
3 23
7 6 5
1 10
1
1 12
5
3 5
1 2 7
3 8
8 3 8
9 29
1 3 1 4 8 10 10 6 4
1 32
3
7 5
5 10 2 6 8 6 10
8 34
6 2 7 7 9 7 10 3
5 7
3 3 10 7 9
9 12
6 4 3 1 5 8 4 8 5
3 33
3 5 3
8 1
1 3 3 1 7 6 6 9
10 37
9 4 1 5 7 6 9 2 5 8
5 3
4 3 6 5 9
1 8
7
8 25
10 5 10 2 7 5 8 10
8 34
9 4 1 8 10 3 7 10
3 28
6 4 3
9 14
10 8 7 2 9 8 4 5 9
2 28
10 10
1 1
6
2 30
8 5
7 4
2 7 1 10 8 5 3
2 13
6 10
1 30
1
5 40
2 1 1 2 6
1 17
9
6 28
3 7 5 6 2 3
6 7
5 3 2 8 2 8
5 31
7 2 5 5 9
8 34
1 5 4 10 7 3 8 9
4 8
5 5 4 10
3 21
3 5 1
6 2
5 9 4 3 3 7
2 18
3 7
6 10
3 6 7 5 10 5
3 31
4 2 4
3 39
6 5 1
10 18
8 1 4 4 10 5 7 6 4 8
1 26
9
8 6
3 9 6 9 6 9 4 3
4 5
7 2 2 6
2 13
5 2
6 26
6 10 5 4 1 9
6 35
2 5 5 8 6 5
1 20
4
4 36
7 10 1 5
10 18
9 8 1 6 4 7 3 10 1 3
10 36
2 10 7 6 1 7 5 6 10 9
8 6
3 7 4 3 4 9 6 1
5 6
1 8 8 6 2
2 4
4 2
9 5
7 2 3 1 3 10 9 6 4
4 9
1 10 9 5
9 14
10 6 9 1 2 3 6 5 10
5 22
8 5 1 8 6
10 9
10 2 6 3 7 2 8 10 2 6
10 35
8 1 2 9 3 6 6 10 3 8
5 19
8 5 6 8 8
3 8
5 4 8
6 32
10 9 1 4 10 2
2 10
3 9
9 7
4 5 9 3 10 7 9 3 10
9 5
8 10 1 5 9 3 9 4 7
7 19
2 7 8 8 5 2 1
6 38
8 7 6 3 5 2
7 6
7 9 5 7 6 9 1
10 8
6 9 7 7 8 9 5 2 2 9
8 34
5 9 5 7 3 1 1 3
7 23
10 3 8 2 6 8 7
2 18
1 7
4 34
4 10 1 5
2 8
10 7
5 14
2 3 6 2 8
5 27
4 8 4 7 10
4 30
4 1 7 1
8 11
9 2 8 6 5 6 4 9
1 26
5
6 14
6 2 9 2 10 9
2 10
1 3
3 7
4 9 2
2 7
1 10
9 27
5 4 6 4 6 4 6 5 5
10 5
3 9 1 8 3 2 8 10 7 1
6 14
9 7 8 10 10 7
1 38
4
5 10
1 5 1 3 7
10 30
5 7 5 7 6 5 1 10 5 6
4 29
3 2 6 5
7 26
8 5 4 6 8 7 3
8 25
6 6 1 10 2 10 4 3
6 36
3 10 9 1 1 6
7 27
1 1 6 4 10 4 8
8 31
7 7 10 6 6 8 7 3
7 20
6 1 4 1 4 1 5
4 25
9 10 5 2
4 20
7 5 10 10
3 22
4 1 1
3 28
7 7 3
3 38
2 5 9
9 23
6 5 3 2 8 1 1 3 2
5 13
4 9 2 9 1
5 18
2 6 7 7 6
4 12
10 7 7 9
8 13
9 9 10 10 6 6 10 3
7 35
4 6 8 7 8 7 1
7 33
1 3 9 6 2 9 1
2 11
9 4
4 9
1 8 9 1
2 23
5 1
4 5
9 3 1 1
9 19
2 7 8 6 10 10 5 6 5
5 32
10 6 4 9 2
4 26
6 10 5 4
1 34
4
7 18
6 1 4 10 9 5 2
3 16
1 4 10
8 40
7 3 9 2 7 2 9 2
2 24
6 2
4 34
10 10 5 5
4 9
4 7 9 6
3 29
8 7 3
9 37
8 9 10 6 5 6 10 5 8
7 14
7 10 10 9 3 10 8
3 12
5 10 8
4 34
3 4 6 8
8 20
6 8 1 7 1 4 6 3
3 19
4 5 2
7 9
1 3 5 6 2 2 10
6 25
7 5 4 4 7 8
6 17
8 4 2 1 1 8
7 25
1 8 5 6 6 7 4
7 12
6 1 3 9 9 1 9
1 32
8
10 36
7 9 10 4 2 8 2 8 6 2
8 20
1 4 1 5 5 1 3 8
7 26
2 1 9 6 4 10 2
7 33
10 9 1 7 9 6 4 9
2 27
4 8
1 29
1
8 40
6 3 4 8 8 1 1 3
4 32
5 2 6 5
5 30
6 9 10 4 9
8 8
2 6 2 1 10 4 1 3
6 14
8 5 7 8 9 4
7 19
4 5 9 1 8 10 9
10 34
1 10 4 6 8 9 8 4 10 7
7 20
7 4 8 8 8 9 9
10 17
7 7 6 7 3 6 5 1 10 7
7 33
3 7 8 7 1 3 8
10 19
10 3 9 7 6 5 3 9 5 10
7 30
9 4 10 3 10 6 9
1 15
1
7 5
1 1 3 6 8 3 9
7 19
10 3 9 2 10 5 7
10 40
8 8 5 4 5 7 10 8 3 9
6 29
9 5 5 3 10 3 8
8 6
1 6 9 10 2 10 9 6
6 5
8 6 8 2 8 10
8 19
7 10 6 6 5 3 6 7
2 11
3 5
7 40
4 2 2 1 1 1 7
6 29
9 10 3 2 1 9 8
4 22
5 1 6 9
6 19
6 6 2 2 4 9
8 31
2 2 4 5 4 2 8 10
4 1
9 5 6 9
9 33
3 8 4 7 8 8 1 8 8
5 37
2 10 9 1 5
8 38
4 5 6 7 4 8 5 7
1 29
10
3 17
3 7 1
8 35
8 1 2 7 9 9 5 9
5 27
5 10 6 7 4
7 39
4 6 7 5 3 9 1
4 38
1 9 6 7
2 36
3 6`

// expected logic from 1185C2.go.
func expected(times []int, M int) []int {
	var cnt [101]int
	res := make([]int, len(times))
	for i, t := range times {
		L := M - t
		sum := 0
		kept := 0
		for tv := 1; tv <= 100; tv++ {
			c := cnt[tv]
			if c == 0 {
				continue
			}
			total := c * tv
			if sum+total <= L {
				sum += total
				kept += c
			} else {
				rem := (L - sum) / tv
				if rem > 0 {
					kept += rem
				}
				break
			}
		}
		res[i] = i - kept
		cnt[t]++
	}
	return res
}

type testCase struct {
	n     int
	M     int
	times []int
}

func parseTestcases() ([]testCase, error) {
	scan := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return nil, fmt.Errorf("invalid test file")
	}
	t, err := strconv.Atoi(scan.Text())
	if err != nil {
		return nil, err
	}
	tests := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			return nil, fmt.Errorf("invalid test file")
		}
		n, _ := strconv.Atoi(scan.Text())
		if !scan.Scan() {
			return nil, fmt.Errorf("invalid test file")
		}
		M, _ := strconv.Atoi(scan.Text())
		times := make([]int, n)
		for j := 0; j < n; j++ {
			if !scan.Scan() {
				return nil, fmt.Errorf("invalid test file")
			}
			v, _ := strconv.Atoi(scan.Text())
			times[j] = v
		}
		tests = append(tests, testCase{n: n, M: M, times: times})
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	return tests, nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.M)
	for i, v := range tc.times {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func runCase(bin string, tc testCase) error {
	input := buildInput(tc)
	gotStr, err := runBinary(bin, input)
	if err != nil {
		return err
	}
	gotFields := strings.Fields(strings.TrimSpace(gotStr))
	if len(gotFields) != len(tc.times) {
		return fmt.Errorf("expected %d outputs got %d", len(tc.times), len(gotFields))
	}
	exp := expected(tc.times, tc.M)
	for i, tok := range gotFields {
		v, err := strconv.Atoi(tok)
		if err != nil || v != exp[i] {
			return fmt.Errorf("position %d expected %d got %q", i+1, exp[i], tok)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTestcases()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
