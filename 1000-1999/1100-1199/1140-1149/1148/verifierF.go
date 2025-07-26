package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func run(bin string, in []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func genCase() (string, []int64, []uint64) {
	n := rand.Intn(6) + 1
	vals := make([]int64, n)
	masks := make([]uint64, n)
	sum := int64(0)
	for i := 0; i < n; i++ {
		vals[i] = int64(rand.Intn(41) - 20)
		masks[i] = uint64(rand.Int63n(1<<10) + 1)
		sum += vals[i]
	}
	if sum == 0 {
		vals[0]++
		sum++
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", vals[i], masks[i]))
	}
	return sb.String(), vals, masks
}

func apply(vals []int64, masks []uint64, s uint64) int64 {
	tot := int64(0)
	for i := 0; i < len(vals); i++ {
		v := vals[i]
		if bits.OnesCount64(s&masks[i])%2 == 1 {
			v = -v
		}
		tot += v
	}
	return tot
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for t := 0; t < 100; t++ {
		input, vals, masks := genCase()
		out, err := run(bin, []byte(input))
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s", t+1, err, input)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		s, err := strconv.ParseUint(out, 10, 64)
		if err != nil || s == 0 {
			fmt.Fprintf(os.Stderr, "bad output on test %d\n", t+1)
			os.Exit(1)
		}
		initial := int64(0)
		for _, v := range vals {
			initial += v
		}
		after := apply(vals, masks, s)
		if initial > 0 && after >= 0 {
			fmt.Fprintf(os.Stderr, "sum sign not changed on test %d\n", t+1)
			os.Exit(1)
		}
		if initial < 0 && after <= 0 {
			fmt.Fprintf(os.Stderr, "sum sign not changed on test %d\n", t+1)
			os.Exit(1)
		}
		if after == 0 {
			fmt.Fprintf(os.Stderr, "sum became zero on test %d\n", t+1)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
