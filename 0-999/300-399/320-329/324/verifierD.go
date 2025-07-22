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

type Segment struct {
	vertical  bool
	fixed     int64
	low, high int64
	headCoord int64
	headDir   byte
}

type queryD struct {
	x, y int64
	dir  byte
	t    int64
}

type testCaseD struct {
	input    string
	expected []string
}

func solveQueriesD(b int64, vert []Segment, hor []Segment, qs []queryD) []string {
	res := make([]string, len(qs))
	for i, q := range qs {
		x, y := q.x, q.y
		dir := q.dir
		t := q.t
		for t > 0 {
			var dx, dy int64
			switch dir {
			case 'R':
				dx = 1
			case 'L':
				dx = -1
			case 'U':
				dy = 1
			case 'D':
				dy = -1
			}
			distBound := int64(0)
			if dx > 0 {
				distBound = b - x
			} else if dx < 0 {
				distBound = x
			} else if dy > 0 {
				distBound = b - y
			} else if dy < 0 {
				distBound = y
			}
			bestDist := distBound + 1
			var seg *Segment
			if dx != 0 {
				for j := range vert {
					s := &vert[j]
					if y < s.low || y > s.high {
						continue
					}
					d := (s.fixed - x) * dx
					if d > 0 && d < bestDist {
						bestDist = d
						seg = s
					}
				}
			} else {
				for j := range hor {
					s := &hor[j]
					if x < s.low || x > s.high {
						continue
					}
					d := (s.fixed - y) * dy
					if d > 0 && d < bestDist {
						bestDist = d
						seg = s
					}
				}
			}
			if seg == nil || bestDist > distBound {
				move := t
				if move > distBound {
					move = distBound
				}
				x += dx * move
				y += dy * move
				break
			}
			if t < bestDist {
				x += dx * t
				y += dy * t
				break
			}
			x += dx * bestDist
			y += dy * bestDist
			t -= bestDist
			if seg.vertical {
				dseg := seg.headCoord - y
				sign := int64(1)
				if dseg < 0 {
					sign = -1
				}
				lenSeg := dseg * sign
				if t < lenSeg {
					y += sign * t
					break
				}
				y = seg.headCoord
				t -= lenSeg
			} else {
				dseg := seg.headCoord - x
				sign := int64(1)
				if dseg < 0 {
					sign = -1
				}
				lenSeg := dseg * sign
				if t < lenSeg {
					x += sign * t
					break
				}
				x = seg.headCoord
				t -= lenSeg
			}
			dir = seg.headDir
		}
		res[i] = fmt.Sprintf("%d %d", x, y)
	}
	return res
}

func genCaseD(rng *rand.Rand) testCaseD {
	b := int64(rng.Intn(10) + 5)
	nSeg := rng.Intn(3)
	var vert []Segment
	var hor []Segment
	usedX := make(map[int64]bool)
	for i := 0; i < nSeg; i++ {
		x := int64(rng.Intn(int(b-1)) + 1)
		for usedX[x] {
			x = int64(rng.Intn(int(b-1)) + 1)
		}
		usedX[x] = true
		y0 := int64(rng.Intn(int(b + 1)))
		y1 := int64(rng.Intn(int(b + 1)))
		for y1 == y0 {
			y1 = int64(rng.Intn(int(b + 1)))
		}
		s := Segment{vertical: true, fixed: x}
		if y0 < y1 {
			s.low = y0
			s.high = y1
			s.headCoord = y1
			s.headDir = 'U'
		} else {
			s.low = y1
			s.high = y0
			s.headCoord = y1
			s.headDir = 'D'
		}
		vert = append(vert, s)
	}
	qcnt := rng.Intn(5) + 1
	var qs []queryD
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", nSeg, b)
	for _, s := range vert {
		if s.vertical {
			if s.headDir == 'U' {
				fmt.Fprintf(&sb, "%d %d %d %d\n", s.fixed, s.low, s.fixed, s.high)
			} else {
				fmt.Fprintf(&sb, "%d %d %d %d\n", s.fixed, s.high, s.fixed, s.low)
			}
		} else {
			// not used
		}
	}
	fmt.Fprintf(&sb, "%d\n", qcnt)
	dirs := []byte{'U', 'D', 'L', 'R'}
	for i := 0; i < qcnt; i++ {
		x := int64(rng.Intn(int(b + 1)))
		y := int64(rng.Intn(int(b + 1)))
		dir := dirs[rng.Intn(len(dirs))]
		t := int64(rng.Intn(10) + 1)
		qs = append(qs, queryD{x: x, y: y, dir: dir, t: t})
		fmt.Fprintf(&sb, "%d %d %c %d\n", x, y, dir, t)
	}
	expected := solveQueriesD(b, vert, hor, qs)
	return testCaseD{input: sb.String(), expected: expected}
}

func runCaseD(bin string, tc testCaseD) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	gotLines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(gotLines) != len(tc.expected) {
		return fmt.Errorf("expected %d lines got %d", len(tc.expected), len(gotLines))
	}
	for i, l := range gotLines {
		if strings.TrimSpace(l) != tc.expected[i] {
			return fmt.Errorf("expected %q got %q", tc.expected[i], l)
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
	for i := 0; i < 100; i++ {
		tc := genCaseD(rng)
		if err := runCaseD(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
