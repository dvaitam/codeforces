package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

const MOD int64 = 1000000007

func modPow(a, e int64) int64 {
	res := int64(1)
	a %= MOD
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func modInv(a int64) int64 { return modPow((a%MOD+MOD)%MOD, MOD-2) }

// genPartitions returns all sorted (desc) partitions of n into >= 2 positive parts.
func genPartitions(n int) [][]int {
	var result [][]int
	var bt func(rem, maxP int, cur []int)
	bt = func(rem, maxP int, cur []int) {
		if rem == 0 {
			if len(cur) >= 2 {
				s := make([]int, len(cur))
				copy(s, cur)
				result = append(result, s)
			}
			return
		}
		for k := maxP; k >= 1; k-- {
			if k > rem {
				continue
			}
			bt(rem-k, k, append(cur, k))
		}
	}
	bt(n, n, nil)
	return result
}

func sliceKey(s []int) string {
	parts := make([]string, len(s))
	for i, v := range s {
		parts[i] = strconv.Itoa(v)
	}
	return strings.Join(parts, ",")
}

// oracle computes the expected time mod MOD using Gaussian elimination over Q.
// Works correctly for small N (tested up to N~10).
func oracle(a []int) int64 {
	N := 0
	for _, x := range a {
		N += x
	}
	// canonical sorted form
	sorted := make([]int, len(a))
	copy(sorted, a)
	sort.Sort(sort.Reverse(sort.IntSlice(sorted)))
	// remove zeros
	end := len(sorted)
	for end > 0 && sorted[end-1] == 0 {
		end--
	}
	sorted = sorted[:end]

	if len(sorted) <= 1 || N <= 1 {
		return 0
	}

	states := genPartitions(N)
	m := len(states)

	stateIdx := make(map[string]int, m)
	for i, s := range states {
		stateIdx[sliceKey(s)] = i
	}

	initIdx, ok := stateIdx[sliceKey(sorted)]
	if !ok {
		return 0
	}

	denom := new(big.Int).SetInt64(int64(N) * int64(N-1))

	// Build (I - Q) * E = 1
	A := make([][]*big.Rat, m)
	b := make([]*big.Rat, m)
	for i := range A {
		A[i] = make([]*big.Rat, m)
		for j := range A[i] {
			A[i][j] = new(big.Rat)
		}
		A[i][i].SetInt64(1)
		b[i] = new(big.Rat).SetInt64(1)
	}

	for i, s := range states {
		k := len(s)

		// Same-color self-loops
		for ci := 0; ci < k; ci++ {
			if s[ci] < 2 {
				continue
			}
			num := big.NewInt(int64(s[ci]) * int64(s[ci]-1))
			prob := new(big.Rat).SetFrac(num, denom)
			A[i][i].Sub(A[i][i], prob)
		}

		// Cross-color transitions
		for ci := 0; ci < k; ci++ {
			for cj := 0; cj < k; cj++ {
				if ci == cj {
					continue
				}
				num := big.NewInt(int64(s[ci]) * int64(s[cj]))
				prob := new(big.Rat).SetFrac(num, denom)

				ns := make([]int, k)
				copy(ns, s)
				ns[ci]++
				ns[cj]--

				var ns2 []int
				for _, v := range ns {
					if v > 0 {
						ns2 = append(ns2, v)
					}
				}
				sort.Sort(sort.Reverse(sort.IntSlice(ns2)))

				nk := sliceKey(ns2)
				if j, ok2 := stateIdx[nk]; ok2 {
					A[i][j].Sub(A[i][j], prob)
				}
				// absorbing → E=0, no contribution
			}
		}
	}

	// Gauss-Jordan elimination
	for col := 0; col < m; col++ {
		pivot := -1
		for row := col; row < m; row++ {
			if A[row][col].Sign() != 0 {
				pivot = row
				break
			}
		}
		if pivot < 0 {
			continue
		}
		A[col], A[pivot] = A[pivot], A[col]
		b[col], b[pivot] = b[pivot], b[col]

		pivInv := new(big.Rat).Inv(A[col][col])
		A[col][col].SetInt64(1)
		for c2 := col + 1; c2 < m; c2++ {
			A[col][c2].Mul(A[col][c2], pivInv)
		}
		b[col].Mul(b[col], pivInv)

		for row := 0; row < m; row++ {
			if row == col || A[row][col].Sign() == 0 {
				continue
			}
			factor := new(big.Rat).Set(A[row][col])
			A[row][col].SetInt64(0)
			for c2 := col + 1; c2 < m; c2++ {
				A[row][c2].Sub(A[row][c2], new(big.Rat).Mul(factor, A[col][c2]))
			}
			b[row].Sub(b[row], new(big.Rat).Mul(factor, b[col]))
		}
	}

	ans := b[initIdx]
	p := new(big.Int).Mod(ans.Num(), big.NewInt(MOD))
	if p.Sign() < 0 {
		p.Add(p, big.NewInt(MOD))
	}
	q := new(big.Int).Mod(ans.Denom(), big.NewInt(MOD))
	if q.Sign() < 0 {
		q.Add(q, big.NewInt(MOD))
	}
	return p.Int64() * modInv(q.Int64()) % MOD
}

func runExe(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genTest(rng *rand.Rand) (string, []int) {
	// Keep N small (<=8) so oracle stays fast
	n := rng.Intn(4) + 1
	a := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(2) + 1
	}
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(n) + "\n")
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String(), a
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Known samples from the problem statement
	type known struct {
		input string
		ans   int64
	}
	knownTests := []known{
		{"1\n1\n", 0},
		{"2\n1 1\n", 1},
		// 83/4 mod M: 83 * modInv(4) = 83 * 250000002 mod M = 750000026
		{"3\n3 2 1\n", 83 * modInv(4) % MOD},
	}
	for i, kt := range knownTests {
		got, err := runExe(bin, kt.input)
		if err != nil {
			fmt.Printf("known test %d: runtime error: %v\ninput: %s", i+1, err, kt.input)
			os.Exit(1)
		}
		gotVal, _ := strconv.ParseInt(got, 10, 64)
		if gotVal != kt.ans {
			fmt.Printf("known test %d failed\nInput: %sExpected: %d\nGot: %d\n", i+1, kt.input, kt.ans, gotVal)
			os.Exit(1)
		}
	}

	// Random tests vs brute-force oracle
	for i := 0; i < 100; i++ {
		input, a := genTest(rng)
		exp := oracle(a)
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		gotVal, _ := strconv.ParseInt(got, 10, 64)
		if gotVal != exp {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%d\nGot:\n%d\n", i+1, input, exp, gotVal)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
