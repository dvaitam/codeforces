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

const maxV = 100000

var isPrime [maxV]bool
var primesByLen [][]string
var prefixMap [][]map[string][]string

func initData() {
	for i := 2; i < maxV; i++ {
		isPrime[i] = true
	}
	for i := 2; i*i < maxV; i++ {
		if isPrime[i] {
			for j := i * i; j < maxV; j += i {
				isPrime[j] = false
			}
		}
	}
	primesByLen = make([][]string, 6)
	for v := 2; v < maxV; v++ {
		if !isPrime[v] {
			continue
		}
		s := fmt.Sprintf("%05d", v)
		for n := 1; n <= 5; n++ {
			primesByLen[n] = append(primesByLen[n], s[5-n:])
		}
	}
	prefixMap = make([][]map[string][]string, 6)
	for n := 1; n <= 5; n++ {
		prefixMap[n] = make([]map[string][]string, n+1)
		for L := 0; L <= n; L++ {
			prefixMap[n][L] = make(map[string][]string)
		}
		for _, p := range primesByLen[n] {
			for L := 1; L < n; L++ {
				pre := p[:L]
				prefixMap[n][L][pre] = append(prefixMap[n][L][pre], p)
			}
		}
	}
}

type Test struct {
	t      int
	primes []string
}

func randomPrime(n int) string {
	lst := primesByLen[n]
	return strings.TrimLeft(lst[rand.Intn(len(lst))], "0")
}

func generateTest() Test {
	t := rand.Intn(5) + 1
	primes := make([]string, t)
	for i := 0; i < t; i++ {
		n := rand.Intn(5) + 1
		// choose prime without leading zero
		for {
			p := randomPrime(n)
			if len(p) == n {
				primes[i] = p
				break
			}
		}
	}
	return Test{t, primes}
}

func (t Test) Input() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t.t)
	for _, p := range t.primes {
		fmt.Fprintf(&sb, "%s\n", p)
	}
	return sb.String()
}

func countMatrices(pi string) int64 {
	n := len(pi)
	var mat [5][5]byte
	for j := 0; j < n; j++ {
		mat[0][j] = pi[j]
		mat[j][0] = pi[j]
	}
	var dfs func(i int) int64
	dfs = func(i int) int64 {
		if i > n {
			return 1
		}
		idx := i - 1
		pre := make([]byte, idx)
		for j := 0; j < idx; j++ {
			pre[j] = mat[idx][j]
		}
		cand := prefixMap[n][idx][string(pre)]
		var cnt int64
		for _, pstr := range cand {
			for j := idx; j < n; j++ {
				mat[idx][j] = pstr[j]
				if j > idx {
					mat[j][idx] = pstr[j]
				}
			}
			cnt += dfs(i + 1)
		}
		return cnt
	}
	if n == 1 {
		return 1
	}
	return dfs(2)
}

func solve(t Test) []int64 {
	res := make([]int64, t.t)
	for i, p := range t.primes {
		res[i] = countMatrices(p)
	}
	return res
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	initData()
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		t := generateTest()
		inp := t.Input()
		exp := solve(t)
		out, err := runBinary(bin, inp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		outs := strings.Fields(out)
		if len(outs) != len(exp) {
			fmt.Fprintf(os.Stderr, "test %d wrong number of outputs\ninput:\n%s\noutput:\n%s\n", i+1, inp, out)
			os.Exit(1)
		}
		for j, valStr := range outs {
			var val int64
			if _, e := fmt.Sscan(valStr, &val); e != nil {
				fmt.Fprintf(os.Stderr, "test %d parse error: %v\n", i+1, e)
				os.Exit(1)
			}
			if val != exp[j] {
				fmt.Fprintf(os.Stderr, "test %d failed: expected %d got %d at line %d\ninput:\n%s\n", i+1, exp[j], val, j+1, inp)
				os.Exit(1)
			}
		}
	}
	fmt.Println("all tests passed")
}
