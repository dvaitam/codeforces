package main

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func extend(states [][2]int, x int) [][2]int {
	next := make([][2]int, 0, len(states)*2)
	for _, st := range states {
		if x > st[0] {
			next = append(next, [2]int{x, st[1]})
		}
		if x < st[1] {
			next = append(next, [2]int{st[0], x})
		}
	}
	res := make([][2]int, 0, 2)
	for _, st := range next {
		dominated := false
		for i := 0; i < len(res); {
			ot := res[i]
			if st[0] >= ot[0] && st[1] <= ot[1] {
				dominated = true
				break
			}
			if st[0] <= ot[0] && st[1] >= ot[1] {
				res[i] = res[len(res)-1]
				res = res[:len(res)-1]
			} else {
				i++
			}
		}
		if !dominated {
			res = append(res, st)
		}
	}
	return res
}

func isDecinc(arr []int) bool {
	states := [][2]int{{math.MinInt64, math.MaxInt64}}
	for _, x := range arr {
		states = extend(states, x)
		if len(states) == 0 {
			return false
		}
	}
	return true
}

func solveD(p []int) int64 {
	n := len(p)
	var ans int64
	for l := 0; l < n; l++ {
		states := [][2]int{{math.MinInt64, math.MaxInt64}}
		for r := l; r < n; r++ {
			states = extend(states, p[r])
			if len(states) == 0 {
				break
			}
			ans++
		}
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(45)
	const t = 100
	var sb strings.Builder
	var exp strings.Builder
	for i := 0; i < t; i++ {
		n := rand.Intn(6) + 1
		perm := rand.Perm(n)
		for j := range perm {
			perm[j]++
		}
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j, v := range perm {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		exp.WriteString(fmt.Sprintf("%d\n", solveD(perm)))
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error running binary: %v\noutput:\n%s", err, out.String())
		os.Exit(1)
	}
	got := strings.TrimSpace(out.String())
	want := strings.TrimSpace(exp.String())
	if got != want {
		fmt.Fprintf(os.Stderr, "wrong answer\nexpected:\n%s\ngot:\n%s\n", want, got)
		os.Exit(1)
	}
	fmt.Println("all tests passed")
}
