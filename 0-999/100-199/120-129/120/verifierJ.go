package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

// ---- Embedded correct solver for 120 J ----
// Uses closest-pair divide-and-conquer on abs(x), abs(y) coordinates.

type pt struct {
	origX, origY int
	x, y         int
	id           int
}

func iabs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func closestPair(pts []pt) (int64, pt, pt) {
	if len(pts) <= 3 {
		var bestD int64 = -1
		var p1, p2 pt
		for i := 0; i < len(pts); i++ {
			for j := i + 1; j < len(pts); j++ {
				dx := int64(pts[i].x - pts[j].x)
				dy := int64(pts[i].y - pts[j].y)
				d := dx*dx + dy*dy
				if bestD == -1 || d < bestD {
					bestD = d
					p1 = pts[i]
					p2 = pts[j]
				}
			}
		}
		return bestD, p1, p2
	}

	mid := len(pts) / 2
	midX := pts[mid].x

	d1, p1a, p1b := closestPair(pts[:mid])
	d2, p2a, p2b := closestPair(pts[mid:])

	d := d1
	pa, pb := p1a, p1b
	if d2 < d {
		d = d2
		pa, pb = p2a, p2b
	}

	var strip []pt
	for _, p := range pts {
		dx := int64(p.x - midX)
		if dx*dx < d {
			strip = append(strip, p)
		}
	}

	sort.Slice(strip, func(i, j int) bool {
		return strip[i].y < strip[j].y
	})

	for i := 0; i < len(strip); i++ {
		for j := i + 1; j < len(strip); j++ {
			dy := int64(strip[j].y - strip[i].y)
			if dy*dy >= d {
				break
			}
			dx := int64(strip[i].x - strip[j].x)
			dist := dx*dx + dy*dy
			if dist < d {
				d = dist
				pa = strip[i]
				pb = strip[j]
			}
		}
	}

	return d, pa, pb
}

func getK(origX, origY, targetX, targetY int) int {
	sx := 1
	if origX*1 != targetX && origX*-1 == targetX {
		sx = -1
	}
	sy := 1
	if origY*1 != targetY && origY*-1 == targetY {
		sy = -1
	}

	if sx == 1 && sy == 1 {
		return 1
	}
	if sx == -1 && sy == 1 {
		return 2
	}
	if sx == 1 && sy == -1 {
		return 3
	}
	return 4
}

func solveJ(input string) string {
	fields := strings.Fields(input)
	idx := 0
	nextInt := func() int {
		v, _ := strconv.Atoi(fields[idx])
		idx++
		return v
	}

	n := nextInt()
	pts := make([]pt, n)
	for i := 0; i < n; i++ {
		x := nextInt()
		y := nextInt()
		pts[i] = pt{
			origX: x,
			origY: y,
			x:     iabs(x),
			y:     iabs(y),
			id:    i + 1,
		}
	}

	sort.Slice(pts, func(i, j int) bool {
		if pts[i].x == pts[j].x {
			return pts[i].y < pts[j].y
		}
		return pts[i].x < pts[j].x
	})

	_, pa, pb := closestPair(pts)

	k1 := getK(pa.origX, pa.origY, iabs(pa.origX), iabs(pa.origY))
	k2 := getK(pb.origX, pb.origY, -iabs(pb.origX), -iabs(pb.origY))

	return fmt.Sprintf("%d %d %d %d", pa.id, k1, pb.id, k2)
}

// ---- Verifier: check that candidate's answer is optimal ----

func applyK(x, y, k int) (int, int) {
	switch k {
	case 1:
		return x, y
	case 2:
		return -x, y
	case 3:
		return x, -y
	case 4:
		return -x, -y
	}
	return x, y
}

func magnitude(x, y int) float64 {
	return math.Sqrt(float64(x)*float64(x) + float64(y)*float64(y))
}

func bruteForceOptimal(xs, ys []int) float64 {
	n := len(xs)
	best := math.Inf(1)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			for k1 := 1; k1 <= 4; k1++ {
				for k2 := 1; k2 <= 4; k2++ {
					x1, y1 := applyK(xs[i], ys[i], k1)
					x2, y2 := applyK(xs[j], ys[j], k2)
					d := magnitude(x1+x2, y1+y2)
					if d < best {
						best = d
					}
				}
			}
		}
	}
	return best
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: go run verifierJ.go /path/to/binary\n")
		os.Exit(1)
	}
	candidate := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		n := rng.Intn(4) + 2
		xs := make([]int, n)
		ys := make([]int, n)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for i := 0; i < n; i++ {
			xs[i] = rng.Intn(21) - 10
			ys[i] = rng.Intn(21) - 10
			fmt.Fprintf(&sb, "%d %d\n", xs[i], ys[i])
		}
		input := sb.String()

		// Run candidate
		candOut, cErr := runBinary(candidate, input)
		if cErr != nil {
			fmt.Fprintf(os.Stderr, "test %d: candidate error: %v\n", t+1, cErr)
			os.Exit(1)
		}

		candFields := strings.Fields(strings.TrimSpace(candOut))
		if len(candFields) != 4 {
			fmt.Fprintf(os.Stderr, "test %d: expected 4 values, got %d: %q\n", t+1, len(candFields), strings.TrimSpace(candOut))
			os.Exit(1)
		}

		ci, _ := strconv.Atoi(candFields[0])
		ck1, _ := strconv.Atoi(candFields[1])
		cj, _ := strconv.Atoi(candFields[2])
		ck2, _ := strconv.Atoi(candFields[3])

		if ci < 1 || ci > n || cj < 1 || cj > n || ci == cj || ck1 < 1 || ck1 > 4 || ck2 < 1 || ck2 > 4 {
			fmt.Fprintf(os.Stderr, "test %d: invalid output: %s\n", t+1, strings.TrimSpace(candOut))
			os.Exit(1)
		}

		x1, y1 := applyK(xs[ci-1], ys[ci-1], ck1)
		x2, y2 := applyK(xs[cj-1], ys[cj-1], ck2)
		candMag := magnitude(x1+x2, y1+y2)

		optMag := bruteForceOptimal(xs, ys)

		if candMag > optMag+1e-9 {
			fmt.Fprintf(os.Stderr, "test %d failed: candidate magnitude %.10f > optimal %.10f\ninput:\n%scandidate output: %s\n", t+1, candMag, optMag, input, strings.TrimSpace(candOut))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}
