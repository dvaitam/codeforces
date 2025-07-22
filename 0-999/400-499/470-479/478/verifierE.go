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

var wavies []int64

func dfs(num int64, p1, p2, length, maxLen int) {
	if length >= maxLen {
		return
	}
	for d := 0; d <= 9; d++ {
		if (p2 > p1 && p2 > d) || (p2 < p1 && p2 < d) {
			newNum := num*10 + int64(d)
			wavies = append(wavies, newNum)
			dfs(newNum, p2, d, length+1, maxLen)
		}
	}
}

func generateWavies() {
	const maxLen = 14
	for d := 1; d <= 9; d++ {
		wavies = append(wavies, int64(d))
	}
	type seed struct {
		num    int64
		p1, p2 int
	}
	var seeds []seed
	for d1 := 1; d1 <= 9; d1++ {
		for d2 := 0; d2 <= 9; d2++ {
			num := int64(d1*10 + d2)
			wavies = append(wavies, num)
			seeds = append(seeds, seed{num, d1, d2})
		}
	}
	for _, s := range seeds {
		dfs(s.num, s.p1, s.p2, 2, maxLen)
	}
	sort.Slice(wavies, func(i, j int) bool { return wavies[i] < wavies[j] })
}

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

func expected(n, k int64) string {
	var cnt int64
	for _, v := range wavies {
		if v%n == 0 {
			cnt++
			if cnt == k {
				return fmt.Sprintf("%d", v)
			}
		}
	}
	return "-1"
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Int63n(1_000_000_000) + 1
	k := rng.Int63n(50) + 1
	input := fmt.Sprintf("%d %d\n", n, k)
	return input, expected(n, k)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	generateWavies()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
