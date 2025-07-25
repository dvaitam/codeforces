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
	n    int
	sher string
	mor  string
}

func solve(n int, sher, mor string) (int, int) {
	freq1 := [10]int{}
	freq2 := [10]int{}
	for i := 0; i < n; i++ {
		d := mor[i] - '0'
		freq1[d]++
		freq2[d]++
	}
	minFlicks := 0
	for i := 0; i < n; i++ {
		d := int(sher[i] - '0')
		j := d
		for j < 10 && freq1[j] == 0 {
			j++
		}
		if j < 10 {
			freq1[j]--
		} else {
			minFlicks++
			for k := 0; k < 10; k++ {
				if freq1[k] > 0 {
					freq1[k]--
					break
				}
			}
		}
	}
	maxSher := 0
	for i := 0; i < n; i++ {
		d := int(sher[i] - '0')
		j := d + 1
		for j < 10 && freq2[j] == 0 {
			j++
		}
		if j < 10 {
			freq2[j]--
			maxSher++
		} else {
			for k := 0; k < 10; k++ {
				if freq2[k] > 0 {
					freq2[k]--
					break
				}
			}
		}
	}
	return minFlicks, maxSher
}

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

func buildInput(tc testCase) string {
	return fmt.Sprintf("%d\n%s\n%s\n", tc.n, tc.sher, tc.mor)
}

func genCases() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{
		{1, "0", "0"},
		{2, "01", "10"},
		{3, "123", "321"},
	}
	for len(cases) < 120 {
		n := rng.Intn(20) + 1
		b1 := make([]byte, n)
		b2 := make([]byte, n)
		for i := 0; i < n; i++ {
			b1[i] = byte('0' + rng.Intn(10))
			b2[i] = byte('0' + rng.Intn(10))
		}
		cases = append(cases, testCase{n: n, sher: string(b1), mor: string(b2)})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases()
	for i, tc := range cases {
		input := buildInput(tc)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		minF, maxS := solve(tc.n, tc.sher, tc.mor)
		exp := fmt.Sprintf("%d\n%d", minF, maxS)
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, exp, out, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
