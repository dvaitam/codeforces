package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type test struct{ input, expected string }

func solveCase(n int, a [][4]int64, b [4]int64) string {
	var sumA int64
	suitSums := make([]int64, 4)
	for i := 0; i < n; i++ {
		for j := 0; j < 4; j++ {
			sumA += a[i][j]
			suitSums[j] += a[i][j]
		}
	}
	var sumB int64
	for j := 0; j < 4; j++ {
		sumB += b[j]
		suitSums[j] += b[j]
	}
	total := sumA + sumB
	T := total / int64(n)

	maxSuit := make([]int64, n)
	for i := 0; i < n; i++ {
		mx := a[i][0]
		for j := 1; j < 4; j++ {
			if a[i][j] > mx {
				mx = a[i][j]
			}
		}
		maxSuit[i] = mx
	}
	prefix := make([]int64, n)
	if n > 0 {
		prefix[0] = maxSuit[0]
		for i := 1; i < n; i++ {
			if maxSuit[i] > prefix[i-1] {
				prefix[i] = maxSuit[i]
			} else {
				prefix[i] = prefix[i-1]
			}
		}
	}
	suffix := make([]int64, n)
	if n > 0 {
		suffix[n-1] = maxSuit[n-1]
		for i := n - 2; i >= 0; i-- {
			if maxSuit[i] > suffix[i+1] {
				suffix[i] = maxSuit[i]
			} else {
				suffix[i] = suffix[i+1]
			}
		}
	}
	results := make([]int64, n)
	for i := 0; i < n; i++ {
		baseOthers := int64(0)
		if i > 0 {
			baseOthers = prefix[i-1]
		}
		if i+1 < n && suffix[i+1] > baseOthers {
			baseOthers = suffix[i+1]
		}
		d := T
		for j := 0; j < 4; j++ {
			d -= a[i][j]
		}
		Sother := make([]int64, 4)
		for j := 0; j < 4; j++ {
			Sother[j] = suitSums[j] - a[i][j]
		}
		bestDiff := int64(0)
		for j0 := 0; j0 < 4; j0++ {
			Btemp := make([]int64, 4)
			copy(Btemp, b[:])
			Stemp := make([]int64, 4)
			copy(Stemp, Sother)
			x := make([]int64, 4)
			take := Btemp[j0]
			if take > d {
				take = d
			}
			x[j0] = take
			Btemp[j0] -= take
			Stemp[j0] -= take
			remaining := d - take
			for remaining > 0 {
				best := -1
				for s := 0; s < 4; s++ {
					if Btemp[s] > 0 {
						if best == -1 || Stemp[s] > Stemp[best] {
							best = s
						}
					}
				}
				if best == -1 {
					break
				}
				give := Btemp[best]
				if give > remaining {
					give = remaining
				}
				x[best] += give
				Btemp[best] -= give
				Stemp[best] -= give
				remaining -= give
			}
			curMax := a[i][0] + x[0]
			for s := 1; s < 4; s++ {
				if a[i][s]+x[s] > curMax {
					curMax = a[i][s] + x[s]
				}
			}
			y := baseOthers
			for s := 0; s < 4; s++ {
				val := (Stemp[s] + int64(n-1) - 1) / int64(n-1)
				if val > y {
					y = val
				}
			}
			diff := curMax - y
			if diff > bestDiff {
				bestDiff = diff
			}
		}
		if bestDiff < 0 {
			bestDiff = 0
		}
		results[i] = bestDiff
	}
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", results[i]))
	}
	return sb.String()
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(47))
	var tests []test
	for len(tests) < 100 {
		n := rng.Intn(3) + 1
		a := make([][4]int64, n)
		var sumA int64
		for i := 0; i < n; i++ {
			for j := 0; j < 4; j++ {
				val := int64(rng.Intn(5))
				a[i][j] = val
				sumA += val
			}
		}
		T := int64(rng.Intn(5) + 1)
		if T*int64(n) < sumA {
			T = (sumA + int64(n) - 1) / int64(n)
		}
		sumB := T*int64(n) - sumA
		b := [4]int64{}
		for j := 0; j < 4; j++ {
			if j == 3 {
				b[j] = sumB
			} else {
				val := int64(rng.Intn(int(sumB + 1)))
				b[j] = val
				sumB -= val
			}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			for j := 0; j < 4; j++ {
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(fmt.Sprintf("%d", a[i][j]))
			}
			sb.WriteByte('\n')
		}
		for j := 0; j < 4; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", b[j]))
		}
		sb.WriteByte('\n')
		tests = append(tests, test{sb.String(), solveCase(n, a, b)})
	}
	return tests
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
