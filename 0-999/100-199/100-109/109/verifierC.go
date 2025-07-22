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

type UF struct {
	parent []int
	size   []int
}

func NewUF(n int) *UF {
	p := make([]int, n)
	s := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i
		s[i] = 1
	}
	return &UF{p, s}
}

func (u *UF) Find(x int) int {
	if u.parent[x] != x {
		u.parent[x] = u.Find(u.parent[x])
	}
	return u.parent[x]
}

func (u *UF) Union(a, b int) {
	a = u.Find(a)
	b = u.Find(b)
	if a == b {
		return
	}
	if u.size[a] > u.size[b] {
		a, b = b, a
	}
	u.parent[a] = b
	u.size[b] += u.size[a]
	u.size[a] = 0
}

func lucky(x int) bool {
	if x == 0 {
		return true
	}
	for x > 0 {
		d := x % 10
		if d != 4 && d != 7 {
			return false
		}
		x /= 10
	}
	return true
}

func expectedOutput(input string) string {
	r := bufio.NewReader(strings.NewReader(input))
	var n int
	fmt.Fscan(r, &n)
	uf := NewUF(n)
	for i := 0; i < n-1; i++ {
		var a, b, c int
		fmt.Fscan(r, &a, &b, &c)
		if !lucky(c) {
			uf.Union(a-1, b-1)
		}
	}
	total := int64(n)
	var ans int64
	for i := 0; i < n; i++ {
		s := uf.size[i]
		if s == 0 {
			continue
		}
		rem := total - int64(s)
		ans += int64(s) * rem * (rem - 1)
	}
	return fmt.Sprintf("%d\n", ans)
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 2; i <= n; i++ {
		a := i
		b := rng.Intn(i-1) + 1
		w := rng.Intn(1000) + 1
		if rng.Intn(5) == 0 {
			if rng.Intn(2) == 0 {
				w = 4
			} else {
				w = 7
			}
		}
		sb.WriteString(fmt.Sprintf("%d %d %d\n", a, b, w))
	}
	return sb.String()
}

func runCase(bin string, input string, expected string) error {
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
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []string{
		"1\n",
		"2\n1 2 4\n",
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}
	for i, tc := range cases {
		exp := expectedOutput(tc)
		if err := runCase(bin, tc, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
