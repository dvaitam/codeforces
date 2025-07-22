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

const mod = 1000000007

func matMul(a, b [][]int, n int) [][]int {
	c := make([][]int, n)
	for i := 0; i < n; i++ {
		c[i] = make([]int, n)
		ai := a[i]
		ci := c[i]
		for k := 0; k < n; k++ {
			if ai[k] == 0 {
				continue
			}
			aik := ai[k]
			bk := b[k]
			for j := 0; j < n; j++ {
				ci[j] = (ci[j] + aik*bk[j]) % mod
			}
		}
	}
	return c
}

func vecMatMul(v []int, m [][]int, n int) []int {
	res := make([]int, n)
	for i := 0; i < n; i++ {
		if v[i] == 0 {
			continue
		}
		vi := v[i]
		mi := m[i]
		for j := 0; j < n; j++ {
			res[j] = (res[j] + vi*mi[j]) % mod
		}
	}
	return res
}

func applyPower(dp []int, mat [][]int, exp int) []int {
	n := len(dp)
	base := make([][]int, n)
	for i := 0; i < n; i++ {
		base[i] = make([]int, n)
		copy(base[i], mat[i])
	}
	res := make([]int, n)
	copy(res, dp)
	first := true
	for exp > 0 {
		if exp&1 != 0 {
			if first {
				res = vecMatMul(res, base, n)
				first = false
			} else {
				res = vecMatMul(res, base, n)
			}
		}
		exp >>= 1
		if exp > 0 {
			base = matMul(base, base, n)
		}
	}
	return res
}

func solveCase(w [8]int) string {
	T := make([][][]int, 8)
	for h := 1; h <= 7; h++ {
		n := 1 << h
		m := make([][]int, n)
		for a := 0; a < n; a++ {
			m[a] = make([]int, n)
		}
		maxH := 1 << (h - 1)
		for a := 0; a < n; a++ {
			for b := 0; b < n; b++ {
				cnt := 0
				for H := 0; H < maxH; H++ {
					ok := true
					for r := 0; r < h; r++ {
						painted := 0
						if r == 0 {
							painted++
						} else if H&(1<<(r-1)) != 0 {
							painted++
						}
						if r == h-1 {
							painted++
						} else if H&(1<<r) != 0 {
							painted++
						}
						if a&(1<<r) != 0 {
							painted++
						}
						if b&(1<<r) != 0 {
							painted++
						}
						if painted == 4 {
							ok = false
							break
						}
					}
					if ok {
						cnt++
					}
				}
				m[a][b] = cnt
			}
		}
		T[h] = m
	}
	var dp []int
	hPrev := 0
	for h := 1; h <= 7; h++ {
		wi := w[h]
		if wi == 0 {
			continue
		}
		if hPrev == 0 {
			size := 1 << h
			dp = make([]int, size)
			dp[size-1] = 1
		} else {
			size2 := 1 << h
			newDp := make([]int, size2)
			min := hPrev
			if h < hPrev {
				min = h
			}
			onesHi := 0
			if h > hPrev {
				onesHi = ((1 << (h - hPrev)) - 1) << hPrev
			}
			onesPrev := 0
			if hPrev > h {
				onesPrev = ((1 << (hPrev - h)) - 1) << h
			}
			for mask1, v := range dp {
				if v == 0 {
					continue
				}
				if hPrev > h && (mask1&onesPrev) != onesPrev {
					continue
				}
				mask2 := mask1 & ((1 << min) - 1)
				mask2 |= onesHi
				newDp[mask2] = (newDp[mask2] + v) % mod
			}
			dp = newDp
		}
		dp = applyPower(dp, T[h], wi)
		hPrev = h
	}
	ans := dp[(1<<hPrev)-1]
	return fmt.Sprintf("%d", ans)
}

func runBinary(bin, input string) (string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func generateCase(rng *rand.Rand) (string, [8]int) {
	var w [8]int
	for i := 1; i <= 7; i++ {
		w[i] = rng.Intn(4)
	}
	var sb strings.Builder
	for i := 1; i <= 7; i++ {
		sb.WriteString(fmt.Sprintf("%d ", w[i]))
	}
	sb.WriteString("\n")
	return sb.String(), w
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	inputs := []string{"1 1 1 1 1 1 1\n"}
	ws := [][8]int{{0, 1, 1, 1, 1, 1, 1, 1}}
	for i := 0; i < 100; i++ {
		in, w := generateCase(rng)
		inputs = append(inputs, in)
		ws = append(ws, w)
	}
	for i := range inputs {
		out, err := runBinary(bin, inputs[i])
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		exp := solveCase(ws[i])
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:%sexpected: %s\nfound: %s\n", i+1, inputs[i], exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
