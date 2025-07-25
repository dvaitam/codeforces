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

type Item struct{ p, t, l, r int64 }

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func check(A []Item, x float64, T int64) bool {
	var max1, max2 float64
	n := len(A)
	i := 0
	for i < n {
		t1 := float64(A[i].p) * (1.0 - x*float64(A[i].r)/float64(T))
		if max1 > t1+1e-12 {
			return false
		}
		t2 := float64(A[i].p) * (1.0 - x*float64(A[i].l)/float64(T))
		if t2 > max2 {
			max2 = t2
		}
		j := i
		for j+1 < n && A[j+1].p == A[i].p {
			j++
			t1 = float64(A[j].p) * (1.0 - x*float64(A[j].r)/float64(T))
			if max1 > t1+1e-12 {
				return false
			}
			t2 = float64(A[j].p) * (1.0 - x*float64(A[j].l)/float64(T))
			if t2 > max2 {
				max2 = t2
			}
		}
		max1 = max2
		i = j + 1
	}
	return true
}

func solve(n int, P, Tarr []int64) float64 {
	A := make([]Item, n)
	for i := 0; i < n; i++ {
		A[i].p = P[i]
		A[i].t = Tarr[i]
	}
	var T int64
	for i := 0; i < n; i++ {
		T += A[i].t
	}
	sort.Slice(A, func(i, j int) bool { return A[i].t*A[j].p < A[j].t*A[i].p })
	var suml, sumr int64
	for i := 0; i < n; i++ {
		sumr += A[i].t
		j := i
		for j+1 < n && A[j].t*A[j+1].p == A[j+1].t*A[j].p {
			j++
			sumr += A[j].t
		}
		for k := i; k <= j; k++ {
			A[k].l = suml + A[k].t
			A[k].r = sumr
		}
		suml = sumr
		i = j
	}
	sort.Slice(A, func(i, j int) bool { return A[i].p < A[j].p })
	lo, hi := 0.0, 1.0
	for it := 0; it < 50; it++ {
		mid := (lo + hi) * 0.5
		if check(A, mid, T) {
			lo = mid
		} else {
			hi = mid
		}
	}
	return lo
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 2
	P := make([]int64, n)
	T := make([]int64, n)
	for i := 0; i < n; i++ {
		P[i] = int64(rng.Intn(10) + 1)
	}
	for i := 0; i < n; i++ {
		T[i] = int64(rng.Intn(10) + 1)
	}
	input := fmt.Sprintf("%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", P[i])
	}
	input += "\n"
	for i := 0; i < n; i++ {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", T[i])
	}
	input += "\n"
	exp := fmt.Sprintf("%.12f", solve(n, P, T))
	return input, exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
