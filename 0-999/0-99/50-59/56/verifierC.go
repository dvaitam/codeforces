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

type Node struct {
	name     string
	children []*Node
}

func randomName(rng *rand.Rand) string {
	l := rng.Intn(5) + 1
	b := make([]byte, l)
	for i := range b {
		b[i] = byte('A' + rng.Intn(26))
	}
	return string(b)
}

func buildTree(rng *rand.Rand, depth int) *Node {
	node := &Node{name: randomName(rng)}
	if depth < 3 {
		c := rng.Intn(3)
		node.children = make([]*Node, c)
		for i := 0; i < c; i++ {
			node.children[i] = buildTree(rng, depth+1)
		}
	}
	return node
}

func writeDesc(b *strings.Builder, n *Node) {
	b.WriteString(n.name)
	if len(n.children) > 0 {
		b.WriteByte(':')
		for i, ch := range n.children {
			if i > 0 {
				b.WriteByte(',')
			}
			writeDesc(b, ch)
		}
	}
	b.WriteByte('.')
}

func expected(desc string) int64 {
	s := desc
	pos := 0
	cnt := map[string]int{}
	var ans int64
	var parse func()
	parse = func() {
		start := pos
		for pos < len(s) && s[pos] >= 'A' && s[pos] <= 'Z' {
			pos++
		}
		name := s[start:pos]
		ans += int64(cnt[name])
		cnt[name]++
		if pos < len(s) && s[pos] == ':' {
			pos++
			for {
				parse()
				if pos < len(s) && s[pos] == ',' {
					pos++
					continue
				}
				break
			}
		}
		if pos < len(s) && s[pos] == '.' {
			pos++
		}
		cnt[name]--
	}
	parse()
	return ans
}

func generateCase(rng *rand.Rand) (string, int64) {
	root := buildTree(rng, 0)
	var b strings.Builder
	writeDesc(&b, root)
	desc := b.String()
	exp := expected(desc)
	return desc + "\n", exp
}

func runCase(bin, input string, expected int64) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	gotStr := strings.TrimSpace(out.String())
	var got int64
	_, err := fmt.Sscan(gotStr, &got)
	if err != nil {
		return fmt.Errorf("invalid output: %s", gotStr)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
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
	for i := 0; i < 100; i++ {
		input, exp := generateCase(rng)
		if err := runCase(bin, input, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
