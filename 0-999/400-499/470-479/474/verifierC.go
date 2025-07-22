package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type Mole struct {
	x, y, a, b int
}

func rotate90(x, y, a, b int) (int, int) {
	dx := x - a
	dy := y - b
	return a - dy, b + dx
}

func isSquare(px, py [4]int) bool {
	dists := make([]int, 0, 6)
	for i := 0; i < 4; i++ {
		for j := i + 1; j < 4; j++ {
			dx := px[i] - px[j]
			dy := py[i] - py[j]
			dists = append(dists, dx*dx+dy*dy)
		}
	}
	sort.Ints(dists)
	if dists[0] == 0 {
		return false
	}
	side := dists[0]
	for i := 1; i < 4; i++ {
		if dists[i] != side {
			return false
		}
	}
	return dists[4] == dists[5] && dists[4] == 2*side
}

func minMoves(moles [4]Mole) int {
	var rx [4][4]int
	var ry [4][4]int
	for j := 0; j < 4; j++ {
		cx, cy := moles[j].x, moles[j].y
		for k := 0; k < 4; k++ {
			rx[j][k] = cx
			ry[j][k] = cy
			cx, cy = rotate90(cx, cy, moles[j].a, moles[j].b)
		}
	}
	best := -1
	for k0 := 0; k0 < 4; k0++ {
		for k1 := 0; k1 < 4; k1++ {
			for k2 := 0; k2 < 4; k2++ {
				for k3 := 0; k3 < 4; k3++ {
					cost := k0 + k1 + k2 + k3
					if best != -1 && cost >= best {
						continue
					}
					px := [4]int{rx[0][k0], rx[1][k1], rx[2][k2], rx[3][k3]}
					py := [4]int{ry[0][k0], ry[1][k1], ry[2][k2], ry[3][k3]}
					if isSquare(px, py) {
						best = cost
					}
				}
			}
		}
	}
	return best
}

func genCase(rng *rand.Rand) []Mole {
	n := rng.Intn(5) + 1
	moles := make([]Mole, 4*n)
	for i := 0; i < 4*n; i++ {
		moles[i] = Mole{
			x: rng.Intn(11) - 5,
			y: rng.Intn(11) - 5,
			a: rng.Intn(11) - 5,
			b: rng.Intn(11) - 5,
		}
	}
	return moles
}

func expectedC(moles []Mole) []int {
	n := len(moles) / 4
	res := make([]int, n)
	for i := 0; i < n; i++ {
		var g [4]Mole
		copy(g[:], moles[i*4:(i+1)*4])
		res[i] = minMoves(g)
	}
	return res
}

func runCase(bin string, moles []Mole, exp []int) error {
	n := len(moles) / 4
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < 4*n; i++ {
		fmt.Fprintf(&sb, "%d %d %d %d\n", moles[i].x, moles[i].y, moles[i].a, moles[i].b)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	lines := strings.Fields(strings.TrimSpace(out.String()))
	if len(lines) != n {
		return fmt.Errorf("expected %d lines got %d", n, len(lines))
	}
	for i, l := range lines {
		var val int
		fmt.Sscan(l, &val)
		if val != exp[i] {
			return fmt.Errorf("regiment %d expected %d got %d", i+1, exp[i], val)
		}
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
		moles := genCase(rng)
		exp := expectedC(moles)
		if err := runCase(bin, moles, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
