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
	p []int
}

func NewDSU(n int) *DSU {
	p := make([]int, n+1)
	for i := 1; i <= n; i++ {
		p[i] = i
	}
	return &DSU{p: p}
}

func (d *DSU) Find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.Find(d.p[x])
	}
	return d.p[x]
}

func (d *DSU) Union(x, y int) int {
	rx := d.Find(x)
	ry := d.Find(y)
	if rx == ry {
		return rx
	}
	d.p[ry] = rx
	return rx
}

type person struct {
	name string
	room int
	keys []int
}

func check(n, m int, edges [][2]int, initOwner, finalOwner []int, initRoom, finalRoom map[string]int) bool {
	dsu := NewDSU(n)
	keysIn := make([][]int, n+1)
	for key := 1; key <= m; key++ {
		r := initOwner[key]
		keysIn[r] = append(keysIn[r], key)
	}
	used := make([]bool, m+1)
	queue := make([]int, 0, n)
	inq := make([]bool, n+1)
	for room := 1; room <= n; room++ {
		if len(keysIn[room]) > 0 {
			queue = append(queue, room)
			inq[room] = true
		}
	}
	for head := 0; head < len(queue); head++ {
		r0 := queue[head]
		inq[r0] = false
		r := dsu.Find(r0)
		for _, key := range keysIn[r] {
			if used[key] {
				continue
			}
			u := edges[key-1][0]
			v := edges[key-1][1]
			if dsu.Find(u) != r && dsu.Find(v) != r {
				continue
			}
			used[key] = true
			oldU := dsu.Find(u)
			oldV := dsu.Find(v)
			newR := dsu.Union(u, v)
			other := oldU
			if newR == oldU {
				other = oldV
			}
			if other != newR {
				keysIn[newR] = append(keysIn[newR], keysIn[other]...)
				keysIn[other] = nil
			}
			if !inq[newR] {
				queue = append(queue, newR)
				inq[newR] = true
			}
		}
	}
	for name, r1 := range initRoom {
		r2 := finalRoom[name]
		if dsu.Find(r1) != dsu.Find(r2) {
			return false
		}
	}
	for key := 1; key <= m; key++ {
		r1 := initOwner[key]
		r2 := finalOwner[key]
		if dsu.Find(r1) != dsu.Find(r2) {
			return false
		}
	}
	return true
}

func generateYesCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 2
	m := rng.Intn(3) + 1
	k := rng.Intn(3) + 1
	edges := make([][2]int, m)
	for i := 0; i < m; i++ {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		for v == u {
			v = rng.Intn(n) + 1
		}
		edges[i] = [2]int{u, v}
	}
	names := make([]string, k)
	for i := 0; i < k; i++ {
		names[i] = fmt.Sprintf("A%d", i+1)
	}
	initRoom := make(map[string]int, k)
	initOwner := make([]int, m+1)
	for _, name := range names {
		initRoom[name] = rng.Intn(n) + 1
	}
	perm := rng.Perm(m)
	for i, idx := range perm {
		initOwner[idx+1] = initRoom[names[i%k]]
	}
	finalRoom := make(map[string]int, k)
	finalOwner := make([]int, m+1)
	for k, v := range initRoom {
		finalRoom[k] = v
	}
	copy(finalOwner, initOwner)
	sb := strings.Builder{}
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)
	for i := 0; i < m; i++ {
		fmt.Fprintf(&sb, "%d %d\n", edges[i][0], edges[i][1])
	}
	for _, name := range names {
		room := initRoom[name]
		keys := []int{}
		for key := 1; key <= m; key++ {
			if initOwner[key] == room {
				keys = append(keys, key)
			}
		}
		fmt.Fprintf(&sb, "%s %d %d", name, room, len(keys))
		for _, key := range keys {
			fmt.Fprintf(&sb, " %d", key)
		}
		sb.WriteByte('\n')
	}
	for _, name := range names {
		room := finalRoom[name]
		keys := []int{}
		for key := 1; key <= m; key++ {
			if finalOwner[key] == room {
				keys = append(keys, key)
			}
		}
		fmt.Fprintf(&sb, "%s %d %d", name, room, len(keys))
		for _, key := range keys {
			fmt.Fprintf(&sb, " %d", key)
		}
		sb.WriteByte('\n')
	}
	return sb.String(), "YES"
}

func generateNoCase() (string, string) {
	sb := strings.Builder{}
	sb.WriteString("2 0 1\n")
	sb.WriteString("A 1 0\n")
	sb.WriteString("A 2 0\n")
	return sb.String(), "NO"
}

func runCase(exe, input, expect string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	res := strings.TrimSpace(out.String())
	if res != expect {
		return fmt.Errorf("expected %s got %s", expect, res)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		var in, exp string
		if i%2 == 0 {
			in, exp = generateYesCase(rng)
		} else {
			in, exp = generateNoCase()
		}
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
