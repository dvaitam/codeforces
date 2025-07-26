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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func unique(a []int64) []int64 {
	if len(a) == 0 {
		return a
	}
	j := 1
	for i := 1; i < len(a); i++ {
		if a[i] != a[i-1] {
			a[j] = a[i]
			j++
		}
	}
	return a[:j]
}

func solve(points [][2]int64) int64 {
	n := len(points)
	xs := make([]int64, n)
	ys := make([]int64, n)
	for i := 0; i < n; i++ {
		xs[i] = points[i][0]
		ys[i] = points[i][1]
	}
	ux := make([]int64, n)
	copy(ux, xs)
	sort.Slice(ux, func(i, j int) bool { return ux[i] < ux[j] })
	ux = unique(ux)
	pts := make([]struct {
		x int
		y int64
	}, n)
	for i := 0; i < n; i++ {
		xi := sort.Search(len(ux), func(j int) bool { return ux[j] >= xs[i] })
		pts[i] = struct {
			x int
			y int64
		}{xi, ys[i]}
	}
	sort.Slice(pts, func(i, j int) bool { return pts[i].y > pts[j].y })
	countX := make([]int, len(ux))
	active := 0
	var ans int64
	for i := 0; i < n; {
		yv := pts[i].y
		j := i
		for j < n && pts[j].y == yv {
			xi := pts[j].x
			if countX[xi] == 0 {
				active++
			}
			countX[xi]++
			j++
		}
		ac := int64(active)
		ans += ac * (ac + 1) / 2
		i = j
	}
	return ans
}

func genTest(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	pts := make([][2]int64, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		x := int64(rng.Intn(10))
		y := int64(rng.Intn(10))
		pts[i] = [2]int64{x, y}
		sb.WriteString(fmt.Sprintf("%d %d\n", x, y))
	}
	return sb.String(), fmt.Sprint(solve(pts))
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expected := genTest(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if out != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, expected, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
