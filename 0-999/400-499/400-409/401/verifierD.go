package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func runCandidate(bin, input string) (string, error) {
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

func solveCase(nStr string, m int) int64 {
	digits := make([]int, len(nStr))
	for i, ch := range nStr {
		digits[i] = int(ch - '0')
	}
	sort.Ints(digits)
	type pair struct{ d, cnt int }
	var vp []pair
	for i := 0; i < len(digits); {
		j := i + 1
		for j < len(digits) && digits[j] == digits[i] {
			j++
		}
		vp = append(vp, pair{digits[i], j - i})
		i = j
	}
	L := len(vp)
	totalStates := 1
	bases := make([]int, L)
	for i := 0; i < L; i++ {
		bases[i] = totalStates
		totalStates *= vp[i].cnt + 1
	}
	totalDigits := len(digits)
	countsState := make([][]int, totalStates)
	sumState := make([]int, totalStates)
	for idx := 0; idx < totalStates; idx++ {
		cnts := make([]int, L)
		rem := idx
		sum := 0
		for i := 0; i < L; i++ {
			c := rem / bases[i] % (vp[i].cnt + 1)
			cnts[i] = c
			sum += c
			rem -= c * bases[i]
		}
		countsState[idx] = cnts
		sumState[idx] = sum
	}
	dp := make([]int64, totalStates*m)
	dp[0*m+0] = 1
	for idx := 0; idx < totalStates; idx++ {
		s := sumState[idx]
		if s >= totalDigits {
			continue
		}
		for rem0 := 0; rem0 < m; rem0++ {
			cur := dp[idx*m+rem0]
			if cur == 0 {
				continue
			}
			for i := 0; i < L; i++ {
				if countsState[idx][i] < vp[i].cnt {
					if s == 0 && vp[i].d == 0 {
						continue
					}
					newIdx := idx + bases[i]
					newRem := (rem0*10 + vp[i].d) % m
					dp[newIdx*m+newRem] += cur
				}
			}
		}
	}
	lastIdx := totalStates - 1
	return dp[lastIdx*m+0]
}

func generateCase(rng *rand.Rand) (string, string) {
	length := rng.Intn(7) + 1
	var sb strings.Builder
	for i := 0; i < length; i++ {
		d := rng.Intn(10)
		if i == 0 && d == 0 {
			d = rng.Intn(9) + 1
		}
		sb.WriteByte(byte('0' + d))
	}
	nStr := sb.String()
	m := rng.Intn(99) + 1
	input := fmt.Sprintf("%s %d\n", nStr, m)
	exp := fmt.Sprintf("%d", solveCase(nStr, m))
	return input, exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
