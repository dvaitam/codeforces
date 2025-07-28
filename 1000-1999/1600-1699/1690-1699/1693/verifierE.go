package main

import (
	"bytes"
	"container/list"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type state struct {
	arr  []int
	dist int
}

func allZero(a []int) bool {
	for _, v := range a[1 : len(a)-1] {
		if v != 0 {
			return false
		}
	}
	return true
}

func serialize(a []int) string {
	var sb strings.Builder
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	return sb.String()
}

func minOps(a []int) int {
	start := append([]int(nil), a...)
	q := list.New()
	q.PushBack(state{start, 0})
	seen := map[string]bool{}
	for q.Len() > 0 {
		e := q.Remove(q.Front()).(state)
		key := serialize(e.arr)
		if seen[key] {
			continue
		}
		seen[key] = true
		if allZero(e.arr) {
			return e.dist
		}
		maxVal := -1
		for _, v := range e.arr {
			if v > maxVal {
				maxVal = v
			}
		}
		left := -1
		for i, v := range e.arr {
			if v == maxVal {
				left = i
				break
			}
		}
		if left > 0 && left < len(e.arr)-1 {
			nb := append([]int(nil), e.arr...)
			m := 0
			for j := 0; j < left; j++ {
				if nb[j] > m {
					m = nb[j]
				}
			}
			nb[left] = m
			q.PushBack(state{nb, e.dist + 1})
		}
		right := -1
		for i := len(e.arr) - 1; i >= 0; i-- {
			if e.arr[i] == maxVal {
				right = i
				break
			}
		}
		if right > 0 && right < len(e.arr)-1 {
			nb := append([]int(nil), e.arr...)
			m := 0
			for j := right + 1; j < len(nb); j++ {
				if nb[j] > m {
					m = nb[j]
				}
			}
			nb[right] = m
			q.PushBack(state{nb, e.dist + 1})
		}
	}
	return -1
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(46)
	const t = 100
	var sb strings.Builder
	var exp strings.Builder
	for i := 0; i < t; i++ {
		n := rand.Intn(3) + 1
		a := make([]int, n+2)
		for j := 1; j <= n; j++ {
			a[j] = rand.Intn(4)
		}
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 1; j <= n; j++ {
			if j > 1 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", a[j]))
		}
		sb.WriteByte('\n')
		exp.WriteString(fmt.Sprintf("%d\n", minOps(a)))
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
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
