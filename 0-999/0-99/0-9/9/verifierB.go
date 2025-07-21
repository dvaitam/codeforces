package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runCandidate(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveCase(n, vb, vs int, xs []int, xu, yu int) int {
	eps := 1e-9
	bestIdx := 2
	x := xs[1]
	bestTime := float64(x)/float64(vb) + math.Hypot(float64(xu-x), float64(yu))/float64(vs)
	bestDist := math.Hypot(float64(xu-x), float64(yu))
	for i := 2; i < n; i++ {
		xi := xs[i]
		tBus := float64(xi) / float64(vb)
		dRun := math.Hypot(float64(xu-xi), float64(yu))
		tTotal := tBus + dRun/float64(vs)
		if tTotal+eps < bestTime {
			bestTime = tTotal
			bestDist = dRun
			bestIdx = i + 1
		} else if math.Abs(tTotal-bestTime) <= eps && dRun < bestDist {
			bestDist = dRun
			bestIdx = i + 1
		}
	}
	return bestIdx
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(9) + 2 // 2..10 stops
	vb := rng.Intn(100) + 1
	vs := rng.Intn(100) + 1
	xs := make([]int, n)
	cur := 0
	for i := 0; i < n; i++ {
		if i == 0 {
			cur = 0
		} else {
			cur += rng.Intn(20) + 1
		}
		xs[i] = cur
	}
	xu := rng.Intn(200) + 1
	yu := rng.Intn(200) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, vb, vs))
	for i, v := range xs {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d %d\n", xu, yu))
	expect := fmt.Sprintf("%d", solveCase(n, vb, vs, xs, xu, yu))
	return sb.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
