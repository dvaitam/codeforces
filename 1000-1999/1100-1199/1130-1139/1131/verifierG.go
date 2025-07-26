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

type block struct {
	hs []int
	cs []int64
}
type query struct{ id, mul int }
type testCase struct {
	blocks  []block
	queries []query
}

func buildRef() (string, error) {
	ref := "./refG.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1131G.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

func runBinary(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCases() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 0, 100)
	// simple single block single query
	cases = append(cases, testCase{blocks: []block{{hs: []int{1}, cs: []int64{1}}}, queries: []query{{id: 1, mul: 1}}})
	for len(cases) < 100 {
		B := rng.Intn(3) + 1
		blocks := make([]block, B)
		for i := 0; i < B; i++ {
			L := rng.Intn(3) + 1
			hs := make([]int, L)
			cs := make([]int64, L)
			for j := 0; j < L; j++ {
				hs[j] = rng.Intn(3) + 1
			}
			for j := 0; j < L; j++ {
				cs[j] = int64(rng.Intn(5) + 1)
			}
			blocks[i] = block{hs: hs, cs: cs}
		}
		q := rng.Intn(3) + 1
		queries := make([]query, q)
		for i := 0; i < q; i++ {
			id := rng.Intn(B) + 1
			mul := rng.Intn(3) + 1
			queries[i] = query{id: id, mul: mul}
		}
		cases = append(cases, testCase{blocks: blocks, queries: queries})
	}
	return cases
}

func runCase(bin, ref string, tc testCase) error {
	// compute n as total lengths
	n := 0
	for _, q := range tc.queries {
		b := tc.blocks[q.id-1]
		n += len(b.hs)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", len(tc.blocks), n)
	for _, b := range tc.blocks {
		fmt.Fprintf(&sb, "%d\n", len(b.hs))
		for i, v := range b.hs {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
		for i, v := range b.cs {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
	}
	fmt.Fprintf(&sb, "%d\n", len(tc.queries))
	for _, q := range tc.queries {
		fmt.Fprintf(&sb, "%d %d\n", q.id, q.mul)
	}
	input := sb.String()
	expect, err := runBinary(ref, input)
	if err != nil {
		return fmt.Errorf("reference failed: %v", err)
	}
	got, err := runBinary(bin, input)
	if err != nil {
		return err
	}
	if got != expect {
		return fmt.Errorf("expected %s got %s", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	cases := genCases()
	for i, tc := range cases {
		if err := runCase(bin, ref, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
