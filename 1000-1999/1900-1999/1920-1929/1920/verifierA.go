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

type constraint struct {
	a int
	x int
}

func solveCase(cons []constraint) int {
	lo := 0
	hi := int(1e9)
	banned := map[int]struct{}{}
	for _, c := range cons {
		switch c.a {
		case 1:
			if c.x > lo {
				lo = c.x
			}
		case 2:
			if c.x < hi {
				hi = c.x
			}
		case 3:
			banned[c.x] = struct{}{}
		}
	}
	if lo > hi {
		return 0
	}
	count := hi - lo + 1
	for x := range banned {
		if x >= lo && x <= hi {
			count--
		}
	}
	if count < 0 {
		count = 0
	}
	return count
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(99) + 2 // 2..100
	cons := make([]constraint, n)
	has1, has2 := false, false
	for i := 0; i < n; i++ {
		a := rng.Intn(3) + 1
		if i == n-1 { // ensure both types appear
			if !has1 {
				a = 1
			} else if !has2 {
				a = 2
			}
		}
		x := rng.Intn(1_000_000_000) + 1
		cons[i] = constraint{a, x}
		if a == 1 {
			has1 = true
		}
		if a == 2 {
			has2 = true
		}
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	fmt.Fprintf(&sb, "%d\n", n)
	for _, c := range cons {
		fmt.Fprintf(&sb, "%d %d\n", c.a, c.x)
	}
	ans := solveCase(cons)
	return sb.String(), fmt.Sprintf("%d", ans)
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
