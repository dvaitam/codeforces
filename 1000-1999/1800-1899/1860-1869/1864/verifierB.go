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

func expectedB(n, k int, s string) string {
	if k%2 == 0 {
		b := []byte(s)
		sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
		return string(b) + "\n"
	}
	odd := make([]byte, 0, (n+1)/2)
	even := make([]byte, 0, n/2)
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			odd = append(odd, s[i])
		} else {
			even = append(even, s[i])
		}
	}
	sort.Slice(odd, func(i, j int) bool { return odd[i] < odd[j] })
	sort.Slice(even, func(i, j int) bool { return even[i] < even[j] })
	res := make([]byte, n)
	oi, ei := 0, 0
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			res[i] = odd[oi]
			oi++
		} else {
			res[i] = even[ei]
			ei++
		}
	}
	return string(res) + "\n"
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	k := rng.Intn(n-1) + 1
	bytes := make([]byte, n)
	for i := range bytes {
		bytes[i] = byte('a' + rng.Intn(26))
	}
	s := string(bytes)
	input := fmt.Sprintf("1\n%d %d\n%s\n", n, k, s)
	expect := expectedB(n, k, s)
	return input, expect
}

func runCase(bin, input, expected string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	expStr := strings.TrimSpace(expected)
	if outStr != expStr {
		return fmt.Errorf("expected %q got %q", expStr, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
