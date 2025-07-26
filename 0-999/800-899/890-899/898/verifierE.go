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

func isSquare(x int) bool {
	r := int(math.Sqrt(float64(x)))
	return r*r == x
}

func costToSquare(x int) int {
	r := int(math.Sqrt(float64(x)))
	c1 := x - r*r
	c2 := (r+1)*(r+1) - x
	if c1 < c2 {
		return c1
	}
	return c2
}

func solveE(a []int) int {
	n := len(a)
	squareCosts := []int{}
	nonsquareCosts := []int{}
	squares := 0
	for _, v := range a {
		if isSquare(v) {
			squares++
			if v == 0 {
				squareCosts = append(squareCosts, 2)
			} else {
				squareCosts = append(squareCosts, 1)
			}
		} else {
			nonsquareCosts = append(nonsquareCosts, costToSquare(v))
		}
	}
	target := n / 2
	if squares == target {
		return 0
	}
	if squares > target {
		sort.Ints(squareCosts)
		need := squares - target
		moves := 0
		for i := 0; i < need; i++ {
			moves += squareCosts[i]
		}
		return moves
	}
	sort.Ints(nonsquareCosts)
	need := target - squares
	moves := 0
	for i := 0; i < need; i++ {
		moves += nonsquareCosts[i]
	}
	return moves
}

func genCaseE(rng *rand.Rand) (string, int) {
	n := (rng.Intn(5) + 1) * 2
	a := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(50)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteString("\n")
	return sb.String(), solveE(a)
}

func runCaseE(bin, in string, exp int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	gotStr := strings.TrimSpace(buf.String())
	got, err := strconv.Atoi(gotStr)
	if err != nil {
		return fmt.Errorf("invalid output: %v", err)
	}
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
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
	for i := 0; i < 100; i++ {
		in, exp := genCaseE(rng)
		if err := runCaseE(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
