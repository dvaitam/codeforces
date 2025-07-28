package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const SQRTN = 1000

var dp [SQRTN][SQRTN]int64
var binom [SQRTN][SQRTN]int64

func h(x int) int64 { return binom[x+3][3] }

func g(num int, siz int, mod int64) int64 {
	if num == 0 {
		return 1
	}
	if siz == SQRTN {
		return 0
	}
	if dp[num][siz] != -1 {
		return dp[num][siz]
	}
	var res int64
	for i := 0; i*siz <= num; i++ {
		res = (res + g(num-i*siz, siz+1, mod)*h(i)) % mod
	}
	dp[num][siz] = res
	return res
}

func f(x, y, n int64, mod int64) int64 {
	num := int(x*y - n)
	return g(num, 1, mod)
}

func expected(n int64) (int64, int64) {
	low, up := int64(0), int64(1500)
	for up-low > 1 {
		mid := (up + low) / 2
		a := mid / 2
		b := mid - a
		if a*b >= n {
			up = mid
		} else {
			low = mid
		}
	}
	a := up / 2
	b := up - a
	return a, b
}

func expectedWays(n, m int64) int64 {
	a, b := expected(n)
	var ans int64
	for i := a; (a+b)-i >= 1 && i*((a+b)-i) >= n; i++ {
		ans = (ans + f(i, (a+b)-i, n, m)) % m
	}
	for i := a - 1; i >= 1 && i*((a+b)-i) >= n; i-- {
		ans = (ans + f(i, (a+b)-i, n, m)) % m
	}
	return ans
}

func generateTests() []struct {
	n int64
	u int
	m int64
} {
	rand.Seed(5)
	t := 100
	res := make([]struct {
		n int64
		u int
		m int64
	}, t)
	for i := 0; i < t; i++ {
		u := 1
		if rand.Intn(2) == 0 {
			u = 1
		} else {
			u = 2
		}
		n := int64(rand.Intn(20) + 1)
		m := int64(1000000007)
		res[i] = struct {
			n int64
			u int
			m int64
		}{n, u, m}
	}
	return res
}

func verifyCase(bin string, tc struct {
	n int64
	u int
	m int64
}) error {
	var input strings.Builder
	input.WriteString(fmt.Sprintf("1 %d\n", tc.u))
	if tc.u == 2 {
		input.WriteString(fmt.Sprintf("%d\n", tc.m))
	}
	input.WriteString(fmt.Sprintf("%d\n", tc.n))

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("execution error: %v", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(strings.TrimSpace(string(out))))
	scanner.Split(bufio.ScanWords)

	a, b := expected(tc.n)
	perim := (a + b) * 2

	if tc.u == 1 {
		if !scanner.Scan() {
			return fmt.Errorf("no output")
		}
		ga, _ := strconv.ParseInt(scanner.Text(), 10, 64)
		if !scanner.Scan() {
			return fmt.Errorf("missing width")
		}
		gb, _ := strconv.ParseInt(scanner.Text(), 10, 64)
		if ga != a || gb != b {
			return fmt.Errorf("expected dims %d %d got %d %d", a, b, ga, gb)
		}
		rows := make([]string, ga)
		for i := int64(0); i < ga; i++ {
			if !scanner.Scan() {
				return fmt.Errorf("missing row")
			}
			rows[i] = scanner.Text()
		}
		countHash := 0
		for _, r := range rows {
			countHash += strings.Count(r, "#")
		}
		if countHash != int(tc.n) {
			return fmt.Errorf("expected %d # got %d", tc.n, countHash)
		}
	} else {
		if !scanner.Scan() {
			return fmt.Errorf("no output")
		}
		gp, _ := strconv.ParseInt(scanner.Text(), 10, 64)
		if gp != perim {
			return fmt.Errorf("expected perimeter %d got %d", perim, gp)
		}
		if !scanner.Scan() {
			return fmt.Errorf("missing ways")
		}
		gw, _ := strconv.ParseInt(scanner.Text(), 10, 64)
		// compute expected ways
		for i := 0; i < SQRTN; i++ {
			for j := 0; j < SQRTN; j++ {
				dp[i][j] = -1
			}
		}
		for i := 0; i < SQRTN; i++ {
			for j := 0; j <= i; j++ {
				if j == 0 || j == i {
					binom[i][j] = 1
				} else {
					binom[i][j] = (binom[i-1][j] + binom[i-1][j-1]) % tc.m
				}
			}
		}
		expWays := expectedWays(tc.n, tc.m)
		if gw != expWays {
			return fmt.Errorf("expected ways %d got %d", expWays, gw)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		if err := verifyCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
