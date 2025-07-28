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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func bruteC(k int, g []int) string {
	n := 1 << (k + 1)
	pref := make([]int, n+1)
	for i := 1; i <= n; i++ {
		pref[i] = pref[i-1] ^ g[i-1]
	}
	for l1 := 1; l1 <= n; l1++ {
		for r1 := l1; r1 <= n; r1++ {
			x1 := pref[r1] ^ pref[l1-1]
			for l2 := 1; l2 <= n; l2++ {
				for r2 := l2; r2 <= n; r2++ {
					if r1 < l2 || r2 < l1 {
						x2 := pref[r2] ^ pref[l2-1]
						if x1 == x2 {
							return fmt.Sprintf("%d %d %d %d", l1, r1, l2, r2)
						}
					}
				}
			}
		}
	}
	return "-1"
}

func generateCase(rng *rand.Rand) (string, string) {
	k := rng.Intn(3) // 0..2
	n := 1 << (k + 1)
	g := make([]int, n)
	maxVal := 1 << (2 * k)
	if maxVal == 0 {
		maxVal = 1
	}
	for i := 0; i < n; i++ {
		g[i] = rng.Intn(maxVal)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d\n", k)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", g[i])
	}
	sb.WriteByte('\n')
	input := sb.String()
	expected := bruteC(k, g)
	return input, expected
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
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
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
