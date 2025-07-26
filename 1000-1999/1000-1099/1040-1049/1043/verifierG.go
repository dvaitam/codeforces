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

const (
	mod1 = 1000000007
	mod2 = 1000000009
	base = 91138233
)

func preprocess(s string) ([]int, []int, []int, []int, [][26]int) {
	n := len(s)
	h1 := make([]int, n+1)
	h2 := make([]int, n+1)
	pow1 := make([]int, n+1)
	pow2 := make([]int, n+1)
	pow1[0], pow2[0] = 1, 1
	for i := 1; i <= n; i++ {
		pow1[i] = int((int64(pow1[i-1]) * base) % mod1)
		pow2[i] = int((int64(pow2[i-1]) * base) % mod2)
	}
	for i := 1; i <= n; i++ {
		c := int(s[i-1])
		h1[i] = int((int64(h1[i-1])*base + int64(c)) % mod1)
		h2[i] = int((int64(h2[i-1])*base + int64(c)) % mod2)
	}
	cnt := make([][26]int, n+1)
	for i := 1; i <= n; i++ {
		for c := 0; c < 26; c++ {
			cnt[i][c] = cnt[i-1][c]
		}
		cnt[i][s[i-1]-'a']++
	}
	return h1, h2, pow1, pow2, cnt
}

func substrHash(h, pow []int, l, r, mod int) int {
	val := (h[r] - int(int64(h[l-1])*int64(pow[r-l])%int64(mod)) + mod) % mod
	return val
}

func expected(n int, s string, qs [][2]int) string {
	h1, h2, pow1, pow2, cnt := preprocess(s)
	var sb strings.Builder
	for _, q := range qs {
		l, r := q[0], q[1]
		k := r - l + 1
		if k < 2 {
			sb.WriteString("-1\n")
			continue
		}
		distinct := 0
		hasRepeat := false
		for c := 0; c < 26; c++ {
			f := cnt[r][c] - cnt[l-1][c]
			if f > 0 {
				distinct++
				if f > 1 {
					hasRepeat = true
				}
			}
		}
		if !hasRepeat {
			sb.WriteString("-1\n")
			continue
		}
		periodic := false
		if distinct == 1 {
			periodic = true
		} else {
			getHashEqual := func(d int) bool {
				h1a := substrHash(h1, pow1, l, r-d, mod1)
				h1b := substrHash(h1, pow1, l+d, r, mod1)
				if h1a != h1b {
					return false
				}
				h2a := substrHash(h2, pow2, l, r-d, mod2)
				h2b := substrHash(h2, pow2, l+d, r, mod2)
				return h2a == h2b
			}
			for di := 1; di*di <= k; di++ {
				if k%di != 0 {
					continue
				}
				if di > 1 && di < k && getHashEqual(di) {
					periodic = true
					break
				}
				d2 := k / di
				if d2 > 1 && d2 < k && d2 != di && getHashEqual(d2) {
					periodic = true
					break
				}
			}
		}
		if periodic {
			sb.WriteString("1\n")
			continue
		}
		if s[l-1] == s[r-1] || (l+1 <= r && s[l-1] == s[l]) || (r-1 >= l && s[r-1] == s[r-2]) {
			sb.WriteString("2\n")
		} else {
			sb.WriteString("3\n")
		}
	}
	return sb.String()
}

func runCase(bin, input, want string) error {
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
	if strings.TrimSpace(want) != got {
		return fmt.Errorf("expected %q got %q", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	type test struct {
		n  int
		s  string
		qs [][2]int
	}

	tests := []test{
		{n: 3, s: "aba", qs: [][2]int{{1, 3}}},
	}
	for i := 0; i < 100; i++ {
		n := rng.Intn(10) + 1
		b := make([]byte, n)
		for j := range b {
			b[j] = byte('a' + rng.Intn(3))
		}
		qn := rng.Intn(5) + 1
		qs := make([][2]int, qn)
		for j := 0; j < qn; j++ {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			qs[j] = [2]int{l, r}
		}
		tests = append(tests, test{n: n, s: string(b), qs: qs})
	}

	for idx, tc := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n%s\n%d\n", tc.n, tc.s, len(tc.qs)))
		for _, p := range tc.qs {
			sb.WriteString(fmt.Sprintf("%d %d\n", p[0], p[1]))
		}
		want := expected(tc.n, tc.s, tc.qs)
		if err := runCase(bin, sb.String(), strings.TrimSpace(want)); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, sb.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
