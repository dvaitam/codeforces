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
	hashtags []string
}

func solve(tc testCase) []string {
	n := len(tc.hashtags)
	hashtags := make([]string, n)
	copy(hashtags, tc.hashtags)
	for i := n - 2; i >= 0; i-- {
		a := hashtags[i]
		b := hashtags[i+1]
		la, lb := len(a), len(b)
		j := 1
		for j < la && j < lb && a[j] == b[j] {
			j++
		}
		if j == la {
			continue
		}
		if j == lb {
			hashtags[i] = a[:lb]
			continue
		}
		if a[j] > b[j] {
			hashtags[i] = a[:j]
		}
	}
	return hashtags
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
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.hashtags)))
	for _, s := range tc.hashtags {
		sb.WriteString(s)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func genCases() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 0, 110)
	for i := 0; i < 110; i++ {
		n := rng.Intn(8) + 1
		h := make([]string, n)
		for j := 0; j < n; j++ {
			l := rng.Intn(6) + 1
			b := make([]byte, l)
			for k := 0; k < l; k++ {
				b[k] = byte('a' + rng.Intn(26))
			}
			h[j] = "#" + string(b)
		}
		cases = append(cases, testCase{hashtags: h})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
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
		expArr := solve(tc)
		exp := strings.Join(expArr, "\n")
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:\n%s", i+1, exp, out, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
