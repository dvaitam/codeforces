package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const mod = int64(1e9 + 7)

type test struct {
	s       string
	q       int
	queries [][2]int
}

func modExp(x, y int64) int64 {
	res := int64(1)
	for y > 0 {
		if y&1 == 1 {
			res = res * x % mod
		}
		x = x * x % mod
		y >>= 1
	}
	return res
}

func pd(c byte) int {
	if c >= 'A' && c <= 'Z' {
		return int(c-'A') + 1
	}
	return int(c-'a') + 27
}

func solve(tc test) string {
	s := tc.s
	n := len(s)
	cnt := make([]int, 53)
	for i := 0; i < n; i++ {
		cnt[pd(s[i])]++
	}
	jc := make([]int64, n+1)
	inv := make([]int64, n+1)
	jc[0] = 1
	for i := 1; i <= n; i++ {
		jc[i] = jc[i-1] * int64(i) % mod
	}
	inv[n] = modExp(jc[n], mod-2)
	for i := n; i > 0; i-- {
		inv[i-1] = inv[i] * int64(i) % mod
	}
	half := n / 2
	g := jc[half] * jc[half] % mod
	for i := 1; i <= 52; i++ {
		if cnt[i] > 0 {
			g = g * inv[cnt[i]] % mod
		}
	}
	p := make([]int64, n+1)
	p[0] = 1
	for i := 1; i <= 52; i++ {
		c := cnt[i]
		if c == 0 {
			continue
		}
		for j := n; j >= c; j-- {
			p[j] = (p[j] + p[j-c]) % mod
		}
	}
	var sres [53][53]int64
	p2 := make([]int64, n+1)
	for i := 1; i <= 52; i++ {
		ci := cnt[i]
		if ci == 0 || ci > half {
			continue
		}
		copy(p2, p)
		for j := ci; j <= n; j++ {
			p2[j] = (p2[j] - p2[j-ci] + mod) % mod
		}
		idx := half - ci
		if idx >= 0 {
			sres[i][i] = p2[idx] * g % mod * 2 % mod
		}
		for j := i + 1; j <= 52; j++ {
			cj := cnt[j]
			if cj == 0 || ci+cj > half {
				continue
			}
			w := half - ci - cj
			var ct int64
			sign := int64(1)
			for ww := w; ww >= 0; ww -= cj {
				ct = (ct + p2[ww]*sign) % mod
				sign = -sign
			}
			sres[i][j] = (ct%mod + mod) % mod * g % mod * 2 % mod
		}
	}
	var out strings.Builder
	for _, q := range tc.queries {
		x := pd(s[q[0]-1])
		y := pd(s[q[1]-1])
		if x > y {
			x, y = y, x
		}
		out.WriteString(fmt.Sprintln(sres[x][y]))
	}
	return strings.TrimSpace(out.String())
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(45))
	tests := make([]test, 0, 100)
	letters := []rune("abAB")
	for i := 0; i < 100; i++ {
		n := 2 * (rng.Intn(2) + 1)
		b := make([]rune, n)
		for j := range b {
			b[j] = letters[rng.Intn(len(letters))]
		}
		q := rng.Intn(3) + 1
		qs := make([][2]int, q)
		for j := 0; j < q; j++ {
			x := rng.Intn(n) + 1
			y := rng.Intn(n) + 1
			for y == x {
				y = rng.Intn(n) + 1
			}
			qs[j] = [2]int{x, y}
		}
		tests = append(tests, test{string(b), q, qs})
	}
	return tests
}

func run(binary string, input string) (string, error) {
	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%s\n%d\n", tc.s, tc.q)
		for _, pq := range tc.queries {
			fmt.Fprintf(&sb, "%d %d\n", pq[0], pq[1])
		}
		input := sb.String()
		want := solve(tc)
		got, err := run(binary, input)
		if err != nil {
			fmt.Printf("Test %d: error running binary: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected:\n%s\nGot:\n%s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
