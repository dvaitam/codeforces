package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Test struct {
	n int
	l *big.Int
	r *big.Int
}

func generateTests() []Test {
	rand.Seed(5)
	tests := make([]Test, 0, 100)
	edge := []Test{{1, big.NewInt(0), big.NewInt(1)}, {5, big.NewInt(3), big.NewInt(31)}, {6, big.NewInt(10), big.NewInt(40)}}
	tests = append(tests, edge...)
	for len(tests) < 100 {
		n := rand.Intn(20) + 1
		maxVal := new(big.Int).Lsh(big.NewInt(1), uint(n))
		r := new(big.Int).Rand(rand.New(rand.NewSource(int64(rand.Int()))), maxVal)
		l := new(big.Int).Rand(rand.New(rand.NewSource(int64(rand.Int()))), maxVal)
		if l.Cmp(r) > 0 {
			l, r = r, l
		}
		tests = append(tests, Test{n, l, r})
	}
	return tests
}

func solve(l, r *big.Int) string {
	k := r.BitLen() - 1
	candidate := new(big.Int)
	found := false
	for ; k >= 0; k-- {
		two := new(big.Int).Lsh(big.NewInt(1), uint(k))
		if r.Cmp(two) < 0 {
			continue
		}
		minusOne := new(big.Int).Sub(two, big.NewInt(1))
		if l.Cmp(minusOne) <= 0 {
			candidate.Lsh(big.NewInt(1), uint(k+1))
			candidate.Sub(candidate, big.NewInt(1))
			found = true
			break
		}
	}
	if !found {
		temp := new(big.Int).Sub(r, big.NewInt(1))
		candidate.Or(r, temp)
	}
	if r.Cmp(candidate) > 0 {
		return r.Text(2)
	}
	return candidate.Text(2)
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
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	binary := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		var in strings.Builder
		fmt.Fprintf(&in, "%d\n%s\n%s\n", t.n, t.l.Text(2), t.r.Text(2))
		expect := solve(t.l, t.r)
		got, err := run(binary, in.String())
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\noutput:\n%s", i+1, err, got)
			os.Exit(1)
		}
		got = strings.TrimSpace(got)
		if got != expect {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:\n%s\ngot:\n%s", i+1, in.String(), expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
	time.Sleep(0)
}
