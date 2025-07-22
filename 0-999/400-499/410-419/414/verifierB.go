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

const mod = 1000000007

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

func solveB(n, k int) int {
	prev := make([]int, n+1)
	for i := 1; i <= n; i++ {
		prev[i] = 1
	}
	cur := make([]int, n+1)
	for length := 2; length <= k; length++ {
		for i := 1; i <= n; i++ {
			cur[i] = 0
		}
		for x := 1; x <= n; x++ {
			v := prev[x]
			if v == 0 {
				continue
			}
			for m := x; m <= n; m += x {
				cur[m] += v
				if cur[m] >= mod {
					cur[m] -= mod
				}
			}
		}
		prev, cur = cur, prev
	}
	ans := 0
	for i := 1; i <= n; i++ {
		ans += prev[i]
		if ans >= mod {
			ans -= mod
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(30) + 1
	k := rng.Intn(10) + 1
	input := fmt.Sprintf("%d %d\n", n, k)
	expected := solveB(n, k)
	return input, expected
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := generateCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: bad output\ninput:\n%soutput:\n%s", i+1, input, out)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
