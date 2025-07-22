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

func solveCase(times []int, T int) int {
	sort.Ints(times)
	best := 0
	r := 0
	for l := 0; l < len(times); l++ {
		for r < len(times) && times[r]-times[l] <= T {
			r++
		}
		if r-l > best {
			best = r - l
		}
	}
	return best
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(100) + 1
	times := make([]int, n)
	for i := 0; i < n; i++ {
		times[i] = rng.Intn(1000) + 1
	}
	T := rng.Intn(1000) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range times {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", T))
	expect := fmt.Sprintf("%d", solveCase(append([]int(nil), times...), T))
	return sb.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := []struct{ in, out string }{}
	// simple case
	cases = append(cases, "1\n5\n0\n", "1")
	cases = append(cases, "3\n1 3 8\n4\n", "2")

	for i := 0; i < 100; i++ {
		in, out := generateCase(rng)
		cases = append(cases, struct{ in, out string }{in, out})
	}

	for i, tc := range cases {
		got, err := runCandidate(bin, tc.in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.in)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(tc.out) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc.out, got, tc.in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
