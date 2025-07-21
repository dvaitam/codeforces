package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(grid [3]string) string {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if grid[i][j] != grid[2-i][2-j] {
				return "NO"
			}
		}
	}
	return "YES"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	total := 0
	for mask := 0; mask < 512; mask++ {
		var g [3]string
		for i := 0; i < 3; i++ {
			var b strings.Builder
			for j := 0; j < 3; j++ {
				idx := i*3 + j
				if mask&(1<<idx) != 0 {
					b.WriteByte('X')
				} else {
					b.WriteByte('.')
				}
			}
			g[i] = b.String()
		}
		input := g[0] + "\n" + g[1] + "\n" + g[2] + "\n"
		exp := expected(g)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", total+1, err, input)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", total+1, exp, out, input)
			os.Exit(1)
		}
		total++
	}
	fmt.Printf("All %d tests passed\n", total)
}
