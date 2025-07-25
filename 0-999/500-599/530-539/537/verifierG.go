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

func abs64(a int64) int64 {
	if a < 0 {
		return -a
	}
	return a
}

type record struct {
	t int64
	x int64
	y int64
}

func expectedAnswerG(l int64, recs []record) string {
	type Occ struct{ c0, x0, y0 int64 }
	occ := make(map[int64]*Occ, len(recs))
	var dxGlobal, dyGlobal int64
	var hasDX, hasDY bool
	for _, r := range recs {
		c := r.t / l
		mod := r.t % l
		if o, ok := occ[mod]; !ok {
			occ[mod] = &Occ{c0: c, x0: r.x, y0: r.y}
		} else {
			d := c - o.c0
			dx := r.x - o.x0
			dy := r.y - o.y0
			if d <= 0 || dx%d != 0 || dy%d != 0 {
				return "NO"
			}
			vx := dx / d
			vy := dy / d
			if !hasDX {
				dxGlobal = vx
				hasDX = true
			} else if dxGlobal != vx {
				return "NO"
			}
			if !hasDY {
				dyGlobal = vy
				hasDY = true
			} else if dyGlobal != vy {
				return "NO"
			}
		}
	}
	if !hasDX {
		dxGlobal = 0
	}
	if !hasDY {
		dyGlobal = 0
	}
	type P struct {
		idx  int
		x, y int64
	}
	fps := []P{{idx: 0, x: 0, y: 0}}
	for r, o := range occ {
		sx := o.x0 - o.c0*dxGlobal
		sy := o.y0 - o.c0*dyGlobal
		if r == 0 {
			if sx != 0 || sy != 0 {
				return "NO"
			}
			continue
		}
		fps = append(fps, P{idx: int(r), x: sx, y: sy})
	}
	fps = append(fps, P{idx: int(l), x: dxGlobal, y: dyGlobal})
	sort.Slice(fps, func(i, j int) bool { return fps[i].idx < fps[j].idx })
	res := make([]byte, l)
	for i := 0; i+1 < len(fps); i++ {
		a := fps[i]
		b := fps[i+1]
		dt := b.idx - a.idx
		dx := b.x - a.x
		dy := b.y - a.y
		D := abs64(dx) + abs64(dy)
		if D > int64(dt) || (int64(dt)-D)%2 != 0 {
			return "NO"
		}
		pos := a.idx
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
		for k := 0; k < rem; k += 2 {
			res[pos] = 'L'
			res[pos+1] = 'R'
			pos += 2
		}
	}
	return string(res)
}

func generateCaseG(rng *rand.Rand) (int, int64, []record) {
	n := rng.Intn(5) + 1
	l := int64(rng.Intn(10) + 1)
	recs := make([]record, n)
	var t int64
	for i := 0; i < n; i++ {
		t += int64(rng.Intn(int(l*3)) + 1)
		x := int64(rng.Intn(11) - 5)
		y := int64(rng.Intn(11) - 5)
		recs[i] = record{t, x, y}
	}
	return n, int64(l), recs
}

func runCaseG(bin string, n int, l int64, recs []record) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, l))
	for _, r := range recs {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", r.t, r.x, r.y))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expected := expectedAnswerG(l, recs)
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, l, recs := generateCaseG(rng)
		if err := runCaseG(bin, n, l, recs); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
