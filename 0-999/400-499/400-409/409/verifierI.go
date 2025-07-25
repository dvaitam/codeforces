package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func run(binary string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout")
		}
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(s string) string {
	n := len(s)
	const INF = 1 << 30
	dist := make([][4]int, n)
	for i := 0; i < n; i++ {
		for d := 0; d < 4; d++ {
			dist[i][d] = INF
		}
	}
	type state struct{ p, d int }
	dq := []state{{0, 0}}
	dist[0][0] = 0
	var best = INF
	pushFront := func(st state) { dq = append([]state{st}, dq...) }
	pushBack := func(st state) { dq = append(dq, st) }
	popFront := func() state { st := dq[0]; dq = dq[1:]; return st }
	for len(dq) > 0 {
		u := popFront()
		p, d := u.p, u.d
		cd := dist[p][d]
		if cd >= best {
			continue
		}
		c := s[p]
		if c == '@' {
			best = cd
			continue
		}
		var ndirs []int
		switch c {
		case '>':
			ndirs = []int{0}
		case 'v':
			ndirs = []int{1}
		case '<':
			ndirs = []int{2}
		case '^':
			ndirs = []int{3}
		case '_':
			ndirs = []int{0}
		case '|':
			ndirs = []int{1}
		case '?':
			ndirs = []int{0, 1, 2, 3}
		default:
			ndirs = []int{d}
		}
		for _, nd := range ndirs {
			add := 0
			if c == '&' {
				add = 1
			}
			np := p
			switch nd {
			case 0:
				np = (p + 1) % n
			case 2:
				np = (p - 1 + n) % n
			case 1, 3:
				np = p
			}
			ndist := cd + add
			if ndist < dist[np][nd] {
				dist[np][nd] = ndist
				st := state{np, nd}
				if add == 1 {
					pushBack(st)
				} else {
					pushFront(st)
				}
			}
		}
	}
	if best == INF {
		return "false"
	}
	if best == 0 {
		return ""
	}
	return strings.Repeat("0", best)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierI.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		var program string
		if i < 50 {
			program = "@"
		} else {
			program = "&@"
		}
		input := program + "\n"
		want := solve(program)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d error: %v\n", i+1, err)
			os.Exit(1)
		}
		if out != want {
			fmt.Printf("test %d failed: expected %q got %q\n", i+1, want, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
