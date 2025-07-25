package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func runCandidate(bin string, input string) (string, error) {
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

func expected(intervals [][2]int) string {
	removed := make([]bool, 1001)
	for _, iv := range intervals {
		a, b := iv[0], iv[1]
		if a < 1 {
			a = 1
		}
		if b > 1000 {
			b = 1000
		}
		for j := a; j <= b; j++ {
			removed[j] = true
		}
	}
	var res []int
	for i := 1; i <= 1000; i++ {
		if !removed[i] {
			res = append(res, i)
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d", len(res))
	for _, v := range res {
		fmt.Fprintf(&sb, " %d", v)
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < 100; i++ {
		n := rng.Intn(10)
		intervals := make([][2]int, n)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for j := 0; j < n; j++ {
			a := rng.Intn(1000) + 1
			b := a + rng.Intn(1000-a+1)
			intervals[j] = [2]int{a, b}
			fmt.Fprintf(&sb, "%d %d\n", a, b)
		}
		input := sb.String()
		want := expected(intervals)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(want) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, want, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
