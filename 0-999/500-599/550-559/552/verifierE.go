package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func evalExpr(digits []*big.Int, ops []byte) *big.Int {
	if len(digits) == 0 {
		return big.NewInt(0)
	}
	cur := new(big.Int).Set(digits[0])
	res := big.NewInt(0)
	for i, op := range ops {
		if op == '+' {
			res.Add(res, cur)
			cur = new(big.Int).Set(digits[i+1])
		} else {
			cur.Mul(cur, digits[i+1])
		}
	}
	res.Add(res, cur)
	return res
}

func evalWithParen(digits []*big.Int, ops []byte, l, r int) *big.Int {
	mid := evalExpr(digits[l:r+1], ops[l:r])
	newDigits := make([]*big.Int, 0, len(digits)-(r-l))
	newDigits = append(newDigits, digits[:l]...)
	newDigits = append(newDigits, mid)
	if r+1 < len(digits) {
		newDigits = append(newDigits, digits[r+1:]...)
	}
	newOps := make([]byte, 0, len(ops)-(r-l))
	newOps = append(newOps, ops[:l]...)
	if r < len(ops) {
		newOps = append(newOps, ops[r])
		if r+1 < len(ops) {
			newOps = append(newOps, ops[r+1:]...)
		}
	}
	return evalExpr(newDigits, newOps)
}

func expected(expr string) string {
	n := (len(expr) + 1) / 2
	digits := make([]*big.Int, n)
	ops := make([]byte, 0, n-1)
	for i := 0; i < len(expr); i++ {
		if i%2 == 0 {
			digits[i/2] = big.NewInt(int64(expr[i] - '0'))
		} else {
			ops = append(ops, expr[i])
		}
	}
	starts := []int{0}
	ends := []int{n - 1}
	for i, op := range ops {
		if op == '*' {
			starts = append(starts, i+1)
			ends = append(ends, i)
		}
	}
	maxVal := big.NewInt(0)
	for _, l := range starts {
		for _, r := range ends {
			if l > r {
				continue
			}
			val := evalWithParen(digits, ops, l, r)
			if val.Cmp(maxVal) > 0 {
				maxVal = val
			}
		}
	}
	return maxVal.String()
}

func genCase(rng *rand.Rand) (string, string) {
	digitsCount := rng.Intn(9) + 1
	var sb strings.Builder
	starCnt := 0
	for i := 0; i < digitsCount; i++ {
		if i > 0 {
			var op byte
			if starCnt >= 15 {
				op = '+'
			} else {
				if rng.Intn(2) == 0 {
					op = '+'
				} else {
					op = '*'
					starCnt++
				}
			}
			sb.WriteByte(op)
		}
		sb.WriteByte(byte(rng.Intn(9) + '1'))
	}
	expr := sb.String()
	return expr + "\n", expected(expr)
}

func runCase(bin, input, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
