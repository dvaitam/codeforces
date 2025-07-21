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

type seg struct{ x1, y1, x2, y2 int }

func check(segs []seg) string {
	var hor, ver []seg
	for _, s := range segs {
		if s.y1 == s.y2 {
			hor = append(hor, s)
		} else if s.x1 == s.x2 {
			ver = append(ver, s)
		} else {
			return "NO"
		}
	}
	if len(hor) != 2 || len(ver) != 2 {
		return "NO"
	}
	y0, y1 := hor[0].y1, hor[1].y1
	if y0 == y1 {
		return "NO"
	}
	var yMin, yMax int
	if y0 < y1 {
		yMin, yMax = y0, y1
	} else {
		yMin, yMax = y1, y0
	}
	h0l, h0r := min(hor[0].x1, hor[0].x2), max(hor[0].x1, hor[0].x2)
	h1l, h1r := min(hor[1].x1, hor[1].x2), max(hor[1].x1, hor[1].x2)
	if h0l != h1l || h0r != h1r || h0l >= h0r {
		return "NO"
	}
	xMin, xMax := h0l, h0r
	xv0, xv1 := ver[0].x1, ver[1].x1
	if xv0 == xv1 {
		return "NO"
	}
	var vMin, vMax int
	if xv0 < xv1 {
		vMin, vMax = xv0, xv1
	} else {
		vMin, vMax = xv1, xv0
	}
	if vMin != xMin || vMax != xMax {
		return "NO"
	}
	v0l, v0r := min(ver[0].y1, ver[0].y2), max(ver[0].y1, ver[0].y2)
	v1l, v1r := min(ver[1].y1, ver[1].y2), max(ver[1].y1, ver[1].y2)
	if v0l != v1l || v0r != v1r {
		return "NO"
	}
	if v0l != yMin || v0r != yMax {
		return "NO"
	}
	return "YES"
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func generateCase(rng *rand.Rand) (string, string) {
	segs := make([]seg, 4)
	var sb strings.Builder
	for i := range segs {
		if rng.Intn(2) == 0 {
			x := rng.Intn(11) - 5
			y1 := rng.Intn(11) - 5
			y2 := rng.Intn(11) - 5
			segs[i] = seg{x, y1, x, y2}
		} else {
			y := rng.Intn(11) - 5
			x1 := rng.Intn(11) - 5
			x2 := rng.Intn(11) - 5
			segs[i] = seg{x1, y, x2, y}
		}
	}
	for _, s := range segs {
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", s.x1, s.y1, s.x2, s.y2))
	}
	exp := check(segs)
	return sb.String(), exp
}

func runCase(bin, input, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
