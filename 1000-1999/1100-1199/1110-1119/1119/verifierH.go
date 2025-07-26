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

const md = 998244353

type testCase struct {
	n, k    int
	x, y, z int
	triples [][3]int
}

func add(x, y int) int {
	x += y
	if x >= md {
		x -= md
	}
	return x
}
func sub(x, y int) int {
	x -= y
	if x < 0 {
		x += md
	}
	return x
}
func mul(x, y int) int { return int((int64(x) * int64(y)) % md) }

func power(x, y int) int {
	res := 1
	for y > 0 {
		if y&1 == 1 {
			res = mul(res, x)
		}
		x = mul(x, x)
		y >>= 1
	}
	return res
}

func fwt(a []int) {
	n := len(a)
	for l := 1; l < n; l <<= 1 {
		for i := 0; i < n; i += l << 1 {
			for j := 0; j < l; j++ {
				u := a[i+j]
				v := a[i+j+l]
				a[i+j] = add(u, v)
				a[i+j+l] = sub(u, v)
			}
		}
	}
}

func referenceSolve(tc testCase) string {
	n, m := tc.n, tc.k
	N := 1 << m
	foo := make([]int, N)
	bar := make([]int, N)
	baz := make([]int, N)
	xorAll := 0
	for _, t := range tc.triples {
		a := t[0] ^ t[2]
		b := t[1] ^ t[2]
		xorAll ^= t[2]
		foo[a]++
		bar[b]++
		baz[a^b]++
	}
	fwt(foo)
	fwt(bar)
	fwt(baz)
	ans := make([]int, N)
	inv2 := (md + 1) / 2
	x := ((tc.x % md) + md) % md
	y := ((tc.y % md) + md) % md
	z := ((tc.z % md) + md) % md
	b1 := add((x+y)%md, z)
	b2 := add((x-y+md)%md, z)
	b3 := add((-x+y+md)%md, z)
	b4 := add(((-x - y + 2*md) % md), z)
	for i := 0; i < N; i++ {
		foo[i] = add(foo[i], n)
		bar[i] = add(bar[i], n)
		baz[i] = add(baz[i], n)
		a1 := mul(foo[i], inv2)
		b1t := mul(bar[i], inv2)
		c1 := mul(baz[i], inv2)
		d := (a1 + b1t + c1 - n) / 2
		e := a1 - d
		f := b1t - d
		g := c1 - d
		res := 1
		res = mul(res, power(b1, d))
		res = mul(res, power(b2, e))
		res = mul(res, power(b3, f))
		res = mul(res, power(b4, g))
		ans[i] = res
	}
	fwt(ans)
	invN := power(N, md-2)
	var sb strings.Builder
	for i := 0; i < N; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		idx := xorAll ^ i
		sb.WriteString(fmt.Sprintf("%d", mul(ans[idx], invN)))
	}
	return sb.String()
}

func buildCase(n, k int, x, y, z int, triples [][3]int) (string, string) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d %d %d\n", n, k, x, y, z))
	for _, t := range triples {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", t[0], t[1], t[2]))
	}
	return sb.String(), referenceSolve(testCase{n: n, k: k, x: x, y: y, z: z, triples: triples})
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	k := rng.Intn(4) + 1
	x := rng.Intn(10)
	y := rng.Intn(10)
	z := rng.Intn(10)
	triples := make([][3]int, n)
	maxVal := 1 << k
	for i := 0; i < n; i++ {
		triples[i][0] = rng.Intn(maxVal)
		triples[i][1] = rng.Intn(maxVal)
		triples[i][2] = rng.Intn(maxVal)
	}
	return buildCase(n, k, x, y, z, triples)
}

func runCase(bin, input, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := []string{}
	exps := []string{}
	in, exp := buildCase(1, 1, 1, 1, 1, [][3]int{{0, 0, 0}})
	cases = append(cases, in)
	exps = append(exps, exp)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(cases) < 102 {
		in, exp := genCase(rng)
		cases = append(cases, in)
		exps = append(exps, exp)
	}
	for i := range cases {
		if err := runCase(bin, cases[i], exps[i]); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, cases[i])
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
