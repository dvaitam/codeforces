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

type pair struct{ a, idx int }

func solveB2(a []int) string {
	n := len(a)
	v := make([]pair, n)
	for i := 0; i < n; i++ {
		v[i] = pair{a[i], i}
	}
	sort.Slice(v, func(i, j int) bool { return v[i].a < v[j].a })
	ansx := make([]int, n)
	ansy := make([]int, n)
	dir := make([]int, n)
	taken := make([]int, n+1)
	takenID := make([]int, n+1)
	if v[0].a == 0 {
		id0 := v[0].idx
		ansx[id0] = 1
		ansy[id0] = 1
		dir[id0] = id0 + 1
		taken[1] = 1
		takenID[1] = id0
	} else {
		i := 0
		for i+1 < n && v[i].a != v[i+1].a {
			i++
		}
		dis := v[i].a
		id1 := v[i].idx
		id2 := v[i+1].idx
		ansx[id1] = 1
		ansy[id1] = 1
		dir[id1] = id2 + 1
		taken[1] = 1
		takenID[1] = id1
		x := 1 + dis
		y := 1
		if x > n {
			y += x - n
			x = n
		}
		ansx[id2] = x
		ansy[id2] = y
		dir[id2] = id1 + 1
		taken[x] = y
		takenID[x] = id2
	}
	curX := 1
	for _, p := range v {
		aVal := p.a
		id := p.idx
		if dir[id] != 0 {
			continue
		}
		for curX <= n && taken[curX] != 0 {
			curX++
		}
		var y int
		if aVal == 0 {
			y = 1
			dir[id] = id + 1
		} else {
			if curX-aVal >= 1 {
				y = taken[curX-aVal]
				dir[id] = takenID[curX-aVal] + 1
			} else {
				y = aVal - curX + 2
				dir[id] = takenID[1] + 1
			}
		}
		ansx[id] = curX
		ansy[id] = y
		taken[curX] = y
		takenID[curX] = id
	}
	var sb strings.Builder
	sb.WriteString("YES\n")
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", ansx[i], ansy[i]))
	}
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d\n", dir[i]))
	}
	return sb.String()
}

func generateCase(rng *rand.Rand) (string, []int) {
	n := rng.Intn(10) + 2
	arr := make([]int, n)
	arr[0] = 0
	for i := 1; i < n; i++ {
		arr[i] = rng.Intn(n + 1)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String(), arr
}

func runCase(bin string, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, arr := generateCase(rng)
		expect := solveB2(arr)
		if err := runCase(bin, in, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
