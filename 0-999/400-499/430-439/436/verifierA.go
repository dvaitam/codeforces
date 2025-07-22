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

type Candy struct {
	t int
	h int
	m int
}

// simulate logic copied from solution
func simulate(startType int, candies0, candies1 []Candy, initialX int) int {
	x := initialX
	next := startType
	eaten := 0
	used0 := make([]bool, len(candies0))
	used1 := make([]bool, len(candies1))
	for {
		bestIdx, bestMass := -1, -1
		if next == 0 {
			for i, c := range candies0 {
				if !used0[i] && c.h <= x && c.m > bestMass {
					bestMass = c.m
					bestIdx = i
				}
			}
			if bestIdx == -1 {
				break
			}
			used0[bestIdx] = true
			x += candies0[bestIdx].m
		} else {
			for i, c := range candies1 {
				if !used1[i] && c.h <= x && c.m > bestMass {
					bestMass = c.m
					bestIdx = i
				}
			}
			if bestIdx == -1 {
				break
			}
			used1[bestIdx] = true
			x += candies1[bestIdx].m
		}
		eaten++
		next ^= 1
	}
	return eaten
}

func expected(n, x int, candies []Candy) string {
	var c0, c1 []Candy
	for _, c := range candies {
		if c.t == 0 {
			c0 = append(c0, c)
		} else {
			c1 = append(c1, c)
		}
	}
	a0 := simulate(0, c0, c1, x)
	a1 := simulate(1, c0, c1, x)
	if a1 > a0 {
		return fmt.Sprintf("%d", a1)
	}
	return fmt.Sprintf("%d", a0)
}

func runCase(bin, input, expected string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
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
		n := rng.Intn(8) + 1
		x := rng.Intn(20) + 1
		candies := make([]Candy, n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, x))
		for j := 0; j < n; j++ {
			t := rng.Intn(2)
			h := rng.Intn(20) + 1
			m := rng.Intn(20) + 1
			candies[j] = Candy{t: t, h: h, m: m}
			sb.WriteString(fmt.Sprintf("%d %d %d\n", t, h, m))
		}
		input := sb.String()
		exp := expected(n, x, candies)
		if err := runCase(bin, input, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
