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

type Rect struct{ x1, y1, x2, y2 int }

func expected(rects []Rect) string {
	if len(rects) == 0 {
		return "NO"
	}
	minX, minY := 1<<30, 1<<30
	maxX, maxY := 0, 0
	var sum int64
	for _, r := range rects {
		if r.x1 < minX {
			minX = r.x1
		}
		if r.y1 < minY {
			minY = r.y1
		}
		if r.x2 > maxX {
			maxX = r.x2
		}
		if r.y2 > maxY {
			maxY = r.y2
		}
		sum += int64(r.x2-r.x1) * int64(r.y2-r.y1)
	}
	dx := maxX - minX
	dy := maxY - minY
	if dx != dy {
		return "NO"
	}
	if int64(dx)*int64(dy) != sum {
		return "NO"
	}
	return "YES"
}

func genSquare(rng *rand.Rand, n int) []Rect {
	L := rng.Intn(10) + 1
	x0 := rng.Intn(10)
	y0 := rng.Intn(10)
	rects := []Rect{{x0, y0, x0 + L, y0 + L}}
	for len(rects) < n {
		idx := rng.Intn(len(rects))
		r := rects[idx]
		if rng.Intn(2) == 0 {
			if r.x2-r.x1 < 2 {
				continue
			}
			cut := rng.Intn(r.x2-r.x1-1) + r.x1 + 1
			rects[idx] = Rect{r.x1, r.y1, cut, r.y2}
			rects = append(rects, Rect{cut, r.y1, r.x2, r.y2})
		} else {
			if r.y2-r.y1 < 2 {
				continue
			}
			cut := rng.Intn(r.y2-r.y1-1) + r.y1 + 1
			rects[idx] = Rect{r.x1, r.y1, r.x2, cut}
			rects = append(rects, Rect{r.x1, cut, r.x2, r.y2})
		}
	}
	return rects
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	rects := genSquare(rng, n)
	if rng.Intn(2) == 0 {
		// make it NO by adding extra rectangle outside
		r := rects[rng.Intn(len(rects))]
		rects = append(rects, Rect{r.x2 + 1, r.y2 + 1, r.x2 + 2, r.y2 + 2})
		n++
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		r := rects[i]
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", r.x1, r.y1, r.x2, r.y2))
	}
	return sb.String(), expected(rects[:n])
}

func runCase(bin string, input, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	ans := strings.TrimSpace(out.String())
	if ans != exp {
		return fmt.Errorf("expected %s got %s", exp, ans)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
