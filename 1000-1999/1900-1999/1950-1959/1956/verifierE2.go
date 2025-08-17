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

func solveCase(n int, a []int) string {
	b := make([]int, n)
	copy(b, a)
	for i := 0; i < 2200; i++ {
		for j := 0; j < n; j++ {
			nxt := (j + 1) % n
			diff := b[nxt] - b[j]
			if diff < 0 {
				diff = 0
			}
			b[nxt] = diff
		}
	}
	ans := []int{}
	for i := 0; i < n; i++ {
		prev := (i - 1 + n) % n
		if b[prev] != 0 || b[i] == 0 {
			continue
		}
		ans = append(ans, i+1)
		next := (i + 1) % n
		next2 := (i + 2) % n
		if b[next] > 0 && b[next2] > 0 {
			x := b[i]
			y := b[next]
			z := b[next2]
			xs := y / x
			if i == n-1 {
				xs++
			}
			delta := y
			if i != n-1 {
				delta -= x
			}
			delta += y % x
			if int64(delta)*int64(xs)/2 < int64(z) {
				ans = append(ans, next+1)
			}
		}
	}
	sort.Ints(ans)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(ans))
	if len(ans) > 0 {
		for _, v := range ans {
			fmt.Fprintf(&sb, "%d ", v)
		}
	}
	sb.WriteByte('\n')
	return sb.String()
}

func genCase(rng *rand.Rand) (string, string) {
	t := rng.Intn(3) + 1
	var in bytes.Buffer
	var out bytes.Buffer
	fmt.Fprintf(&in, "%d\n", t)
	for i := 0; i < t; i++ {
		n := rng.Intn(5) + 1
		fmt.Fprintf(&in, "%d\n", n)
		a := make([]int, n)
		for j := 0; j < n; j++ {
			a[j] = rng.Intn(11)
			if j > 0 {
				in.WriteByte(' ')
			}
			fmt.Fprintf(&in, "%d", a[j])
		}
		in.WriteByte('\n')
		out.WriteString(solveCase(n, a))
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
		fmt.Println("usage: go run verifierE2.go /path/to/binary")
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
