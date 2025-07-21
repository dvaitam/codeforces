package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solveD(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, m, k int
	fmt.Fscan(reader, &n, &m, &k)
	xs := make([]int64, n)
	ys := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &xs[i], &ys[i])
	}
	rs := make([]int64, m)
	cx := make([]int64, m)
	cy := make([]int64, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &rs[i], &cx[i], &cy[i])
	}
	wlen := (m + 63) / 64
	bs := make([][]uint64, n)
	for i := 0; i < n; i++ {
		b := make([]uint64, wlen)
		xi, yi := xs[i], ys[i]
		for j := 0; j < m; j++ {
			dx := xi - cx[j]
			dy := yi - cy[j]
			if dx*dx+dy*dy < rs[j]*rs[j] {
				idx := j >> 6
				off := uint(j & 63)
				b[idx] |= 1 << off
			}
		}
		bs[i] = b
	}
	var sb strings.Builder
	for qi := 0; qi < k; qi++ {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		a--
		b--
		cnt := 0
		ba, bb := bs[a], bs[b]
		for i := 0; i < wlen; i++ {
			cnt += bits.OnesCount64(ba[i] ^ bb[i])
		}
		sb.WriteString(fmt.Sprintf("%d\n", cnt))
	}
	return sb.String()
}

func genTestD() (string, string) {
	n := rand.Intn(3) + 1
	m := rand.Intn(3) + 1
	k := rand.Intn(3) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)
	for i := 0; i < n; i++ {
		x := rand.Intn(11) - 5
		y := rand.Intn(11) - 5
		fmt.Fprintf(&sb, "%d %d\n", x, y)
	}
	for i := 0; i < m; i++ {
		r := rand.Intn(5) + 1
		x := rand.Intn(11) - 5
		y := rand.Intn(11) - 5
		fmt.Fprintf(&sb, "%d %d %d\n", r, x, y)
	}
	for i := 0; i < k; i++ {
		a := rand.Intn(n) + 1
		b := rand.Intn(n) + 1
		fmt.Fprintf(&sb, "%d %d\n", a, b)
	}
	in := sb.String()
	out := solveD(in)
	return in, out
}

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierD.go <binary>")
		return
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, exp := genTestD()
		got, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\nInput:\n%s\nOutput:\n%s\n", i+1, err, in, got)
			return
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s", i+1, in, exp, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
