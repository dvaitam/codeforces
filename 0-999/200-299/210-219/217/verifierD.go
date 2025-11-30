package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const mod = 1000000007

const rawTestcasesData = `
3 3 2 10 8
2 1 3 5 9 4
3 8 4 9 8 7 11
3 4 5 3 9 7 12 1
2 3 4 1 5 1 5
5 7 5 7 7 12 10 8
3 6 0
2 3 3 4 5 11
5 5 3 9 7 10
4 7 4 4 6 11 1
4 3 5 6 9 10 10 2
3 5 2 2 2
5 8 0
4 2 3 3 1 5
5 7 0
2 1 3 12 10 6
4 4 0
4 1 0
2 1 1 7
4 5 1 12
2 6 2 6 3
5 7 3 9 7 11
2 5 3 11 12 12
3 5 3 5 9 5
4 1 3 10 6 1
5 3 0
4 8 2 11 6
4 8 0
2 1 2 5 11
5 5 4 10 6 3 6
3 6 2 10 5
4 7 0
2 3 2 9 4
4 4 2 3 11
5 2 0
4 6 5 4 8 3 2 6
3 8 2 4 2
2 4 2 10 3
4 6 5 2 10 6 10 3
5 5 4 5 8 6 11
5 5 3 10 7 1
5 3 1 1
5 7 4 12 4 1 12
5 5 4 6 4 2 10
4 2 1 1
2 4 3 10 1 1
5 2 1 9
4 4 5 1 9 9 7 1
2 6 1 5
5 1 2 4 4
2 2 1 4
4 3 0
5 7 0
4 4 2 10 9
5 1 3 6 1 1
3 1 0
2 2 3 1 12 2
5 6 1 6
2 6 3 11 7 10
4 6 2 4 6
5 2 1 9
2 7 0
3 1 2 8 10
5 1 4 7 1 6 11
5 6 3 12 7 8
2 4 1 9
4 2 3 4 7 3
2 6 2 9 5
2 8 5 2 12 11 9 7
2 6 4 9 2 10 12
2 8 1 4
5 1 4 2 10 2 11
5 3 0
4 2 0
2 8 5 5 10 5 2 1
3 2 4 12 2 9 1
4 3 0
3 3 5 4 8 10 12 7
4 6 4 7 6 9 7
2 7 4 4 7 12 3
5 8 1 11
5 3 1 2
5 8 5 9 8 10 12 3
3 5 1 3
4 4 5 9 5 11 12 7
4 4 2 1 5
5 7 1 3
4 4 2 8 3
5 8 5 10 4 8 10 11
2 8 5 2 7 12 1 8
3 4 5 12 11 2 4 5
3 4 2 3 3
2 5 1 1
4 3 3 2 12 2
2 2 2 5 1
4 8 4 12 11 6 1
2 6 2 7 7
5 2 1 11
5 7 1 9
4 2 2 2 11
`

var rawTestcases = strings.Split(strings.TrimSpace(rawTestcasesData), "\n")

type orbit struct {
	r, cr, cmr int
}

func solve217D(n, m, t int, values []int) (int, error) {
	// Logic directly mirrored from 0-999/200-299/210-219/217/217D.go so the verifier
	// computes expected answers without depending on the standalone binary.
	if len(values) != t {
		return 0, fmt.Errorf("expected %d residues got %d", t, len(values))
	}
	_ = n // kept to match the original signature from the solution.

	cnt := make([]int, m)
	for _, x := range values {
		r := x % m
		cnt[r]++
	}

	orbits := make([]orbit, 0, m/2+1)
	for r := 1; r*2 < m; r++ {
		cr := cnt[r]
		cmr := cnt[m-r]
		if cr+cmr > 0 {
			orbits = append(orbits, orbit{r, cr, cmr})
		}
	}
	if m%2 == 0 {
		r := m / 2
		if cnt[r] > 0 {
			orbits = append(orbits, orbit{r, cnt[r], 0})
		}
	}

	start := make([]byte, m)
	for i := range start {
		start[i] = '0'
	}
	start[0] = '1'
	dp := map[string]int{
		string(start): 1,
	}

	for _, o := range orbits {
		dp2 := make(map[string]int)
		r := o.r
		tot := o.cr + o.cmr
		for maskStr, w := range dp {
			dp2[maskStr] = (dp2[maskStr] + w) % mod
			if tot == 0 {
				continue
			}
			mask := []byte(maskStr)
			if mask[r] == '1' || mask[(m-r)%m] == '1' {
				continue
			}
			newMask := make([]byte, m)
			copy(newMask, mask)
			for j := 0; j < m; j++ {
				if mask[j] == '1' {
					newMask[(j+r)%m] = '1'
					newMask[(j+m-r)%m] = '1'
				}
			}
			nmStr := string(newMask)
			dp2[nmStr] = (dp2[nmStr] + int((int64(w)*int64(tot))%mod)) % mod
		}
		dp = dp2
	}

	ans := 0
	for _, v := range dp {
		ans = (ans + v) % mod
	}
	return ans, nil
}

func expectedFromLine(line string) (string, error) {
	fields := strings.Fields(line)
	if len(fields) < 3 {
		return "", fmt.Errorf("invalid input")
	}
	nums := make([]int, len(fields))
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return "", fmt.Errorf("bad integer %q: %w", f, err)
		}
		nums[i] = v
	}
	n, m, t := nums[0], nums[1], nums[2]
	if len(nums) != 3+t {
		return "", fmt.Errorf("expected %d residues got %d", t, len(nums)-3)
	}
	ans, err := solve217D(n, m, t, nums[3:])
	if err != nil {
		return "", err
	}
	return strconv.Itoa(ans), nil
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v, output: %s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for idx, line := range rawTestcases {
		expect, err := expectedFromLine(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d invalid: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := run(bin, line+"\n")
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(rawTestcases))
}
