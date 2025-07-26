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

func minMoves(a, l, r []int) int {
	n := len(a)
	type state struct {
		arr []int
		d   int
	}
	start := make([]int, n)
	copy(start, a)
	key := func(x []int) string { return fmt.Sprint(x) }
	vis := map[string]bool{key(start): true}
	q := []state{{start, 0}}
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		ok := true
		for i := 0; i < n; i++ {
			if cur.arr[i] < l[i] || cur.arr[i] > r[i] {
				ok = false
				break
			}
		}
		if ok {
			return cur.d
		}
		for i := 0; i < n; i++ {
			if cur.arr[i] == 0 {
				continue
			}
			for _, j := range []int{(i + n - 1) % n, (i + 1) % n} {
				nxt := append([]int(nil), cur.arr...)
				nxt[i]--
				nxt[j]++
				k := key(nxt)
				if !vis[k] {
					vis[k] = true
					q = append(q, state{nxt, cur.d + 1})
				}
			}
		}
	}
	return -1
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 3
	a := make([]int, n)
	l := make([]int, n)
	r := make([]int, n)
	total := 0
	for i := 0; i < n; i++ {
		l[i] = rng.Intn(3)
		r[i] = l[i] + rng.Intn(3)
		a[i] = l[i] + rng.Intn(r[i]-l[i]+1)
		total += a[i]
	}
	// ensure solvable
	sumL, sumR := 0, 0
	for i := 0; i < n; i++ {
		sumL += l[i]
		sumR += r[i]
	}
	if total < sumL {
		a[0] += sumL - total
		total = sumL
	}
	if total > sumR {
		a[0] -= total - sumR
		total = sumR
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", a[i], l[i], r[i]))
	}
	ans := minMoves(a, l, r)
	return sb.String(), fmt.Sprintf("%d", ans)
}

func run(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, want := genCase(rng)
		got, err := run(cand, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
