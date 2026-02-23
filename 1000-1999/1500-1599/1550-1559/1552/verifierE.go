package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesE.txt so the verifier is self-contained.
const testcasesRaw = `
3 3 1 1 1 2 3 2 3 3 2
1 1 1
3 1 3 3 1
2 2 1 2 2 1
1 1 1
1 1 1
3 3 3 3 3 1 2 3 3 1 2
2 2 1 1 1 1
3 2 2 1 1 2 2 1
2 2 1 1 2 2
1 1 1
2 2 1 2 2 2
1 1 1
1 1 1
1 1 1
3 3 3 2 3 3 1 3 3 2 2
2 1 1 2
2 2 2 2 1 1
3 1 1 2 2
3 2 3 3 2 3 1 1
1 1 1
3 3 3 2 1 2 3 1 1 1 2
1 1 1
3 2 1 2 1 3 3 3
1 1 1
3 3 3 1 2 1 3 3 1 3 2
3 2 3 1 3 3 1 3
1 1 1
3 3 3 2 3 2 2 3 3 1 2
2 2 2 2 1 1
2 2 2 1 1 1
3 3 2 3 2 3 1 3 3 1 3
3 3 3 2 3 1 2 3 2 3 3
3 3 2 3 3 3 2 2 1 3 1
1 1 1
1 1 1
2 1 2 2
1 1 1
1 1 1
2 1 1 2
1 1 1
2 2 2 2 1 1
3 2 1 3 2 1 3 3
3 3 2 3 3 3 3 1 3 1 2
2 1 1 1
1 1 1
1 1 1
3 2 1 1 3 2 1 2
2 1 1 1
2 1 2 2
3 3 2 1 3 2 3 1 1 3 1
3 1 2 3 3
2 1 2 1
1 1 1
2 2 1 2 1 1
3 2 2 1 3 3 3 1
1 1 1
1 1 1
3 3 2 1 3 2 3 3 2 2 1
2 1 2 1
1 1 1
3 1 1 3 3
1 1 1
1 1 1
3 3 1 1 2 3 2 1 1 1 1
1 1 1
3 1 2 2 2
2 2 2 2 1 2
3 1 3 1 3
3 3 1 3 3 1 3 1 2 3 1
2 2 1 2 2 2
2 1 2 1
3 1 1 3 2
3 2 3 1 2 3 3 3
2 1 1 1
1 1 1
3 3 1 2 1 3 3 3 2 1 2
3 3 1 1 1 1 3 2 1 1 1
3 3 1 3 2 3 3 1 2 3 1
2 1 1 2
2 2 1 1 1 2
3 1 1 2 2
2 2 2 1 2 2
2 1 1 2
1 1 1
2 2 1 2 2 2
2 2 2 2 1 1
3 3 1 2 1 3 1 2 3 3 1
3 1 1 3 3
1 1 1
1 1 1
3 1 2 2 1
3 3 2 1 2 3 3 1 2 2 3
1 1 1
3 1 2 1 1
3 3 3 1 3 3 2 2 1 1 2
1 1 1
2 1 1 1
3 2 2 2 2 1 1 3
2 1 2 2

`

type testCase struct {
	n   int
	k   int
	seq []int
}

func solve(tc testCase) [][2]int {
	n, k := tc.n, tc.k
	L := make([]int, n+1)
	R := make([]int, n+1)
	if k == 1 {
		for i := 1; i <= n; i++ {
			L[i] = i
			R[i] = i
		}
	} else {
		for i := 1; i <= n; i += k - 1 {
			end := i + k - 2
			if end > n {
				end = n
			}
			for j := i; j <= end; j++ {
				L[j] = i
				R[j] = end
			}
		}
	}
	p := make([]int, n+1)
	vis := make([]bool, n+1)
	a0 := make([]int, n+1)
	a1 := make([]int, n+1)
	total := n * k
	for idx := 1; idx <= total; idx++ {
		x := tc.seq[idx-1]
		if vis[x] {
			continue
		}
		if p[x] != 0 {
			a0[x] = p[x]
			a1[x] = idx
			vis[x] = true
			for j := L[x]; j <= R[x]; j++ {
				p[j] = 0
			}
		} else {
			p[x] = idx
		}
	}
	res := make([][2]int, n)
	for i := 1; i <= n; i++ {
		res[i-1] = [2]int{a0[i], a1[i]}
	}
	return res
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	cases := make([]testCase, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			continue
		}
		k, err := strconv.Atoi(fields[1])
		if err != nil || k < 2 {
			continue
		}
		if len(fields) != 2+n*k {
			continue
		}
		seq := make([]int, n*k)
		counts := make(map[int]int)
		for i := 0; i < n*k; i++ {
			v, err := strconv.Atoi(fields[2+i])
			if err != nil {
				continue
			}
			seq[i] = v
			counts[v]++
		}
		valid := true
		for i := 1; i <= n; i++ {
			if counts[i] != k {
				valid = false
				break
			}
		}
		if !valid {
			continue
		}
		cases = append(cases, testCase{n: n, k: k, seq: seq})
	}
	return cases, nil
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for i, tc := range cases {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
		for idx, v := range tc.seq {
			if idx > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')

		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		gotLines := strings.Split(strings.TrimSpace(got), "\n")
		if len(gotLines) != tc.n {
			fmt.Printf("case %d: expected %d lines, got %d\n", i+1, tc.n, len(gotLines))
			os.Exit(1)
		}
		
		limit := (tc.n + tc.k - 2) / (tc.k - 1)
		cov := make([]int, tc.n*tc.k+1)

		for idx, line := range gotLines {
			fields := strings.Fields(line)
			if len(fields) != 2 {
				fmt.Printf("case %d line %d: expected 2 numbers, got %q\n", i+1, idx+1, line)
				os.Exit(1)
			}
			a, err1 := strconv.Atoi(fields[0])
			b, err2 := strconv.Atoi(fields[1])
			if err1 != nil || err2 != nil {
				fmt.Printf("case %d line %d: non-integer output\n", i+1, idx+1)
				os.Exit(1)
			}
			if a >= b {
				fmt.Printf("case %d line %d: invalid interval %d %d\n", i+1, idx+1, a, b)
				os.Exit(1)
			}
			if a < 1 || b > tc.n*tc.k {
				fmt.Printf("case %d line %d: out of bounds %d %d\n", i+1, idx+1, a, b)
				os.Exit(1)
			}
			color := idx + 1
			if tc.seq[a-1] != color || tc.seq[b-1] != color {
				fmt.Printf("case %d line %d: ends do not match color %d\n", i+1, idx+1, color)
				os.Exit(1)
			}
			for j := a; j <= b; j++ {
				cov[j]++
				if cov[j] > limit {
					fmt.Printf("case %d: point %d covered more than %d times\n", i+1, j, limit)
					os.Exit(1)
				}
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
