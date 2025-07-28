package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Query struct {
	idx int
	x   int
}

type Test struct {
	n  int
	q  int
	a  []int
	qs []Query
}

func generateTests() []Test {
	rand.Seed(4)
	tests := make([]Test, 0, 100)
	for len(tests) < 100 {
		n := rand.Intn(10) + 1
		q := rand.Intn(10) + 1
		a := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = rand.Intn(20) + 1
		}
		qs := make([]Query, q)
		for i := 0; i < q; i++ {
			qs[i] = Query{rand.Intn(n) + 1, rand.Intn(20) + 1}
		}
		tests = append(tests, Test{n, q, a, qs})
	}
	return tests
}

const MOD int = 1000000007
const maxVal int = 200000

var spf [maxVal + 1]int

func initSieve() {
	for i := 2; i <= maxVal; i++ {
		if spf[i] == 0 {
			for j := i; j <= maxVal; j += i {
				if spf[j] == 0 {
					spf[j] = i
				}
			}
		}
	}
}

func factorize(x int) map[int]int {
	res := make(map[int]int)
	for x > 1 {
		p := spf[x]
		if p == 0 {
			p = x
		}
		cnt := 0
		for x%p == 0 {
			x /= p
			cnt++
		}
		res[p] += cnt
	}
	return res
}

type primeData struct {
	indexExp map[int]int
	cnt      map[int]int
	minExp   int
	missing  int
}

func modPow(a, e int) int {
	res := 1
	b := a % MOD
	for e > 0 {
		if e&1 == 1 {
			res = res * b % MOD
		}
		b = b * b % MOD
		e >>= 1
	}
	return res
}

func minKey(m map[int]int) int {
	min := -1
	for k := range m {
		if min == -1 || k < min {
			min = k
		}
	}
	if min == -1 {
		return 0
	}
	return min
}

func solve(test Test) string {
	primes := make(map[int]*primeData)
	n := test.n
	curGCD := 1
	for i, val := range test.a {
		factors := factorize(val)
		for p, e := range factors {
			pd := primes[p]
			if pd == nil {
				pd = &primeData{indexExp: make(map[int]int), cnt: make(map[int]int), missing: n}
				primes[p] = pd
			}
			pd.indexExp[i] = e
			pd.cnt[e]++
			pd.missing--
		}
	}
	for p, pd := range primes {
		if pd.missing == 0 {
			pd.minExp = minKey(pd.cnt)
			curGCD = curGCD * modPow(p, pd.minExp) % MOD
		}
	}
	var b strings.Builder
	for _, q := range test.qs {
		idx := q.idx - 1
		factors := factorize(q.x)
		for p, add := range factors {
			pd := primes[p]
			if pd == nil {
				pd = &primeData{indexExp: make(map[int]int), cnt: make(map[int]int), missing: n}
				primes[p] = pd
			}
			oldExp := pd.indexExp[idx]
			if oldExp == 0 {
				pd.missing--
			} else {
				pd.cnt[oldExp]--
				if pd.cnt[oldExp] == 0 {
					delete(pd.cnt, oldExp)
				}
			}
			newExp := oldExp + add
			pd.indexExp[idx] = newExp
			pd.cnt[newExp]++
			oldMin := pd.minExp
			var newMin int
			if pd.missing > 0 {
				newMin = 0
			} else {
				if (oldExp == oldMin && pd.cnt[oldExp] == 0) || oldMin == 0 {
					newMin = minKey(pd.cnt)
				} else {
					newMin = oldMin
				}
			}
			if newMin > oldMin {
				curGCD = curGCD * modPow(p, newMin-oldMin) % MOD
				pd.minExp = newMin
			} else {
				pd.minExp = newMin
			}
		}
		fmt.Fprintf(&b, "%d\n", curGCD)
	}
	return b.String()
}

func run(binary string, input string) (string, error) {
	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	initSieve()
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	binary := os.Args[1]
	tests := generateTests()
	for i, test := range tests {
		var in strings.Builder
		fmt.Fprintf(&in, "%d %d\n", test.n, test.q)
		for _, v := range test.a {
			fmt.Fprintf(&in, "%d ", v)
		}
		in.WriteByte('\n')
		for _, q := range test.qs {
			fmt.Fprintf(&in, "%d %d\n", q.idx, q.x)
		}
		expect := solve(test)
		got, err := run(binary, in.String())
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\noutput:\n%s", i+1, err, got)
			os.Exit(1)
		}
		got = strings.ReplaceAll(strings.TrimSpace(got), "\r\n", "\n")
		expect = strings.ReplaceAll(strings.TrimSpace(expect), "\r\n", "\n")
		if got != expect {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:\n%s\ngot:\n%s", i+1, in.String(), expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
	time.Sleep(0)
}
