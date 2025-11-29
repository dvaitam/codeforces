package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `5 3 3 5 1 15 8 2 6
1 3 4 2 13
5 1 5 2 1 7 14 9 6
4 2 1 2 20 20 15 5
2 1 1 2 7 6
2 3 3 2 18 7
2 2 4 3 1 12
4 2 2 3 3 11 10 20
5 1 5 3 3 10 12 10 16
3 2 4 4 6 2 9
1 3 4 1 18
4 3 4 5 1 15 2 6
5 2 1 2 15 12 17 12 17
3 4 1 5 12 10 2
4 1 2 3 17 20 12 5
3 3 5 1 10 11 10
2 1 2 3 16 6
1 1 5 5 13
1 2 5 3 9
4 4 2 1 2 16 11 7
2 5 2 4 4 6
4 3 2 1 14 10 5 15
5 2 5 4 16 11 16 9 10
4 4 2 1 13 18 6 16
3 2 1 4 9 17 18
5 3 1 3 19 2 10 12 18
3 4 3 3 11 6 19
1 4 5 3 11
3 4 3 5 12 12 9
3 4 3 2 15 12 11
5 2 5 2 7 12 16 10 3
4 2 5 5 17 14 10 20
5 3 1 2 6 19 15 20 6
2 2 1 4 8 6
1 2 1 3 6
4 2 5 1 14 15 12 13
5 1 5 2 8 12 1 12 13
3 4 1 5 12 2 18
5 3 1 3 18 17 11 19 10
3 2 4 4 19 18 12
4 2 2 5 13 19 16 7
2 5 1 3 1 13
1 3 5 5 18
2 3 5 4 14 14
2 4 3 4 13 13
2 5 5 3 10 16
3 4 1 3 10 16 10
2 4 1 1 20 15
2 3 1 2 13 1
4 5 5 3 8 16 2 8
4 3 2 3 10 16 20 16
5 5 1 1 5 10 10 18 11
5 3 5 1 15 12 12 19 5
1 1 3 5 15
1 5 2 1 14
4 5 5 4 13 16 13 7
3 4 1 3 1 14 19
3 4 3 2 6 16 18
4 3 5 2 14 19 18 2
1 2 3 1 3
1 3 4 1 13
4 1 1 1 8 20 4 5
3 4 2 2 20 6 14
2 1 5 2 2 18
1 4 1 3 2
5 5 1 4 20 5 1 14 3
3 5 4 4 12 12 2
2 3 2 5 17 10
5 5 5 2 9 3 18 8 9
3 5 2 2 12 15 13
2 2 1 3 3 19
1 1 1 5 20
4 2 4 4 16 11 4 17
1 5 4 1 5
4 2 1 1 16 7 5 20
4 3 2 3 11 20 12 13
4 2 3 3 14 12 17 2
5 5 2 2 13 3 4 2 2
2 2 2 1 16 16
3 1 4 4 10 20 14
3 4 4 1 7 5 6
1 3 4 4 5
5 3 1 3 6 10 8 2 16
1 3 3 3 2
1 4 4 2 4
3 3 4 2 6 2 7
1 5 2 1 20
4 3 3 2 15 12 10 3
4 2 2 2 7 2 20 14
3 1 4 1 15 14 6
1 1 5 5 19
3 1 1 2 16 3 16
1 2 1 4 13
1 1 3 4 15
3 1 1 2 6 17 13
2 5 2 1 11 4
1 5 2 5 3
2 5 1 4 18 12
4 4 5 4 17 5 1 10
4 2 5 1 5 10 16 3
`

// solveCase mirrors 1312F.go logic and returns winning move count.
func solveCase(n, x, y, z int, a []int64) int {
	const H = 500
	var g [H + 1][3]int
	for h := 1; h <= H; h++ {
		for t := 0; t < 3; t++ {
			var used [4]bool
			h2 := h - x
			if h2 < 0 {
				h2 = 0
			}
			used[g[h2][0]] = true
			if t != 1 {
				h2 = h - y
				if h2 < 0 {
					h2 = 0
				}
				used[g[h2][1]] = true
			}
			if t != 2 {
				h2 = h - z
				if h2 < 0 {
					h2 = 0
				}
				used[g[h2][2]] = true
			}
			mex := 0
			for used[mex] {
				mex++
			}
			g[h][t] = mex
		}
	}

	base := 200
	period := -1
	for p := 1; p <= 100; p++ {
		ok := true
		for h := base; h+p <= H; h++ {
			for t := 0; t < 3; t++ {
				if g[h][t] != g[h+p][t] {
					ok = false
					break
				}
			}
			if !ok {
				break
			}
		}
		if ok {
			period = p
			break
		}
	}
	if period == -1 {
		period = 1
	}

	getg := func(h0 int64, t int) int {
		if h0 <= H {
			return g[h0][t]
		}
		if h0 < int64(base) {
			return g[h0][t]
		}
		idx := base + int((h0-int64(base))%int64(period))
		return g[idx][t]
	}

	xor := 0
	for _, v := range a {
		xor ^= getg(v, 0)
	}
	if xor == 0 {
		return 0
	}
	cnt := 0
	for _, ai := range a {
		gi := getg(ai, 0)
		h2 := ai - int64(x)
		if h2 < 0 {
			h2 = 0
		}
		if xor^gi^getg(h2, 0) == 0 {
			cnt++
		}
		h2 = ai - int64(y)
		if h2 < 0 {
			h2 = 0
		}
		if xor^gi^getg(h2, 1) == 0 {
			cnt++
		}
		h2 = ai - int64(z)
		if h2 < 0 {
			h2 = 0
		}
		if xor^gi^getg(h2, 2) == 0 {
			cnt++
		}
	}
	return cnt
}

type testCase struct {
	n int
	x int
	y int
	z int
	a []int64
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(testcasesRaw, "\n")
	var tests []testCase
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 4 {
			return nil, fmt.Errorf("invalid line")
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}
		if len(parts) != 4+n {
			return nil, fmt.Errorf("line expects %d numbers got %d", 4+n, len(parts))
		}
		x, _ := strconv.Atoi(parts[1])
		y, _ := strconv.Atoi(parts[2])
		z, _ := strconv.Atoi(parts[3])
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			v, err := strconv.ParseInt(parts[4+i], 10, 64)
			if err != nil {
				return nil, err
			}
			arr[i] = v
		}
		tests = append(tests, testCase{n: n, x: x, y: y, z: z, a: arr})
	}
	return tests, nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d %d %d\n", tc.n, tc.x, tc.y, tc.z)
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
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
	expected := strconv.Itoa(solveCase(tc.n, tc.x, tc.y, tc.z, tc.a))
	got, err := runBinary(bin, input)
	if err != nil {
		return err
	}
	if strings.TrimSpace(got) != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
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
