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

type Node struct {
	covered bool
	c       [4]*Node
}

func (node *Node) update(x1, x2, y1, y2, xl, xr, yl, yr int) {
	if node.covered || x2 <= xl || xr <= x1 || y2 <= yl || yr <= y1 {
		return
	}
	if x1 <= xl && xr <= x2 && y1 <= yl && yr <= y2 {
		node.covered = true
		node.c = [4]*Node{}
		return
	}
	if xl+1 == xr && yl+1 == yr {
		node.covered = true
		node.c = [4]*Node{}
		return
	}
	mx := (xl + xr) / 2
	my := (yl + yr) / 2
	if node.c[0] == nil {
		node.c[0] = &Node{}
		node.c[1] = &Node{}
		node.c[2] = &Node{}
		node.c[3] = &Node{}
	}
	node.c[0].update(x1, x2, y1, y2, xl, mx, yl, my)
	node.c[1].update(x1, x2, y1, y2, xl, mx, my, yr)
	node.c[2].update(x1, x2, y1, y2, mx, xr, yl, my)
	node.c[3].update(x1, x2, y1, y2, mx, xr, my, yr)
	if node.c[0].covered && node.c[1].covered && node.c[2].covered && node.c[3].covered {
		node.covered = true
		node.c = [4]*Node{}
	}
}

func (node *Node) query(x1, x2, y1, y2, xl, xr, yl, yr int) bool {
	if x2 <= xl || xr <= x1 || y2 <= yl || yr <= y1 {
		return true
	}
	if node.covered {
		return true
	}
	if xl+1 == xr && yl+1 == yr {
		return node.covered
	}
	mx := (xl + xr) / 2
	my := (yl + yr) / 2
	if node.c[0] == nil {
		return false
	}
	return node.c[0].query(x1, x2, y1, y2, xl, mx, yl, my) &&
		node.c[1].query(x1, x2, y1, y2, xl, mx, my, yr) &&
		node.c[2].query(x1, x2, y1, y2, mx, xr, yl, my) &&
		node.c[3].query(x1, x2, y1, y2, mx, xr, my, yr)
}

type Rect struct{ x1, y1, x2, y2 int }

func uniqueInts(a []int) []int {
	if len(a) == 0 {
		return a
	}
	j := 1
	for i := 1; i < len(a); i++ {
		if a[i] != a[i-1] {
			a[j] = a[i]
			j++
		}
	}
	return a[:j]
}

func solveD(rects []Rect) int {
	n := len(rects)
	xs := make([]int, 0, 2*n)
	ys := make([]int, 0, 2*n)
	for _, r := range rects {
		xs = append(xs, r.x1, r.x2)
		ys = append(ys, r.y1, r.y2)
	}
	sort.Ints(xs)
	xs = uniqueInts(xs)
	sort.Ints(ys)
	ys = uniqueInts(ys)
	toIndex := func(arr []int, v int) int { return sort.SearchInts(arr, v) }
	root := &Node{}
	visible := 1
	seen := make([]bool, n)
	for i := n - 1; i >= 0; i-- {
		r := rects[i]
		x1 := toIndex(xs, r.x1)
		x2 := toIndex(xs, r.x2)
		y1 := toIndex(ys, r.y1)
		y2 := toIndex(ys, r.y2)
		if !root.query(x1, x2, y1, y2, 0, len(xs)-1, 0, len(ys)-1) {
			visible++
			seen[i] = true
		}
		root.update(x1, x2, y1, y2, 0, len(xs)-1, 0, len(ys)-1)
	}
	_ = seen
	return visible
}

func generateD(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 1
	rects := make([]Rect, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		x1 := rng.Intn(10)
		y1 := rng.Intn(10)
		x2 := x1 + rng.Intn(5) + 1
		y2 := y1 + rng.Intn(5) + 1
		rects[i] = Rect{x1, y1, x2, y2}
		fmt.Fprintf(&sb, "%d %d %d %d\n", x1, y1, x2, y2)
	}
	res := solveD(rects)
	return sb.String(), fmt.Sprint(res)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(45))
	for i := 0; i < 100; i++ {
		in, exp := generateD(rng)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(in)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			fmt.Printf("case %d runtime error: %v\n%s", i+1, err, out.String())
			return
		}
		got := strings.TrimSpace(out.String())
		if got != exp {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, in, exp, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
