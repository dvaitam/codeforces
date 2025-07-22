package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

func expectedAnswers(n int, a []int64, ks []int) []int64 {
	arr := make([]int64, n)
	copy(arr, a)
	sort.Slice(arr, func(i, j int) bool { return arr[i] > arr[j] })
	P := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		P[i] = P[i-1] + arr[i-1]
	}
	var ans1 int64
	for i := 2; i <= n; i++ {
		ans1 += arr[i-1] * int64(i-1)
	}
	uniq := make(map[int]struct{})
	for _, k := range ks {
		uniq[k] = struct{}{}
	}
	ansMap := make(map[int]int64, len(uniq))
	rem := n - 1
	for k := range uniq {
		if k <= 1 {
			ansMap[k] = ans1
		} else {
			curCount := 0
			kPow := int64(1)
			d := 1
			var ans int64
			for curCount < rem {
				kPow *= int64(k)
				level := kPow
				use := level
				if use > int64(rem-curCount) {
					use = int64(rem - curCount)
				}
				lIdx := curCount + 1 + 1
				rIdx := curCount + int(use) + 1
				if lIdx < 2 {
					lIdx = 2
				}
				if rIdx > n {
					rIdx = n
				}
				sumB := P[rIdx] - P[lIdx-1]
				ans += sumB * int64(d)
				curCount += int(use)
				d++
			}
			ansMap[k] = ans
		}
	}
	res := make([]int64, len(ks))
	for i, k := range ks {
		res[i] = ansMap[k]
	}
	return res
}

func runCase(bin string, n int, a []int64, ks []int) error {
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range a {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(fmt.Sprintf("%d", v))
	}
	input.WriteByte('\n')
	q := len(ks)
	input.WriteString(fmt.Sprintf("%d\n", q))
	for i, k := range ks {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(strconv.Itoa(k))
	}
	input.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) != q {
		return fmt.Errorf("expected %d numbers got %d", q, len(fields))
	}
	exp := expectedAnswers(n, a, ks)
	for i, f := range fields {
		got, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return fmt.Errorf("bad output: %v", err)
		}
		if got != exp[i] {
			return fmt.Errorf("at index %d expected %d got %d", i, exp[i], got)
		}
	}
	return nil
}

func parseLine(line string) (int, []int64, []int, error) {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return 0, nil, nil, fmt.Errorf("invalid line")
	}
	n, _ := strconv.Atoi(parts[0])
	q, _ := strconv.Atoi(parts[1])
	need := 2 + n + q
	if len(parts) != need {
		return 0, nil, nil, fmt.Errorf("expected %d numbers got %d", need, len(parts))
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		v, _ := strconv.ParseInt(parts[2+i], 10, 64)
		a[i] = v
	}
	ks := make([]int, q)
	for i := 0; i < q; i++ {
		v, _ := strconv.Atoi(parts[2+n+i])
		ks[i] = v
	}
	return n, a, ks, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		n, a, ks, err := parseLine(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", idx, err)
			os.Exit(1)
		}
		if err := runCase(bin, n, a, ks); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
