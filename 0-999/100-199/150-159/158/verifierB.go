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

func expectedB(groups []int) int {
	counts := make([]int, 5)
	for _, g := range groups {
		if g >= 1 && g <= 4 {
			counts[g]++
		}
	}
	taxis := counts[4]
	taxis += counts[3]
	if counts[1] > counts[3] {
		counts[1] -= counts[3]
	} else {
		counts[1] = 0
	}
	taxis += counts[2] / 2
	if counts[2]%2 != 0 {
		taxis++
		if counts[1] > 2 {
			counts[1] -= 2
		} else {
			counts[1] = 0
		}
	}
	if counts[1] > 0 {
		taxis += (counts[1] + 3) / 4
	}
	return taxis
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
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		n := rng.Intn(1000) + 1
		groups := make([]int, n)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for i := 0; i < n; i++ {
			g := rng.Intn(4) + 1
			groups[i] = g
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", g)
		}
		sb.WriteByte('\n')
		input := sb.String()
		expected := expectedB(groups)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s", t+1, err, out)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if out != fmt.Sprint(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:\n%s", t+1, expected, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
