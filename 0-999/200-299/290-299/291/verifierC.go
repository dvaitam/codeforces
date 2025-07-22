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

func solveC(n, k int, ips []uint32) string {
	arr := append([]uint32(nil), ips...)
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	for t := 1; t <= 31; t++ {
		shift := 32 - t
		cnt := 0
		var last uint32 = 0xffffffff
		for i, ip := range arr {
			v := ip >> shift
			if i == 0 || v != last {
				cnt++
				last = v
			}
			if cnt > k {
				break
			}
		}
		if cnt == k {
			mask := ^uint32(0) << shift
			return fmt.Sprintf("%d.%d.%d.%d\n", (mask>>24)&0xff, (mask>>16)&0xff, (mask>>8)&0xff, mask&0xff)
		}
	}
	return "-1\n"
}

func genValidCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	t := rng.Intn(30) + 1
	if t > 31 {
		t = 31
	}
	kMax := 1 << t
	if kMax > n {
		kMax = n
	}
	k := rng.Intn(kMax) + 1
	groups := rng.Perm(1 << t)[:k]
	ips := make([]uint32, n)
	for i := 0; i < n; i++ {
		g := groups[rng.Intn(k)]
		suffix := rng.Uint32() & ((1 << (32 - t)) - 1)
		ips[i] = uint32(g)<<(32-t) | suffix
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, k)
	for _, ip := range ips {
		fmt.Fprintf(&sb, "%d.%d.%d.%d\n", ip>>24, (ip>>16)&0xff, (ip>>8)&0xff, ip&0xff)
	}
	return sb.String(), solveC(n, k, ips)
}

func genRandomCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	k := rng.Intn(n) + 1
	ips := make([]uint32, n)
	for i := 0; i < n; i++ {
		ips[i] = rng.Uint32()
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, k)
	for _, ip := range ips {
		fmt.Fprintf(&sb, "%d.%d.%d.%d\n", ip>>24, (ip>>16)&0xff, (ip>>8)&0xff, ip&0xff)
	}
	return sb.String(), solveC(n, k, ips)
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
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		var in, expect string
		if rand.Intn(2) == 0 {
			in, expect = genValidCase(rand.New(rand.NewSource(time.Now().UnixNano() + int64(i))))
		} else {
			in, expect = genRandomCase(rand.New(rand.NewSource(time.Now().UnixNano() + int64(i))))
		}
		got, err := run(bin, in)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\ninput:\n%soutput:\n%s", i+1, err, in, got)
			return
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%sbut got:\n%s", i+1, in, expect, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
