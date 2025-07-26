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

func solve(sticks []int) string {
	sort.Ints(sticks)
	n := len(sticks) / 4
	area := sticks[0] * sticks[len(sticks)-1]
	for i := 0; i < n; i++ {
		if sticks[2*i] != sticks[2*i+1] || sticks[len(sticks)-2*i-2] != sticks[len(sticks)-2*i-1] || sticks[2*i]*sticks[len(sticks)-2*i-1] != area {
			return "NO"
		}
	}
	return "YES"
}

func genValid(rng *rand.Rand) ([]int, string) {
	n := rng.Intn(10) + 1
	area := rng.Intn(20) + 1
	sticks := make([]int, 0, 4*n)
	for i := 0; i < n; i++ {
		a := rng.Intn(area) + 1
		b := area / a
		// ensure area divides
		for a*b != area {
			a = rng.Intn(area) + 1
			if area%a == 0 {
				b = area / a
			}
		}
		sticks = append(sticks, a, a, b, b)
	}
	rng.Shuffle(len(sticks), func(i, j int) { sticks[i], sticks[j] = sticks[j], sticks[i] })
	return sticks, "YES"
}

func genRandom(rng *rand.Rand) ([]int, string) {
	n := rng.Intn(10) + 1
	sticks := make([]int, 4*n)
	for i := range sticks {
		sticks[i] = rng.Intn(20) + 1
	}
	return sticks, solve(sticks)
}

func generateCase(rng *rand.Rand) (string, string) {
	if rng.Intn(2) == 0 {
		sticks, exp := genValid(rng)
		nums := make([]string, len(sticks))
		for i, v := range sticks {
			nums[i] = fmt.Sprintf("%d", v)
		}
		input := fmt.Sprintf("1\n%d\n%s\n", len(sticks)/4, strings.Join(nums, " "))
		return input, exp
	}
	sticks, exp := genRandom(rng)
	nums := make([]string, len(sticks))
	for i, v := range sticks {
		nums[i] = fmt.Sprintf("%d", v)
	}
	input := fmt.Sprintf("1\n%d\n%s\n", len(sticks)/4, strings.Join(nums, " "))
	return input, exp
}

func runCase(bin, input, expect string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(strings.Split(out.String(), "\n")[0])
	if got != expect {
		return fmt.Errorf("expected %s got %s", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expect := generateCase(rng)
		if err := runCase(bin, input, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
