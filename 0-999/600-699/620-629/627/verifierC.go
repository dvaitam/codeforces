package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type station struct {
	x int
	p int
}

func expectedC(d, n int, stations []station) int64 {
	st := make([]station, 0, len(stations)+2)
	st = append(st, station{0, 0})
	st = append(st, stations...)
	st = append(st, station{d, 0})
	sort.Slice(st, func(i, j int) bool { return st[i].x < st[j].x })
	for i := 0; i < len(st)-1; i++ {
		if st[i+1].x-st[i].x > n {
			return -1
		}
	}
	next := make([]int, len(st))
	stack := []int{}
	for i := len(st) - 1; i >= 0; i-- {
		for len(stack) > 0 && st[stack[len(stack)-1]].p >= st[i].p {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			next[i] = -1
		} else {
			next[i] = stack[len(stack)-1]
		}
		stack = append(stack, i)
	}
	fuel := n
	var cost int64
	for i := 0; i < len(st)-1; i++ {
		dist := st[i+1].x - st[i].x
		target := n
		if j := next[i]; j != -1 {
			if st[j].x-st[i].x <= n {
				target = st[j].x - st[i].x
			}
		}
		if target < dist {
			target = dist
		}
		if fuel < target {
			add := target - fuel
			cost += int64(add * st[i].p)
			fuel += add
		}
		fuel -= dist
	}
	return cost
}

func generateC(rng *rand.Rand) (string, int64) {
	d := rng.Intn(100) + 1
	n := rng.Intn(d) + 1
	m := rng.Intn(5) + 1
	st := make([]station, m)
	posSet := map[int]bool{}
	for i := 0; i < m; i++ {
		for {
			x := rng.Intn(d-1) + 1
			if !posSet[x] {
				posSet[x] = true
				st[i] = station{x, rng.Intn(10) + 1}
				break
			}
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", d, n, m)
	for _, s := range st {
		fmt.Fprintf(&sb, "%d %d\n", s.x, s.p)
	}
	cost := expectedC(d, n, st)
	return sb.String(), cost
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		return
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(44))
	for i := 0; i < 100; i++ {
		input, exp := generateC(rng)
		cmd := exec.Command(exe)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			fmt.Printf("case %d runtime error: %v\n%s", i+1, err, out.String())
			return
		}
		got := strings.TrimSpace(out.String())
		if got != fmt.Sprint(exp) {
			fmt.Printf("case %d failed: expected %d got %s\ninput:\n%s", i+1, exp, got, input)
			return
		}
	}
	fmt.Println("All tests passed")
}
