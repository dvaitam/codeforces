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

type CaseC struct {
	input string
	alice int64
	bob   int64
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func expectedC(s []int64) (int64, int64) {
	n := len(s)
	total := int64(0)
	for _, v := range s {
		total += v
	}
	dp := make([][2]int64, n+1)
	for i := n - 1; i >= 0; i-- {
		val := s[i]
		dp[i][0] = max64(val+dp[i+1][1], -val+dp[i+1][0])
		dp[i][1] = min64(-val+dp[i+1][0], val+dp[i+1][1])
	}
	diff := dp[0][1]
	alice := (total + diff) / 2
	bob := total - alice
	return alice, bob
}

func generateCaseC(rng *rand.Rand) CaseC {
	n := rng.Intn(50) + 1
	slices := make([]int64, n)
	for i := range slices {
		slices[i] = rng.Int63n(100_000) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range slices {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	a, b := expectedC(slices)
	return CaseC{sb.String(), a, b}
}

func runCase(exe string, input string, alice, bob int64) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(exe, ".go") {
		cmd = exec.Command("go", "run", exe)
	} else {
		cmd = exec.Command(exe)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	var ga, gb int64
	if _, err := fmt.Sscan(outStr, &ga, &gb); err != nil {
		return fmt.Errorf("cannot parse output: %v\n%s", err, outStr)
	}
	if ga != alice || gb != bob {
		return fmt.Errorf("expected %d %d got %d %d", alice, bob, ga, gb)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	slices := []int64{1, 2, 3}
	a, b := expectedC(slices)
	cases := []CaseC{
		{"3\n1 2 3\n", a, b},
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCaseC(rng))
	}
	for i, tc := range cases {
		if err := runCase(exe, tc.input, tc.alice, tc.bob); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
