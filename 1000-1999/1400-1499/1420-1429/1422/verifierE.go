package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const LOG = 17
const hashB1 uint64 = 911382323
const hashB2 uint64 = 972663749

func solveE(input string) string {
	str := strings.TrimSpace(input)
	s := []byte(str)
	n := len(s)
	if n == 0 {
		return ""
	}

	chArr := make([]byte, n+1)
	nxtArr := make([]int32, n+1)
	lnsArr := make([]int32, n+1)
	ansArr := make([]int32, n+2)
	pow1 := make([]uint64, n+1)
	pow2 := make([]uint64, n+1)

	pow1[0], pow2[0] = 1, 1
	for i := 0; i < n; i++ {
		pow1[i+1] = pow1[i] * hashB1
		pow2[i+1] = pow2[i] * hashB2
	}

	var up [LOG][]int32
	var h1 [LOG][]uint64
	var h2 [LOG][]uint64
	for k := 0; k < LOG; k++ {
		up[k] = make([]int32, n+1)
		h1[k] = make([]uint64, n+1)
		h2[k] = make([]uint64, n+1)
	}

	build := func(u int32, c byte, nx int32) {
		ui := int(u)
		chArr[ui] = c
		nxtArr[ui] = nx
		lnsArr[ui] = lnsArr[int(nx)] + 1
		up[0][ui] = nx
		v := uint64(c-'a') + 1
		h1[0][ui] = v
		h2[0][ui] = v
		for k := 1; k < LOG; k++ {
			if int(lnsArr[ui]) >= 1<<k {
				mid := up[k-1][ui]
				mi := int(mid)
				up[k][ui] = up[k-1][mi]
				h1[k][ui] = h1[k-1][ui]*pow1[1<<(k-1)] + h1[k-1][mi]
				h2[k][ui] = h2[k-1][ui]*pow2[1<<(k-1)] + h2[k-1][mi]
			}
		}
	}

	cmp := func(a, b int32) int {
		if a == b {
			return 0
		}
		x, y := a, b
		for k := LOG - 1; k >= 0; k-- {
			if int(lnsArr[int(x)]) >= 1<<k && int(lnsArr[int(y)]) >= 1<<k &&
				h1[k][int(x)] == h1[k][int(y)] &&
				h2[k][int(x)] == h2[k][int(y)] {
				x = up[k][int(x)]
				y = up[k][int(y)]
			}
		}
		if x == 0 {
			if y == 0 {
				return 0
			}
			return -1
		}
		if y == 0 {
			return 1
		}
		if chArr[int(x)] < chArr[int(y)] {
			return -1
		}
		return 1
	}

	jump := func(u int32, d int32) int32 {
		x := u
		rem := int(d)
		for k := LOG - 1; k >= 0; k-- {
			if rem >= 1<<k {
				x = up[k][int(x)]
				rem -= 1 << k
			}
		}
		return x
	}

	collect := func(u int32, k int) string {
		buf := make([]byte, 0, k)
		x := u
		for k > 0 && x != 0 {
			buf = append(buf, chArr[int(x)])
			x = nxtArr[int(x)]
			k--
		}
		return string(buf)
	}

	lastK := func(u int32, k int) string {
		start := jump(u, lnsArr[int(u)]-int32(k))
		return collect(start, k)
	}

	for i := n - 1; i >= 0; i-- {
		u := int32(i + 1)
		build(u, s[i], ansArr[i+1])
		if i+1 < n && s[i] == s[i+1] {
			if cmp(u, ansArr[i+2]) <= 0 {
				ansArr[i] = u
			} else {
				ansArr[i] = ansArr[i+2]
			}
		} else {
			ansArr[i] = u
		}
	}

	var out strings.Builder
	for i := 0; i < n; i++ {
		u := ansArr[i]
		l := int(lnsArr[int(u)])
		if l == 0 {
			out.WriteString("0\n")
			continue
		}
		var res string
		if l <= 10 {
			res = collect(u, l)
		} else {
			res = collect(u, 5) + "..." + lastK(u, 2)
		}
		out.WriteString(strconv.Itoa(l))
		out.WriteByte(' ')
		out.WriteString(res)
		out.WriteByte('\n')
	}
	return strings.TrimSpace(out.String())
}

func generateInputs() []string {
	rng := rand.New(rand.NewSource(46))
	var inputs []string
	fixed := []string{"a", "aaaa", "abab"}
	inputs = append(inputs, fixed...)
	for len(inputs) < 100 {
		n := rng.Intn(8) + 1
		var sb strings.Builder
		for i := 0; i < n; i++ {
			sb.WriteByte(byte('a' + rng.Intn(3)))
		}
		inputs = append(inputs, sb.String())
	}
	return inputs
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input + "\n")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierE /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	inputs := generateInputs()
	for i, input := range inputs {
		expected := solveE(input)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Printf("Wrong answer on test %d\nInput:%s\nExpected:%s\nGot:%s\n", i+1, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(inputs))
}
