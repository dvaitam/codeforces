package main

import (
	"bufio"
	"fmt"
	"os"
)

type sensor struct {
	col   int   // expected color
	cells []int // cell ids along the ray in order from the sensor
	ptr   int   // index of the current first not-removed cell
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m, k int
	if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
		return
	}

	// Read sensor grids.
	l := make([][]int, m)
	r := make([][]int, m)
	for i := 0; i < m; i++ {
		l[i] = make([]int, k)
		for j := 0; j < k; j++ {
			fmt.Fscan(in, &l[i][j])
		}
	}
	for i := 0; i < m; i++ {
		r[i] = make([]int, k)
		for j := 0; j < k; j++ {
			fmt.Fscan(in, &r[i][j])
		}
	}

	f := make([][]int, n)
	b := make([][]int, n)
	for i := 0; i < n; i++ {
		f[i] = make([]int, k)
		for j := 0; j < k; j++ {
			fmt.Fscan(in, &f[i][j])
		}
	}
	for i := 0; i < n; i++ {
		b[i] = make([]int, k)
		for j := 0; j < k; j++ {
			fmt.Fscan(in, &b[i][j])
		}
	}

	d := make([][]int, n)
	u := make([][]int, n)
	for i := 0; i < n; i++ {
		d[i] = make([]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(in, &d[i][j])
		}
	}
	for i := 0; i < n; i++ {
		u[i] = make([]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(in, &u[i][j])
		}
	}

	N := n * m * k
	idx := func(x, y, z int) int {
		return (x*m+y)*k + z
	}

	// Helper to abort with impossible.
	fail := func() {
		fmt.Println(-1)
		os.Exit(0)
	}

	removed := make([]bool, N) // true => cell is forced to be empty

	// Mark empty lines and detect immediate contradictions.
	xEmpty := make([][]bool, m) // by (y, z)
	for y := 0; y < m; y++ {
		xEmpty[y] = make([]bool, k)
		for z := 0; z < k; z++ {
			if (l[y][z] == 0) != (r[y][z] == 0) {
				fail()
			}
			xEmpty[y][z] = l[y][z] == 0
		}
	}
	yEmpty := make([][]bool, n) // by (x, z)
	for x := 0; x < n; x++ {
		yEmpty[x] = make([]bool, k)
		for z := 0; z < k; z++ {
			if (f[x][z] == 0) != (b[x][z] == 0) {
				fail()
			}
			yEmpty[x][z] = f[x][z] == 0
		}
	}
	zEmpty := make([][]bool, n) // by (x, y)
	for x := 0; x < n; x++ {
		zEmpty[x] = make([]bool, m)
		for y := 0; y < m; y++ {
			if (d[x][y] == 0) != (u[x][y] == 0) {
				fail()
			}
			zEmpty[x][y] = d[x][y] == 0
		}
	}

	for x := 0; x < n; x++ {
		for y := 0; y < m; y++ {
			for z := 0; z < k; z++ {
				if xEmpty[y][z] || yEmpty[x][z] || zEmpty[x][y] {
					removed[idx(x, y, z)] = true
				}
			}
		}
	}

	sensors := make([]sensor, 0, 2*(m*k+n*k+n*m))

	addSensor := func(col int, cells []int) {
		if col == 0 {
			return
		}
		if len(cells) == 0 {
			fail()
		}
		sensors = append(sensors, sensor{col: col, cells: cells})
	}

	// Build sensor lists for all six sides, skipping removed cells.
	for y := 0; y < m; y++ {
		for z := 0; z < k; z++ {
			if xEmpty[y][z] {
				continue
			}
			cells := make([]int, 0, n)
			for x := 0; x < n; x++ {
				id := idx(x, y, z)
				if !removed[id] {
					cells = append(cells, id)
				}
			}
			addSensor(l[y][z], cells) // from x = 0 side
			// from x = n+1 side (reverse order)
			rev := make([]int, 0, len(cells))
			for i := len(cells) - 1; i >= 0; i-- {
				rev = append(rev, cells[i])
			}
			addSensor(r[y][z], rev)
		}
	}

	for x := 0; x < n; x++ {
		for z := 0; z < k; z++ {
			if yEmpty[x][z] {
				continue
			}
			cells := make([]int, 0, m)
			for y := 0; y < m; y++ {
				id := idx(x, y, z)
				if !removed[id] {
					cells = append(cells, id)
				}
			}
			addSensor(f[x][z], cells)
			rev := make([]int, 0, len(cells))
			for i := len(cells) - 1; i >= 0; i-- {
				rev = append(rev, cells[i])
			}
			addSensor(b[x][z], rev)
		}
	}

	for x := 0; x < n; x++ {
		for y := 0; y < m; y++ {
			if zEmpty[x][y] {
				continue
			}
			cells := make([]int, 0, k)
			for z := 0; z < k; z++ {
				id := idx(x, y, z)
				if !removed[id] {
					cells = append(cells, id)
				}
			}
			addSensor(d[x][y], cells)
			rev := make([]int, 0, len(cells))
			for i := len(cells) - 1; i >= 0; i-- {
				rev = append(rev, cells[i])
			}
			addSensor(u[x][y], rev)
		}
	}

	// who[cell] contains sensor ids that currently point to this cell.
	who := make([][]int, N)
	state := make([]uint8, N) // 0-none, 1-one color, 2-conflict
	color := make([]int, N)   // valid when state == 1
	inQueue := make([]bool, N)
	queue := make([]int, 0)

	advance := func(sid int) bool {
		s := &sensors[sid]
		for s.ptr < len(s.cells) && removed[s.cells[s.ptr]] {
			s.ptr++
		}
		return s.ptr < len(s.cells)
	}

	pushConflict := func(cid int) {
		if !removed[cid] && !inQueue[cid] {
			inQueue[cid] = true
			queue = append(queue, cid)
		}
	}

	addToCell := func(sid int, cid int) {
		col := sensors[sid].col
		if state[cid] == 0 {
			state[cid] = 1
			color[cid] = col
		} else if state[cid] == 1 {
			if color[cid] != col {
				state[cid] = 2
				pushConflict(cid)
			}
		}
	}

	// Initialize pointers and cells.
	for sid := range sensors {
		if !advance(sid) {
			fail()
		}
		cid := sensors[sid].cells[sensors[sid].ptr]
		addToCell(sid, cid)
		who[cid] = append(who[cid], sid)
	}

	// Resolve conflicts by removing cells that would need different colors simultaneously.
	for head := 0; head < len(queue); head++ {
		cid := queue[head]
		inQueue[cid] = false
		if removed[cid] {
			continue
		}
		removed[cid] = true
		curSensors := who[cid]
		who[cid] = nil
		for _, sid := range curSensors {
			if !advance(sid) {
				fail()
			}
			newCid := sensors[sid].cells[sensors[sid].ptr]
			addToCell(sid, newCid)
			who[newCid] = append(who[newCid], sid)
		}
		state[cid] = 0
	}

	// Build final assignment.
	ans := make([]int, N)
	for sid := range sensors {
		cid := sensors[sid].cells[sensors[sid].ptr]
		if removed[cid] {
			fail()
		}
		if ans[cid] == 0 {
			ans[cid] = sensors[sid].col
		} else if ans[cid] != sensors[sid].col {
			fail()
		}
	}

	// Output in required order.
	out := bufio.NewWriter(os.Stdout)
	for x := 0; x < n; x++ {
		for y := 0; y < m; y++ {
			for z := 0; z < k; z++ {
				id := idx(x, y, z)
				if removed[id] {
					fmt.Fprint(out, 0)
				} else {
					fmt.Fprint(out, ans[id])
				}
				if !(x == n-1 && y == m-1 && z == k-1) {
					fmt.Fprint(out, " ")
				}
			}
		}
	}
	out.WriteByte('\n')
	out.Flush()
}
