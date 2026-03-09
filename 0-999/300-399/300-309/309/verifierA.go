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

// reference solution: count pairs of frogs that can meet on a circle
func solve(n, l, t int, a []int) float64 {
	sorted := make([]int, n)
	copy(sorted, a)
	sort.Ints(sorted)

	// duplicate for circular handling
	b := make([]int, 2*n)
	for i := 0; i < n; i++ {
		b[i] = sorted[i]
		b[n+i] = sorted[i] + l
	}

	t2 := 2 * t
	x := t2 / l
	t2 %= l

	var r float64
	j := 0
	for i := 0; i < n; i++ {
		if j < i {
			j = i
		}
		for j < 2*n && b[j]-b[i] <= t2 {
			j++
		}
		r += float64(j - i - 1)
	}

	return 0.25 * (r + float64(x)*float64(n)*float64(n-1))
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(10) + 2
	l := rng.Intn(1000) + n
	t := rng.Intn(500) + 1

	positions := make([]int, n)
	used := make(map[int]bool)
	for i := 0; i < n; i++ {
		for {
			p := rng.Intn(l)
			if !used[p] {
				positions[i] = p
				used[p] = true
				break
			}
		}
	}
	sort.Ints(positions)

	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, l, t)
	for i, p := range positions {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", p)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func parseInput(input string) (n, l, t int, a []int) {
	fields := strings.Fields(input)
	n, _ = strconv.Atoi(fields[0])
	l, _ = strconv.Atoi(fields[1])
	t, _ = strconv.Atoi(fields[2])
	a = make([]int, n)
	for i := 0; i < n; i++ {
		a[i], _ = strconv.Atoi(fields[3+i])
	}
	return
}

func runCase(exe, input string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}

	outStr := strings.TrimSpace(out.String())
	got, err := strconv.ParseFloat(outStr, 64)
	if err != nil {
		return fmt.Errorf("bad output: %v", err)
	}

	n, l, t, a := parseInput(input)
	expected := solve(n, l, t, a)

	if math.Abs(got-expected) > 1e-4 {
		return fmt.Errorf("expected %.8f got %.8f", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		in := generateCase(rng)
		if err := runCase(exe, in); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
