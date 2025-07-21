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

type move struct {
	c byte
	k int
}

func expected(a, b, n, m int, chips [][2]int, moves []move) (int, int) {
	chipSet := make(map[[2]int]struct{})
	for _, p := range chips {
		chipSet[p] = struct{}{}
	}
	r1, r2 := 1, a
	c1, c2 := 1, b
	alice, bob := 0, 0
	turn := 0
	for _, mv := range moves {
		var removed [][2]int
		switch mv.c {
		case 'U':
			for i := r1; i < r1+mv.k; i++ {
				for j := c1; j <= c2; j++ {
					p := [2]int{i, j}
					if _, ok := chipSet[p]; ok {
						removed = append(removed, p)
					}
				}
			}
			r1 += mv.k
		case 'D':
			for i := r2 - mv.k + 1; i <= r2; i++ {
				for j := c1; j <= c2; j++ {
					p := [2]int{i, j}
					if _, ok := chipSet[p]; ok {
						removed = append(removed, p)
					}
				}
			}
			r2 -= mv.k
		case 'L':
			for j := c1; j < c1+mv.k; j++ {
				for i := r1; i <= r2; i++ {
					p := [2]int{i, j}
					if _, ok := chipSet[p]; ok {
						removed = append(removed, p)
					}
				}
			}
			c1 += mv.k
		case 'R':
			for j := c2 - mv.k + 1; j <= c2; j++ {
				for i := r1; i <= r2; i++ {
					p := [2]int{i, j}
					if _, ok := chipSet[p]; ok {
						removed = append(removed, p)
					}
				}
			}
			c2 -= mv.k
		}
		if turn%2 == 0 {
			alice += len(removed)
		} else {
			bob += len(removed)
		}
		for _, p := range removed {
			delete(chipSet, p)
		}
		turn++
	}
	return alice, bob
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	dirs := []byte{'U', 'D', 'L', 'R'}
	for caseNum := 0; caseNum < 100; caseNum++ {
		a := rng.Intn(4) + 3
		b := rng.Intn(4) + 3
		n := rng.Intn(4) + 1
		m := rng.Intn(4) + 1
		chipMap := make(map[[2]int]struct{})
		chips := make([][2]int, 0, n)
		for len(chips) < n {
			x := rng.Intn(a) + 1
			y := rng.Intn(b) + 1
			p := [2]int{x, y}
			if _, ok := chipMap[p]; ok {
				continue
			}
			chipMap[p] = struct{}{}
			chips = append(chips, p)
		}
		moves := make([]move, m)
		r1, r2 := 1, a
		c1, c2 := 1, b
		for i := 0; i < m; i++ {
			dir := dirs[rng.Intn(4)]
			var limit int
			if dir == 'U' || dir == 'D' {
				limit = r2 - r1 + 1
			} else {
				limit = c2 - c1 + 1
			}
			if limit <= 1 {
				moves[i] = move{dir, 0}
				continue
			}
			k := rng.Intn(limit-1) + 1
			moves[i] = move{dir, k}
			switch dir {
			case 'U':
				r1 += k
			case 'D':
				r2 -= k
			case 'L':
				c1 += k
			case 'R':
				c2 -= k
			}
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", a, b, n, m))
		for _, p := range chips {
			sb.WriteString(fmt.Sprintf("%d %d\n", p[0], p[1]))
		}
		for _, mv := range moves {
			sb.WriteString(fmt.Sprintf("%c %d\n", mv.c, mv.k))
		}
		input := sb.String()
		exa, exb := expected(a, b, n, m, chips, moves)
		expectedStr := fmt.Sprintf("%d %d", exa, exb)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", caseNum+1, err, input)
			os.Exit(1)
		}
		if got != expectedStr {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", caseNum+1, expectedStr, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
