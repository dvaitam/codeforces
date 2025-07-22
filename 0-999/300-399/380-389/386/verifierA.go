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

func solveCase(bids []int) (int, int) {
	maxVal, secondVal := -1, -1
	maxIdx := -1
	for i, v := range bids {
		if v > maxVal {
			secondVal = maxVal
			maxVal = v
			maxIdx = i + 1
		} else if v > secondVal {
			secondVal = v
		}
	}
	return maxIdx, secondVal
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(999) + 2 // 2..1000
	seen := make(map[int]bool)
	bids := make([]int, 0, n)
	for len(bids) < n {
		x := rng.Intn(10000) + 1
		if !seen[x] {
			seen[x] = true
			bids = append(bids, x)
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range bids {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	idx, price := solveCase(bids)
	expect := fmt.Sprintf("%d %d", idx, price)
	return sb.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// some deterministic edge cases
	cases := []struct{ in, out string }{}
	// smallest
	cases = append(cases, func() (string, string) {
		bids := []int{1, 2}
		in := "2\n1 2\n"
		idx, price := solveCase(bids)
		return in, fmt.Sprintf("%d %d", idx, price)
	}())
	// large n
	cases = append(cases, func() (string, string) {
		n := 1000
		bids := make([]int, n)
		for i := 0; i < n; i++ {
			bids[i] = i + 1
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := range bids {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", bids[i]))
		}
		sb.WriteByte('\n')
		idx, price := solveCase(bids)
		return sb.String(), fmt.Sprintf("%d %d", idx, price)
	}())

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
