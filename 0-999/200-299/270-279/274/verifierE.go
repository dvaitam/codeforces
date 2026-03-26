package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Correct solver for 274E (billiard ball in grid with blocked cells).

type Point struct{ x, y int }
type State struct{ x, y, dx, dy int }
type Segment struct{ x1, y1, x2, y2, dx, dy int }
type Interval struct{ min, max int }
type Event struct {
	U   int
	typ int
	V1  int
	V2  int
}

var bit []int

func bitAdd(idx, val int) {
	for ; idx < len(bit); idx += idx & -idx {
		bit[idx] += val
	}
}

func bitQuery(idx int) int {
	sum := 0
	for ; idx > 0; idx -= idx & -idx {
		sum += bit[idx]
	}
	return sum
}

func mergeIntervals(intervals []Interval) []Interval {
	if len(intervals) == 0 {
		return nil
	}
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i].min < intervals[j].min
	})
	var merged []Interval
	curr := intervals[0]
	for i := 1; i < len(intervals); i++ {
		if intervals[i].min <= curr.max {
			if intervals[i].max > curr.max {
				curr.max = intervals[i].max
			}
		} else {
			merged = append(merged, curr)
			curr = intervals[i]
		}
	}
	merged = append(merged, curr)
	return merged
}

func solve274E(input string) string {
	scanner := bufio.NewScanner(strings.NewReader(input))
	const maxCapacity = 2 * 1024 * 1024
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)
	scanner.Split(bufio.ScanWords)

	scanInt := func() int {
		scanner.Scan()
		n, _ := strconv.Atoi(scanner.Text())
		return n
	}

	scanString := func() string {
		scanner.Scan()
		return scanner.Text()
	}

	if !scanner.Scan() {
		return "0"
	}
	n, _ := strconv.Atoi(scanner.Text())
	m := scanInt()
	k := scanInt()

	blockedSet := make(map[Point]bool)
	blocks0 := make(map[int][]int)
	blocks1 := make(map[int][]int)

	for i := 0; i < k; i++ {
		bx := scanInt()
		by := scanInt()
		blockedSet[Point{bx, by}] = true
		blocks0[bx-by] = append(blocks0[bx-by], bx)
		blocks1[bx+by] = append(blocks1[bx+by], bx)
	}

	for _, arr := range blocks0 {
		sort.Ints(arr)
	}
	for _, arr := range blocks1 {
		sort.Ints(arr)
	}

	xs := scanInt()
	ys := scanInt()
	dirStr := scanString()

	dx, dy := 0, 0
	switch dirStr {
	case "NE":
		dx, dy = -1, 1
	case "NW":
		dx, dy = -1, -1
	case "SE":
		dx, dy = 1, 1
	case "SW":
		dx, dy = 1, -1
	}

	isBlocked := func(cx, cy int) bool {
		if cx < 1 || cx > n || cy < 1 || cy > m {
			return true
		}
		return blockedSet[Point{cx, cy}]
	}

	getMinT := func(x, y, dx, dy int) int {
		dir := dx * dy
		var blocks map[int][]int
		if dir == 1 {
			blocks = blocks0
		} else {
			blocks = blocks1
		}

		C := x - dir*y
		tMin := n + m + 5

		tbx := x
		if dx == 1 {
			tbx = n + 1 - x
		}
		tby := y
		if dy == 1 {
			tby = m + 1 - y
		}
		if tbx < tMin {
			tMin = tbx
		}
		if tby < tMin {
			tMin = tby
		}

		lines := []int{C, C + dx, C - dx}
		for i, L := range lines {
			arr, ok := blocks[L]
			if !ok {
				continue
			}

			if dx == 1 {
				idx := sort.Search(len(arr), func(j int) bool {
					if i == 2 {
						return arr[j] >= x
					}
					return arr[j] > x
				})
				if idx < len(arr) {
					u := arr[idx]
					t := u - x
					if i == 2 {
						t++
					}
					if t < tMin {
						tMin = t
					}
				}
			} else {
				idx := sort.Search(len(arr), func(j int) bool {
					if i == 2 {
						return arr[j] > x
					}
					return arr[j] >= x
				})
				idx--
				if idx >= 0 {
					u := arr[idx]
					t := x - u
					if i == 2 {
						t++
					}
					if t < tMin {
						tMin = t
					}
				}
			}
		}
		return tMin
	}

	visitedStates := make(map[State]bool)
	var segments []Segment

	x, y := xs, ys
	for {
		state := State{x, y, dx, dy}
		if visitedStates[state] {
			break
		}
		visitedStates[state] = true

		t := getMinT(x, y, dx, dy)
		nx := x + (t-1)*dx
		ny := y + (t-1)*dy

		segments = append(segments, Segment{x, y, nx, ny, dx, dy})

		bx := isBlocked(nx+dx, ny)
		by := isBlocked(nx, ny+dy)
		bxy := isBlocked(nx+dx, ny+dy)

		var ndx, ndy int
		if bx && by {
			ndx, ndy = -dx, -dy
		} else if bx {
			ndx, ndy = -dx, dy
		} else if by {
			ndx, ndy = dx, -dy
		} else if bxy {
			ndx, ndy = -dx, -dy
		} else {
			ndx, ndy = dx, dy
		}

		x, y, dx, dy = nx, ny, ndx, ndy
	}

	horizMap := make(map[int][]Interval)
	vertMap := make(map[int][]Interval)

	for _, seg := range segments {
		U1 := seg.x1 + seg.y1
		U2 := seg.x2 + seg.y2
		V1 := seg.x1 - seg.y1
		V2 := seg.x2 - seg.y2

		if seg.dx*seg.dy == 1 {
			if U1 > U2 {
				U1, U2 = U2, U1
			}
			horizMap[V1] = append(horizMap[V1], Interval{U1, U2})
		} else {
			if V1 > V2 {
				V1, V2 = V2, V1
			}
			vertMap[U1] = append(vertMap[U1], Interval{V1, V2})
		}
	}

	var totalCells int64 = 0
	mergedHoriz := make(map[int][]Interval)
	for V, intervals := range horizMap {
		mr := mergeIntervals(intervals)
		mergedHoriz[V] = mr
		for _, iv := range mr {
			totalCells += int64((iv.max-iv.min)/2 + 1)
		}
	}

	mergedVert := make(map[int][]Interval)
	for U, intervals := range vertMap {
		mr := mergeIntervals(intervals)
		mergedVert[U] = mr
		for _, iv := range mr {
			totalCells += int64((iv.max-iv.min)/2 + 1)
		}
	}

	var events []Event
	for V, intervals := range mergedHoriz {
		for _, iv := range intervals {
			events = append(events, Event{U: iv.min, typ: 1, V1: V + 100005})
			events = append(events, Event{U: iv.max, typ: 3, V1: V + 100005})
		}
	}
	for U, intervals := range mergedVert {
		for _, iv := range intervals {
			events = append(events, Event{U: U, typ: 2, V1: iv.min + 100005, V2: iv.max + 100005})
		}
	}

	sort.Slice(events, func(i, j int) bool {
		if events[i].U != events[j].U {
			return events[i].U < events[j].U
		}
		return events[i].typ < events[j].typ
	})

	bit = make([]int, 300005)
	var intersections int64 = 0

	for _, ev := range events {
		if ev.typ == 1 {
			bitAdd(ev.V1, 1)
		} else if ev.typ == 3 {
			bitAdd(ev.V1, -1)
		} else if ev.typ == 2 {
			intersections += int64(bitQuery(ev.V2) - bitQuery(ev.V1-1))
		}
	}

	return fmt.Sprintf("%d", totalCells-intersections)
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	m := rng.Intn(10) + 1
	maxCells := n*m - 1
	k := 0
	if maxCells > 0 {
		k = rng.Intn(minInt(maxCells, 4) + 1)
	}
	blocks := map[Point]bool{}
	for len(blocks) < k {
		x := rng.Intn(n) + 1
		y := rng.Intn(m) + 1
		blocks[Point{x, y}] = true
	}
	xs := rng.Intn(n) + 1
	ys := rng.Intn(m) + 1
	for blocks[Point{xs, ys}] {
		xs = rng.Intn(n) + 1
		ys = rng.Intn(m) + 1
	}
	dirs := []string{"NE", "NW", "SE", "SW"}
	dir := dirs[rng.Intn(4)]
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	for p := range blocks {
		sb.WriteString(fmt.Sprintf("%d %d\n", p.x, p.y))
	}
	sb.WriteString(fmt.Sprintf("%d %d %s\n", xs, ys, dir))
	return sb.String()
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func runCase(exe string, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := generateCase(rng)
		exp := solve274E(in)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
