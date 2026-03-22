package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const MOD = 1000000007

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

// expected computes the answer using the efficient DP approach embedded
// directly (matching the accepted solution logic).
func expected(n, k, m int) int64 {
	const SIZE = 1 << 21
	dp := make([]int, SIZE)
	newDp := make([]int, SIZE)

	active := make([]int, 0, 10000)
	newActive := make([]int, 0, 10000)

	dp[0] = 1
	active = append(active, 0)

	var c [4]int
	var nextC [4]int

	for i := 1; i <= n; i++ {
		newActive = newActive[:0]

		for _, state := range active {
			ways := dp[state]
			dp[state] = 0
			if ways == 0 {
				continue
			}

			j := state & 0xF
			f := (state >> 4) & 1
			activeComps := 0
			for d := 0; d < m; d++ {
				val := (state >> (5 + 4*d)) & 0xF
				c[d] = val
				activeComps += val
			}

			remExp := c[m-1]

			addState := func(nj, nf int, nc *[4]int) {
				ntotal := nf
				for d := 0; d < m; d++ {
					ntotal += nc[d]
				}
				if ntotal > k-nj+1 {
					return
				}
				enc := nj | (nf << 4)
				for d := 0; d < m; d++ {
					enc |= (nc[d] << (5 + 4*d))
				}
				if newDp[enc] == 0 {
					newActive = append(newActive, enc)
				}
				newDp[enc] = (newDp[enc] + ways) % MOD
			}

			// Action 1: Skip node i
			{
				valid := true
				nextF := f
				if remExp > 0 {
					if remExp > 1 || f == 1 {
						valid = false
					} else {
						nextF = 1
					}
				}
				if valid {
					for d := 0; d < m-1; d++ {
						nextC[d+1] = c[d]
					}
					nextC[0] = 0
					addState(j, nextF, &nextC)
				}
			}

			if j < k {
				// Action 2a: Start new component
				{
					valid := true
					nextF := f
					if remExp > 0 {
						if remExp > 1 || f == 1 {
							valid = false
						} else {
							nextF = 1
						}
					}
					if valid {
						for d := 0; d < m-1; d++ {
							nextC[d+1] = c[d]
						}
						nextC[0] = 1
						addState(j+1, nextF, &nextC)
					}
				}

				// Action 2b: Append to existing component
				for src := 0; src < m; src++ {
					count := c[src]
					if count == 0 {
						continue
					}
					currRemExp := remExp
					if src == m-1 {
						currRemExp--
					}
					valid := true
					nextF := f
					if currRemExp > 0 {
						if currRemExp > 1 || f == 1 {
							valid = false
						} else {
							nextF = 1
						}
					}
					if valid {
						for d := 0; d < m-1; d++ {
							cnt := c[d]
							if d == src {
								cnt--
							}
							nextC[d+1] = cnt
						}
						nextC[0] = 1
						oldWays := ways
						ways = int((int64(ways) * int64(count)) % int64(MOD))
						addState(j+1, nextF, &nextC)
						ways = oldWays
					}
				}

				// Action 2c: Prepend to existing component
				totalComps := activeComps + f
				if totalComps > 0 {
					valid := true
					nextF := f
					if remExp > 0 {
						if remExp > 1 || f == 1 {
							valid = false
						} else {
							nextF = 1
						}
					}
					if valid {
						for d := 0; d < m-1; d++ {
							nextC[d+1] = c[d]
						}
						nextC[0] = 0
						oldWays := ways
						ways = int((int64(ways) * int64(totalComps)) % int64(MOD))
						addState(j+1, nextF, &nextC)
						ways = oldWays
					}
				}

				// Action 2d: Merge two components
				if totalComps >= 2 {
					for src := 0; src < m; src++ {
						count := c[src]
						if count == 0 {
							continue
						}
						currRemExp := remExp
						if src == m-1 {
							currRemExp--
						}
						valid := true
						nextF := f
						if currRemExp > 0 {
							if currRemExp > 1 || f == 1 {
								valid = false
							} else {
								nextF = 1
							}
						}
						if valid {
							for d := 0; d < m-1; d++ {
								cnt := c[d]
								if d == src {
									cnt--
								}
								nextC[d+1] = cnt
							}
							nextC[0] = 0
							oldWays := ways
							mult := int64(count) * int64(totalComps-1)
							ways = int((int64(ways) * mult) % int64(MOD))
							addState(j+1, nextF, &nextC)
							ways = oldWays
						}
					}
				}
			}
		}

		dp, newDp = newDp, dp
		active, newActive = newActive, active
	}

	ans := 0
	for _, state := range active {
		j := state & 0xF
		if j != k {
			continue
		}
		f := (state >> 4) & 1
		activeComps := 0
		for d := 0; d < m; d++ {
			activeComps += (state >> (5 + 4*d)) & 0xF
		}
		if activeComps+f == 1 {
			ans = (ans + dp[state]) % MOD
		}
	}
	return int64(ans)
}

func generateCase(rng *rand.Rand) (string, int64) {
	n := rng.Intn(20) + 1
	k := rng.Intn(n) + 1
	if k > 12 {
		k = 12
	}
	m := rng.Intn(4) + 1
	input := fmt.Sprintf("%d %d %d\n", n, k, m)
	return input, expected(n, k, m)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if strings.HasSuffix(bin, ".go") {
		tmp, err := os.CreateTemp("", "verifierF1-bin-*")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create temp file: %v\n", err)
			os.Exit(1)
		}
		tmp.Close()
		defer os.Remove(tmp.Name())
		out, err := exec.Command("go", "build", "-o", tmp.Name(), bin).CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "compile error: %v\n%s", err, out)
			os.Exit(1)
		}
		bin = tmp.Name()
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	fixed := [][3]int{
		{3, 3, 1},
		{4, 2, 2},
	}
	idx := 0
	for ; idx < len(fixed); idx++ {
		n := fixed[idx][0]
		k := fixed[idx][1]
		m := fixed[idx][2]
		inp := fmt.Sprintf("%d %d %d\n", n, k, m)
		exp := strconv.FormatInt(expected(n, k, m), 10)
		out, err := run(bin, inp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", idx+1, err, inp)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", idx+1, exp, out, inp)
			os.Exit(1)
		}
	}
	for ; idx < 100; idx++ {
		inp, expVal := generateCase(rng)
		out, err := run(bin, inp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", idx+1, err, inp)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strconv.FormatInt(expVal, 10) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:%s", idx+1, expVal, out, inp)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
