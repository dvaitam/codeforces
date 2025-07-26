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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(3) + 3
	m := rng.Intn(3) + 3
	used1 := make(map[[2]int]bool)
	used2 := make(map[[2]int]bool)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < n; i++ {
		for {
			x := rng.Intn(11)
			y := rng.Intn(11)
			if !used1[[2]int{x, y}] {
				used1[[2]int{x, y}] = true
				fmt.Fprintf(&sb, "%d %d\n", x, y)
				break
			}
		}
	}
	for i := 0; i < m; i++ {
		for {
			x := rng.Intn(11)
			y := rng.Intn(11)
			if !used2[[2]int{x, y}] {
				used2[[2]int{x, y}] = true
				fmt.Fprintf(&sb, "%d %d\n", x, y)
				break
			}
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		low := strings.ToLower(strings.TrimSpace(out))
		if low != "yes" && low != "no" {
			fmt.Fprintf(os.Stderr, "case %d failed: expected YES or NO got %q\ninput:\n%s", i+1, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
