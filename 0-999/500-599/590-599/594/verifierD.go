package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const mod int64 = 1000000007

func phi(x int64) int64 {
	if x == 0 {
		return 0
	}
	res := x
	i := int64(2)
	for i*i <= x {
		if x%i == 0 {
			for x%i == 0 {
				x /= i
			}
			res = res / i * (i - 1)
		}
		i++
	}
	if x > 1 {
		res = res / x * (x - 1)
	}
	return res
}

func solveD(a []int, queries [][2]int) []int64 {
	res := make([]int64, len(queries))
	for qi, q := range queries {
		l, r := q[0]-1, q[1]-1
		prod := int64(1)
		factors := make(map[int64]int)
		for i := l; i <= r; i++ {
			x := a[i]
			v := int64(x)
			prod = (prod * v) % mod
			xx := x
			for p := 2; p*p <= xx; p++ {
				if xx%p == 0 {
					c := 0
					for xx%p == 0 {
						xx /= p
						c++
					}
					factors[int64(p)] += c
				}
			}
			if xx > 1 {
				factors[int64(xx)]++
			}
		}
		phiVal := int64(1)
		for p, e := range factors {
			pe := int64(1)
			for i := 0; i < e-1; i++ {
				pe = (pe * p) % mod
			}
			phiVal = phiVal * ((pe * (p - 1)) % mod) % mod
		}
		res[qi] = phiVal % mod
	}
	return res
}

func run(binary, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(4)
	for t := 0; t < 100; t++ {
		n := rand.Intn(5) + 1
		a := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = rand.Intn(10) + 1
		}
		qnum := rand.Intn(5) + 1
		queries := make([][2]int, qnum)
		for i := 0; i < qnum; i++ {
			l := rand.Intn(n) + 1
			r := rand.Intn(n-l+1) + l
			queries[i] = [2]int{l, r}
		}
		input := fmt.Sprintf("%d\n", n)
		for i, v := range a {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", v)
		}
		input += "\n"
		input += fmt.Sprintf("%d\n", qnum)
		for i, q := range queries {
			input += fmt.Sprintf("%d %d", q[0], q[1])
			if i+1 < qnum {
				input += "\n"
			} else {
				input += "\n"
			}
		}
		expectedVals := solveD(a, queries)
		expectedParts := make([]string, len(expectedVals))
		for i, v := range expectedVals {
			expectedParts[i] = fmt.Sprintf("%d", v%mod)
		}
		expected := strings.Join(expectedParts, "\n")
		got, err := run(bin, input)
		if err != nil {
			fmt.Println("test", t, "runtime error:", err)
			fmt.Println("output:\n" + got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Println("test", t, "failed")
			fmt.Println("input:\n" + input)
			fmt.Println("expected:\n" + expected)
			fmt.Println("got:\n" + got)
			os.Exit(1)
		}
	}
	fmt.Println("OK")
}
