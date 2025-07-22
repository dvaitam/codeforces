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

func generateCase(rng *rand.Rand) (int, []int) {
	n := rng.Intn(19) + 2 // 2..20
	a := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(20) + 1
	}
	return n, a
}

func checkOutput(n int, skills []int, out string) error {
	fields := strings.Fields(out)
	if len(fields) < 2 {
		return fmt.Errorf("not enough output")
	}
	x, err := strconv.Atoi(fields[0])
	if err != nil || x < 0 || x > n {
		return fmt.Errorf("invalid x")
	}
	if len(fields) < 1+x+1 {
		return fmt.Errorf("not enough ids for team1")
	}
	team1 := make([]int, x)
	used := make([]bool, n)
	for i := 0; i < x; i++ {
		id, err := strconv.Atoi(fields[1+i])
		if err != nil || id < 1 || id > n || used[id-1] {
			return fmt.Errorf("invalid team1 id")
		}
		used[id-1] = true
		team1[i] = id - 1
	}
	idx := 1 + x
	if idx >= len(fields) {
		return fmt.Errorf("missing y")
	}
	y, err := strconv.Atoi(fields[idx])
	if err != nil || y < 0 || x+y != n {
		return fmt.Errorf("invalid y")
	}
	if len(fields) != idx+1+y {
		return fmt.Errorf("wrong number of ids for team2")
	}
	team2 := make([]int, y)
	for i := 0; i < y; i++ {
		id, err := strconv.Atoi(fields[idx+1+i])
		if err != nil || id < 1 || id > n || used[id-1] {
			return fmt.Errorf("invalid team2 id")
		}
		used[id-1] = true
		team2[i] = id - 1
	}
	if x+y != n {
		return fmt.Errorf("counts don't add up")
	}
	maxSkill := 0
	for _, v := range skills {
		if v > maxSkill {
			maxSkill = v
		}
	}
	sum1 := 0
	for _, id := range team1 {
		sum1 += skills[id]
	}
	sum2 := 0
	for _, id := range team2 {
		sum2 += skills[id]
	}
	if abs(sum1-sum2) > maxSkill {
		return fmt.Errorf("skill difference too large")
	}
	if abs(x-y) > 1 {
		return fmt.Errorf("team size difference too large")
	}
	return nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func runCase(bin string, n int, a []int) error {
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range a {
		input.WriteString(fmt.Sprintf("%d", v))
		if i+1 < n {
			input.WriteByte(' ')
		}
	}
	input.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if err := checkOutput(n, a, strings.TrimSpace(out.String())); err != nil {
		return err
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		n, a := generateCase(rng)
		if err := runCase(bin, n, a); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput n=%d a=%v\n", i+1, err, n, a)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
