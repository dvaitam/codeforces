package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCaseC struct {
	n       int
	m       int
	arr     []int64
	queries [][2]int64
}

func expectedC(tc testCaseC) []int64 {
	n := tc.n
	a := append([]int64(nil), tc.arr...)
	total := int64(n) * int64(n+1) / 2
	for i := 0; i+1 < n; i++ {
		if a[i] != a[i+1] {
			total += int64(i+1) * int64(n-i-1)
		}
	}
	res := make([]int64, tc.m)
	for qi, q := range tc.queries {
		idx := int(q[0]) - 1
		x := q[1]
		if a[idx] != x {
			if idx-1 >= 0 {
				oldDiff := a[idx-1] != a[idx]
				newDiff := a[idx-1] != x
				if oldDiff != newDiff {
					weight := int64(idx) * int64(n-idx)
					if newDiff {
						total += weight
					} else {
						total -= weight
					}
				}
			}
			if idx+1 < n {
				oldDiff := a[idx] != a[idx+1]
				newDiff := x != a[idx+1]
				if oldDiff != newDiff {
					weight := int64(idx+1) * int64(n-idx-1)
					if newDiff {
						total += weight
					} else {
						total -= weight
					}
				}
			}
			a[idx] = x
		}
		res[qi] = total
	}
	return res
}

func genCaseC(rng *rand.Rand) testCaseC {
	n := rng.Intn(8) + 1
	m := rng.Intn(8) + 1
	arr := make([]int64, n)
	for i := range arr {
		arr[i] = int64(rng.Intn(5) + 1)
	}
	queries := make([][2]int64, m)
	for i := 0; i < m; i++ {
		queries[i][0] = int64(rng.Intn(n) + 1)
		queries[i][1] = int64(rng.Intn(5) + 1)
	}
	return testCaseC{n: n, m: m, arr: arr, queries: queries}
}

func runCaseC(bin string, tc testCaseC) error {
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for i, v := range tc.arr {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(fmt.Sprintf("%d", v))
	}
	input.WriteByte('\n')
	for _, q := range tc.queries {
		input.WriteString(fmt.Sprintf("%d %d\n", q[0], q[1]))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != tc.m {
		return fmt.Errorf("expected %d lines got %d", tc.m, len(lines))
	}
	expect := expectedC(tc)
	for i, line := range lines {
		val, err := strconv.ParseInt(strings.TrimSpace(line), 10, 64)
		if err != nil {
			return fmt.Errorf("invalid int %q", line)
		}
		if val != expect[i] {
			return fmt.Errorf("line %d expected %d got %d", i+1, expect[i], val)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		tc := genCaseC(rng)
		if err := runCaseC(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n", t+1, err)
			var inp strings.Builder
			inp.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
			for i, v := range tc.arr {
				if i > 0 {
					inp.WriteByte(' ')
				}
				inp.WriteString(fmt.Sprintf("%d", v))
			}
			inp.WriteByte('\n')
			for _, q := range tc.queries {
				inp.WriteString(fmt.Sprintf("%d %d\n", q[0], q[1]))
			}
			fmt.Fprint(os.Stderr, inp.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
