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

func winners(n int, a []int) int {
	for n >= a[0] {
		cnt := 0
		for _, v := range a {
			if v <= n {
				cnt++
			} else {
				break
			}
		}
		if cnt == 0 {
			break
		}
		n -= cnt
	}
	return n
}

func genCase(rng *rand.Rand) (string, string) {
	t := rng.Intn(3) + 1
	var in bytes.Buffer
	var out bytes.Buffer
	fmt.Fprintf(&in, "%d\n", t)
	for ; t > 0; t-- {
		k := rng.Intn(5) + 1
		q := rng.Intn(5) + 1
		fmt.Fprintf(&in, "%d %d\n", k, q)
		a := make([]int, k)
		used := make(map[int]bool)
		for i := 0; i < k; i++ {
			val := rng.Intn(100) + 1
			for used[val] {
				val = rng.Intn(100) + 1
			}
			used[val] = true
			a[i] = val
		}
		// sort a
		for i := 0; i < k; i++ {
			for j := i + 1; j < k; j++ {
				if a[j] < a[i] {
					a[i], a[j] = a[j], a[i]
				}
			}
		}
		for i, v := range a {
			if i > 0 {
				in.WriteByte(' ')
			}
			fmt.Fprintf(&in, "%d", v)
		}
		in.WriteByte('\n')
		ns := make([]int, q)
		for i := 0; i < q; i++ {
			ns[i] = rng.Intn(100) + 1
			if i > 0 {
				in.WriteByte(' ')
			}
			fmt.Fprintf(&in, "%d", ns[i])
		}
		in.WriteByte('\n')
		for i, v := range ns {
			if i > 0 {
				out.WriteByte(' ')
			}
			fmt.Fprintf(&out, "%d", winners(v, a))
		}
		out.WriteByte('\n')
	}
	return in.String(), strings.TrimSpace(out.String())
}

func run(bin, input string) (string, error) {
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, in, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
