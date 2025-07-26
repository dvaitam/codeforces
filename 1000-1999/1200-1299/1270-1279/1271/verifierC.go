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

func expected(n, sx, sy int, pts [][2]int) string {
	x1, x2, y1, y2 := 0, 0, 0, 0
	for _, p := range pts {
		x := p[0]
		y := p[1]
		if x < sx {
			x1++
		}
		if x > sx {
			x2++
		}
		if y < sy {
			y1++
		}
		if y > sy {
			y2++
		}
	}
	ans := x1
	ansx, ansy := sx, sy
	if x2 > ans {
		ans = x2
	}
	if y1 > ans {
		ans = y1
	}
	if y2 > ans {
		ans = y2
	}
	if ans == x1 {
		ansx = sx - 1
		ansy = sy
	} else if ans == x2 {
		ansx = sx + 1
		ansy = sy
	} else if ans == y1 {
		ansx = sx
		ansy = sy - 1
	} else if ans == y2 {
		ansx = sx
		ansy = sy + 1
	}
	return fmt.Sprintf("%d\n%d %d", ans, ansx, ansy)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(1000) + 1
	sx := rng.Intn(1000)
	sy := rng.Intn(1000)
	pts := make([][2]int, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, sx, sy)
	for i := 0; i < n; i++ {
		x := rng.Intn(1000)
		y := rng.Intn(1000)
		if x == sx && y == sy {
			if x > 0 {
				x--
			} else {
				x++
			}
		}
		pts[i] = [2]int{x, y}
		fmt.Fprintf(&sb, "%d %d\n", x, y)
	}
	return sb.String(), expected(n, sx, sy, pts)
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expected = strings.TrimSpace(expected)
	if got != expected {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expect := generateCase(rng)
		if err := runCase(exe, input, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
