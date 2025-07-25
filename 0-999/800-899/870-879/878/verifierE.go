package main

import (
	"bytes"
	"fmt"
	"math/big"
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return out.String() + stderr.String(), err
	}
	return out.String(), nil
}

type Query struct{ l, r int }

type Test struct {
	n       int
	q       int
	arr     []int64
	queries []Query
	input   string
}

func genTest(rng *rand.Rand) Test {
	n := rng.Intn(6) + 1
	q := rng.Intn(4) + 1
	arr := make([]int64, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	for i := 0; i < n; i++ {
		arr[i] = int64(rng.Intn(11) - 5)
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(arr[i], 10))
	}
	sb.WriteByte('\n')
	queries := make([]Query, q)
	for i := 0; i < q; i++ {
		l := rng.Intn(n) + 1
		r := rng.Intn(n-l+1) + l
		queries[i] = Query{l: l, r: r}
		sb.WriteString(fmt.Sprintf("%d %d\n", l, r))
	}
	return Test{n: n, q: q, arr: arr, queries: queries, input: sb.String()}
}

const mod = 1000000007

func maxValue(a []int64, l, r int, memo map[[2]int]*big.Int) *big.Int {
	if l == r {
		return big.NewInt(a[l-1])
	}
	key := [2]int{l, r}
	if v, ok := memo[key]; ok {
		return new(big.Int).Set(v)
	}
	best := big.NewInt(-1 << 60)
	for k := l; k < r; k++ {
		left := maxValue(a, l, k, memo)
		right := maxValue(a, k+1, r, memo)
		cand := new(big.Int).Lsh(right, 1)
		cand.Add(cand, left)
		if cand.Cmp(best) > 0 {
			best = cand
		}
	}
	memo[key] = new(big.Int).Set(best)
	return best
}

func solve(t Test) string {
	var sb strings.Builder
	for _, qu := range t.queries {
		memo := make(map[[2]int]*big.Int)
		val := maxValue(t.arr, qu.l, qu.r, memo)
		val.Mod(val, big.NewInt(mod))
		if val.Sign() < 0 {
			val.Add(val, big.NewInt(mod))
		}
		sb.WriteString(val.String())
		sb.WriteByte('\n')
	}
	return strings.TrimSpace(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		t := genTest(rng)
		expected := solve(t)
		out, err := run(bin, t.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\noutput:\n%s", i+1, err, out)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if out != expected {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", i+1, t.input, expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("ok 100 tests")
}
