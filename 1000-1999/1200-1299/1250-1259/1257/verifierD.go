package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func expected(monsters []int, heroes [][2]int) int {
	n := len(monsters)
	maxPower := make([]int, n+2)
	for _, h := range heroes {
		p, s := h[0], h[1]
		if s <= n && maxPower[s] < p {
			maxPower[s] = p
		}
	}
	for i := n - 1; i >= 1; i-- {
		if maxPower[i] < maxPower[i+1] {
			maxPower[i] = maxPower[i+1]
		}
	}
	for _, m := range monsters {
		if m > maxPower[1] {
			return -1
		}
	}
	days := 0
	i := 0
	for i < n {
		days++
		maxA := 0
		j := i
		for j < n {
			if monsters[j] > maxA {
				maxA = monsters[j]
			}
			length := j - i + 1
			if maxPower[length] < maxA {
				break
			}
			j++
		}
		i = j
	}
	return days
}

func generateCase(rng *rand.Rand) (string, []string) {
	t := rng.Intn(3) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	exp := make([]string, t)
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		n := rng.Intn(5) + 1
		monsters := make([]int, n)
		for i := 0; i < n; i++ {
			monsters[i] = rng.Intn(9) + 1
		}
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", monsters[i]))
		}
		sb.WriteByte('\n')
		m := rng.Intn(5) + 1
		heroes := make([][2]int, m)
		for i := 0; i < m; i++ {
			p := rng.Intn(9) + 1
			s := rng.Intn(n) + 1
			heroes[i] = [2]int{p, s}
			sb.WriteString(fmt.Sprintf("%d %d\n", p, s))
		}
		exp[caseIdx] = fmt.Sprintf("%d", expected(monsters, heroes))
	}
	return sb.String(), exp
}

func runCase(bin, input string, exp []string) error {
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
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != len(exp) {
		return fmt.Errorf("expected %d lines got %d", len(exp), len(lines))
	}
	for i, line := range lines {
		if strings.TrimSpace(line) != exp[i] {
			return fmt.Errorf("line %d expected %s got %s", i+1, exp[i], strings.TrimSpace(line))
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
