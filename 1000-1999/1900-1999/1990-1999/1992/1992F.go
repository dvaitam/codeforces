package main

import (
	"bufio"
	"fmt"
	"os"
)

type factor struct {
	p int
	e int
}

func factorize(x int) []factor {
	res := []factor{}
	for i := 2; i*i <= x; i++ {
		if x%i == 0 {
			c := 0
			for x%i == 0 {
				x /= i
				c++
			}
			res = append(res, factor{i, c})
		}
	}
	if x > 1 {
		res = append(res, factor{x, 1})
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, x int
		fmt.Fscan(in, &n, &x)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		fac := factorize(x)
		m := len(fac)
		base := make([]int, m)
		base[0] = 1
		for i := 1; i < m; i++ {
			base[i] = base[i-1] * (fac[i-1].e + 1)
		}
		states := base[m-1] * (fac[m-1].e + 1)
		decode := make([][]int, states)
		for s := 0; s < states; s++ {
			decode[s] = make([]int, m)
			for j := 0; j < m; j++ {
				decode[s][j] = (s / base[j]) % (fac[j].e + 1)
			}
		}
		exps := make([]int, m)
		for j := 0; j < m; j++ {
			exps[j] = fac[j].e
		}
		target := states - 1
		dp := make([]bool, states)
		dp[0] = true
		segments := 1
		for _, val := range arr {
			vec := make([]int, m)
			tmp := val
			valid := true
			for j, f := range fac {
				c := 0
				for tmp%f.p == 0 {
					tmp /= f.p
					c++
				}
				if c > f.e {
					c = f.e
				}
				vec[j] = c
			}
			if tmp != 1 {
				valid = false
			}
			newdp := make([]bool, states)
			copy(newdp, dp)
			if valid {
				for s := 0; s < states; s++ {
					if dp[s] {
						idx := 0
						for j := 0; j < m; j++ {
							v := decode[s][j] + vec[j]
							if v > exps[j] {
								v = exps[j]
							}
							idx += v * base[j]
						}
						if !newdp[idx] {
							newdp[idx] = true
						}
					}
				}
			}
			if newdp[target] {
				segments++
				dp = make([]bool, states)
				dp[0] = true
				if valid {
					for s := 0; s < states; s++ {
						if dp[s] {
							idx := 0
							for j := 0; j < m; j++ {
								v := decode[s][j] + vec[j]
								if v > exps[j] {
									v = exps[j]
								}
								idx += v * base[j]
							}
							if !dp[idx] {
								dp[idx] = true
							}
						}
					}
				}
			} else {
				dp = newdp
			}
		}
		fmt.Fprintln(out, segments)
	}
}
