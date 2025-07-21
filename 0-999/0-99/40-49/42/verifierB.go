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

type Pos struct{ x, y int }

func posFrom(s string) Pos {
	return Pos{int(s[0] - 'a'), int(s[1] - '1')}
}

func pathClear(a, b Pos, occ [][]int) bool {
	dx, dy := 0, 0
	if b.x > a.x {
		dx = 1
	} else if b.x < a.x {
		dx = -1
	}
	if b.y > a.y {
		dy = 1
	} else if b.y < a.y {
		dy = -1
	}
	x, y := a.x+dx, a.y+dy
	for x != b.x || y != b.y {
		if occ[y][x] != 0 {
			return false
		}
		x += dx
		y += dy
	}
	return true
}

func isAttacked(occ [][]int, rooks []Pos, wk Pos, p Pos) bool {
	for _, r := range rooks {
		if r.x == p.x || r.y == p.y {
			if pathClear(r, p, occ) {
				return true
			}
		}
	}
	dx := wk.x - p.x
	if dx < 0 {
		dx = -dx
	}
	dy := wk.y - p.y
	if dy < 0 {
		dy = -dy
	}
	if dx <= 1 && dy <= 1 {
		return true
	}
	return false
}

func solve(s1, s2, s3, s4 string) string {
	rk1 := posFrom(s1)
	rk2 := posFrom(s2)
	wk := posFrom(s3)
	bk := posFrom(s4)
	occ := make([][]int, 8)
	for i := range occ {
		occ[i] = make([]int, 8)
	}
	occ[rk1.y][rk1.x] = 1
	occ[rk2.y][rk2.x] = 1
	occ[wk.y][wk.x] = 2
	occ[bk.y][bk.x] = 3
	rooks := []Pos{rk1, rk2}
	if !isAttacked(occ, rooks, wk, bk) {
		return "OTHER"
	}
	dirs := []int{-1, 0, 1}
	for _, dx := range dirs {
		for _, dy := range dirs {
			if dx == 0 && dy == 0 {
				continue
			}
			nx, ny := bk.x+dx, bk.y+dy
			if nx < 0 || nx > 7 || ny < 0 || ny > 7 {
				continue
			}
			if nx == wk.x && ny == wk.y {
				continue
			}
			newRooks := make([]Pos, 0, 2)
			for _, r := range rooks {
				if r.x == nx && r.y == ny {
					continue
				}
				newRooks = append(newRooks, r)
			}
			occ2 := make([][]int, 8)
			for i := range occ2 {
				occ2[i] = make([]int, 8)
			}
			for _, r := range newRooks {
				occ2[r.y][r.x] = 1
			}
			occ2[wk.y][wk.x] = 2
			occ2[ny][nx] = 3
			if !isAttacked(occ2, newRooks, wk, Pos{nx, ny}) {
				return "OTHER"
			}
		}
	}
	return "CHECKMATE"
}

func randomPos(rng *rand.Rand) string {
	x := rng.Intn(8)
	y := rng.Intn(8)
	return fmt.Sprintf("%c%d", 'a'+x, y+1)
}

func generateCase(rng *rand.Rand) (string, string) {
	for {
		s1 := randomPos(rng)
		s2 := randomPos(rng)
		s3 := randomPos(rng)
		s4 := randomPos(rng)
		if s1 == s2 || s1 == s3 || s1 == s4 || s2 == s3 || s2 == s4 || s3 == s4 {
			continue
		}
		// kings cannot attack each other
		wk := posFrom(s3)
		bk := posFrom(s4)
		dx := wk.x - bk.x
		if dx < 0 {
			dx = -dx
		}
		dy := wk.y - bk.y
		if dy < 0 {
			dy = -dy
		}
		if dx <= 1 && dy <= 1 {
			continue
		}
		input := fmt.Sprintf("%s %s %s %s\n", s1, s2, s3, s4)
		return input, solve(s1, s2, s3, s4)
	}
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	outStr = strings.ToUpper(outStr)
	if outStr != expected {
		return fmt.Errorf("expected %s got %s", expected, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
