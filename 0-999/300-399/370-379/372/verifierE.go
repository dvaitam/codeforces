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

const mod = 1000000007

type point struct{ a, b, c, d int }

type pair struct{ dot, cross int64 }

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func solveCase(pts []point) int64 {
	n := len(pts)
	xs := make([]int64, n)
	ys := make([]int64, n)
	for i, p := range pts {
		xs[i] = int64(p.a) * int64(p.d)
		ys[i] = int64(p.c) * int64(p.b)
	}
	counts := make(map[pair]int)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			dot := xs[i]*xs[j] - ys[i]*ys[j]
			cross := xs[i]*ys[j] + ys[i]*xs[j]
			g := gcd(dot, cross)
			if g != 0 {
				dot /= g
				cross /= g
			}
			counts[pair{dot, cross}]++
		}
	}
	var ans int64
	for _, m := range counts {
		if m >= 2 {
			add := modPow(2, int64(m)) - 1 - int64(m)
			add %= mod
			if add < 0 {
				add += mod
			}
			ans = (ans + add) % mod
		}
	}
	return ans
}

func modPow(a, e int64) int64 {
	res := int64(1)
	a %= mod
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func runCase(bin string, pts []point) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(pts)))
	for _, p := range pts {
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", p.a, p.b, p.c, p.d))
	}
	input := sb.String()
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	expected := solveCase(pts)
	var got int64
	fmt.Sscan(strings.TrimSpace(out.String()), &got)
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const tests = 100
	for i := 0; i < tests; i++ {
		n := rng.Intn(5) + 2
		pts := make([]point, n)
		for j := 0; j < n; j++ {
			for {
				a := rng.Intn(5) - 2
				c := rng.Intn(5) - 2
				if a != 0 || c != 0 {
					pts[j].a = a
					pts[j].b = rng.Intn(3) + 1
					pts[j].c = c
					pts[j].d = rng.Intn(3) + 1
					break
				}
			}
		}
		if err := runCase(bin, pts); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", tests)
}
