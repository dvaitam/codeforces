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

func oracleSolve(input string) string {
	words := strings.Fields(input)
	idx := 0
	nextInt := func() int {
		v := 0
		s := words[idx]
		idx++
		for _, ch := range s {
			v = v*10 + int(ch-'0')
		}
		return v
	}
	n := nextInt()
	k := nextInt()
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = nextInt()
	}
	sort.Ints(a)

	mod := 998244353
	ans := 0

	prevDp := make([]int, n)
	currDp := make([]int, n)

	if k < 2 || n < k {
		return "0"
	}

	maxX := (a[n-1] - a[0]) / (k - 1)

	for x := 1; x <= maxX; x++ {
		for i := 0; i < n; i++ {
			prevDp[i] = 1
		}

		for j := 2; j <= k; j++ {
			sum := 0
			p := 0
			for i := 0; i < n; i++ {
				for p < i && a[i]-a[p] >= x {
					sum += prevDp[p]
					if sum >= mod {
						sum -= mod
					}
					p++
				}
				currDp[i] = sum
			}
			prevDp, currDp = currDp, prevDp
		}

		sumK := 0
		for i := 0; i < n; i++ {
			sumK += prevDp[i]
			if sumK >= mod {
				sumK -= mod
			}
		}

		if sumK == 0 {
			break
		}
		ans += sumK
		if ans >= mod {
			ans -= mod
		}
	}

	return fmt.Sprint(ans)
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(8) + 2
	k := rng.Intn(n-1) + 2
	if k > n {
		k = n
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(100)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, k)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", arr[i])
	}
	sb.WriteByte('\n')
	return sb.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := generateCase(rng)
		exp := oracleSolve(in)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed:\nexpected:\n%s\n\ngot:\n%s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
