package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func FInt(L, R int, v, primes []int) int {
	size := R - L
	ok := make([]bool, size)
	primeArr := make([]bool, size)
	for i := 0; i < size; i++ {
		primeArr[i] = true
	}
	V := len(v)
	for _, vi := range v {
		j0 := sort.Search(V, func(j int) bool { return vi+v[j] >= L })
		for j := j0; j < V; j++ {
			x := vi + v[j]
			if x >= R {
				break
			}
			ok[x-L] = true
		}
	}
	for _, p := range primes {
		x := (L / p) * p
		if x < L || x <= p {
			x += p
		}
		for x < R {
			primeArr[x-L] = false
			x += p
		}
	}
	cnt := 0
	for i := 0; i < size; i++ {
		if ok[i] && primeArr[i] {
			cnt++
		}
	}
	return cnt
}

func solveC(L, R int) int {
	v := make([]int, 0)
	for i := 1; i*i <= R; i++ {
		v = append(v, i*i)
	}
	NN := int(math.Sqrt(float64(R))) + 1
	isP := make([]bool, NN+1)
	primes := make([]int, 0)
	for i := 2; i <= NN; i++ {
		isP[i] = true
	}
	for i := 2; i <= NN; i++ {
		if isP[i] {
			primes = append(primes, i)
			for j := i * 2; j <= NN; j += i {
				isP[j] = false
			}
		}
	}
	return FInt(L, R+1, v, primes)
}

func generateCase(rng *rand.Rand) (string, string) {
	l := rng.Intn(5000) + 1
	r := l + rng.Intn(200)
	input := fmt.Sprintf("%d %d\n", l, r)
	exp := fmt.Sprintf("%d", solveC(l, r))
	return input, exp
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	edge := [][2]int{{1, 1}, {2, 5}, {10, 10}, {100, 150}, {2000, 2200}}
	for i, p := range edge {
		input := fmt.Sprintf("%d %d\n", p[0], p[1])
		if err := runCase(exe, input, fmt.Sprintf("%d", solveC(p[0], p[1]))); err != nil {
			fmt.Fprintf(os.Stderr, "edge case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
	}
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
