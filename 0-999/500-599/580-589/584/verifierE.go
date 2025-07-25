package main

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func run(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("time limit")
		}
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func compute(n int, p, s []int) (int64, [][2]int) {
	pp := make([]int, n+1)
	ps := make([]int, n+1)
	for i := 1; i <= n; i++ {
		pp[p[i-1]] = i
		ps[s[i-1]] = i
	}
	ppos := make([]int, n+1)
	for i := 1; i <= n; i++ {
		ppos[i] = ps[p[i-1]]
	}
	var mc int64
	for x := 1; x <= n; x++ {
		diff := int64(pp[x] - ps[x])
		if diff < 0 {
			diff = -diff
		}
		mc += diff
	}
	swaps := make([][2]int, 0)
	for {
		flag := false
		for x := 1; x <= n; x++ {
			if ppos[x] >= x+1 {
				for y := ppos[x]; y > x; y-- {
					if ppos[y] <= x {
						ppos[x], ppos[y] = ppos[y], ppos[x]
						swaps = append(swaps, [2]int{x, y})
						flag = true
						break
					}
				}
			}
		}
		if !flag {
			break
		}
	}
	return mc / 2, swaps
}

func apply(n int, p []int, swaps [][2]int) []int {
	arr := make([]int, n)
	copy(arr, p)
	for _, sw := range swaps {
		i, j := sw[0]-1, sw[1]-1
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

func parseOutput(out string, kExpected int) (int64, [][2]int, error) {
	reader := strings.NewReader(out)
	var cost int64
	var k int
	if _, err := fmt.Fscan(reader, &cost); err != nil {
		return 0, nil, fmt.Errorf("cannot parse cost: %v", err)
	}
	if _, err := fmt.Fscan(reader, &k); err != nil {
		return 0, nil, fmt.Errorf("cannot parse k: %v", err)
	}
	swaps := make([][2]int, k)
	for i := 0; i < k; i++ {
		if _, err := fmt.Fscan(reader, &swaps[i][0], &swaps[i][1]); err != nil {
			return 0, nil, fmt.Errorf("cannot parse swap: %v", err)
		}
	}
	return cost, swaps, nil
}

func permutations(n int, rng *rand.Rand) ([]int, []int) {
	p := rng.Perm(n)
	s := rng.Perm(n)
	for i := range p {
		p[i]++
		s[i]++
	}
	return p, s
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	type tc struct {
		n    int
		p, s []int
	}
	cases := make([]tc, 0, 100)
	cases = append(cases, tc{2, []int{1, 2}, []int{2, 1}})
	for len(cases) < 100 {
		n := rng.Intn(10) + 2
		p, s := permutations(n, rng)
		cases = append(cases, tc{n, p, s})
	}

	for idx, tc := range cases {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i, v := range tc.p {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		for i, v := range tc.s {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expCost, _ := compute(tc.n, tc.p, tc.s)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", idx+1, err, input)
			os.Exit(1)
		}
		cost, swaps, err := parseOutput(out, -1)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d parse error: %v\noutput:%s", idx+1, err, out)
			os.Exit(1)
		}
		var swapCost int64
		for _, sw := range swaps {
			if sw[0] < 1 || sw[0] > tc.n || sw[1] < 1 || sw[1] > tc.n || sw[0] == sw[1] {
				fmt.Fprintf(os.Stderr, "case %d invalid swap %v\n", idx+1, sw)
				os.Exit(1)
			}
			swapCost += int64(math.Abs(float64(sw[0] - sw[1])))
		}
		final := apply(tc.n, tc.p, swaps)
		for i, v := range final {
			if v != tc.s[i] {
				fmt.Fprintf(os.Stderr, "case %d permutation mismatch\n", idx+1)
				os.Exit(1)
			}
		}
		if swapCost != cost {
			fmt.Fprintf(os.Stderr, "case %d cost mismatch reported %d actual %d\n", idx+1, cost, swapCost)
			os.Exit(1)
		}
		if cost != expCost {
			fmt.Fprintf(os.Stderr, "case %d minimal cost mismatch expected %d got %d\n", idx+1, expCost, cost)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
