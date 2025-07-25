package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveD(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return ""
	}
	var tmp string
	fmt.Fscan(reader, &tmp)
	s := []byte(tmp)
	z := make([]int, n)
	z[0] = n
	l, r := 0, 0
	for i := 1; i < n; i++ {
		if i <= r {
			k0 := r - i + 1
			if z[i-l] < k0 {
				z[i] = z[i-l]
			} else {
				z[i] = k0
			}
		}
		for i+z[i] < n && s[z[i]] == s[i+z[i]] {
			z[i]++
		}
		if i+z[i]-1 > r {
			l = i
			r = i + z[i] - 1
		}
	}
	diff := make([]int, n+2)
	for L := 1; L < n; L++ {
		l64 := int64(k) * int64(L)
		if l64 <= int64(n) && int64(z[L]) >= l64 {
			r64 := int64(z[L]) + int64(L)
			maxR := int64(k+1) * int64(L)
			if r64 > maxR {
				r64 = maxR
			}
			if r64 > int64(n) {
				r64 = int64(n)
			}
			lpos := int(l64)
			rpos := int(r64)
			if lpos <= rpos {
				diff[lpos]++
				diff[rpos+1]--
			}
		}
		if int64(k-1)*int64(L) <= int64(z[L]) {
			pos := k * L
			if pos <= n {
				diff[pos]++
				diff[pos+1]--
			}
		}
	}
	res := make([]byte, n)
	cur := 0
	for i := 1; i <= n; i++ {
		cur += diff[i]
		if cur > 0 {
			res[i-1] = '1'
		} else {
			res[i-1] = '0'
		}
	}
	return string(res)
}

func genTestD(rng *rand.Rand) string {
	n := rng.Intn(50) + 1
	k := rng.Intn(10) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte('a' + rng.Intn(2))
	}
	sb.Write(b)
	sb.WriteByte('\n')
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
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	for i := 1; i <= 100; i++ {
		in := genTestD(rng)
		expect := solveD(in)
		got, err := run(exe, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i, in, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
