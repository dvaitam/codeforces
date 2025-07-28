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

type testCaseE struct {
	n int
	l int
	r int
	s int
}

func parseTestcasesE(path string) ([]testCaseE, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var cases []testCaseE
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 4 {
			return nil, fmt.Errorf("bad line: %s", line)
		}
		n, _ := strconv.Atoi(fields[0])
		l, _ := strconv.Atoi(fields[1])
		r, _ := strconv.Atoi(fields[2])
		s, _ := strconv.Atoi(fields[3])
		cases = append(cases, testCaseE{n, l, r, s})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

func solveE(n, l, r, s int) string {
	l--
	r--
	k := r - l + 1
	for first := 1; first+k-1 <= n; first++ {
		sum := k*first + k*(k-1)/2
		if sum > s {
			break
		}
		extra := s - sum
		if extra <= k {
			ans := make([]int, n)
			used := make([]bool, n+1)
			needAdd := r - extra + 1
			ok := true
			for i := l; i <= r; i++ {
				val := first + (i - l)
				if i >= needAdd {
					val++
				}
				if val > n {
					ok = false
					break
				}
				ans[i] = val
				used[val] = true
			}
			if !ok {
				continue
			}
			cur := 1
			for i := 0; i < n; i++ {
				if i >= l && i <= r {
					continue
				}
				for cur <= n && used[cur] {
					cur++
				}
				ans[i] = cur
				used[cur] = true
			}
			var sb strings.Builder
			for i := 0; i < n; i++ {
				if i > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(strconv.Itoa(ans[i]))
			}
			return sb.String()
		}
	}
	return "-1"
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
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcasesE("testcasesE.txt")
	if err != nil {
		fmt.Println("failed to parse testcases:", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", tc.n, tc.l, tc.r, tc.s))
		expected := solveE(tc.n, tc.l, tc.r, tc.s)
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
