package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type testCaseF struct {
	n, m, k int
	a       []int64
	d       []int64
	f       []int64
	exp     string
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func existsSumInRange(D, F []int64, L, R int64) bool {
	i, j := 0, len(F)-1
	for i < len(D) && j >= 0 {
		s := D[i] + F[j]
		if s < L {
			i++
		} else if s > R {
			j--
		} else {
			return true
		}
	}
	return false
}

func can(a, D, F []int64, limit int64) bool {
	pos := -1
	for i := 1; i < len(a); i++ {
		gap := a[i] - a[i-1]
		if gap > limit {
			if pos != -1 {
				return false
			}
			pos = i
		}
	}
	if pos == -1 {
		return true
	}
	gap := a[pos] - a[pos-1]
	if gap > 2*limit {
		return false
	}
	left := maxInt64(a[pos-1], a[pos]-limit)
	right := minInt64(a[pos], a[pos-1]+limit)
	return existsSumInRange(D, F, left, right)
}

func solveF(n, m, k int, a, d, f []int64) string {
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	sort.Slice(d, func(i, j int) bool { return d[i] < d[j] })
	sort.Slice(f, func(i, j int) bool { return f[i] < f[j] })
	low, high := int64(0), int64(4_000_000_000)
	for low < high {
		mid := (low + high) / 2
		if can(a, d, f, mid) {
			high = mid
		} else {
			low = mid + 1
		}
	}
	return fmt.Sprint(low)
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateTests() []testCaseF {
	rng := rand.New(rand.NewSource(6))
	cases := make([]testCaseF, 100)
	for i := range cases {
		n := rng.Intn(4) + 2
		m := rng.Intn(4) + 1
		k := rng.Intn(4) + 1
		a := make([]int64, n)
		cur := int64(0)
		for j := 0; j < n; j++ {
			cur += int64(rng.Intn(5) + 1)
			a[j] = cur
		}
		d := make([]int64, m)
		for j := 0; j < m; j++ {
			d[j] = int64(rng.Intn(10) + 1)
		}
		f := make([]int64, k)
		for j := 0; j < k; j++ {
			f[j] = int64(rng.Intn(10) + 1)
		}
		cases[i] = testCaseF{n: n, m: m, k: k, a: a, d: d, f: f, exp: solveF(n, m, k, a, d, f)}
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateTests()
	for i, tc := range cases {
		var sb strings.Builder
		fmt.Fprintln(&sb, 1)
		fmt.Fprintf(&sb, "%d %d %d\n", tc.n, tc.m, tc.k)
		for j := 0; j < tc.n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(tc.a[j]))
		}
		sb.WriteByte('\n')
		for j := 0; j < tc.m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(tc.d[j]))
		}
		sb.WriteByte('\n')
		for j := 0; j < tc.k; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(tc.f[j]))
		}
		sb.WriteByte('\n')
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != tc.exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, tc.exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
