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

type bitset []uint64

func newBitset(nbits int) bitset {
	n := (nbits + 63) >> 6
	return make(bitset, n)
}

func (b bitset) setAll(nbits int) {
	for i := range b {
		b[i] = ^uint64(0)
	}
	tail := nbits & 63
	if tail != 0 {
		b[len(b)-1] &= (uint64(1) << tail) - 1
	}
}

func (b bitset) any() bool {
	for _, w := range b {
		if w != 0 {
			return true
		}
	}
	return false
}

func (dst bitset) shl(src bitset, k int) {
	n := len(src)
	wordShift := k >> 6
	offset := uint(k & 63)
	for i := n - 1; i >= 0; i-- {
		var v uint64
		j := i - wordShift
		if j >= 0 {
			v = src[j] << offset
			if offset != 0 && j-1 >= 0 {
				v |= src[j-1] >> (64 - offset)
			}
		}
		dst[i] = v
	}
}

func (dst bitset) shr(src bitset, k int) {
	n := len(src)
	wordShift := k >> 6
	offset := uint(k & 63)
	for i := 0; i < n; i++ {
		var v uint64
		j := i + wordShift
		if j < n {
			v = src[j] >> offset
			if offset != 0 && j+1 < n {
				v |= src[j+1] << (64 - offset)
			}
		}
		dst[i] = v
	}
}

func feasible(a []int, d int) bool {
	bsize := d + 1
	dpPrev := newBitset(bsize)
	dpPrev.setAll(bsize)
	dpCur := newBitset(bsize)
	tmp := newBitset(bsize)
	for _, ai := range a {
		tmp.shl(dpPrev, ai)
		for i := range dpCur {
			dpCur[i] = tmp[i]
		}
		tmp.shr(dpPrev, ai)
		for i := range dpCur {
			dpCur[i] |= tmp[i]
		}
		tail := bsize & 63
		if tail != 0 {
			mask := uint64(1)<<uint(tail) - 1
			dpCur[len(dpCur)-1] &= mask
		}
		if !dpCur.any() {
			return false
		}
		dpPrev, dpCur = dpCur, dpPrev
	}
	return true
}

func solveG(a []int) string {
	maxa := 0
	for _, v := range a {
		if v > maxa {
			maxa = v
		}
	}
	lo, hi := maxa, 2*maxa
	ans := hi
	for lo <= hi {
		mid := (lo + hi) >> 1
		if feasible(a, mid) {
			ans = mid
			hi = mid - 1
		} else {
			lo = mid + 1
		}
	}
	return fmt.Sprintf("%d\n", ans)
}

func genCaseG(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	a := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(5) + 1
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	input := sb.String()
	expect := solveG(a)
	return input, expect
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expect := genCaseG(rng)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\ninput:\n%soutput:\n%s", i+1, err, input, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%sbut got:\n%s", i+1, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
