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

type testCase struct {
	input    string
	expected string
}

func solveCase(k []int, a []int64) string {
	var ans int64
	for i := range k {
		cap := int64(1)
		d := int64(0)
		for cap < a[i] {
			cap <<= 2
			d++
		}
		if int64(k[i])+d > ans {
			ans = int64(k[i]) + d
		}
	}
	return fmt.Sprintf("%d", ans)
}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(10) + 1
	ks := make([]int, n)
	as := make([]int64, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		ks[i] = rng.Intn(20)
		as[i] = int64(rng.Intn(20) + 1)
		sb.WriteString(fmt.Sprintf("%d %d\n", ks[i], as[i]))
	}
	expect := solveCase(ks, as)
	return testCase{input: sb.String(), expected: expect}
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != tc.expected {
		return fmt.Errorf("expected %s got %s", tc.expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{{input: "1\n0 1\n", expected: "0"}}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
