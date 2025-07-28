package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solve(x string) string {
	bs := []byte(x)
	n := len(bs)
	for i := n - 1; i >= 1; i-- {
		d1 := int(bs[i-1] - '0')
		d2 := int(bs[i] - '0')
		sum := d1 + d2
		if sum >= 10 {
			res := make([]byte, 0, n)
			res = append(res, bs[:i-1]...)
			res = append(res, byte('0'+sum/10))
			res = append(res, byte('0'+sum%10))
			res = append(res, bs[i+1:]...)
			return string(res)
		}
	}
	sum := int(bs[0] - '0' + bs[1] - '0')
	res := make([]byte, 0, n-1)
	res = append(res, byte('0'+sum))
	res = append(res, bs[2:]...)
	return string(res)
}

func runCase(bin, s string) error {
	input := fmt.Sprintf("1\n%s\n", s)
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
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	expect := solve(s)
	if got != expect {
		return fmt.Errorf("expected %s got %s", expect, got)
	}
	return nil
}

func randomCase(rng *rand.Rand) string {
	length := rng.Intn(18) + 2
	b := make([]byte, length)
	b[0] = '1' + byte(rng.Intn(9))
	for i := 1; i < length; i++ {
		b[i] = byte(rng.Intn(10)) + '0'
	}
	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []string{"10057", "10", "99", "1234", "1010", "90876"}
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCase(rng))
	}
	for idx, s := range cases {
		if err := runCase(bin, s); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
