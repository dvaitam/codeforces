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

var colorMap = map[rune]int{'R': 0, 'O': 1, 'Y': 2, 'G': 3, 'B': 4, 'V': 5}

func compose(p, g []int) []int {
	r := make([]int, len(p))
	for i := range p {
		r[i] = g[p[i]]
	}
	return r
}

func key(p []int) string {
	b := make([]byte, len(p))
	for i, v := range p {
		b[i] = byte(v)
	}
	return string(b)
}

var perms = genGroup()

func genGroup() [][]int {
	Rx := []int{1, 5, 2, 0, 4, 3}
	Ry := []int{2, 1, 5, 3, 0, 4}
	Rz := []int{0, 2, 3, 4, 1, 5}
	perms := [][]int{}
	seen := map[string]bool{}
	id := []int{0, 1, 2, 3, 4, 5}
	queue := [][]int{id}
	seen[key(id)] = true
	gens := [][]int{Rx, Ry, Rz}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		perms = append(perms, cur)
		for _, g := range gens {
			nxt := compose(cur, g)
			k := key(nxt)
			if !seen[k] {
				seen[k] = true
				queue = append(queue, nxt)
			}
		}
	}
	return perms
}

func cycles(perm []int) []int {
	vis := make([]bool, len(perm))
	var cyc []int
	for i := 0; i < len(perm); i++ {
		if !vis[i] {
			j := i
			cnt := 0
			for !vis[j] {
				vis[j] = true
				cnt++
				j = perm[j]
			}
			if cnt > 0 {
				cyc = append(cyc, cnt)
			}
		}
	}
	return cyc
}

func dfsAssign(cyc []int, avail, used [6]int, idx int) int {
	if idx == len(cyc) {
		for i := 0; i < 6; i++ {
			if used[i] != avail[i] {
				return 0
			}
		}
		return 1
	}
	cnt := 0
	l := cyc[idx]
	for c := 0; c < 6; c++ {
		if used[c]+l <= avail[c] {
			used[c] += l
			cnt += dfsAssign(cyc, avail, used, idx+1)
			used[c] -= l
		}
	}
	return cnt
}

func solve(input string) string {
	s := strings.TrimSpace(input)
	var avail [6]int
	for _, ch := range s {
		if idx, ok := colorMap[ch]; ok {
			avail[idx]++
		}
	}

	total := 0
	for _, p := range perms {
		cyc := cycles(p)
		total += dfsAssign(cyc, avail, [6]int{}, 0)
	}
	return fmt.Sprintf("%d", total/len(perms))
}

func genCase(rng *rand.Rand) string {
	chars := []byte{'R', 'O', 'Y', 'G', 'B', 'V'}
	b := make([]byte, 6)
	for i := 0; i < 6; i++ {
		b[i] = chars[rng.Intn(len(chars))]
	}
	return string(b) + "\n"
}

func runCase(bin string, input, expected string) error {
	cmd := exec.Command(bin)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		expected := solve(input)
		if err := runCase(bin, input, expected); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
