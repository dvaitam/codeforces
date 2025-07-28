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
	size   []int
}

func NewDSU(n int) *DSU {
	p := make([]int, n)
	sz := make([]int, n)
	for i := range p {
		p[i] = i
		sz[i] = 1
	}
	return &DSU{parent: p, size: sz}
}

func (d *DSU) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) union(a, b int) {
	ra := d.find(a)
	rb := d.find(b)
	if ra == rb {
		return
	}
	if d.size[ra] < d.size[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	d.size[ra] += d.size[rb]
}

func solveE2(n, k int, s, t string) string {
	dsu := NewDSU(n)
	for i := 0; i < n; i++ {
		if i+k < n {
			dsu.union(i, i+k)
		}
		if i+k+1 < n {
			dsu.union(i, i+k+1)
		}
	}
	counts := make(map[int][]int)
	for i := 0; i < n; i++ {
		r := dsu.find(i)
		arr, ok := counts[r]
		if !ok {
			arr = make([]int, 26)
			counts[r] = arr
		}
		arr[s[i]-'a']++
		arr[t[i]-'a']--
	}
	for _, arr := range counts {
		for j := 0; j < 26; j++ {
			if arr[j] != 0 {
				return "NO"
			}
		}
	}
	return "YES"
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	k := rng.Intn(n) + 1
	letters := "abcdefghijklmnopqrstuvwxyz"
	var sb1, sb2 strings.Builder
	for i := 0; i < n; i++ {
		sb1.WriteByte(letters[rng.Intn(len(letters))])
		sb2.WriteByte(letters[rng.Intn(len(letters))])
	}
	s := sb1.String()
	t := sb2.String()
	input := fmt.Sprintf("1\n%d %d\n%s\n%s\n", n, k, s, t)
	expected := solveE2(n, k, s, t)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierE2.go /path/to/binary")
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
