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

type vec struct{ dx, dy int64 }

func solveE(n int, xs, ys []int64) string {
	calc := func(order []int) int64 {
		var best int64
		for t := 0; t < 2; t++ {
			seq := make([]int, 0, n)
			l, r := 0, n-1
			for i := 0; i < n; i++ {
				if (i%2 == 0) == (t == 0) {
					seq = append(seq, order[l])
					l++
				} else {
					seq = append(seq, order[r])
					r--
				}
			}
			var sum int64
			for i := 0; i < n; i++ {
				j := seq[i]
				k := seq[(i+1)%n]
				dx := xs[j] - xs[k]
				if dx < 0 {
					dx = -dx
				}
				dy := ys[j] - ys[k]
				if dy < 0 {
					dy = -dy
				}
				sum += dx + dy
			}
			if sum > best {
				best = sum
			}
		}
		return best
	}
	dirs := []vec{{1, 0}, {0, 1}, {1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
	idx := make([]int, n)
	var ans int64
	for _, v := range dirs {
		for i := range idx {
			idx[i] = i
		}
		sort.Slice(idx, func(i, j int) bool {
			return v.dx*xs[idx[i]]+v.dy*ys[idx[i]] < v.dx*xs[idx[j]]+v.dy*ys[idx[j]]
		})
		if val := calc(idx); val > ans {
			ans = val
		}
	}
	return fmt.Sprintf("%d", ans)
}

func generateCaseE(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 3
	xs := make([]int64, n)
	ys := make([]int64, n)
	used := make(map[[2]int64]bool)
	for i := 0; i < n; i++ {
		for {
			x := rng.Int63n(50)
			y := rng.Int63n(50)
			if !used[[2]int64{x, y}] {
				used[[2]int64{x, y}] = true
				xs[i] = x
				ys[i] = y
				break
			}
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", xs[i], ys[i]))
	}
	out := solveE(n, xs, ys)
	return sb.String(), out
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseE(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
