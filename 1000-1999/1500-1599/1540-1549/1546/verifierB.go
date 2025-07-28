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
	n      int
	m      int
	orig   []string
	final  []string
	stolen string
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String() + errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func genString(rng *rand.Rand, m int) string {
	b := make([]byte, m)
	for i := 0; i < m; i++ {
		b[i] = byte('a' + rng.Intn(26))
	}
	return string(b)
}

func genCase(rng *rand.Rand) testCase {
	n := rng.Intn(5)*2 + 1 // 1,3,5,7,9
	m := rng.Intn(5) + 1   // length 1..5
	orig := make([]string, n)
	for i := 0; i < n; i++ {
		orig[i] = genString(rng, m)
	}
	stolenIndex := rng.Intn(n)

	// Copy originals to tmp for swapping
	tmp := make([]string, n)
	copy(tmp, orig)
	indices := []int{}
	for i := 0; i < n; i++ {
		if i != stolenIndex {
			indices = append(indices, i)
		}
	}
	rng.Shuffle(len(indices), func(i, j int) { indices[i], indices[j] = indices[j], indices[i] })
	for i := 0; i < len(indices); i += 2 {
		a := indices[i]
		b := indices[i+1]
		k := rng.Intn(m) + 1 // at least one position
		pos := rng.Perm(m)[:k]
		sa := []byte(tmp[a])
		sb := []byte(tmp[b])
		for _, p := range pos {
			sa[p], sb[p] = sb[p], sa[p]
		}
		tmp[a] = string(sa)
		tmp[b] = string(sb)
	}
	final := []string{}
	for i := 0; i < n; i++ {
		if i != stolenIndex {
			final = append(final, tmp[i])
		}
	}
	rng.Shuffle(len(final), func(i, j int) { final[i], final[j] = final[j], final[i] })
	return testCase{n: n, m: m, orig: orig, final: final, stolen: orig[stolenIndex]}
}

func compute(orig []string, final []string, m int) string {
	res := make([]byte, m)
	for _, s := range orig {
		for j := 0; j < m; j++ {
			res[j] ^= s[j]
		}
	}
	for _, s := range final {
		for j := 0; j < m; j++ {
			res[j] ^= s[j]
		}
	}
	return string(res)
}

func main() {
	var bin string
	if len(os.Args) == 2 {
		bin = os.Args[1]
	} else if len(os.Args) == 3 && os.Args[1] == "--" {
		bin = os.Args[2]
	} else {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []testCase
	// simple edge cases
	cases = append(cases, testCase{n: 1, m: 1, orig: []string{"a"}, final: []string{}, stolen: "a"})
	for len(cases) < 100 {
		cases = append(cases, genCase(rng))
	}
	for i, tc := range cases {
		var input strings.Builder
		fmt.Fprintf(&input, "1\n%d %d\n", tc.n, tc.m)
		for _, s := range tc.orig {
			fmt.Fprintln(&input, s)
		}
		for _, s := range tc.final {
			fmt.Fprintln(&input, s)
		}
		out, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(out)
		expect := compute(tc.orig, tc.final, tc.m)
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, input.String())
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
