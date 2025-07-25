package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type pair struct{ first, second int }

type TestE struct {
	n   int
	arr []int
}

func (t TestE) Input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t.n))
	for i, v := range t.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func maxPair(a, b pair) pair {
	if a.first > b.first {
		return a
	}
	return b
}

func expectedE(t TestE) string {
	n := t.n
	a := make([]pair, n+5)
	for i := 1; i <= n; i++ {
		a[i] = pair{-t.arr[i-1], i}
	}
	arr := a[1 : n+1]
	sort.Slice(arr, func(i, j int) bool { return arr[i].first < arr[j].first })
	const K = 12
	f := make([][]pair, K)
	for k := 0; k < K; k++ {
		f[k] = make([]pair, n+5)
	}
	for i := 1; i <= n; i++ {
		f[0][i] = pair{a[i+1].first - a[i].first, i}
	}
	for k := 1; k < K; k++ {
		shift := 1 << (k - 1)
		for i := 1; i <= n; i++ {
			if i+shift <= n {
				f[k][i] = maxPair(f[k-1][i], f[k-1][i+shift])
			} else {
				f[k][i] = f[k-1][i]
			}
		}
	}
	lg := make([]int, n+5)
	for i := 3; i <= n; i++ {
		lg[i] = lg[(i+1)>>1] + 1
	}
	rmq := func(l, r int) pair {
		length := r - l + 1
		k := lg[length]
		b1 := f[k][l]
		b2 := f[k][r-(1<<k)+1]
		if b1.first > b2.first {
			return b1
		}
		return b2
	}
	const INF = 1000000000
	d1, d2, d3 := -INF, -INF, -INF
	p1, p2, p3 := 0, 0, 0
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			len1 := j - i
			if i > 2*len1 || len1 > 2*i {
				continue
			}
			m := i
			if len1 > m {
				m = len1
			}
			l := (m + 1) >> 1
			r := m * 2
			rem := n - j
			if r > rem {
				r = rem
			}
			if l <= r && j+l <= n {
				x := rmq(j+l, j+r)
				d1cand := a[i+1].first - a[i].first
				d2cand := a[j+1].first - a[j].first
				if d1cand > d1 || (d1cand == d1 && (d2cand > d2 || (d2cand == d2 && x.first > d3))) {
					d1 = d1cand
					p1 = i
					d2 = d2cand
					p2 = j
					d3 = x.first
					p3 = x.second
				}
			}
		}
	}
	ans := make([]int, n+1)
	for i := 1; i <= p1; i++ {
		ans[a[i].second] = 1
	}
	for i := p1 + 1; i <= p2; i++ {
		ans[a[i].second] = 2
	}
	for i := p2 + 1; i <= p3; i++ {
		ans[a[i].second] = 3
	}
	for i := p3 + 1; i <= n; i++ {
		ans[a[i].second] = -1
	}
	var sb strings.Builder
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(ans[i]))
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func genTests() []TestE {
	rand.Seed(5)
	tests := make([]TestE, 0, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(10) + 3
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rand.Intn(5000) + 1
		}
		tests = append(tests, TestE{n: n, arr: arr})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		exp := strings.TrimSpace(expectedE(tc))
		gotRaw, err := run(bin, tc.Input())
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n%s", i+1, err, gotRaw)
			os.Exit(1)
		}
		got := strings.TrimSpace(gotRaw)
		if got != exp {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, tc.Input(), exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
