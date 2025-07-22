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

type testCase struct{ input string }

func solveCase(in string) string {
	in = strings.TrimSpace(in)
	cnt := [10]int{}
	for i := 0; i < len(in); i++ {
		d := in[i] - '0'
		cnt[d]++
	}
	cnt[1]--
	cnt[6]--
	cnt[8]--
	cnt[9]--
	perms := []string{"1869", "6198", "1896", "1689", "1986", "1968", "1698"}
	m := 0
	var sb strings.Builder
	for d := 9; d >= 1; d-- {
		for cnt[d] > 0 {
			m = (m*10 + d) % 7
			sb.WriteByte(byte('0' + d))
			cnt[d]--
		}
	}
	sb.WriteString(perms[m])
	for cnt[0] > 0 {
		sb.WriteByte('0')
		cnt[0]--
	}
	return sb.String()
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := solveCase(tc.input)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(20) + 4
	digits := []byte{'1', '6', '8', '9'}
	for len(digits) < n {
		digits = append(digits, byte('0'+rng.Intn(10)))
	}
	rng.Shuffle(len(digits), func(i, j int) { digits[i], digits[j] = digits[j], digits[i] })
	return testCase{input: string(digits)}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{{input: "1689"}, {input: "1869"}, {input: "196889"}}
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCase(rng))
	}
	for idx, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: %s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
