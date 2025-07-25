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

const MOD int64 = 1000000007
const MAXN = 200

var fact [MAXN + 1]int64
var invfact [MAXN + 1]int64
var pow25 [MAXN + 1]int64

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

func initPrecalc() {
	fact[0] = 1
	for i := 1; i <= MAXN; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invfact[MAXN] = modPow(fact[MAXN], MOD-2)
	for i := MAXN; i >= 1; i-- {
		invfact[i-1] = invfact[i] * int64(i) % MOD
	}
	pow25[0] = 1
	for i := 1; i <= MAXN; i++ {
		pow25[i] = pow25[i-1] * 25 % MOD
	}
}

func C(n, k int) int64 {
	if k < 0 || k > n {
		return 0
	}
	return fact[n] * invfact[k] % MOD * invfact[n-k] % MOD
}

type event struct {
	t   int
	str string
	n   int
}

type testCase struct {
	start  string
	events []event
}

func genCase(rng *rand.Rand) testCase {
	m := rng.Intn(5) + 5
	startLen := rng.Intn(3) + 1
	b := make([]byte, startLen)
	for i := 0; i < startLen; i++ {
		b[i] = byte('a' + rng.Intn(3))
	}
	s := string(b)
	evs := make([]event, 0, m)
	for i := 0; i < m; i++ {
		if rng.Intn(2) == 0 {
			l := rng.Intn(3) + 1
			bb := make([]byte, l)
			for j := 0; j < l; j++ {
				bb[j] = byte('a' + rng.Intn(3))
			}
			evs = append(evs, event{t: 1, str: string(bb)})
			s = string(bb)
		} else {
			n := len(s) + rng.Intn(4)
			evs = append(evs, event{t: 2, n: n})
		}
	}
	return testCase{start: s, events: evs}
}

func solve(tc testCase) string {
	cur := tc.start
	curLen := len(cur)
	dp := map[int][]int64{}
	arrInit := make([]int64, curLen+1)
	arrInit[curLen] = 1
	dp[curLen] = arrInit
	var out strings.Builder
	for _, ev := range tc.events {
		if ev.t == 1 {
			cur = ev.str
			curLen = len(cur)
			if _, ok := dp[curLen]; !ok {
				arr := make([]int64, curLen+1)
				arr[curLen] = 1
				dp[curLen] = arr
			}
		} else {
			n := ev.n
			if n < curLen {
				out.WriteString("0\n")
				continue
			}
			arr := dp[curLen]
			for len(arr) <= n {
				i := len(arr)
				val := (arr[i-1]*26 + C(i-1, curLen-1)*pow25[i-curLen]) % MOD
				arr = append(arr, val)
			}
			dp[curLen] = arr
			out.WriteString(fmt.Sprintf("%d\n", arr[n]%MOD))
		}
	}
	return strings.TrimSpace(out.String())
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.events)))
	sb.WriteString(tc.start + "\n")
	for _, ev := range tc.events {
		if ev.t == 1 {
			sb.WriteString(fmt.Sprintf("1 %s\n", ev.str))
		} else {
			sb.WriteString(fmt.Sprintf("2 %d\n", ev.n))
		}
	}
	return sb.String()
}

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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append([]string{os.Args[0]}, os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	initPrecalc()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		input := buildInput(tc)
		exp := solve(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%s\nexp:\n%s\n---\ngot:\n%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
