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
	parent []int
	size   []int
}

func NewDSU(n int) *DSU {
	d := &DSU{parent: make([]int, n), size: make([]int, n)}
	for i := 0; i < n; i++ {
		d.parent[i] = i
		d.size[i] = 1
	}
	return d
}

func (d *DSU) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) union(x, y int) {
	fx := d.find(x)
	fy := d.find(y)
	if fx == fy {
		return
	}
	if d.size[fx] < d.size[fy] {
		fx, fy = fy, fx
	}
	d.parent[fy] = fx
	d.size[fx] += d.size[fy]
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String() + errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

type event struct {
	id   int
	left bool
}

func solveD(n int, segs [][2]int) string {
	events := make([]event, 2*n+2)
	for i := 0; i < n; i++ {
		l := segs[i][0]
		r := segs[i][1]
		events[l] = event{id: i + 1, left: true}
		events[r] = event{id: i + 1, left: false}
	}
	dsu := NewDSU(n)
	stack := []int{}
	edges := 0
	for pos := 1; pos <= 2*n; pos++ {
		e := events[pos]
		if e.id == 0 {
			continue
		}
		id := e.id - 1
		if e.left {
			stack = append(stack, id)
		} else {
			temp := []int{}
			for len(stack) > 0 && stack[len(stack)-1] != id {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				dsu.union(id, top)
				edges++
				temp = append(temp, top)
				if edges >= n {
					return "NO"
				}
			}
			if len(stack) == 0 {
				return "NO"
			}
			stack = stack[:len(stack)-1]
			for i := len(temp) - 1; i >= 0; i-- {
				stack = append(stack, temp[i])
			}
		}
	}
	if edges != n-1 {
		return "NO"
	}
	root := dsu.find(0)
	for i := 1; i < n; i++ {
		if dsu.find(i) != root {
			return "NO"
		}
	}
	return "YES"
}

func generateCase(r *rand.Rand) (string, string) {
	n := r.Intn(7) + 1
	perm := r.Perm(2 * n)
	segs := make([][2]int, n)
	for i := 0; i < n; i++ {
		a := perm[2*i] + 1
		b := perm[2*i+1] + 1
		if a > b {
			a, b = b, a
		}
		segs[i] = [2]int{a, b}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", segs[i][0], segs[i][1]))
	}
	expect := solveD(n, segs)
	return sb.String(), expect
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
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
