package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCase struct {
	input    string
	expected string
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func extgcd(a, b int64) (g, x, y int64) {
	if b == 0 {
		return a, 1, 0
	}
	g, x1, y1 := extgcd(b, a%b)
	return g, y1, x1 - (a/b)*y1
}

func modInverse(a, mod int64) int64 {
	g, x, _ := extgcd(a, mod)
	if g != 1 {
		return 0
	}
	x %= mod
	if x < 0 {
		x += mod
	}
	return x
}

func solveCase(n int, ab [][2]int64, queries [][2]int64) []int64 {
	diffs := make([]int64, n)
	var bsum int64
	for i := 0; i < n; i++ {
		diffs[i] = ab[i][0] - ab[i][1]
		bsum += ab[i][1]
	}
	sort.Slice(diffs, func(i, j int) bool { return diffs[i] > diffs[j] })
	prefix := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		prefix[i] = prefix[i-1] + diffs[i-1]
	}
	values := make([]int64, n+1)
	for i := 0; i <= n; i++ {
		values[i] = bsum + prefix[i]
	}
	S := int64(1)
	for S*S <= int64(n) {
		S++
	}
	pre := make([][]int64, S+1)
	for s := int64(1); s <= S; s++ {
		pre[s] = make([]int64, s)
		for i := 0; i < int(s); i++ {
			pre[s][i] = -1 << 63
		}
		for r := int64(0); r < s; r++ {
			maxv := int64(-1 << 63)
			for x := r; x <= int64(n); x += s {
				if v := values[x]; v > maxv {
					maxv = v
				}
			}
			pre[s][r] = maxv
		}
	}
	res := make([]int64, len(queries))
	for idx, q := range queries {
		xj, yj := q[0], q[1]
		g := gcd(xj, yj)
		if int64(n)%g != 0 {
			res[idx] = -1
			continue
		}
		lcm := xj / g * yj
		x1 := xj / g
		y1 := yj / g
		n1 := int64(n) / g
		inv := modInverse(x1%y1, y1)
		t0 := (n1 % y1) * inv % y1
		r0 := xj * t0
		if r0 > int64(n) {
			res[idx] = -1
			continue
		}
		if lcm <= S {
			val := pre[lcm][r0%lcm]
			if val == -1<<63 {
				res[idx] = -1
			} else {
				res[idx] = val
			}
		} else {
			maxv := int64(-1 << 63)
			for R := r0; R <= int64(n); R += lcm {
				if v := values[R]; v > maxv {
					maxv = v
				}
			}
			if maxv == -1<<63 {
				res[idx] = -1
			} else {
				res[idx] = maxv
			}
		}
	}
	return res
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(5))
	tests := []testCase{}
	// small fixed example
	tests = append(tests, testCase{
		input:    "2\n5 4\n3 6\n2\n1 1\n2 2\n",
		expected: "11\n11"})
	for len(tests) < 100 {
		n := rng.Intn(5) + 1
		ab := make([][2]int64, n)
		for i := 0; i < n; i++ {
			ab[i][0] = int64(rng.Intn(10) + 1)
			ab[i][1] = int64(rng.Intn(10) + 1)
		}
		m := rng.Intn(3) + 1
		qs := make([][2]int64, m)
		for j := 0; j < m; j++ {
			qs[j][0] = int64(rng.Intn(n) + 1)
			qs[j][1] = int64(rng.Intn(n) + 1)
		}
		res := solveCase(n, ab, qs)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			sb.WriteString(fmt.Sprintf("%d %d\n", ab[i][0], ab[i][1]))
		}
		sb.WriteString(fmt.Sprintf("%d\n", m))
		for i, q := range qs {
			sb.WriteString(fmt.Sprintf("%d %d\n", q[0], q[1]))
			res[i] = res[i]
		}
		expSb := strings.Builder{}
		for i, v := range res {
			if i > 0 {
				expSb.WriteByte('\n')
			}
			expSb.WriteString(fmt.Sprint(v))
		}
		tests = append(tests, testCase{input: sb.String(), expected: expSb.String()})
	}
	return tests
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
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	_ = rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := generateTests()
	for i, tc := range tests {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		expFields := strings.Fields(tc.expected)
		gotFields := strings.Fields(got)
		if len(expFields) != len(gotFields) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc.expected, got, tc.input)
			os.Exit(1)
		}
		ok := true
		for idx := range expFields {
			if expFields[idx] != gotFields[idx] {
				ok = false
				break
			}
		}
		if !ok {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc.expected, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
