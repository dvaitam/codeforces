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

type DSU struct {
	parent []int
}

func newDSU(n int) *DSU {
	d := &DSU{parent: make([]int, n)}
	for i := range d.parent {
		d.parent[i] = i
	}
	return d
}

func (d *DSU) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) union(a, b int) {
	pa := d.find(a)
	pb := d.find(b)
	if pa != pb {
		d.parent[pa] = pb
	}
}

func run(bin, input string) (string, error) {
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

func solveD(words []string) int {
	dsu := newDSU(26)
	used := make([]bool, 26)
	for _, s := range words {
		first := -1
		for _, ch := range s {
			idx := int(ch - 'a')
			used[idx] = true
			if first == -1 {
				first = idx
			} else {
				dsu.union(first, idx)
			}
		}
	}
	comps := make(map[int]bool)
	for i := 0; i < 26; i++ {
		if used[i] {
			root := dsu.find(i)
			comps[root] = true
		}
	}
	return len(comps)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	words := make([]string, n)
	for i := 0; i < n; i++ {
		l := rng.Intn(6) + 1
		b := make([]byte, l)
		for j := 0; j < l; j++ {
			b[j] = byte('a' + rng.Intn(5))
		}
		words[i] = string(b)
	}
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d\n", n))
	for _, w := range words {
		input.WriteString(w)
		input.WriteByte('\n')
	}
	expect := fmt.Sprintf("%d", solveD(words))
	return input.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
