package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func runProgram(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func solve(n int, k int64, xs, ys, rs []float64) int64 {
	M := 1
	pi := math.Pi
	for i := 0; i < n; i++ {
		type ev struct {
			a float64
			d int
		}
		events := make([]ev, 0, 2*(n-1))
		for j := 0; j < n; j++ {
			if j == i {
				continue
			}
			dx := xs[j] - xs[i]
			dy := ys[j] - ys[i]
			dist := math.Hypot(dx, dy)
			if dist <= rs[j] {
				continue
			}
			phi := math.Atan2(dy, dx)
			ang := math.Asin(rs[j] / dist)
			width := ang * 2.0
			L := math.Mod(phi-ang, pi)
			if L < 0 {
				L += pi
			}
			R := L + width
			if R < pi {
				events = append(events, ev{L, 1}, ev{R, -1})
			} else {
				events = append(events, ev{L, 1}, ev{pi, -1}, ev{0, 1}, ev{R - pi, -1})
			}
		}
		if len(events) == 0 {
			continue
		}
		sort.Slice(events, func(a, b int) bool {
			if events[a].a == events[b].a {
				return events[a].d > events[b].d
			}
			return events[a].a < events[b].a
		})
		cnt, best := 0, 0
		for _, e := range events {
			cnt += e.d
			if cnt > best {
				best = cnt
			}
		}
		if best+1 > M {
			M = best + 1
		}
	}
	return int64(n) + k*int64(M)
}

func randomTest(rng *rand.Rand) (int, int64, []float64, []float64, []float64) {
	n := rng.Intn(4) + 1
	k := int64(rng.Intn(6))
	xs := make([]float64, n)
	ys := make([]float64, n)
	rs := make([]float64, n)
	for i := 0; i < n; i++ {
		for {
			xs[i] = rng.Float64()*10 - 5
			ys[i] = rng.Float64()*10 - 5
			rs[i] = rng.Float64()*2 + 0.5
			ok := true
			for j := 0; j < i; j++ {
				dx := xs[i] - xs[j]
				dy := ys[i] - ys[j]
				d := math.Hypot(dx, dy)
				if d <= rs[i]+rs[j] {
					ok = false
					break
				}
			}
			if ok {
				break
			}
		}
	}
	return n, k, xs, ys, rs
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 100; i++ {
		n, k, xs, ys, rs := randomTest(rng)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for j := 0; j < n; j++ {
			sb.WriteString(fmt.Sprintf("%.2f %.2f %.2f\n", xs[j], ys[j], rs[j]))
		}
		expected := solve(n, k, xs, ys, rs)
		out, err := runProgram(bin, sb.String())
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			return
		}
		if strings.TrimSpace(out) != fmt.Sprintf("%d", expected) {
			fmt.Printf("Test %d failed: expected %d got %s\n", i+1, expected, strings.TrimSpace(out))
			return
		}
	}
	fmt.Println("All tests passed")
}
