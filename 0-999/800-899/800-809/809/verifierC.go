package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const mod int64 = 1000000007

var A, B, limitK int64
var memo [31][2][2][2]struct {
	cnt int64
	sum int64
}
var used [31][2][2][2]bool

func dfs(pos int, la, lb, lk bool) (int64, int64) {
	if pos < 0 {
		return 1, 0
	}
	iLa, iLb, iLk := 0, 0, 0
	if la {
		iLa = 1
	}
	if lb {
		iLb = 1
	}
	if lk {
		iLk = 1
	}
	if used[pos][iLa][iLb][iLk] {
		res := memo[pos][iLa][iLb][iLk]
		return res.cnt, res.sum
	}
	bitA := (A >> pos) & 1
	bitB := (B >> pos) & 1
	bitK := (limitK >> pos) & 1
	var resCnt, resSum int64
	for da := int64(0); da <= 1; da++ {
		if la && da > bitA {
			continue
		}
		nLa := la && da == bitA
		for db := int64(0); db <= 1; db++ {
			if lb && db > bitB {
				continue
			}
			nLb := lb && db == bitB
			xr := da ^ db
			if lk && xr > bitK {
				continue
			}
			nLk := lk && xr == bitK
			c, s := dfs(pos-1, nLa, nLb, nLk)
			resCnt += c
			resSum = (resSum + s + (c%mod)*((xr<<pos)%mod)) % mod
		}
	}
	used[pos][iLa][iLb][iLk] = true
	memo[pos][iLa][iLb][iLk] = struct{ cnt, sum int64 }{resCnt, resSum}
	return resCnt, resSum
}

func calc(a, b, k int64) int64 {
	if a < 0 || b < 0 || k <= 0 {
		return 0
	}
	A, B = a, b
	limitK = k - 1
	for i := range used {
		for j := range used[i] {
			for l := range used[i][j] {
				used[i][j][l][0] = false
				used[i][j][l][1] = false
			}
		}
	}
	cnt, sum := dfs(30, true, true, true)
	return (sum + cnt) % mod
}

func query(x1, y1, x2, y2, k int64) int64 {
	res := calc(x2-1, y2-1, k)
	res -= calc(x1-2, y2-1, k)
	res -= calc(x2-1, y1-2, k)
	res += calc(x1-2, y1-2, k)
	res %= mod
	if res < 0 {
		res += mod
	}
	return res
}

func generateCaseC(rng *rand.Rand) string {
	q := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for i := 0; i < q; i++ {
		x1 := rng.Intn(20) + 1
		x2 := rng.Intn(20-x1+1) + x1
		y1 := rng.Intn(20) + 1
		y2 := rng.Intn(20-y1+1) + y1
		k := rng.Intn(30) + 1
		sb.WriteString(fmt.Sprintf("%d %d %d %d %d\n", x1, y1, x2, y2, k))
	}
	return sb.String()
}

func runCase(bin string, input string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	// compute expected
	rdr := strings.NewReader(input)
	var q int
	fmt.Fscan(rdr, &q)
	expOut := make([]string, q)
	for i := 0; i < q; i++ {
		var x1, y1, x2, y2, k int64
		fmt.Fscan(rdr, &x1, &y1, &x2, &y2, &k)
		expOut[i] = fmt.Sprint(query(x1, y1, x2, y2, k))
	}
	gotLines := strings.Fields(strings.TrimSpace(out.String()))
	if len(gotLines) != q {
		return fmt.Errorf("expected %d lines got %d", q, len(gotLines))
	}
	for i := 0; i < q; i++ {
		if gotLines[i] != expOut[i] {
			return fmt.Errorf("case line %d: expected %s got %s", i+1, expOut[i], gotLines[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCaseC(rng)
		if err := runCase(bin, input); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
