package main

import (
	"bytes"
	"context"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// ---- Embedded correct solver for 2056 F2 ----

type Pre struct {
	d          uint
	P          uint64
	fullParity uint8
	fullXor    uint64
	prefParity []uint8
	prefXor    []uint64
}

func xorUpto(n uint64) uint64 {
	switch n & 3 {
	case 0:
		return n
	case 1:
		return 1
	case 2:
		return n + 1
	default:
		return 0
	}
}

func build(p int) *Pre {
	d := bits.Len(uint(p - 1))
	P := 1 << d
	a := make([]uint8, P)
	for k := 0; k < p; k++ {
		if ((p - k - 1) & (k >> 1)) == 0 {
			a[k] = 1
		}
	}
	for bit := 1; bit < P; bit <<= 1 {
		step := bit << 1
		for start := 0; start < P; start += step {
			for i := start + bit; i < start+step; i++ {
				a[i] ^= a[i-bit]
			}
		}
	}
	prefParity := make([]uint8, P+1)
	prefXor := make([]uint64, P+1)
	for i := 0; i < P; i++ {
		prefParity[i+1] = prefParity[i] ^ a[i]
		prefXor[i+1] = prefXor[i]
		if a[i] == 1 {
			prefXor[i+1] ^= uint64(i)
		}
	}
	return &Pre{
		d:          uint(d),
		P:          uint64(P),
		fullParity: prefParity[P],
		fullXor:    prefXor[P],
		prefParity: prefParity,
		prefXor:    prefXor,
	}
}

func solveReference(input []byte) string {
	idx := 0
	n := len(input)
	nextInt := func() int {
		for idx < n && (input[idx] < '0' || input[idx] > '9') {
			idx++
		}
		val := 0
		for idx < n && input[idx] >= '0' && input[idx] <= '9' {
			val = val*10 + int(input[idx]-'0')
			idx++
		}
		return val
	}
	nextString := func() string {
		for idx < n && input[idx] <= ' ' {
			idx++
		}
		start := idx
		for idx < n && input[idx] > ' ' {
			idx++
		}
		return string(input[start:idx])
	}

	t := nextInt()
	cache := make(map[int]*Pre)
	var out bytes.Buffer
	for ; t > 0; t-- {
		_ = nextInt()
		m := uint64(nextInt())
		s := nextString()
		p := 0
		for i := 0; i < len(s); i++ {
			if s[i] == '1' {
				p++
			}
		}
		pre := cache[p]
		if pre == nil {
			pre = build(p)
			cache[p] = pre
		}
		cnt := m / pre.P
		rem := int(m % pre.P)
		var ans uint64
		if cnt&1 == 1 {
			ans ^= pre.fullXor
		}
		if pre.fullParity == 1 && cnt > 0 {
			ans ^= xorUpto(cnt-1) << pre.d
		}
		ans ^= pre.prefXor[rem]
		if pre.prefParity[rem] == 1 {
			ans ^= cnt << pre.d
		}
		out.WriteString(strconv.FormatUint(ans, 10))
		out.WriteByte('\n')
	}
	return out.String()
}

// ---- Verifier infrastructure ----

func runProgram(bin string, input []byte) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout")
		}
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func randBinary(k int, rng *rand.Rand) string {
	b := make([]byte, k)
	b[0] = '1'
	for i := 1; i < k; i++ {
		if rng.Intn(2) == 0 {
			b[i] = '0'
		} else {
			b[i] = '1'
		}
	}
	return string(b)
}

func buildTests() []byte {
	rng := rand.New(rand.NewSource(2056))
	type tc struct {
		k int
		m int64
		s string
	}

	var cases []tc

	cases = append(cases, tc{1, 1, "1"})
	cases = append(cases, tc{1, 5, "1"})
	cases = append(cases, tc{2, 7, "10"})
	cases = append(cases, tc{3, 6, "101"})

	cases = append(cases, tc{10, 123456789, randBinary(10, rng)})

	long1 := strings.Repeat("10", 50000)
	cases = append(cases, tc{len(long1), 999_999_937, long1})

	long2 := "1" + strings.Repeat("0", 99999)
	cases = append(cases, tc{len(long2), 1_000_000_000, long2})

	var buf bytes.Buffer
	fmt.Fprintln(&buf, len(cases))
	for _, c := range cases {
		fmt.Fprintf(&buf, "%d %d\n", c.k, c.m)
		fmt.Fprintln(&buf, c.s)
	}
	return buf.Bytes()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF2.go /path/to/binary")
		return
	}

	target, cleanTarget, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanTarget()

	input := buildTests()

	expOut := solveReference(input)
	gotOut, err := runProgram(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target failed: %v\n", err)
		os.Exit(1)
	}

	exp := strings.Fields(expOut)
	got := strings.Fields(gotOut)
	if len(exp) != len(got) {
		fmt.Fprintf(os.Stderr, "output length mismatch: expected %d lines, got %d\n", len(exp), len(got))
		os.Exit(1)
	}
	for i := range exp {
		if exp[i] != got[i] {
			fmt.Fprintf(os.Stderr, "mismatch at case %d: expected %s got %s\n", i+1, exp[i], got[i])
			os.Exit(1)
		}
	}

	fmt.Println("all tests passed")
}
