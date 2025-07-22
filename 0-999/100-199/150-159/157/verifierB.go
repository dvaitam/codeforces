package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func expected(radii []int) float64 {
	sort.Ints(radii)
	n := len(radii)
	var s float64
	for i, r := range radii {
		area := math.Pi * float64(r*r)
		if n%2 == 1 {
			if i%2 == 0 {
				s += area
			} else {
				s -= area
			}
		} else {
			if i%2 == 0 {
				s -= area
			} else {
				s += area
			}
		}
	}
	return s
}

func runCase(bin string, radii []int) error {
	var input strings.Builder
	n := len(radii)
	input.WriteString(fmt.Sprintf("%d\n", n))
	for i, r := range radii {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(fmt.Sprintf("%d", r))
	}
	input.WriteByte('\n')

	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input.String())
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	var got float64
	if _, err := fmt.Fscan(&buf, &got); err != nil {
		return fmt.Errorf("failed to parse output: %v", err)
	}
	exp := expected(append([]int(nil), radii...))
	if math.Abs(got-exp) > 1e-4*math.Max(1, math.Abs(exp)) {
		return fmt.Errorf("expected %.5f got %.5f\ninput:\n%s", exp, got, input.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(100) + 1
		radii := make([]int, n)
		used := make(map[int]bool)
		for j := 0; j < n; j++ {
			for {
				r := rng.Intn(1000) + 1
				if !used[r] {
					used[r] = true
					radii[j] = r
					break
				}
			}
		}
		if err := runCase(bin, radii); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
