package main

import (
	"bytes"
	"container/list"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func bfsSolve(s, t string) int {
	if len(s) != len(t) {
		return -1
	}
	n := len(s)
	count := make(map[rune]int)
	for _, ch := range s {
		count[ch]++
	}
	for _, ch := range t {
		count[ch]--
	}
	for _, v := range count {
		if v != 0 {
			return -1
		}
	}
	visited := map[string]bool{s: true}
	type node struct {
		str  string
		dist int
	}
	q := list.New()
	q.PushBack(node{s, 0})
	for q.Len() > 0 {
		front := q.Remove(q.Front()).(node)
		if front.str == t {
			return front.dist
		}
		curr := front.str
		for l := 0; l < n; l++ {
			for r := l + 1; r < n; r++ {
				bs := []byte(curr)
				ch := bs[r]
				copy(bs[l+1:r+1], bs[l:r])
				bs[l] = ch
				ns := string(bs)
				if !visited[ns] {
					visited[ns] = true
					q.PushBack(node{ns, front.dist + 1})
				}
			}
		}
	}
	return -1
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	sb1 := make([]byte, n)
	sb2 := make([]byte, n)
	for i := 0; i < n; i++ {
		ch := byte('a' + rng.Intn(3))
		sb1[i] = ch
	}
	for i := 0; i < n; i++ {
		sb2[i] = sb1[i]
	}
	// permute sb2
	rng.Shuffle(n, func(i, j int) { sb2[i], sb2[j] = sb2[j], sb2[i] })
	s := string(sb1)
	t := string(sb2)
	input := fmt.Sprintf("1\n%s\n%s\n", s, t)
	expect := bfsSolve(s, t)
	return input, fmt.Sprintf("%d", expect)
}

func runCase(bin, in, expect string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expect {
		return fmt.Errorf("expected %q got %q", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := generateCase(rng)
		if err := runCase(bin, in, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
