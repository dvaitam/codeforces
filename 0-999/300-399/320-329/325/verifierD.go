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

type DSU struct {
	parent      []int
	left, right []bool
}

func NewDSU(n int) *DSU {
	d := &DSU{parent: make([]int, n), left: make([]bool, n), right: make([]bool, n)}
	return d
}

func (d *DSU) Find(x int) int {
	for d.parent[x] != x {
		d.parent[x] = d.parent[d.parent[x]]
		x = d.parent[x]
	}
	return x
}

func (d *DSU) Union(x, y int) {
	rx := d.Find(x)
	ry := d.Find(y)
	if rx == ry {
		return
	}
	d.parent[ry] = rx
	d.left[rx] = d.left[rx] || d.left[ry]
	d.right[rx] = d.right[rx] || d.right[ry]
}

func expected(r, c int, ops [][2]int) int {
	rc := r * c
	idx := make([]int, rc)
	dsu := NewDSU(len(ops) + 5)
	curID := 0
	ans := 0
	for _, op := range ops {
		i := op[0] - 1
		j := op[1] - 1
		pos := i*c + j
		newLeft := j == 0
		newRight := j == c-1
		orLeft, orRight := false, false
		neigh := make([]int, 0, 4)
		// up/down
		for _, di := range []int{-1, 1} {
			ni := i + di
			if ni < 0 || ni >= r {
				continue
			}
			nj := j
			p2 := ni*c + nj
			id2 := idx[p2]
			if id2 == 0 {
				continue
			}
			root := dsu.Find(id2)
			dup := false
			for _, v := range neigh {
				if v == root {
					dup = true
					break
				}
			}
			if dup {
				continue
			}
			neigh = append(neigh, root)
			if dsu.left[root] {
				orLeft = true
			}
			if dsu.right[root] {
				orRight = true
			}
		}
		// left/right wrap
		for _, dj := range []int{-1, 1} {
			ni := i
			nj := j + dj
			if nj < 0 {
				nj = c - 1
			} else if nj >= c {
				nj = 0
			}
			p2 := ni*c + nj
			id2 := idx[p2]
			if id2 == 0 {
				continue
			}
			root := dsu.Find(id2)
			dup := false
			for _, v := range neigh {
				if v == root {
					dup = true
					break
				}
			}
			if dup {
				continue
			}
			neigh = append(neigh, root)
			if dsu.left[root] {
				orLeft = true
			}
			if dsu.right[root] {
				orRight = true
			}
		}
		totalLeft := newLeft || orLeft
		totalRight := newRight || orRight
		if totalLeft && totalRight {
			continue
		}
		curID++
		ans++
		id := curID
		idx[pos] = id
		dsu.parent[id] = id
		dsu.left[id] = newLeft
		dsu.right[id] = newRight
		for _, v := range neigh {
			dsu.Union(id, v)
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, int) {
	r := rng.Intn(4) + 2
	c := rng.Intn(4) + 2
	n := rng.Intn(r*c) + 1
	ops := make([][2]int, n)
	for i := 0; i < n; i++ {
		ops[i] = [2]int{rng.Intn(r) + 1, rng.Intn(c) + 1}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", r, c, n))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", ops[i][0], ops[i][1]))
	}
	return sb.String(), expected(r, c, ops)
}

func runCase(bin, input string, exp int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
