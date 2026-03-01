package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

func nextPermutation(a []int) bool {
	i := len(a) - 2
	for i >= 0 && a[i] >= a[i+1] {
		i--
	}
	if i < 0 {
		return false
	}
	j := len(a) - 1
	for a[j] <= a[i] {
		j--
	}
	a[i], a[j] = a[j], a[i]
	for l, r := i+1, len(a)-1; l < r; l, r = l+1, r-1 {
		a[l], a[r] = a[r], a[l]
	}
	return true
}

func lexLess(x, y []int) bool {
	for i := 0; i < len(x); i++ {
		if x[i] != y[i] {
			return x[i] < y[i]
		}
	}
	return false
}

func equalSlices(x, y []int) bool {
	if len(x) != len(y) {
		return false
	}
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

func simulate(p, b []int, lessons int) []int {
	n := len(b)
	inside := append([]int(nil), b...)
	nextOutside := n + 1
	for t := 0; t < lessons; t++ {
		desk := make([][]int, n)
		for i := 0; i < n; i++ {
			desk[p[i]] = append(desk[p[i]], inside[i])
		}
		nextInside := make([]int, n)
		empty := make([]int, 0)
		for i := 0; i < n; i++ {
			if len(desk[i]) == 0 {
				empty = append(empty, i)
				continue
			}
			sort.Ints(desk[i])
			nextInside[i] = desk[i][0]
		}
		for _, idx := range empty {
			nextInside[idx] = nextOutside
			nextOutside++
		}
		inside = nextInside
	}
	return inside
}

func bruteSolve(p, a []int) ([]int, bool) {
	n := len(p)
	perm := make([]int, n)
	for i := 0; i < n; i++ {
		perm[i] = i + 1
	}

	best := []int(nil)
	for {
		for lessons := 0; lessons <= n+2; lessons++ {
			res := simulate(p, perm, lessons)
			if equalSlices(res, a) {
				cand := append([]int(nil), perm...)
				if best == nil || lexLess(cand, best) {
					best = cand
				}
				break
			}
		}
		if !nextPermutation(perm) {
			break
		}
	}
	if best == nil {
		return nil, false
	}
	return best, true
}

func formatSlice(a []int) string {
	var sb strings.Builder
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 2
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = rng.Intn(n)
	}

	// Enforce at least one duplicate in p (problem guarantee).
	i := rng.Intn(n)
	j := rng.Intn(n)
	for j == i {
		j = rng.Intn(n)
	}
	p[j] = p[i]

	perm := make([]int, n)
	for i := 0; i < n; i++ {
		perm[i] = i + 1
	}
	rng.Shuffle(n, func(i, j int) {
		perm[i], perm[j] = perm[j], perm[i]
	})
	lessons := rng.Intn(n + 2)
	a := simulate(p, perm, lessons)

	expected, ok := bruteSolve(p, a)
	if !ok {
		panic("generated unsolvable case")
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(p[i] + 1))
	}
	sb.WriteByte('\n')
	sb.WriteString(formatSlice(a))

	return sb.String(), formatSlice(expected)
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
	outStr := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if outStr != exp {
		return fmt.Errorf("expected %q got %q", exp, outStr)
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
	for i := 0; i < 300; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
