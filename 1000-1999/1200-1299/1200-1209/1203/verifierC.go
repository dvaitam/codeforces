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

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func solve(nums []int64) string {
	g := nums[0]
	for _, x := range nums[1:] {
		g = gcd(g, x)
	}
	ans := int64(1)
	for p := int64(2); p*p <= g; p++ {
		if g%p == 0 {
			cnt := int64(0)
			for g%p == 0 {
				g /= p
				cnt++
			}
			ans *= cnt + 1
		}
	}
	if g > 1 {
		ans *= 2
	}
	return fmt.Sprintf("%d", ans)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	nums := make([]int64, n)
	for i := range nums {
		nums[i] = int64(rng.Intn(100000) + 1)
	}
	parts := make([]string, n)
	for i, v := range nums {
		parts[i] = fmt.Sprintf("%d", v)
	}
	input := fmt.Sprintf("%d\n%s\n", n, strings.Join(parts, " "))
	expect := solve(nums)
	return input, expect
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
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
