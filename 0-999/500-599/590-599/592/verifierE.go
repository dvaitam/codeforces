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

func expected(n int, c, d int64, arr [][2]int64) int64 {
	angles := make([]float64, n)
	twoPi := 2 * math.Pi
	for i := 0; i < n; i++ {
		x := float64(arr[i][0] - c)
		y := float64(arr[i][1] - d)
		ang := math.Atan2(y, x)
		if ang < 0 {
			ang += twoPi
		}
		angles[i] = ang
	}
	sort.Float64s(angles)
	ext := make([]float64, 2*n)
	copy(ext, angles)
	for i := 0; i < n; i++ {
		ext[i+n] = angles[i] + twoPi
	}
	var bad int64
	j := 0
	eps := 1e-12
	for i := 0; i < n; i++ {
		if j < i+1 {
			j = i + 1
		}
		for j < i+n && ext[j]-ext[i] <= math.Pi+eps {
			j++
		}
		m := j - i - 1
		if m >= 2 {
			bad += int64(m*(m-1)) / 2
		}
	}
	total := int64(n) * int64(n-1) * int64(n-2) / 6
	good := total - bad
	return good
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(30) + 3
	c := int64(rng.Intn(1000) + 1)
	d := int64(rng.Intn(1000) + 1)
	data := make([][2]int64, n)
	for i := 0; i < n; i++ {
		r := int64(rng.Intn(1000) + 1)
		w := int64(rng.Intn(1000) + 1)
		if r == c && w == d {
			r++
		}
		data[i] = [2]int64{r, w}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, c, d))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", data[i][0], data[i][1]))
	}
	ans := expected(n, c, d, data)
	return sb.String(), fmt.Sprintf("%d", ans)
}

func runCase(exe, input, expected string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(exe, ".go") {
		cmd = exec.Command("go", "run", exe)
	} else {
		cmd = exec.Command(exe)
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
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
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
