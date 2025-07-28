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

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func pow2(exp int64) int64 {
	res := int64(1)
	base := int64(2)
	for exp > 0 {
		if exp&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
		exp >>= 1
	}
	return res
}

func expected(n int64, segs [][3]int64) string {
	var total int64
	for _, s := range segs {
		total |= s[2]
	}
	ans := total % mod * pow2(n-1) % mod
	return strconv.FormatInt(ans, 10)
}

func genCase(rng *rand.Rand) (string, string) {
	n := int64(rng.Intn(50) + 1)
	m := rng.Intn(int(n)) + 1
	segs := make([][3]int64, 0, m)
	cover := make([]bool, n)
	for i := 0; i < m; i++ {
		l := int64(rng.Intn(int(n)) + 1)
		r := int64(rng.Intn(int(n-l+1)) + int(l))
		x := rng.Int63n(1 << 30)
		segs = append(segs, [3]int64{l, r, x})
		for j := l; j <= r; j++ {
			cover[j-1] = true
		}
	}
	for idx, ok := range cover {
		if !ok {
			x := rng.Int63n(1 << 30)
			segs = append(segs, [3]int64{int64(idx + 1), int64(idx + 1), x})
		}
	}
	m = len(segs)

	var sb strings.Builder
	sb.WriteString("1\n")
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for _, s := range segs {
		fmt.Fprintf(&sb, "%d %d %d\n", s[0], s[1], s[2])
	}
	inp := sb.String()
	exp := expected(n, segs)
	return inp, exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		got, err := runBinary(exe, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, in, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
