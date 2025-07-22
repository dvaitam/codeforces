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

// DSU structure for solveA
type DSU struct{ p []int }

func newDSU(n int) *DSU {
	p := make([]int, n+1)
	for i := 0; i <= n; i++ {
		p[i] = i
	}
	return &DSU{p: p}
}

func (d *DSU) find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.find(d.p[x])
	}
	return d.p[x]
}

func (d *DSU) union(x, y int) {
	rx := d.find(x)
	ry := d.find(y)
	if rx != ry {
		d.p[ry] = rx
	}
}

func solveA(s string) string {
	n := len(s)
	isPrime := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		isPrime[i] = true
	}
	for i := 2; i*i <= n; i++ {
		if isPrime[i] {
			for j := i * i; j <= n; j += i {
				isPrime[j] = false
			}
		}
	}
	dsu := newDSU(n)
	for p := 2; p <= n; p++ {
		if !isPrime[p] {
			continue
		}
		for k := 2; p*k <= n; k++ {
			dsu.union(p, p*k)
		}
	}
	groups := make(map[int][]int)
	for i := 1; i <= n; i++ {
		r := dsu.find(i)
		groups[r] = append(groups[r], i)
	}
	grpList := make([][]int, 0, len(groups))
	for _, g := range groups {
		grpList = append(grpList, g)
	}
	for i := 1; i < len(grpList); i++ {
		j := i
		for j > 0 && len(grpList[j-1]) < len(grpList[j]) {
			grpList[j-1], grpList[j] = grpList[j], grpList[j-1]
			j--
		}
	}
	freq := make([]int, 26)
	for i := 0; i < n; i++ {
		freq[int(s[i]-'a')]++
	}
	res := make([]byte, n)
	for _, g := range grpList {
		size := len(g)
		idx := -1
		for j := 0; j < 26; j++ {
			if freq[j] >= size {
				idx = j
				break
			}
		}
		if idx == -1 {
			return "NO"
		}
		for _, pos := range g {
			res[pos-1] = byte('a' + idx)
		}
		freq[idx] -= size
	}
	return "YES\n" + string(res)
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte('a' + rng.Intn(26))
	}
	s := string(b)
	return s + "\n", s
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, s := genCase(rng)
		expect := strings.TrimSpace(solveA(s))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
