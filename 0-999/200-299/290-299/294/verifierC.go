package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

const modC = 1000000007

func modpow(a, e int64) int64 {
	res := int64(1)
	a %= modC
	for e > 0 {
		if e&1 == 1 {
			res = res * a % modC
		}
		a = a * a % modC
		e >>= 1
	}
	return res
}

func solveC(n int, pos []int) int64 {
	sort.Ints(pos)
	m := len(pos)
	totalOff := n - m
	fact := make([]int64, n+1)
	invFact := make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % modC
	}
	invFact[n] = modpow(fact[n], modC-2)
	for i := n; i >= 1; i-- {
		invFact[i-1] = invFact[i] * int64(i) % modC
	}
	ans := fact[totalOff]
	first := pos[0] - 1
	ans = ans * invFact[first] % modC
	for i := 1; i < m; i++ {
		gap := pos[i] - pos[i-1] - 1
		ans = ans * invFact[gap] % modC
		if gap > 0 {
			ans = ans * modpow(2, int64(gap-1)) % modC
		}
	}
	last := n - pos[m-1]
	ans = ans * invFact[last] % modC
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	m := rng.Intn(n) + 1
	posMap := make(map[int]struct{})
	positions := make([]int, 0, m)
	for len(positions) < m {
		p := rng.Intn(n) + 1
		if _, ok := posMap[p]; !ok {
			posMap[p] = struct{}{}
			positions = append(positions, p)
		}
	}
	ans := solveC(n, append([]int(nil), positions...))
	var in bytes.Buffer
	fmt.Fprintf(&in, "%d %d\n", n, m)
	for i, p := range positions {
		if i > 0 {
			in.WriteByte(' ')
		}
		fmt.Fprintf(&in, "%d", p)
	}
	in.WriteByte('\n')
	out := fmt.Sprintf("%d\n", ans)
	return in.String(), out
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := generateCase(rng)
		got, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("test %d: execution failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, in, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
