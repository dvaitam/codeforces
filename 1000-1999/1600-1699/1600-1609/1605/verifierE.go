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

type prefixHelper struct {
	arr  []int64
	pref []int64
}

func newPrefix(arr []int64) *prefixHelper {
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	pref := make([]int64, len(arr)+1)
	for i, v := range arr {
		pref[i+1] = pref[i] + v
	}
	return &prefixHelper{arr, pref}
}

func (p *prefixHelper) sumPlus(t int64) int64 {
	if len(p.arr) == 0 {
		return 0
	}
	idx := sort.Search(len(p.arr), func(i int) bool { return p.arr[i] >= -t })
	total := p.pref[len(p.arr)]
	left := p.pref[idx]
	return (total - 2*left) + int64(len(p.arr)-2*idx)*t
}

func (p *prefixHelper) sumMinus(t int64) int64 {
	if len(p.arr) == 0 {
		return 0
	}
	idx := sort.Search(len(p.arr), func(i int) bool { return p.arr[i] >= t })
	total := p.pref[len(p.arr)]
	left := p.pref[idx]
	return (total - 2*left) + int64(2*idx-len(p.arr))*t
}

func mobius(n int) []int {
	mu := make([]int, n+1)
	prime := make([]int, 0)
	isComp := make([]bool, n+1)
	mu[1] = 1
	for i := 2; i <= n; i++ {
		if !isComp[i] {
			prime = append(prime, i)
			mu[i] = -1
		}
		for _, p := range prime {
			if i*p > n {
				break
			}
			isComp[i*p] = true
			if i%p == 0 {
				mu[i*p] = 0
				break
			} else {
				mu[i*p] = -mu[i]
			}
		}
	}
	return mu
}

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func solveE(n int, a, b []int64, queries []int64) []int64 {
	mu := mobius(n)
	diffBase := make([]int64, n+1)
	for i := 2; i <= n; i++ {
		diffBase[i] = b[i] - a[i]
	}
	base := make([]int64, n+1)
	for d := 1; d <= n; d++ {
		if mu[d] == 0 {
			continue
		}
		for m := d; m <= n; m += d {
			base[m] += int64(mu[d]) * diffBase[m/d]
		}
	}
	var plusArr, minusArr []int64
	constArrSum := int64(0)
	for i := 1; i <= n; i++ {
		if mu[i] == 1 {
			plusArr = append(plusArr, base[i])
		} else if mu[i] == -1 {
			minusArr = append(minusArr, base[i])
		} else {
			constArrSum += abs64(base[i])
		}
	}
	plusHelper := newPrefix(plusArr)
	minusHelper := newPrefix(minusArr)
	ans := make([]int64, len(queries))
	for idx, x := range queries {
		t := x - a[1]
		ans[idx] = constArrSum + plusHelper.sumPlus(t) + minusHelper.sumMinus(t)
	}
	return ans
}

type testCaseE struct {
	n  int
	a  []int64
	b  []int64
	qs []int64
}

func generateTestsE() []testCaseE {
	r := rand.New(rand.NewSource(1))
	tests := []testCaseE{}
	for len(tests) < 120 {
		n := r.Intn(10) + 1
		a := make([]int64, n+1)
		b := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			a[i] = int64(r.Intn(20))
			b[i] = int64(r.Intn(20))
		}
		q := r.Intn(5) + 1
		qs := make([]int64, q)
		for i := 0; i < q; i++ {
			qs[i] = int64(r.Intn(40))
		}
		tests = append(tests, testCaseE{n, a, b, qs})
	}
	return tests
}

func runCandidate(bin, input string) (string, error) {
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTestsE()
	for i, tc := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i := 1; i <= tc.n; i++ {
			if i > 1 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", tc.a[i]))
		}
		sb.WriteByte('\n')
		for i := 1; i <= tc.n; i++ {
			if i > 1 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", tc.b[i]))
		}
		sb.WriteByte('\n')
		sb.WriteString(fmt.Sprintf("%d\n", len(tc.qs)))
		for i, v := range tc.qs {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		input := sb.String()
		exp := solveE(tc.n, tc.a, tc.b, tc.qs)
		gotStr, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		fields := strings.Fields(gotStr)
		if len(fields) != len(exp) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d numbers got %d\ninput:\n%s\n", i+1, len(exp), len(fields), input)
			os.Exit(1)
		}
		for j, f := range fields {
			var val int64
			if _, err := fmt.Sscan(f, &val); err != nil || val != exp[j] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%s\nexpected:%v\ngot:%v\n", i+1, input, exp, fields)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
