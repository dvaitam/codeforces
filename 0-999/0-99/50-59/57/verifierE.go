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

const MOD = 1000000007
const KMAX = 200

func solve(k int64, blocks [][2]int) int64 {
	off := KMAX*2 + 5
	size := off*2 + 1
	dist := make([][]int, size)
	for i := range dist {
		dist[i] = make([]int, size)
		for j := range dist[i] {
			dist[i][j] = -1
		}
	}
	moves := [8][2]int{{1, 2}, {2, 1}, {2, -1}, {1, -2}, {-1, -2}, {-2, -1}, {-2, 1}, {-1, 2}}
	type P struct{ x, y int }
	q := make([]P, 0, 500000)
	sx, sy := off, off
	dist[sx][sy] = 0
	q = append(q, P{sx, sy})
	cnt := make([]int64, KMAX+1)
	cnt[0] = 1
	for head := 0; head < len(q); head++ {
		p := q[head]
		d := dist[p.x][p.y]
		if d >= KMAX {
			continue
		}
		nd := d + 1
		for _, mv := range moves {
			nx, ny := p.x+mv[0], p.y+mv[1]
			if nx < 0 || ny < 0 || nx >= size || ny >= size {
				continue
			}
			if dist[nx][ny] != -1 {
				continue
			}
			dist[nx][ny] = nd
			cnt[nd]++
			q = append(q, P{nx, ny})
		}
	}
	T := make([]int64, KMAX+1)
	T[0] = cnt[0]
	for i := 1; i <= KMAX; i++ {
		T[i] = T[i-1] + cnt[i]
	}
	D1 := make([]int64, KMAX+1)
	for i := 1; i <= KMAX; i++ {
		D1[i] = T[i] - T[i-1]
	}
	D2 := make([]int64, KMAX+1)
	for i := 2; i <= KMAX; i++ {
		D2[i] = D1[i] - D1[i-1]
	}
	thr := 50
	D3 := make([]int64, KMAX+1)
	for i := 3; i <= KMAX; i++ {
		D3[i] = D2[i] - D2[i-1]
	}
	baseD3 := D3[thr]
	A3 := baseD3 / 6
	B3 := (D2[thr] - 6*A3*int64(thr-1)) / 2
	C3 := D1[thr] - A3*(3*int64(thr)*int64(thr)-3*int64(thr)+1) - B3*(2*int64(thr)-1)
	D0 := T[0]
	var total int64
	if k <= KMAX {
		total = T[k]
	} else {
		kk := k % MOD
		total = A3 % MOD * kk % MOD * kk % MOD * kk % MOD
		total = (total + B3%MOD*kk%MOD*kk%MOD) % MOD
		total = (total + C3%MOD*kk%MOD) % MOD
		total = (total + D0%MOD) % MOD
		if total < 0 {
			total += MOD
		}
	}
	var sub int64
	for _, b := range blocks {
		x := b[0] + off
		y := b[1] + off
		if x >= 0 && x < size && y >= 0 && y < size {
			d := int64(dist[x][y])
			if d >= 0 && d <= k {
				sub++
			}
		}
	}
	ans := (total - sub) % MOD
	if ans < 0 {
		ans += MOD
	}
	return ans
}

func genCase(rng *rand.Rand) (string, string) {
	k := rng.Int63n(1000000) + 1
	n := rng.Intn(10)
	blocks := make([][2]int, n)
	used := make(map[[2]int]bool)
	for i := 0; i < n; i++ {
		for {
			x := rng.Intn(21) - 10
			y := rng.Intn(21) - 10
			if x == 0 && y == 0 {
				continue
			}
			key := [2]int{x, y}
			if !used[key] {
				used[key] = true
				blocks[i] = key
				break
			}
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", k, n)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d %d\n", blocks[i][0], blocks[i][1])
	}
	exp := fmt.Sprintf("%d\n", solve(k, blocks))
	return sb.String(), exp
}

func runCase(bin, input, exp string) error {
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
	if got != strings.TrimSpace(exp) {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
