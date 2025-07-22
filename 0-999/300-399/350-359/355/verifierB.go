package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func expectedCost(c1, c2, c3, c4 int, buses, trolleys []int) int {
	busCost := 0
	for _, a := range buses {
		cost := a * c1
		if cost > c2 {
			cost = c2
		}
		busCost += cost
	}
	if busCost > c3 {
		busCost = c3
	}
	trolleyCost := 0
	for _, b := range trolleys {
		cost := b * c1
		if cost > c2 {
			cost = c2
		}
		trolleyCost += cost
	}
	if trolleyCost > c3 {
		trolleyCost = c3
	}
	ans := busCost + trolleyCost
	if ans > c4 {
		ans = c4
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, int) {
	c1 := rng.Intn(1000) + 1
	c2 := rng.Intn(1000) + 1
	c3 := rng.Intn(1000) + 1
	c4 := rng.Intn(1000) + 1
	n := rng.Intn(10) + 1
	m := rng.Intn(10) + 1
	buses := make([]int, n)
	trolleys := make([]int, m)
	for i := range buses {
		buses[i] = rng.Intn(1001)
	}
	for i := range trolleys {
		trolleys[i] = rng.Intn(1001)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d %d\n", c1, c2, c3, c4))
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i, v := range buses {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, v := range trolleys {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	exp := expectedCost(c1, c2, c3, c4, buses, trolleys)
	return sb.String(), exp
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
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		outVal, err2 := strconv.Atoi(strings.TrimSpace(out))
		if err2 != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: output not integer\ninput:\n%soutput:\n%s", i+1, in, out)
			os.Exit(1)
		}
		if outVal != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", i+1, exp, outVal, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
