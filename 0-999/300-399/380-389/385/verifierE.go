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

func matMul(a, b [4][4]uint64, mod uint64) [4][4]uint64 {
	var c [4][4]uint64
	for i := 0; i < 4; i++ {
		for k := 0; k < 4; k++ {
			if a[i][k] == 0 {
				continue
			}
			for j := 0; j < 4; j++ {
				c[i][j] += a[i][k] * b[k][j]
			}
		}
		for j := 0; j < 4; j++ {
			c[i][j] %= mod
		}
	}
	return c
}

func matPow(mat [4][4]uint64, p, mod uint64) [4][4]uint64 {
	var res [4][4]uint64
	for i := 0; i < 4; i++ {
		res[i][i] = 1 % mod
	}
	for p > 0 {
		if p&1 != 0 {
			res = matMul(res, mat, mod)
		}
		mat = matMul(mat, mat, mod)
		p >>= 1
	}
	return res
}

func expected(n, sx, sy, dx, dy, t int64) (int64, int64) {
	if t == 0 {
		return sx, sy
	}
	mod := uint64(n * 2)
	X0 := uint64((sx - 1) % n)
	Y0 := uint64((sy - 1) % n)
	S0 := (X0 + Y0) % mod
	DS0 := dx + dy
	Sneg1 := (int64(S0) - DS0) % int64(mod)
	if Sneg1 < 0 {
		Sneg1 += int64(mod)
	}
	V0 := [4]uint64{uint64(S0), uint64(Sneg1), 0, 1}
	M := [4][4]uint64{
		{4 % mod, (mod - 1) % mod, 2 % mod, 4 % mod},
		{1, 0, 0, 0},
		{0, 0, 1, 1},
		{0, 0, 0, 1},
	}
	Mt := matPow(M, uint64(t), mod)
	var Vt [4]uint64
	for i := 0; i < 4; i++ {
		var acc uint64
		for j := 0; j < 4; j++ {
			acc += Mt[i][j] * V0[j]
		}
		Vt[i] = acc % mod
	}
	St := Vt[0]
	diff0 := int64(X0) - int64(Y0)
	delta := dx - dy
	tmod := uint64(t % int64(mod))
	deltaMod := (delta%int64(mod) + int64(mod)) % int64(mod)
	diffT := (int64(diff0)%int64(mod) + (int64(tmod)*deltaMod)%int64(mod)) % int64(mod)
	if diffT < 0 {
		diffT += int64(mod)
	}
	d := uint64(diffT)
	sumd := (St + d) % mod
	Xt := (sumd / 2) % uint64(n)
	diff := (St + mod - d) % mod
	Yt := (diff / 2) % uint64(n)
	return int64(Xt + 1), int64(Yt + 1)
}

func generateCase(rng *rand.Rand) (string, [2]int64) {
	n := rng.Int63n(1000) + 1
	sx := rng.Int63n(n) + 1
	sy := rng.Int63n(n) + 1
	dx := rng.Int63n(21) - 10
	dy := rng.Int63n(21) - 10
	t := rng.Int63n(10000)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d %d %d %d\n", n, sx, sy, dx, dy, t))
	x, y := expected(n, sx, sy, dx, dy, t)
	return sb.String(), [2]int64{x, y}
}

func runCase(bin, input string, expected [2]int64) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var x, y int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &x, &y); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if x != expected[0] || y != expected[1] {
		return fmt.Errorf("expected %d %d got %d %d", expected[0], expected[1], x, y)
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
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
