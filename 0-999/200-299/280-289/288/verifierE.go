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

const mod int64 = 1000000007

func isLucky(x int64) bool {
	if x == 0 {
		return false
	}
	for x > 0 {
		d := x % 10
		if d != 4 && d != 7 {
			return false
		}
		x /= 10
	}
	return true
}

func genLucky(l, r int64) []int64 {
	res := []int64{}
	for i := l; i <= r; i++ {
		if isLucky(i) {
			res = append(res, i)
		}
	}
	return res
}

func expected(l, r int64) int64 {
	arr := genLucky(l, r)
	var sum int64
	for i := 0; i+1 < len(arr); i++ {
		sum = (sum + (arr[i]%mod)*(arr[i+1]%mod)) % mod
	}
	return sum % mod
}

func runCase(bin string, l, r int64) error {
	input := fmt.Sprintf("%d\n%d\n", l, r)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	gotStr := strings.TrimSpace(out.String())
	var got int64
	if _, err := fmt.Sscan(gotStr, &got); err != nil {
		return fmt.Errorf("cannot parse output: %v", err)
	}
	exp := expected(l, r)
	if got != exp {
		return fmt.Errorf("expected %d got %d (l=%d r=%d)", exp, got, l, r)
	}
	return nil
}

func randLucky(rng *rand.Rand) int64 {
	digits := rng.Intn(5) + 1
	var sb strings.Builder
	for i := 0; i < digits; i++ {
		if rng.Intn(2) == 0 {
			sb.WriteByte('4')
		} else {
			sb.WriteByte('7')
		}
	}
	v, _ := strconv.ParseInt(sb.String(), 10, 64)
	return v
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := [][2]int64{
		{4, 7},
		{44, 77},
		{4, 4},
		{7, 44},
	}
	for len(cases) < 100 {
		a := randLucky(rng)
		b := randLucky(rng)
		if a > b {
			a, b = b, a
		}
		if a == b {
			b += 4
		}
		cases = append(cases, [2]int64{a, b})
	}
	for i, tc := range cases {
		if err := runCase(bin, tc[0], tc[1]); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
