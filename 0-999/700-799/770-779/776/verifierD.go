package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type DSU struct {
	parent, rank, diff []int
}

func NewDSU(n int) *DSU {
	d := &DSU{
		parent: make([]int, n),
		rank:   make([]int, n),
		diff:   make([]int, n),
	}
	for i := 0; i < n; i++ {
		d.parent[i] = i
	}
	return d
}

func (d *DSU) find(x int) (int, int) {
	if d.parent[x] == x {
		return x, 0
	}
	root, parity := d.find(d.parent[x])
	parity ^= d.diff[x]
	d.parent[x] = root
	d.diff[x] = parity
	return root, parity
}

func (d *DSU) unite(x, y, val int) bool {
	rx, px := d.find(x)
	ry, py := d.find(y)
	if rx == ry {
		return (px ^ py) == val
	}
	if d.rank[rx] < d.rank[ry] {
		rx, ry = ry, rx
		px, py = py, px
	}
	d.parent[ry] = rx
	d.diff[ry] = px ^ py ^ val
	if d.rank[rx] == d.rank[ry] {
		d.rank[rx]++
	}
	return true
}

func expected(n, m int, door []int, switches [][2]int) string {
	dsu := NewDSU(m + 1)
	possible := true
	for i := 0; i < n && possible; i++ {
		a := switches[i][0]
		b := switches[i][1]
		val := 1 - door[i]
		if !dsu.unite(a, b, val) {
			possible = false
			break
		}
	}
	if possible {
		return "YES"
	} else {
		return "NO"
	}
}

func run(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for t := 1; t <= 100; t++ {
		n := (t % 3) + 2 // 2..4
		m := (t % 3) + 2 // 2..4
		door := make([]int, n)
		for i := 0; i < n; i++ {
			door[i] = (t + i) % 2
		}
		switches := make([][2]int, n)
		for i := 0; i < n; i++ {
			switches[i][0] = i%m + 1
			switches[i][1] = (i+1)%m + 1
		}
		// build switch inputs
		roomsFor := make([][]int, m+1)
		for i := 0; i < n; i++ {
			a := switches[i][0]
			b := switches[i][1]
			roomsFor[a] = append(roomsFor[a], i+1)
			roomsFor[b] = append(roomsFor[b], i+1)
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", door[i])
		}
		sb.WriteByte('\n')
		for sw := 1; sw <= m; sw++ {
			fmt.Fprintf(&sb, "%d", len(roomsFor[sw]))
			for _, r := range roomsFor[sw] {
				fmt.Fprintf(&sb, " %d", r)
			}
			sb.WriteByte('\n')
		}
		inp := sb.String()
		exp := expected(n, m, door, switches)
		got, err := run(bin, inp)
		if err != nil {
			fmt.Printf("Test %d: execution error: %v\nOutput:\n%s\n", t, err, got)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected:%s Got:%s\n", t, inp, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
