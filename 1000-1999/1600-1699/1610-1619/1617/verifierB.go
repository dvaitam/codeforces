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

func run(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		return out.String() + errb.String(), err
	}
	return out.String(), nil
}

func solve(n int) string {
	if n%2 == 0 {
		return fmt.Sprintf("%d %d 1", n/2-1, n/2)
	}
	if (n-1)%4 == 0 {
		k := (n - 1) / 2
		return fmt.Sprintf("%d %d 1", k-1, k+1)
	}
	k := (n - 1) / 2
	return fmt.Sprintf("%d %d 1", k-2, k+2)
}

func genTest(rng *rand.Rand) int {
	return rng.Intn(1_000_000_000-10) + 10
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const tests = 100
	var input bytes.Buffer
	input.WriteString(fmt.Sprintf("%d\n", tests))
	expected := make([]string, tests)
	nums := make([]int, tests)
	for i := 0; i < tests; i++ {
		n := genTest(rng)
		nums[i] = n
		expected[i] = solve(n)
		input.WriteString(fmt.Sprintf("%d\n", n))
	}
	out, err := run(bin, input.String())
	if err != nil {
		fmt.Printf("runtime error: %v\n%s", err, out)
		os.Exit(1)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != tests {
		fmt.Printf("expected %d lines of output got %d\n", tests, len(lines))
		os.Exit(1)
	}
	for i := 0; i < tests; i++ {
		if strings.TrimSpace(lines[i]) != expected[i] {
			fmt.Printf("test %d failed (n=%d) expected:%s got:%s\n", i+1, nums[i], expected[i], strings.TrimSpace(lines[i]))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
