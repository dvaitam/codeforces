package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type point struct{ x, y int }

func generateCase(rng *rand.Rand) (int, []point, [9]int, float64, float64, float64, float64) {
	n := rng.Intn(40) + 9
	pts := make([]point, n)
	for i := range pts {
		pts[i] = point{rng.Intn(21) - 10, rng.Intn(21) - 10}
	}
	x1 := float64(rng.Intn(19)-9) + 0.5
	x2 := float64(rng.Intn(19)-9) + 0.5
	for x2 == x1 {
		x2 = float64(rng.Intn(19)-9) + 0.5
	}
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	y1 := float64(rng.Intn(19)-9) + 0.5
	y2 := float64(rng.Intn(19)-9) + 0.5
	for y2 == y1 {
		y2 = float64(rng.Intn(19)-9) + 0.5
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	var cnt [9]int
	for _, p := range pts {
		region := 0
		if float64(p.x) < x1 {
			if float64(p.y) < y1 {
				region = 0
			} else if float64(p.y) < y2 {
				region = 3
			} else {
				region = 6
			}
		} else if float64(p.x) < x2 {
			if float64(p.y) < y1 {
				region = 1
			} else if float64(p.y) < y2 {
				region = 4
			} else {
				region = 7
			}
		} else {
			if float64(p.y) < y1 {
				region = 2
			} else if float64(p.y) < y2 {
				region = 5
			} else {
				region = 8
			}
		}
		cnt[region]++
	}
	return n, pts, cnt, x1, x2, y1, y2
}

func checkOutput(x1, x2, y1, y2 float64, pts []point, want [9]int) error {
	if math.Abs(x1-x2) <= 1e-6 || math.Abs(y1-y2) <= 1e-6 {
		return fmt.Errorf("lines coincide")
	}
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	var got [9]int
	for _, p := range pts {
		if math.Abs(float64(p.x)-x1) <= 1e-6 || math.Abs(float64(p.x)-x2) <= 1e-6 || math.Abs(float64(p.y)-y1) <= 1e-6 || math.Abs(float64(p.y)-y2) <= 1e-6 {
			return fmt.Errorf("point on line")
		}
		region := 0
		if float64(p.x) < x1 {
			if float64(p.y) < y1 {
				region = 0
			} else if float64(p.y) < y2 {
				region = 3
			} else {
				region = 6
			}
		} else if float64(p.x) < x2 {
			if float64(p.y) < y1 {
				region = 1
			} else if float64(p.y) < y2 {
				region = 4
			} else {
				region = 7
			}
		} else {
			if float64(p.y) < y1 {
				region = 2
			} else if float64(p.y) < y2 {
				region = 5
			} else {
				region = 8
			}
		}
		got[region]++
	}
	for i := 0; i < 9; i++ {
		if got[i] != want[i] {
			return fmt.Errorf("region %d expected %d got %d", i+1, want[i], got[i])
		}
	}
	return nil
}

func runCase(bin string, n int, pts []point, cnt [9]int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, p := range pts {
		sb.WriteString(fmt.Sprintf("%d %d\n", p.x, p.y))
	}
	for i := 0; i < 9; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(cnt[i]))
	}
	sb.WriteByte('\n')

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	parts := strings.Fields(strings.TrimSpace(out.String()))
	if len(parts) != 4 {
		return fmt.Errorf("expected 4 numbers")
	}
	x1, _ := strconv.ParseFloat(parts[0], 64)
	x2, _ := strconv.ParseFloat(parts[1], 64)
	y1, _ := strconv.ParseFloat(parts[2], 64)
	y2, _ := strconv.ParseFloat(parts[3], 64)
	return checkOutput(x1, x2, y1, y2, pts, cnt)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, pts, cnt, _, _, _, _ := generateCase(rng)
		if err := runCase(bin, n, pts, cnt); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
