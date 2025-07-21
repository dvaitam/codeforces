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

type node struct {
	x, y, t int
	p       float64
}

func solve(nodes []node) float64 {
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].t < nodes[j].t })
	n := len(nodes)
	dp := make([]float64, n)
	res := 0.0
	eps := 1e-9
	for i := 0; i < n; i++ {
		dp[i] = nodes[i].p
		for j := 0; j < i; j++ {
			dt := float64(nodes[i].t - nodes[j].t)
			dx := float64(nodes[i].x - nodes[j].x)
			dy := float64(nodes[i].y - nodes[j].y)
			if dt*dt-dx*dx-dy*dy > -eps {
				if dp[j]+nodes[i].p > dp[i] {
					dp[i] = dp[j] + nodes[i].p
				}
			}
		}
		if dp[i] > res {
			res = dp[i]
		}
	}
	return res
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	nodes := make([]node, n)
	used := map[[2]int]bool{}
	for i := 0; i < n; i++ {
		for {
			x := rng.Intn(11) - 5
			y := rng.Intn(11) - 5
			if !used[[2]int{x, y}] {
				used[[2]int{x, y}] = true
				nodes[i].x = x
				nodes[i].y = y
				break
			}
		}
		nodes[i].t = rng.Intn(51)
		nodes[i].p = rng.Float64()
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d %d %.6f\n", nodes[i].x, nodes[i].y, nodes[i].t, nodes[i].p))
	}
	expected := fmt.Sprintf("%.9f", solve(nodes))
	return sb.String(), expected
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
	if valOut, err := strconv.ParseFloat(outStr, 64); err == nil {
		valExp, _ := strconv.ParseFloat(expected, 64)
		if math.Abs(valOut-valExp) > 1e-6 {
			return fmt.Errorf("expected %.9f got %s", valExp, outStr)
		}
	} else {
		return fmt.Errorf("invalid output %s", outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
