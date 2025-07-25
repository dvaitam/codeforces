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

func buildRef() (string, error) {
	ref := "./refA.bin"
	cmd := exec.Command("go", "build", "-o", ref, "555A.go")
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type Case struct {
	n      int
	k      int
	chains [][]int
	input  string
}

func genCases() []Case {
	rng := rand.New(rand.NewSource(42))
	cases := make([]Case, 100)
	for idx := range cases {
		n := rng.Intn(7) + 3 // 3..9
		k := rng.Intn(n) + 1
		perm := rng.Perm(n)
		chains := make([][]int, k)
		pos := 0
		for i := 0; i < k; i++ {
			remain := n - pos
			groupsLeft := k - i
			m := 1
			if remain-groupsLeft > 0 {
				m += rng.Intn(remain - groupsLeft + 1)
			}
			chain := make([]int, m)
			for j := 0; j < m; j++ {
				chain[j] = perm[pos] + 1
				pos++
			}
			sort.Ints(chain)
			chains[i] = chain
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, k)
		for _, ch := range chains {
			fmt.Fprintf(&sb, "%d", len(ch))
			for _, v := range ch {
				fmt.Fprintf(&sb, " %d", v)
			}
			sb.WriteByte('\n')
		}
		cases[idx] = Case{n, k, chains, sb.String()}
	}
	return cases
}

func runCase(bin, ref string, c Case) error {
	expected, err := runBinary(ref, c.input)
	if err != nil {
		return fmt.Errorf("reference failed: %v", err)
	}
	got, err := runBinary(bin, c.input)
	if err != nil {
		return err
	}
	if expected != strings.TrimSpace(got) {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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
	for i, c := range cases {
		if err := runCase(bin, ref, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, c.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
