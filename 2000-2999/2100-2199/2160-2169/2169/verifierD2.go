package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testInput struct {
	text string
}

// Embedded reference solver for 2169 D2.
func solveD2(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var buf bytes.Buffer
	out := bufio.NewWriter(&buf)

	var t int
	fmt.Fscan(in, &t)

	for i := 0; i < t; i++ {
		var x, y, k int64
		fmt.Fscan(in, &x, &y, &k)

		if y == 1 {
			fmt.Fprintln(out, -1)
			continue
		}

		D := y - 1
		V := k - 1
		limit := int64(1000000000000)

		possible := true

		for x > 0 {
			if V >= limit {
				possible = false
				break
			}

			q := V / D
			if q == 0 {
				break
			}

			var steps int64 = 1
			if q < D {
				rem := (q+1)*D - 1 - V
				steps = rem / q
				if steps == 0 {
					steps = 1
				}
			}

			if steps > x {
				steps = x
			}

			V += steps * q
			x -= steps
		}

		if V >= limit {
			possible = false
		}

		if possible {
			fmt.Fprintln(out, V+1)
		} else {
			fmt.Fprintln(out, -1)
		}
	}

	out.Flush()
	return buf.String()
}

func commandForPath(path string) *exec.Cmd {
	switch {
	case strings.HasSuffix(path, ".go"):
		return exec.Command("go", "run", path)
	case strings.HasSuffix(path, ".py"):
		return exec.Command("python3", path)
	case strings.HasSuffix(path, ".js"):
		return exec.Command("node", path)
	default:
		return exec.Command(path)
	}
}

func runBinary(path, input string) (string, error) {
	cmd := commandForPath(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

func normalizeOutput(s string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	return strings.Join(lines, "\n")
}

func fixedTests() []testInput {
	return []testInput{
		{"6\n2 3 5\n2 5 1\n20 2 1000000000000\n175 10 28\n1000000000 998244353 1\n99 1 1\n"},
		{"3\n1 1 1\n1 2 1\n1000000000000 1000000000000 1000000000000\n"},
		{"1\n999999999999 999999999999 1\n"},
	}
}

func randomValue(rng *rand.Rand, limit int64) int64 {
	return rng.Int63n(limit) + 1
}

func randomTests() []testInput {
	tests := fixedTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 80 {
		t := rng.Intn(10) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", t))
		for i := 0; i < t; i++ {
			x := randomValue(rng, 1_000_000_000_000)
			y := randomValue(rng, 1_000_000_000_000)
			k := randomValue(rng, 1_000_000_000_000)
			sb.WriteString(fmt.Sprintf("%d %d %d\n", x, y, k))
		}
		tests = append(tests, testInput{text: sb.String()})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests := randomTests()
	for idx, input := range tests {
		expect := solveD2(input.text)

		got, err := runBinary(bin, input.text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, input.text)
			os.Exit(1)
		}
		if normalizeOutput(expect) != normalizeOutput(got) {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, input.text, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
