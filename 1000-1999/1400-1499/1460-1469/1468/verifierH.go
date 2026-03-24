package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

// Embedded reference solver for 1468H
func refSolve(n, k, m int, b []int) string {
	if (n-m)%(k-1) != 0 {
		return "NO"
	}
	req := (k - 1) / 2
	for i := 0; i < m; i++ {
		L := b[i] - (i + 1)
		R := (n - b[i]) - (m - (i + 1))
		if L >= req && R >= req {
			return "YES"
		}
	}
	return "NO"
}

func runBinary(bin string, input []byte) ([]byte, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return out.Bytes(), nil
}

type Case struct {
	input    []byte
	expected string
}

func genCases() []Case {
	rng := rand.New(rand.NewSource(1468))
	cases := make([]Case, 100)
	for i := range cases {
		n := rng.Intn(20) + 3 // n >= 3 per problem constraints
		k := rng.Intn(n-2) + 3 // k >= 3, k <= n, and k is odd
		if k%2 == 0 {
			k--
			if k < 3 {
				k = 3
			}
		}
		if k > n {
			k = n
			if k%2 == 0 {
				k--
			}
		}
		// m >= 1, m < n
		maxM := n - 1
		if maxM < 1 {
			maxM = 1
		}
		m := rng.Intn(maxM) + 1
		if m >= n {
			m = n - 1
		}
		// Generate m distinct sorted values in [1, n]
		perm := rng.Perm(n)
		chosen := make([]int, m)
		for j := 0; j < m; j++ {
			chosen[j] = perm[j] + 1
		}
		sort.Ints(chosen)

		exp := refSolve(n, k, m, chosen)

		var buf strings.Builder
		fmt.Fprintf(&buf, "1\n%d %d %d\n", n, k, m)
		for j, v := range chosen {
			if j > 0 {
				buf.WriteByte(' ')
			}
			fmt.Fprintf(&buf, "%d", v)
		}
		buf.WriteByte('\n')
		cases[i] = Case{input: []byte(buf.String()), expected: exp}
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases := genCases()
	for i, c := range cases {
		out, err := runBinary(bin, c.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if !strings.EqualFold(got, c.expected) {
			fmt.Printf("wrong answer on case %d\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, string(c.input), c.expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
