package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func powMod(a, e, m int) int {
	res := 1 % m
	a %= m
	for i := 0; i < e; i++ {
		res = res * a % m
	}
	return res
}

func solveB(r io.Reader) string {
	in := bufio.NewReader(r)
	var n string
	if _, err := fmt.Fscan(in, &n); err != nil {
		return ""
	}
	if n == "0" {
		return "4\n"
	}
	mod2, mod4 := 0, 0
	for i := 0; i < len(n); i++ {
		d := int(n[i] - '0')
		mod2 = (mod2*10 + d) % 2
		mod4 = (mod4*10 + d) % 4
	}
	e := mod4
	if e == 0 {
		e = 4
	}
	e4 := mod2
	if e4 == 0 {
		e4 = 2
	}
	res := (1 + powMod(2, e, 5) + powMod(3, e, 5) + powMod(4, e4, 5)) % 5
	return fmt.Sprintf("%d\n", res)
}

func randomBigInt(rng *rand.Rand) string {
	length := rng.Intn(100) + 1
	b := make([]byte, length)
	b[0] = byte(rng.Intn(9)+1) + '0'
	for i := 1; i < length; i++ {
		b[i] = byte(rng.Intn(10)) + '0'
	}
	return string(b)
}

func generateCase(rng *rand.Rand) (string, string) {
	if rng.Intn(10) == 0 {
		nums := []string{"0", "1", "2", "3", "4", "10", "100", "12345678901234567890"}
		n := nums[rng.Intn(len(nums))]
		input := n + "\n"
		expected := solveB(strings.NewReader(input))
		return input, expected
	}
	n := randomBigInt(rng)
	input := n + "\n"
	expected := solveB(strings.NewReader(input))
	return input, expected
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
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
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
	cases := []struct{ in, exp string }{}
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		cases = append(cases, struct{ in, exp string }{in, exp})
	}
	for i, c := range cases {
		if err := runCase(bin, c.in, c.exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, c.in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
