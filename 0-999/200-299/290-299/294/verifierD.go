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

type Fenwick struct {
	n int
	f []int
}

func NewFenwick(n int) *Fenwick { return &Fenwick{n, make([]int, n+1)} }
func (ft *Fenwick) Update(i, v int) {
	for x := i; x <= ft.n; x += x & -x {
		ft.f[x] += v
	}
}
func (ft *Fenwick) Query(i int) int {
	if i > ft.n {
		i = ft.n
	}
	if i < 1 {
		return 0
	}
	s := 0
	for x := i; x > 0; x -= x & -x {
		s += ft.f[x]
	}
	return s
}
func (ft *Fenwick) Range(l, r int) int {
	if l > r {
		return 0
	}
	if l < 1 {
		l = 1
	}
	if r > ft.n {
		r = ft.n
	}
	return ft.Query(r) - ft.Query(l-1)
}

type Segment struct{ x, y, dx, dy, l int }

func solveD(n, m, xs, ys int, dir string) string {
	dx, dy := 0, 0
	switch dir {
	case "DR":
		dx, dy = 1, 1
	case "DL":
		dx, dy = 1, -1
	case "UR":
		dx, dy = -1, 1
	case "UL":
		dx, dy = -1, -1
	}
	p := (xs + ys) & 1
	total := int64(n) * int64(m)
	var need int64
	if p == 0 {
		need = (total + 1) / 2
	} else {
		need = total / 2
	}
	var segs []Segment
	x0, y0 := xs, ys
	dx0, dy0 := dx, dy
	for {
		sx := n - x0
		if dx < 0 {
			sx = x0 - 1
		}
		sy := m - y0
		if dy < 0 {
			sy = y0 - 1
		}
		steps := sx
		if sy < steps {
			steps = sy
		}
		segs = append(segs, Segment{x0, y0, dx, dy, steps})
		x1 := x0 + dx*steps
		y1 := y0 + dy*steps
		if x1 == n || x1 == 1 {
			dx = -dx
		}
		if y1 == m || y1 == 1 {
			dy = -dy
		}
		x0, y0 = x1, y1
		if x0 == xs && y0 == ys && dx == dx0 && dy == dy0 {
			break
		}
	}
	maxD := n + m + 5
	ftMinus := NewFenwick(maxD)
	off := m + 2
	maxV := n + m + 5
	ftPlus := NewFenwick(maxV)
	var unique int64 = 1
	var stepsCount int64 = 1
	if dx0*dy0 > 0 {
		v0 := xs - ys
		ftPlus.Update(v0+off, 1)
	} else {
		d0 := xs + ys
		ftMinus.Update(d0, 1)
	}
	for _, s := range segs {
		typ := 0
		if s.dx*s.dy < 0 {
			typ = 1
		}
		segLen := int64(s.l)
		var inter int
		if typ == 0 {
			sum0 := s.x + s.y
			if s.dy > 0 {
				l := sum0 + 2
				r := sum0 + 2*s.l
				inter = ftMinus.Range(l, r)
			} else {
				l := sum0 - 2*s.l
				r := sum0 - 2
				inter = ftMinus.Range(l, r)
			}
		} else {
			diff0 := s.x - s.y
			if s.dx > 0 {
				l := diff0 + 2
				r := diff0 + 2*s.l
				inter = ftPlus.Range(l+off, r+off)
			} else {
				l := diff0 - 2*s.l
				r := diff0 - 2
				inter = ftPlus.Range(l+off, r+off)
			}
		}
		newAll := segLen - int64(inter)
		if unique+newAll < need {
			unique += newAll
			stepsCount += segLen
			if typ == 0 {
				v := s.x - s.y
				ftPlus.Update(v+off, 1)
			} else {
				d := s.x + s.y
				ftMinus.Update(d, 1)
			}
			continue
		}
		needLeft := need - unique
		lo, hi := int64(1), segLen
		best := segLen
		for lo <= hi {
			mid := (lo + hi) / 2
			var inter2 int
			if typ == 0 {
				sum0 := s.x + s.y
				if s.dy > 0 {
					r := sum0 + 2*int(mid)
					inter2 = ftMinus.Range(sum0+2, r)
				} else {
					l := sum0 - 2*int(mid)
					inter2 = ftMinus.Range(l, sum0-2)
				}
			} else {
				diff0 := s.x - s.y
				if s.dx > 0 {
					r := diff0 + 2*int(mid)
					inter2 = ftPlus.Range(diff0+2+off, r+off)
				} else {
					l := diff0 - 2*int(mid)
					inter2 = ftPlus.Range(l+off, diff0-2+off)
				}
			}
			visited := mid - int64(inter2)
			if visited >= needLeft {
				best = mid
				hi = mid - 1
			} else {
				lo = mid + 1
			}
		}
		return fmt.Sprintf("%d\n", stepsCount+best)
	}
	return "-1\n"
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 2
	m := rng.Intn(8) + 2
	side := rng.Intn(4)
	var xs, ys int
	switch side {
	case 0:
		xs = 1
		ys = rng.Intn(m) + 1
	case 1:
		xs = n
		ys = rng.Intn(m) + 1
	case 2:
		ys = 1
		xs = rng.Intn(n) + 1
	default:
		ys = m
		xs = rng.Intn(n) + 1
	}
	dirs := []string{"UL", "UR", "DL", "DR"}
	dir := dirs[rng.Intn(4)]
	ans := solveD(n, m, xs, ys, dir)
	in := fmt.Sprintf("%d %d\n%d %d %s\n", n, m, xs, ys, dir)
	return in, ans
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := generateCase(rng)
		got, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("test %d: execution failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, in, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
