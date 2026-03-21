package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type station struct {
	x int64
	p int64
}

type segment struct {
	amount int64
	price  int64
}

func solveC(d, n int64, stations []station) int64 {
	st := make([]station, len(stations))
	copy(st, stations)
	sort.Slice(st, func(i, j int) bool { return st[i].x < st[j].x })
	st = append(st, station{x: d, p: 0})

	deque := make([]segment, 0, len(st)+1)
	deque = append(deque, segment{amount: n, price: 0})

	var pos int64
	var totalCost int64

	for _, s := range st {
		dist := s.x - pos
		if dist > n {
			return -1
		}
		rem := dist
		for rem > 0 {
			front := &deque[0]
			if front.amount <= rem {
				totalCost += front.amount * front.price
				rem -= front.amount
				deque = deque[1:]
			} else {
				totalCost += rem * front.price
				front.amount -= rem
				rem = 0
			}
		}
		if s.x == d {
			break
		}
		addAmount := dist
		for len(deque) > 0 {
			back := deque[len(deque)-1]
			if back.price >= s.p {
				addAmount += back.amount
				deque = deque[:len(deque)-1]
			} else {
				break
			}
		}
		deque = append(deque, segment{amount: addAmount, price: s.p})
		pos = s.x
	}
	return totalCost
}

func generateC(rng *rand.Rand) (string, int64) {
	d := int64(rng.Intn(100) + 1)
	n := int64(rng.Intn(int(d)) + 1)
	m := rng.Intn(5) + 1
	st := make([]station, m)
	posSet := map[int64]bool{}
	for i := 0; i < m; i++ {
		for {
			x := int64(rng.Intn(int(d)-1) + 1)
			if !posSet[x] {
				posSet[x] = true
				st[i] = station{x, int64(rng.Intn(10) + 1)}
				break
			}
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", d, n, m)
	for _, s := range st {
		fmt.Fprintf(&sb, "%d %d\n", s.x, s.p)
	}
	cost := solveC(d, n, st)
	return sb.String(), cost
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		return
	}
	exe := os.Args[1]
	if !filepath.IsAbs(exe) {
		abs, err := filepath.Abs(exe)
		if err == nil {
			exe = abs
		}
	}
	rng := rand.New(rand.NewSource(44))
	for i := 0; i < 30; i++ {
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
