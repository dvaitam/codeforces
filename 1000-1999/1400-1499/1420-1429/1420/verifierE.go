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
const testcasesRaw = `12 1 0 1 0 0 0 0 1 1 1 0 1
8 1 1 1 0 0 0 1 0
11 1 1 0 0 0 0 0 1 0 0 1
9 1 1 0 1 1 0 0 1 1
12 1 1 0 1 0 0 0 0 1 0 0 0
9 0 0 0 0 0 0 1 1 0
6 1 1 0 0 0 1
6 1 0 1 1 0 1
4 1 1 0 0
9 0 0 0 1 0 0 1 1 1
5 1 0 1 1 1
2 1 0
5 0 1 1 0 1
6 1 1 0 1 0 0
3 1 1 0
1 1
3 0 1 1
10 0 0 1 1 0 1 0 0 1 1
8 1 1 1 0 1 0 0 1
2 1 0
10 1 0 0 1 0 1 1 1 1 0
5 0 1 0 1 0
6 0 1 1 0 0 0
10 1 0 1 0 0 1 1 1 1 0
1 1
11 1 1 1 0 1 1 1 1 0 1 0
12 0 1 1 0 1 0 1 1 0 1 1 0
3 0 1 0
3 1 0 0
1 0
3 1 1 1
3 0 1 1
2 1 1
6 1 0 0 0 1 1
10 1 1 0 0 1 1 0 1 1 0
1 0
3 1 1 0
6 1 1 0 0 0 0
6 1 0 0 0 0 1
4 1 0 1 0
6 0 0 0 1 0 0
4 1 1 0 0
1 1
7 0 0 1 1 1 0 1
10 0 0 0 1 1 0 0 0 1 1
2 1 0
4 1 0 0 1
8 0 1 0 0 0 0 0 0
8 1 0 1 0 0 0 1 0
12 1 0 1 1 0 1 1 1 0 1 0 0
10 1 1 1 1 1 0 0 1 0 0
10 1 1 0 1 1 0 1 0 1 1
1 1
1 1
12 0 0 0 0 1 1 1 1 0 1 0 1
5 0 1 0 0 1
11 0 0 0 0 0 1 1 0 0 0 0
3 0 0 1
10 1 1 0 0 1 0 0 0 1 1
2 0 0
2 0 1
11 1 1 1 0 0 1 0 1 0 0 1
7 1 1 1 1 0 0 1
9 0 1 1 0 1 1 1 1 1
12 1 0 0 1 1 0 0 0 1 0 1 0
7 1 0 1 1 1 0 1
8 0 0 0 1 1 0 1 1
4 0 0 0 1
7 1 1 1 0 1 1 0
2 0 0
10 0 0 1 1 0 1 1 0 0 1
3 0 0 0
1 0
8 1 1 0 1 1 0 0 0
11 0 1 1 0 1 0 1 0 1 1 0
7 1 1 0 0 0 0 0
11 1 1 1 0 0 1 1 0 0 0 1
1 1
6 1 1 0 0 0 0
7 0 1 0 0 0 0 1
5 1 1 1 1 1
3 0 1 0
12 0 0 0 1 1 0 1 1 1 1 0 0
11 1 0 1 1 0 0 1 0 1 0 1
9 0 1 1 1 1 0 0 1 1
10 0 0 0 1 1 0 0 0 1 0
5 0 1 1 1 0
1 0
10 1 0 1 1 0 0 1 1 0 0
11 0 1 1 0 1 1 1 0 0 1 0
11 0 1 1 1 1 0 0 0 1 0 0
5 0 0 1 0 0
2 0 1
7 1 0 0 1 1 0 0
5 1 1 0 1 0
11 1 1 0 0 0 1 1 0 1 1 0
1 0
11 1 0 1 0 1 1 1 0 0 0 1
2 1 1
12 0 0 1 1 0 1 1 0 0 1 1 1`

const MOD = 998244353

type testCase struct {
	n   int
	arr []int
}

func add(a, b int) int {
	s := a + b
	if s >= MOD {
		return s - MOD
	}
	return s
}

func mul(a, b int) int {
	return int((int64(a) * int64(b)) % MOD)
}

func powmod(a, e int) int {
	res, base := 1, a
	for e > 0 {
		if e&1 == 1 {
			res = mul(res, base)
		}
		base = mul(base, base)
		e >>= 1
	}
	return res
}

func solveCase(tc testCase) string {
	n := tc.n
	a := tc.arr
	P := make([]int, 0, n)
	for i, v := range a {
		if v == 1 {
			P = append(P, i+1)
		}
	}
	m := len(P)
	Z := n - m
	M := n * (n - 1) / 2
	if Z < 2 {
		out := make([]string, M+1)
		for i := range out {
			out[i] = "0"
		}
		return strings.Join(out, " ")
	}
	totalPairsZeros := Z * (Z - 1) / 2
	maxC := M
	const INF = 1000000000
	dpPrev := make([][]int, n+1)
	dpCurr := make([][]int, n+1)
	for j := range dpPrev {
		dpPrev[j] = make([]int, maxC+1)
		dpCurr[j] = make([]int, maxC+1)
		for k := 0; k <= maxC; k++ {
			dpPrev[j][k] = INF
		}
	}
	dpPrev[0][0] = 0
	usedPrev := []int{0}
	prevMaxK := 0
	for i := 1; i <= m; i++ {
		for j := range dpCurr {
			for k := 0; k <= maxC; k++ {
				dpCurr[j][k] = INF
			}
		}
		usedCurr := make([]int, 0, n)
		usedMap := make([]bool, n+1)
		currMaxK := 0
		minPos := i
		maxPos := n - (m - i)
		for _, j1 := range usedPrev {
			for j2 := minPos; j2 <= maxPos; j2++ {
				move := P[i-1] - j2
				if move < 0 {
					move = -move
				}
				d := j2 - j1 - 1
				segCost := d * (d - 1) / 2
				for k1 := 0; k1 <= prevMaxK; k1++ {
					prevCost := dpPrev[j1][k1]
					if prevCost >= INF {
						continue
					}
					k2 := k1 + move
					if k2 > maxC {
						continue
					}
					cost := prevCost + segCost
					if cost < dpCurr[j2][k2] {
						dpCurr[j2][k2] = cost
						if !usedMap[j2] {
							usedMap[j2] = true
							usedCurr = append(usedCurr, j2)
						}
						if k2 > currMaxK {
							currMaxK = k2
						}
					}
				}
			}
		}
		dpPrev, dpCurr = dpCurr, dpPrev
		usedPrev = usedCurr
		prevMaxK = currMaxK
	}

	bestCost := make([]int, maxC+1)
	for k := range bestCost {
		bestCost[k] = INF
	}
	for _, j := range usedPrev {
		tailZeros := n - j
		tailCost := tailZeros * (tailZeros - 1) / 2
		for k := 0; k <= prevMaxK; k++ {
			c := dpPrev[j][k]
			if c >= INF {
				continue
			}
			tot := c + tailCost
			if tot < bestCost[k] {
				bestCost[k] = tot
			}
		}
	}
	prefixMin := make([]int, maxC+1)
	curMin := INF
	for k := 0; k <= maxC; k++ {
		if bestCost[k] < curMin {
			curMin = bestCost[k]
		}
		prefixMin[k] = curMin
	}
	out := make([]string, M+1)
	for k := 0; k <= M; k++ {
		if k <= maxC && prefixMin[k] < INF {
			out[k] = strconv.Itoa(totalPairsZeros - prefixMin[k])
		} else {
			out[k] = strconv.Itoa(totalPairsZeros - prefixMin[maxC])
		}
	}
	return strings.Join(out, " ")
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(testcasesRaw, "\n")
	var cases []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid n: %v", idx+1, err)
		}
		if len(fields) != 1+n {
			return nil, fmt.Errorf("line %d: expected %d numbers got %d", idx+1, 1+n, len(fields))
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[1+i])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse a[%d]: %v", idx+1, i, err)
			}
			arr[i] = v
		}
		cases = append(cases, testCase{n: n, arr: arr})
	}
	if len(cases) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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
		var input strings.Builder
		fmt.Fprintf(&input, "%d\n", tc.n)
		for idx, v := range tc.arr {
			if idx > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')

		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Printf("test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		expect := solveCase(tc)
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", i+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
