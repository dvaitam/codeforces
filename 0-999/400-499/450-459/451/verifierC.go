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

func runCandidate(bin, input string) (string, error) {
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

func solveCase(n, k, d1, d2 int64) string {
	if n%3 != 0 {
		return "no"
	}
	target := n / 3
	for _, s1 := range []int64{1, -1} {
		for _, s2 := range []int64{1, -1} {
			num := k - s1*d1 + s2*d2
			if num%3 != 0 {
				continue
			}
			w2 := num / 3
			w1 := w2 + s1*d1
			w3 := w2 - s2*d2
			if w1 < 0 || w2 < 0 || w3 < 0 {
				continue
			}
			if w1 > target || w2 > target || w3 > target {
				continue
			}
			return "yes"
		}
	}
	return "no"
}

func generateCase(rng *rand.Rand) (string, string) {
	n := int64(rng.Intn(300) + 3)
	k := int64(rng.Intn(int(n + 1)))
	d1 := int64(rng.Intn(int(k + 1)))
	d2 := int64(rng.Intn(int(k + 1)))
	input := fmt.Sprintf("1\n%d %d %d %d\n", n, k, d1, d2)
	expect := solveCase(n, k, d1, d2)
	return input, expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
