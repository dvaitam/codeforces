package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func abs64(a int64) int64 {
	if a < 0 {
		return -a
	}
	return a
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	var l int64
	if _, err := fmt.Fscan(reader, &n, &l); err != nil {
		return
	}
	type Occ struct {
		c0, x0, y0 int64
	}
	occ := make(map[int64]*Occ, n)
	var dxGlobal, dyGlobal int64
	var hasDXGlobal, hasDYGlobal bool
	for i := 0; i < n; i++ {
		var t, x, y int64
		fmt.Fscan(reader, &t, &x, &y)
		c := t / l
		r := t % l
		if o, ok := occ[r]; !ok {
			occ[r] = &Occ{c0: c, x0: x, y0: y}
		} else {
			d := c - o.c0
			dx := x - o.x0
			dy := y - o.y0
			if d <= 0 || dx%d != 0 || dy%d != 0 {
				fmt.Fprintln(writer, "NO")
				return
			}
			vx := dx / d
			vy := dy / d
			// enforce global dx, dy consistency
			if !hasDXGlobal {
				dxGlobal = vx
				hasDXGlobal = true
			} else if dxGlobal != vx {
				fmt.Fprintln(writer, "NO")
				return
			}
			if !hasDYGlobal {
				dyGlobal = vy
				hasDYGlobal = true
			} else if dyGlobal != vy {
				fmt.Fprintln(writer, "NO")
				return
			}
		}
	}
	if !hasDXGlobal {
		dxGlobal = 0
	}
	if !hasDYGlobal {
		dyGlobal = 0
	}
	// collect fixed points
	type P struct {
		idx  int
		x, y int64
	}
	fps := make([]P, 0, len(occ)+2)
	// start point
	fps = append(fps, P{idx: 0, x: 0, y: 0})
	for r, o := range occ {
		sx := o.x0 - o.c0*dxGlobal
		sy := o.y0 - o.c0*dyGlobal
		if r == 0 {
			if sx != 0 || sy != 0 {
				fmt.Fprintln(writer, "NO")
				return
			}
			continue
		}
		fps = append(fps, P{idx: int(r), x: sx, y: sy})
	}
	// end point at l
	fps = append(fps, P{idx: int(l), x: dxGlobal, y: dyGlobal})
	sort.Slice(fps, func(i, j int) bool { return fps[i].idx < fps[j].idx })
	// build moves
	res := make([]byte, l)
	for i := 0; i+1 < len(fps); i++ {
		a := fps[i]
		b := fps[i+1]
		dt := b.idx - a.idx
		dx := b.x - a.x
		dy := b.y - a.y
		D := abs64(dx) + abs64(dy)
		if D > int64(dt) || (int64(dt)-D)%2 != 0 {
			fmt.Fprintln(writer, "NO")
			return
		}
		pos := a.idx
		// x moves
		if dx > 0 {
			for k := int64(0); k < dx; k++ {
				res[pos] = 'R'
				pos++
			}
		} else if dx < 0 {
			for k := int64(0); k < -dx; k++ {
				res[pos] = 'L'
				pos++
			}
		}
		// y moves
		if dy > 0 {
			for k := int64(0); k < dy; k++ {
				res[pos] = 'U'
				pos++
			}
		} else if dy < 0 {
			for k := int64(0); k < -dy; k++ {
				res[pos] = 'D'
				pos++
			}
		}
		rem := dt - int(D)
		// fill with LR pairs
		for k := 0; k < rem; k += 2 {
			res[pos] = 'L'
			res[pos+1] = 'R'
			pos += 2
		}
	}
	// output
	fmt.Fprintln(writer, string(res))
}
