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

const mod int64 = 998244353

func modPow(a, e int64) int64 {
	res := int64(1)
	a %= mod
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func expected(w []int64) string {
	n := len(w)
	m := n / 3
	ans := int64(1)
	for i := 0; i < m; i++ {
		a := w[3*i]
		b := w[3*i+1]
		c := w[3*i+2]
		s1 := a + b
		s2 := a + c
		s3 := b + c
		mx := s1
		if s2 > mx {
			mx = s2
		}
		if s3 > mx {
			mx = s3
		}
		cnt := 0
		if s1 == mx {
			cnt++
		}
		if s2 == mx {
			cnt++
		}
		if s3 == mx {
			cnt++
		}
		ans = ans * int64(cnt) % mod
	}
	fact := make([]int64, m+1)
	fact[0] = 1
	for i := 1; i <= m; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	invFact := make([]int64, m+1)
	invFact[m] = modPow(fact[m], mod-2)
	for i := m; i > 0; i-- {
		invFact[i-1] = invFact[i] * int64(i) % mod
	}
	half := m / 2
	comb := fact[m] * invFact[half] % mod * invFact[m-half] % mod
	ans = ans * comb % mod
	return fmt.Sprintf("%d", ans)
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for tc := 0; tc < 100; tc++ {
		m := rng.Intn(5) + 2 // number of triangles 2..6 => n divisible by 3*? but multiple of 6; adjust
		if m%2 != 0 {
			m++ // ensure divisible by 2 to satisfy n divisible by 6
		}
		n := m * 3
		w := make([]int64, n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			w[i] = int64(rng.Intn(10) + 1)
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", w[i]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expectedOut := expected(w)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", tc+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expectedOut {
			fmt.Printf("case %d failed: expected %s got %s\ninput:\n%s", tc+1, expectedOut, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
