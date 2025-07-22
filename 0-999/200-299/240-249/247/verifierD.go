package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solveD(input string) string {
	r := bufio.NewReader(strings.NewReader(input))
	var n, m, a, b int
	if _, err := fmt.Fscan(r, &n, &m, &a, &b); err != nil {
		return ""
	}
	ys := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(r, &ys[i])
	}
	yps := make([]int, m)
	for j := 0; j < m; j++ {
		fmt.Fscan(r, &yps[j])
	}
	ls := make([]int, m)
	for j := 0; j < m; j++ {
		fmt.Fscan(r, &ls[j])
	}
	D := make([]float64, n)
	aa := float64(a)
	for i := 0; i < n; i++ {
		yi := float64(ys[i])
		D[i] = math.Hypot(aa, yi)
	}
	deltaX := float64(b - a)
	bestI := 0
	bestJ := 0
	minCost := math.Inf(1)
	for j := 0; j < m; j++ {
		yj := float64(yps[j])
		for bestI+1 < n {
			dy0 := float64(ys[bestI]) - yj
			cost0 := D[bestI] + math.Hypot(deltaX, dy0)
			dy1 := float64(ys[bestI+1]) - yj
			cost1 := D[bestI+1] + math.Hypot(deltaX, dy1)
			if cost1 <= cost0 {
				bestI++
			} else {
				break
			}
		}
		dy := float64(ys[bestI]) - yj
		bridge := math.Hypot(deltaX, dy)
		total := D[bestI] + bridge + float64(ls[j])
		if total < minCost {
			minCost = total
			bestJ = j
		}
	}
	return fmt.Sprintf("%d %d\n", bestI+1, bestJ+1)
}

func generateCaseD(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	a := rng.Intn(5) + 1
	b := a + rng.Intn(5) + 1
	ys := make([]int, n)
	yps := make([]int, m)
	for i := range ys {
		if i == 0 {
			ys[i] = rng.Intn(21) - 10
		} else {
			ys[i] = ys[i-1] + rng.Intn(3) + 1
		}
	}
	for j := range yps {
		if j == 0 {
			yps[j] = rng.Intn(21) - 10
		} else {
			yps[j] = yps[j-1] + rng.Intn(3) + 1
		}
	}
	ls := make([]int, m)
	for j := range ls {
		ls[j] = rng.Intn(10) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d %d\n", n, m, a, b))
	for i, v := range ys {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	for j, v := range yps {
		if j > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	for j, v := range ls {
		if j > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]string, 100)
	for i := 0; i < 100; i++ {
		cases[i] = generateCaseD(rng)
	}
	for i, tc := range cases {
		expect := solveD(tc)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%sexpected:%sq\ngot:%sq\n", i+1, tc, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
