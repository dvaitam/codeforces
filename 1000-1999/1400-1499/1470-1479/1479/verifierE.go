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

func clone(x []int) []int { return append([]int(nil), x...) }

func removeOne(c []int, idx int) []int {
	res := clone(c)
	res[idx]--
	if res[idx] == 0 {
		return append(res[:idx], res[idx+1:]...)
	}
	return res
}

func removeOneAddNew(c []int, idx int) []int {
	res := removeOne(c, idx)
	res = append(res, 1)
	sort.Ints(res)
	return res
}

func removeOneAddToOther(c []int, i, j int) []int {
	res := clone(c)
	res[i]--
	if res[i] == 0 {
		res = append(res[:i], res[i+1:]...)
		if j > i {
			j--
		}
	}
	res[j]++
	sort.Ints(res)
	return res
}

var memo map[string]*big.Rat
var half = big.NewRat(1, 2)

func expected(c []int) *big.Rat {
	sort.Ints(c)
	if len(c) == 1 {
		return big.NewRat(0, 1)
	}
	key := fmt.Sprint(c)
	if v, ok := memo[key]; ok {
		return new(big.Rat).Set(v)
	}
	n := 0
	for _, v := range c {
		n += v
	}
	res := big.NewRat(1, 1)
	for i, cnt := range c {
		pSel := new(big.Rat).SetFrac64(int64(cnt), int64(n))
		// create new club
		st1 := removeOneAddNew(c, i)
		t1 := expected(st1)
		term1 := new(big.Rat).Mul(pSel, half)
		term1.Mul(term1, t1)
		sum.Add(sum, term1)
		// join existing club
		for j, cntj := range c {
			pJoin := new(big.Rat).SetFrac64(int64(cntj), int64(n))
			pJoin.Mul(pJoin, pSel)
			pJoin.Mul(pJoin, half)
			st2 := removeOneAddToOther(c, i, j)
			t2 := expected(st2)
			term2 := new(big.Rat).Mul(pJoin, t2)
			sum.Add(sum, term2)
		}
	}
	res.Add(res, sum)
	memo[key] = res
	return new(big.Rat).Set(res)
}

func modResult(r *big.Rat) int {
	mod := big.NewInt(998244353)
	num := new(big.Int).Mod(r.Num(), mod)
	den := new(big.Int).Mod(r.Denom(), mod)
	inv := new(big.Int).Exp(den, new(big.Int).Sub(mod, big.NewInt(2)), mod)
	num.Mul(num, inv)
	num.Mod(num, mod)
	return int(num.Int64())
}

func genTestsE() [][]int {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([][]int, 100)
	for i := range tests {
		for {
			m := rng.Intn(3) + 1
			counts := make([]int, m)
			total := 0
			for j := range counts {
				counts[j] = rng.Intn(3) + 1
				total += counts[j]
			}
			if total <= 5 {
				tests[i] = counts
				break
			}
		}
	}
	return tests
}

func buildInputE(c []int) string {
	var sb strings.Builder
	fmt.Fprintln(&sb, len(c))
	for i, v := range c {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprint(&sb, v)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(bin string, c []int, expect int) error {
	input := buildInputE(c)
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) == 0 {
		return fmt.Errorf("no output")
	}
	val, err := strconv.Atoi(fields[0])
	if err != nil || val != expect {
		return fmt.Errorf("expected %d got %s", expect, fields[0])
	}
	if len(fields) > 1 {
		return fmt.Errorf("extra output")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genTestsE()
	for i, c := range cases {
		memo = make(map[string]*big.Rat)
		exp := modResult(expected(c))
		if err := runCase(bin, c, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
