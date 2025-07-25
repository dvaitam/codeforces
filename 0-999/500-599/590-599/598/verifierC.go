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

type vec struct {
	idx   int
	x, y  int
	angle float64
}

func expectedC(v []vec) (int, int) {
	n := len(v)
	for i := range v {
		v[i].angle = math.Atan2(float64(v[i].y), float64(v[i].x))
	}
	sort.Slice(v, func(i, j int) bool { return v[i].angle < v[j].angle })
	ans1, ans2 := v[0].idx, v[n-1].idx
	minDiff := v[0].angle + 2*math.Pi - v[n-1].angle
	for i := 0; i < n-1; i++ {
		diff := v[i+1].angle - v[i].angle
		if diff < minDiff {
			minDiff = diff
			ans1 = v[i].idx
			ans2 = v[i+1].idx
		}
	}
	return ans1, ans2
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(18) + 2
	vectors := make([]vec, n)
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(n) + "\n")
	for i := 0; i < n; i++ {
		x := rng.Intn(2001) - 1000
		y := rng.Intn(2001) - 1000
		if x == 0 && y == 0 {
			x = 1
		}
		vectors[i] = vec{idx: i + 1, x: x, y: y}
		sb.WriteString(fmt.Sprintf("%d %d\n", x, y))
	}
	a, b := expectedC(append([]vec(nil), vectors...))
	out := fmt.Sprintf("%d %d\n", a, b)
	return sb.String(), out
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
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
