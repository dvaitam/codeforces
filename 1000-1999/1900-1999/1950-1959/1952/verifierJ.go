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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expectedLines(n int) []string {
	switch n {
	case 1:
		return []string{"yoink a", "yoink b", "*slaps a on top of b*", "yeet b", "go touch some grass"}
	case 2:
		return []string{"yoink a", "bruh b is lowkey just 0", "rip this b fell off by a", "vibe check a ratios b", "simp for 7", "bruh a is lowkey just b", "yeet a", "go touch some grass"}
	case 3:
		return []string{"yoink n", "yoink a", "bruh m is lowkey just a[0]", "bruh i is lowkey just 1", "vibe check n ratios i", "simp for 9", "yeet m", "go touch some grass", "vibe check a[i] ratios m", "bruh m is lowkey just a[i]", "*slaps 1 on top of i*", "simp for 5"}
	default:
		return []string{"yoink n", "yoink a", "yoink k", "bruh i is lowkey just 0", "bruh c is lowkey just 1", "bruh j is lowkey just 0", "vibe check n ratios j", "simp for 17", "vibe check k ratios c", "simp for 15", "vibe check c ratios k", "simp for 15", "yeet a[i]", "go touch some grass", "*slaps 1 on top of i*", "simp for 5", "vibe check a[j] ratios a[i]", "*slaps 1 on top of c*", "*slaps 1 on top of j*", "simp for 7"}
	}
}

func generateCase(rng *rand.Rand) (string, []string) {
	n := rng.Intn(6) + 1
	return fmt.Sprintf("%d\n", n), expectedLines(n)
}

func check(out string, expect []string) error {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != len(expect) {
		return fmt.Errorf("expected %d lines got %d", len(expect), len(lines))
	}
	for i, e := range expect {
		if strings.TrimSpace(lines[i]) != e {
			return fmt.Errorf("line %d: expected %q got %q", i+1, e, lines[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierJ.go /path/to/binary")
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
		if err := check(out, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
