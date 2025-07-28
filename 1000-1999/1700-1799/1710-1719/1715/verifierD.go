package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type EdgeD struct {
	u int
	v int
	x int
}

type DSUD struct {
	parent []int
	size   []int
	mx     []int
	need   []bool
	hasOne []bool
}

func NewDSUD(n int) *DSUD {
	d := &DSUD{parent: make([]int, n), size: make([]int, n), mx: make([]int, n), need: make([]bool, n), hasOne: make([]bool, n)}
	for i := 0; i < n; i++ {
		d.parent[i] = i
		d.size[i] = 1
		d.mx[i] = i
	}
	return d
}

func (d *DSUD) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSUD) union(a, b int) int {
	ra := d.find(a)
	rb := d.find(b)
	if ra == rb {
		return ra
	}
	if d.size[ra] < d.size[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	d.size[ra] += d.size[rb]
	if d.mx[rb] > d.mx[ra] {
		d.mx[ra] = d.mx[rb]
	}
	d.need[ra] = d.need[ra] || d.need[rb]
	d.hasOne[ra] = d.hasOne[ra] || d.hasOne[rb]
	return ra
}

func solveD(n int, edges []EdgeD) []int {
	ans := make([]int, n)
	for bit := 0; bit < 30; bit++ {
		forced0 := make([]bool, n)
		for _, e := range edges {
			if ((e.x >> bit) & 1) == 0 {
				forced0[e.u] = true
				forced0[e.v] = true
			}
		}
		dsu := NewDSUD(n)
		forced1 := make([]bool, n)
		for _, e := range edges {
			if ((e.x >> bit) & 1) == 1 {
				fu := forced0[e.u]
				fv := forced0[e.v]
				if fu && fv {
				} else if fu && !fv {
					forced1[e.v] = true
				} else if fv && !fu {
					forced1[e.u] = true
				} else {
					root := dsu.union(e.u, e.v)
					dsu.need[root] = true
				}
			}
		}
		for i := 0; i < n; i++ {
			if forced0[i] {
				continue
			}
			r := dsu.find(i)
			if dsu.mx[r] < i {
				dsu.mx[r] = i
			}
		}
		for i := 0; i < n; i++ {
			if forced1[i] {
				r := dsu.find(i)
				dsu.hasOne[r] = true
			}
		}
		seen := make(map[int]bool)
		for i := 0; i < n; i++ {
			if forced0[i] {
				continue
			}
			r := dsu.find(i)
			if seen[r] {
				continue
			}
			seen[r] = true
			if dsu.need[r] && !dsu.hasOne[r] {
				idx := dsu.mx[r]
				forced1[idx] = true
				dsu.hasOne[r] = true
			}
		}
		for i := 0; i < n; i++ {
			if forced1[i] {
				ans[i] |= (1 << bit)
			}
		}
	}
	return ans
}

type testCaseD struct {
	n     int
	q     int
	edges []EdgeD
}

func genCaseD(rng *rand.Rand) testCaseD {
	n := rng.Intn(5) + 1
	q := rng.Intn(7) + 1
	edges := make([]EdgeD, q)
	for i := 0; i < q; i++ {
		u := rng.Intn(n)
		v := rng.Intn(n)
		x := rng.Intn(16)
		edges[i] = EdgeD{u, v, x}
	}
	return testCaseD{n: n, q: q, edges: edges}
}

func runCaseD(bin string, tc testCaseD) error {
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.q))
	for _, e := range tc.edges {
		input.WriteString(fmt.Sprintf("%d %d %d\n", e.u+1, e.v+1, e.x))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	if len(fields) != tc.n {
		return fmt.Errorf("expected %d numbers got %d", tc.n, len(fields))
	}
	expect := solveD(tc.n, tc.edges)
	for i, f := range fields {
		val, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("invalid int %q", f)
		}
		if val != expect[i] {
			return fmt.Errorf("mismatch expected %v got %v", expect, fields)
		}
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
	for t := 0; t < 100; t++ {
		tc := genCaseD(rng)
		if err := runCaseD(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n", t+1, err)
			var inp strings.Builder
			inp.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.q))
			for _, e := range tc.edges {
				inp.WriteString(fmt.Sprintf("%d %d %d\n", e.u+1, e.v+1, e.x))
			}
			fmt.Fprint(os.Stderr, inp.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
