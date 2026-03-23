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
	"time"
)

// Embedded correct solver for 243C (CF-accepted)

type Rect struct {
	x1, x2, y1, y2 int64
}

func uniqueSorted(a []int64) []int64 {
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	n := 0
	for _, v := range a {
		if n == 0 || a[n-1] != v {
			a[n] = v
			n++
		}
	}
	return a[:n]
}

func solveC(input string) string {
	data := []byte(input)
	p := 0

	skip := func() {
		for p < len(data) {
			c := data[p]
			if c != ' ' && c != '\n' && c != '\r' && c != '\t' {
				break
			}
			p++
		}
	}

	nextByte := func() byte {
		skip()
		b := data[p]
		p++
		return b
	}

	nextInt := func() int64 {
		skip()
		var v int64
		for p < len(data) {
			c := data[p]
			if c < '0' || c > '9' {
				break
			}
			v = v*10 + int64(c-'0')
			p++
		}
		return v
	}

	n := int(nextInt())

	recs := make([]Rect, 0, n)
	xs := make([]int64, 0, 2*n+2)
	ys := make([]int64, 0, 2*n+2)

	const INF int64 = 1 << 60
	minX, maxX := INF, -INF
	minY, maxY := INF, -INF

	var x, y int64

	for i := 0; i < n; i++ {
		d := nextByte()
		dist := nextInt()

		nx, ny := x, y
		switch d {
		case 'L':
			nx -= dist
		case 'R':
			nx += dist
		case 'U':
			ny += dist
		case 'D':
			ny -= dist
		}

		var r Rect
		if x == nx {
			if y < ny {
				r = Rect{2*x - 1, 2*x + 1, 2*y - 1, 2*ny + 1}
			} else {
				r = Rect{2*x - 1, 2*x + 1, 2*ny - 1, 2*y + 1}
			}
		} else {
			if x < nx {
				r = Rect{2*x - 1, 2*nx + 1, 2*y - 1, 2*y + 1}
			} else {
				r = Rect{2*nx - 1, 2*x + 1, 2*y - 1, 2*y + 1}
			}
		}

		recs = append(recs, r)
		xs = append(xs, r.x1, r.x2)
		ys = append(ys, r.y1, r.y2)

		if r.x1 < minX {
			minX = r.x1
		}
		if r.x2 > maxX {
			maxX = r.x2
		}
		if r.y1 < minY {
			minY = r.y1
		}
		if r.y2 > maxY {
			maxY = r.y2
		}

		x, y = nx, ny
	}

	xs = append(xs, minX-2, maxX+2)
	ys = append(ys, minY-2, maxY+2)

	xs = uniqueSorted(xs)
	ys = uniqueSorted(ys)

	wp, hp := len(xs), len(ys)
	w, h := wp-1, hp-1

	xid := make(map[int64]int, wp)
	yid := make(map[int64]int, hp)
	for i, v := range xs {
		xid[v] = i
	}
	for i, v := range ys {
		yid[v] = i
	}

	diff := make([]int32, wp*hp)

	for _, r := range recs {
		l := xid[r.x1]
		rr := xid[r.x2]
		b := yid[r.y1]
		t := yid[r.y2]

		diff[l*hp+b]++
		diff[rr*hp+b]--
		diff[l*hp+t]--
		diff[rr*hp+t]++
	}

	for i := 0; i < wp; i++ {
		base := i * hp
		for j := 1; j < hp; j++ {
			diff[base+j] += diff[base+j-1]
		}
	}
	for i := 1; i < wp; i++ {
		base := i * hp
		prev := base - hp
		for j := 0; j < hp; j++ {
			diff[base+j] += diff[prev+j]
		}
	}

	widths := make([]int64, w)
	heights := make([]int64, h)
	for i := 0; i < w; i++ {
		widths[i] = (xs[i+1] - xs[i]) / 2
	}
	for j := 0; j < h; j++ {
		heights[j] = (ys[j+1] - ys[j]) / 2
	}

	totalCells := w * h
	state := make([]uint8, totalCells)
	var totalFree int64

	for i := 0; i < w; i++ {
		wi := widths[i]
		baseDiff := i * hp
		baseState := i * h
		for j := 0; j < h; j++ {
			if diff[baseDiff+j] > 0 {
				state[baseState+j] = 1
			} else {
				totalFree += wi * heights[j]
			}
		}
	}

	queue := make([]int32, 0, totalCells)
	var reachable int64

	mark := func(i, j int) {
		idx := i*h + j
		if state[idx] == 0 {
			state[idx] = 2
			reachable += widths[i] * heights[j]
			queue = append(queue, int32(idx))
		}
	}

	for i := 0; i < w; i++ {
		mark(i, 0)
		mark(i, h-1)
	}
	for j := 0; j < h; j++ {
		mark(0, j)
		mark(w-1, j)
	}

	for head := 0; head < len(queue); head++ {
		v := int(queue[head])
		i := v / h
		j := v - i*h

		if j > 0 {
			u := v - 1
			if state[u] == 0 {
				state[u] = 2
				reachable += widths[i] * heights[j-1]
				queue = append(queue, int32(u))
			}
		}
		if j+1 < h {
			u := v + 1
			if state[u] == 0 {
				state[u] = 2
				reachable += widths[i] * heights[j+1]
				queue = append(queue, int32(u))
			}
		}
		if i > 0 {
			u := v - h
			if state[u] == 0 {
				state[u] = 2
				reachable += widths[i-1] * heights[j]
				queue = append(queue, int32(u))
			}
		}
		if i+1 < w {
			u := v + h
			if state[u] == 0 {
				state[u] = 2
				reachable += widths[i+1] * heights[j]
				queue = append(queue, int32(u))
			}
		}
	}

	return fmt.Sprintf("%d", totalFree-reachable)
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func genTest() string {
	n := rand.Intn(20) + 1
	dirs := []byte{'L', 'R', 'U', 'D'}
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", n)
	for i := 0; i < n; i++ {
		d := dirs[rand.Intn(4)]
		x := rand.Intn(5) + 1
		fmt.Fprintf(&b, "%c %d\n", d, x)
	}
	return b.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go [--] /path/to/binary_or_source.go")
		os.Exit(1)
	}
	candidate, _ := filepath.Abs(os.Args[1])
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		input := genTest()
		expected := solveC(input)
		actOut, err := runBinary(candidate, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		actual := strings.TrimSpace(actOut)
		if actual != expected {
			fmt.Printf("test %d failed:\ninput:\n%s\nexpected:%s\ngot:%s\n", i+1, input, expected, actual)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
