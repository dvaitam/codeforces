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

func NewDSU(n int) *DSU {
	p := make([]int, n)
	for i := range p {
		p[i] = i
	}
	return &DSU{parent: p}
}

func (d *DSU) Find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.Find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) Union(a, b int) {
	ra := d.Find(a)
	rb := d.Find(b)
	if ra != rb {
		d.parent[ra] = rb
	}
}

func solveE1(n int, s, t string) string {
	k := 3
	dsu := NewDSU(n)
	for i := 0; i+k < n; i++ {
		dsu.Union(i, i+k)
	}
	for i := 0; i+k+1 < n; i++ {
		dsu.Union(i, i+k+1)
	}
	countsS := make(map[int][26]int)
	countsT := make(map[int][26]int)
	for i := 0; i < n; i++ {
		r := dsu.Find(i)
		cs := countsS[r]
		cs[int(s[i]-'a')]++
		countsS[r] = cs
		ct := countsT[r]
		ct[int(t[i]-'a')]++
		countsT[r] = ct
	}
	for r, cs := range countsS {
		ct := countsT[r]
		for j := 0; j < 26; j++ {
			if cs[j] != ct[j] {
				return "NO"
			}
		}
	}
	return "YES"
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	letters := "abcdefghijklmnopqrstuvwxyz"
	var sb1, sb2 strings.Builder
	for i := 0; i < n; i++ {
		sb1.WriteByte(letters[rng.Intn(len(letters))])
		sb2.WriteByte(letters[rng.Intn(len(letters))])
	}
	s := sb1.String()
	t := sb2.String()
	input := fmt.Sprintf("1\n%d 3\n%s\n%s\n", n, s, t)
	expected := solveE1(n, s, t)
	return input, expected
}

func runCase(bin, input, expected string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
