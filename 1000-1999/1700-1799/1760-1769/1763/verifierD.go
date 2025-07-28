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

const MOD int64 = 1_000_000_007

var fact [105]int64
var invFact [105]int64

func modPow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func init() {
	fact[0] = 1
	for i := 1; i < len(fact); i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invFact[len(fact)-1] = modPow(fact[len(fact)-1], MOD-2)
	for i := len(fact) - 1; i > 0; i-- {
		invFact[i-1] = invFact[i] * int64(i) % MOD
	}
}

func comb(n, k int) int64 {
	if k < 0 || k > n {
		return 0
	}
	return fact[n] * invFact[k] % MOD * invFact[n-k] % MOD
}

func solveOne(n, i, j, x, y int) int64 {
	ans := int64(0)
	for k := 2; k <= n-1; k++ {
		if x == n && i != k {
			continue
		}
		if y == n && j != k {
			continue
		}
		if i == k && x != n {
			continue
		}
		if j == k && y != n {
			continue
		}
		var posX, posY byte
		if i == k {
			posX = 'P'
		} else if i < k {
			posX = 'L'
		} else {
			posX = 'R'
		}
		if j == k {
			posY = 'P'
		} else if j < k {
			posY = 'L'
		} else {
			posY = 'R'
		}
		L := k - 1
		switch {
		case posX == 'L' && posY == 'L':
			if x >= n || y >= n || x >= y || j > k-1 {
				continue
			}
			a := x - 1
			b := y - x - 1
			c := n - 1 - y
			if a < i-1 || b < j-i-1 || c < L-j {
				continue
			}
			val := comb(a, i-1) * comb(b, j-i-1) % MOD * comb(c, L-j) % MOD
			ans = (ans + val) % MOD
		case posX == 'L' && posY == 'P':
			if x >= n || j != k || y != n || L < i {
				continue
			}
			a := x - 1
			b := n - 1 - x
			if a < i-1 || b < L-i {
				continue
			}
			val := comb(a, i-1) * comb(b, L-i) % MOD
			ans = (ans + val) % MOD
		case posX == 'L' && posY == 'R':
			if x >= n || y >= n || L < i || j <= k {
				continue
			}
			p := L - i
			if x < y {
				a := x - 1
				b := y - x - 1
				c := n - 1 - y
				s := n - y - j + k
				if s < 0 || s > p || s > c || p-s > b || i-1 > a {
					continue
				}
				val := comb(a, i-1) * comb(b, p-s) % MOD * comb(c, s) % MOD
				ans = (ans + val) % MOD
			} else if x > y {
				a := y - 1
				b := x - y - 1
				c := n - 1 - x
				r := n - j - y + i
				if r < 0 || r > b || r > i-1 || i-1-r > a || p > c {
					continue
				}
				val := comb(a, i-1-r) * comb(b, r) % MOD * comb(c, p) % MOD
				ans = (ans + val) % MOD
			}
		case posX == 'P' && posY == 'R':
			if x != n || j <= k {
				continue
			}
			a := y - 1
			b := n - 1 - y
			aa := a - (n - j)
			bb := b - (j - k - 1)
			if aa < 0 || aa > a || bb < 0 || bb > b || aa+bb != L {
				continue
			}
			val := comb(a, aa) * comb(b, bb) % MOD
			ans = (ans + val) % MOD
		case posX == 'R' && posY == 'R':
			if x >= n || y >= n || k >= i || i >= j || x <= y {
				continue
			}
			a := y - 1
			b := x - y - 1
			c := n - 1 - x
			aa := a - (n - j)
			bb := b - (j - i - 1)
			cc := c - (i - k - 1)
			if aa < 0 || aa > a || bb < 0 || bb > b || cc < 0 || cc > c || aa+bb+cc != L {
				continue
			}
			val := comb(a, aa) * comb(b, bb) % MOD * comb(c, cc) % MOD
			ans = (ans + val) % MOD
		}
	}
	return ans % MOD
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 3
	i := rng.Intn(n-2) + 1
	j := rng.Intn(n-i-1) + i + 1
	x := rng.Intn(n) + 1
	y := rng.Intn(n) + 1
	for y == x {
		y = rng.Intn(n) + 1
	}
	var input strings.Builder
	input.WriteString("1\n")
	input.WriteString(fmt.Sprintf("%d %d %d %d %d\n", n, i, j, x, y))
	ans := solveOne(n, i, j, x, y)
	return input.String(), fmt.Sprintf("%d", ans)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
