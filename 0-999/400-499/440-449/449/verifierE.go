package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const MOD = 1000000007

func solve(input string) string {
	rdr := bufio.NewReader(strings.NewReader(input))
	var t int
	fmt.Fscan(rdr, &t)
	ns := make([]int, t)
	ms := make([]int, t)
	maxK := 0
	for i := 0; i < t; i++ {
		fmt.Fscan(rdr, &ns[i], &ms[i])
		k := ns[i]
		if ms[i] < k {
			k = ms[i]
		}
		if k > maxK {
			maxK = k
		}
	}
	phi := make([]int, maxK+1)
	for i := 0; i <= maxK; i++ {
		phi[i] = i
	}
	for i := 2; i <= maxK; i++ {
		if phi[i] == i {
			for j := i; j <= maxK; j += i {
				phi[j] -= phi[j] / i
			}
		}
	}
	f2 := make([]int64, maxK+1)
	for i := 1; i <= maxK; i++ {
		pi := int64(phi[i])
		for k0 := i; k0 <= maxK; k0 += i {
			d := int64(k0 / i)
			f2[k0] += d * pi
		}
	}
	G := make([]int64, maxK+1)
	A := make([]int64, maxK+1)
	B := make([]int64, maxK+1)
	C := make([]int64, maxK+1)
	inv2 := int64((MOD + 1) / 2)
	inv6 := int64(166666668)
	for k := 1; k <= maxK; k++ {
		kk := int64(k)
		t1 := kk * (kk + 1) % MOD * (2*kk + 1) % MOD * inv6 % MOD
		t3 := kk * (kk + 1) % MOD * inv2 % MOD
		sg := (f2[k] + kk) % MOD
		gk := (2*t1%MOD - 4*t3%MOD + 2*sg%MOD) % MOD
		gk = (gk - (kk * kk % MOD) + MOD) % MOD
		G[k] = gk
		A[k] = (A[k-1] + gk) % MOD
		B[k] = (B[k-1] + gk*kk) % MOD
		C[k] = (C[k-1] + gk*kk%MOD*kk%MOD) % MOD
	}
	var out strings.Builder
	idx := 0
	for idx < t {
		n, m := ns[idx], ms[idx]
		K := n
		if m < K {
			K = m
		}
		if K <= 0 {
			fmt.Fprintln(&out, 0)
			idx++
			continue
		}
		np := int64(n + 1)
		mp := int64(m + 1)
		s0 := A[K]
		s1 := B[K]
		s2 := C[K]
		ans := (np*mp%MOD*s0%MOD - (np+mp)%MOD*s1%MOD + s2) % MOD
		if ans < 0 {
			ans += MOD
		}
		fmt.Fprintln(&out, ans)
		idx++
	}
	return strings.TrimSpace(out.String())
}

func generateTests() []test {
	rand.Seed(453)
	var tests []test
	fixed := []string{
		"1\n1 1\n",
		"2\n1 2\n2 1\n",
		"3\n5 5\n10 1\n2 2\n",
	}
	for _, f := range fixed {
		tests = append(tests, test{f, solve(f)})
	}
	for len(tests) < 100 {
		t := rand.Intn(3) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", t)
		for i := 0; i < t; i++ {
			n := rand.Intn(50) + 1
			m := rand.Intn(50) + 1
			fmt.Fprintf(&sb, "%d %d\n", n, m)
		}
		inp := sb.String()
		tests = append(tests, test{inp, solve(inp)})
	}
	return tests
}

type test struct {
	input, expected string
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:%s\nGot:%s\n", i+1, t.input, strings.TrimSpace(t.expected), got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
