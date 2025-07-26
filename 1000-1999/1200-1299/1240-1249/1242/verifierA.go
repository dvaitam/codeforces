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

func run(bin, input string) (string, error) {
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(n int64) int64 {
	if n == 1 {
		return 1
	}
	orig := n
	var p int64 = -1
	for i := int64(2); i*i <= n; i++ {
		if n%i == 0 {
			p = i
			break
		}
	}
	if p == -1 {
		return n
	}
	for orig%p == 0 {
		orig /= p
	}
	if orig == 1 {
		return p
	}
	return 1
}

func genCase(rng *rand.Rand) int64 {
	return rng.Int63n(1_000_000_000_000) + 1
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 12, 16, 17, 25, 27, 30, 36, 49, 60, 97}
	for len(cases) < 100 {
		cases = append(cases, genCase(rng))
	}

	for i, n := range cases {
		in := fmt.Sprintf("%d\n", n)
		exp := fmt.Sprintf("%d", solve(n))
		out, err := run(bin, in)
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
