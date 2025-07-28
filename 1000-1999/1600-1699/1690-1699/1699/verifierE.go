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

func run(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

const maxM = 1000000

var spf [maxM + 1]int

func init() {
	for i := 2; i <= maxM; i++ {
		if spf[i] == 0 {
			spf[i] = i
			if i <= maxM/i {
				for j := i * i; j <= maxM; j += i {
					if spf[j] == 0 {
						spf[j] = i
					}
				}
			}
		}
	}
	spf[1] = 1
}

func maxPrimeFactor(x int) int {
	res := 1
	for x > 1 {
		p := spf[x]
		if p > res {
			res = p
		}
		for x%p == 0 {
			x /= p
		}
	}
	return res
}

func solveE(n, m int, arr []int) int {
	maxVal := 0
	maxPF := 1
	for _, v := range arr {
		if v > maxVal {
			maxVal = v
		}
		pf := maxPrimeFactor(v)
		if pf > maxPF {
			maxPF = pf
		}
	}
	if maxVal == 1 {
		return 0
	}
	freq := make([]int, maxVal+1)
	for _, v := range arr {
		freq[v]++
	}
	dp := make([]int, maxVal+1)
	cnt := make([]int, maxVal+1)
	for val, c := range freq {
		if c > 0 {
			if val == 1 {
				dp[val] = 1
			}
			cnt[dp[val]] += c
		}
	}
	curMin := 0
	for curMin < len(cnt) && cnt[curMin] == 0 {
		curMin++
	}
	best := int(1 << 60)
	for r := 2; r <= maxVal; r++ {
		if r > dp[r] {
			old := dp[r]
			dp[r] = r
			if freq[r] > 0 {
				cnt[old] -= freq[r]
				cnt[r] += freq[r]
				if old == curMin && cnt[old] == 0 {
					for curMin < len(cnt) && cnt[curMin] == 0 {
						curMin++
					}
				}
			}
		}
		for j := r * 2; j <= maxVal; j += r {
			cand := dp[j/r]
			if cand == 0 {
				continue
			}
			if cand > r {
				cand = r
			}
			if cand > dp[j] {
				old := dp[j]
				dp[j] = cand
				if freq[j] > 0 {
					cnt[old] -= freq[j]
					cnt[cand] += freq[j]
					if old == curMin && cnt[old] == 0 {
						for curMin < len(cnt) && cnt[curMin] == 0 {
							curMin++
						}
					}
				}
			}
		}
		if r >= maxPF {
			if curMin < len(cnt) {
				diff := r - curMin
				if diff < best {
					best = diff
					if best == 0 {
						break
					}
				}
			}
		}
	}
	if best == int(1<<60) {
		best = 0
	}
	return best
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(5) + 1
	m := rng.Intn(1000) + 2
	arr := make([]int, n)
	input := fmt.Sprintf("1\n%d %d\n", n, m)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(m-1) + 1
		if i > 0 {
			input += " "
		}
		input += strconv.Itoa(arr[i])
	}
	input += "\n"
	exp := solveE(n, m, arr)
	return input, exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		val, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: bad output\ninput:\n%soutput:\n%s", i+1, in, out)
			os.Exit(1)
		}
		if val != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", i+1, exp, val, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
