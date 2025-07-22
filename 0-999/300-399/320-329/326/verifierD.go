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

type DSU struct {
	p, r []int
}

func newDSU(n int) *DSU {
	p := make([]int, n)
	r := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i
	}
	return &DSU{p: p, r: r}
}

func (d *DSU) find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.find(d.p[x])
	}
	return d.p[x]
}

func (d *DSU) union(x, y int) {
	rx := d.find(x)
	ry := d.find(y)
	if rx == ry {
		return
	}
	if d.r[rx] < d.r[ry] {
		d.p[rx] = ry
	} else if d.r[ry] < d.r[rx] {
		d.p[ry] = rx
	} else {
		d.p[ry] = rx
		d.r[rx]++
	}
}

type Edge struct{ low, high, id int }

func expected(rects [][4]int) (bool, []int) {
	n := len(rects)
	left := make(map[int][]Edge)
	right := make(map[int][]Edge)
	bottom := make(map[int][]Edge)
	top := make(map[int][]Edge)
	for i, r := range rects {
		x1, y1, x2, y2 := r[0], r[1], r[2], r[3]
		left[x1] = append(left[x1], Edge{y1, y2, i})
		right[x2] = append(right[x2], Edge{y1, y2, i})
		bottom[y1] = append(bottom[y1], Edge{x1, x2, i})
		top[y2] = append(top[y2], Edge{x1, x2, i})
	}
	dsu := newDSU(n)
	process := func(a, b map[int][]Edge) {
		for coord, edgesA := range a {
			edgesB, ok := b[coord]
			if !ok {
				continue
			}
			sort.Slice(edgesA, func(i, j int) bool { return edgesA[i].low < edgesA[j].low })
			sort.Slice(edgesB, func(i, j int) bool { return edgesB[i].low < edgesB[j].low })
			j := 0
			for _, ea := range edgesA {
				for j < len(edgesB) && edgesB[j].high <= ea.low {
					j++
				}
				for k := j; k < len(edgesB) && edgesB[k].low < ea.high; k++ {
					dsu.union(ea.id, edgesB[k].id)
				}
			}
		}
	}
	process(left, right)
	process(right, left)
	process(bottom, top)
	process(top, bottom)

	minX := make([]int, n)
	maxX := make([]int, n)
	minY := make([]int, n)
	maxY := make([]int, n)
	area := make([]int, n)
	for i := 0; i < n; i++ {
		minX[i] = 1 << 30
		minY[i] = 1 << 30
	}
	for i, r := range rects {
		f := dsu.find(i)
		x1, y1, x2, y2 := r[0], r[1], r[2], r[3]
		if x1 < minX[f] {
			minX[f] = x1
		}
		if y1 < minY[f] {
			minY[f] = y1
		}
		if x2 > maxX[f] {
			maxX[f] = x2
		}
		if y2 > maxY[f] {
			maxY[f] = y2
		}
		area[f] += (x2 - x1) * (y2 - y1)
	}
	for i := 0; i < n; i++ {
		if dsu.find(i) != i {
			continue
		}
		dx := maxX[i] - minX[i]
		dy := maxY[i] - minY[i]
		if dx > 0 && dx == dy && area[i] == dx*dy {
			res := []int{}
			for j := 0; j < n; j++ {
				if dsu.find(j) == i {
					res = append(res, j+1)
				}
			}
			return true, res
		}
	}
	return false, nil
}

func runCase(bin string, rects [][4]int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(rects)))
	for _, r := range rects {
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", r[0], r[1], r[2], r[3]))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	gotLines := strings.Split(strings.TrimSpace(out.String()), "\n")
	ok, subset := expected(rects)
	if !ok {
		if len(gotLines) != 1 || strings.TrimSpace(gotLines[0]) != "NO" {
			return fmt.Errorf("expected NO got %q", out.String())
		}
		return nil
	}
	if len(gotLines) < 2 {
		return fmt.Errorf("expected YES")
	}
	if strings.TrimSpace(gotLines[0]) != "YES" {
		return fmt.Errorf("expected YES got %s", gotLines[0])
	}
	m, err := fmt.Sscanf(strings.TrimSpace(gotLines[1]), "%d", new(int))
	_ = m
	_ = err
	// parse second line using fmt.Fscan to get m. Better: using fmt.Sscan.
	var count int
	_, err = fmt.Sscan(strings.TrimSpace(gotLines[1]), &count)
	if err != nil {
		return fmt.Errorf("can't parse count: %v", err)
	}
	if len(gotLines) != 2 {
		arr := strings.Fields(strings.Join(gotLines[1:], " "))
		if len(arr) != 1+count {
			// there may be spaces after numbers etc
		}
	}
	// Flatten all numbers from remaining lines
	vals := []int{}
	for _, ln := range gotLines[2:] {
		fields := strings.Fields(ln)
		for _, f := range fields {
			var x int
			if _, err := fmt.Sscan(f, &x); err == nil {
				vals = append(vals, x)
			}
		}
	}
	// Some solutions may print numbers on same second line separated by spaces
	if len(vals) == 0 {
		fields := strings.Fields(strings.TrimSpace(gotLines[1]))
		if len(fields) > 1 {
			fmtvals := fields[1:]
			vals = make([]int, len(fmtvals))
			for i, f := range fmtvals {
				fmt.Sscan(f, &vals[i])
			}
		}
	}
	if len(vals) != count {
		return fmt.Errorf("expected %d ids got %d", count, len(vals))
	}
	if count != len(subset) {
		return fmt.Errorf("expected subset size %d got %d", len(subset), count)
	}
	// compare as sets
	sort.Ints(vals)
	sort.Ints(subset)
	for i := range vals {
		if vals[i] != subset[i] {
			return fmt.Errorf("expected ids %v got %v", subset, vals)
		}
	}
	return nil
}

func generateRectangles(rng *rand.Rand) [][4]int {
	n := rng.Intn(4) + 1
	rects := make([][4]int, 0, n)
	occupied := map[[2]int]bool{}
	for len(rects) < n {
		x1 := rng.Intn(5)
		y1 := rng.Intn(5)
		x2 := x1 + rng.Intn(3) + 1
		y2 := y1 + rng.Intn(3) + 1
		if x2 > 9 {
			x2 = 9
		}
		if y2 > 9 {
			y2 = 9
		}
		overlap := false
		for x := x1; x < x2; x++ {
			for y := y1; y < y2; y++ {
				if occupied[[2]int{x, y}] {
					overlap = true
					break
				}
			}
			if overlap {
				break
			}
		}
		if overlap {
			continue
		}
		for x := x1; x < x2; x++ {
			for y := y1; y < y2; y++ {
				occupied[[2]int{x, y}] = true
			}
		}
		rects = append(rects, [4]int{x1, y1, x2, y2})
	}
	return rects
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// simple positive case
	rects := [][4]int{{0, 0, 1, 2}, {1, 0, 2, 1}, {1, 1, 2, 2}, {0, 1, 1, 2}}
	if err := runCase(bin, rects); err != nil {
		fmt.Fprintf(os.Stderr, "predefined case failed: %v\n", err)
		os.Exit(1)
	}

	for i := 0; i < 100; i++ {
		rc := generateRectangles(rng)
		if err := runCase(bin, rc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
