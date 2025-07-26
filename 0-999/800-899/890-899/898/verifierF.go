package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func buildInputF(rng *rand.Rand) (string, string) {
	a := rng.Intn(1000)
	b := rng.Intn(1000)
	sa := fmt.Sprintf("%d", a)
	sb := fmt.Sprintf("%d", b)
	sc := fmt.Sprintf("%d", a+b)
	s := sa + sb + sc + "\n"
	return s, sa + "+" + sb + "=" + sc + "\n"
}

func checkOutputF(in, out string) error {
	out = strings.TrimSpace(out)
	plus := strings.IndexByte(out, '+')
	eq := strings.IndexByte(out, '=')
	if plus == -1 || eq == -1 || plus > eq {
		return fmt.Errorf("output format")
	}
	a := out[:plus]
	b := out[plus+1 : eq]
	c := out[eq+1:]
	if len(a) == 0 || len(b) == 0 || len(c) == 0 {
		return fmt.Errorf("empty parts")
	}
	if (len(a) > 1 && a[0] == '0') || (len(b) > 1 && b[0] == '0') || (len(c) > 1 && c[0] == '0') {
		return fmt.Errorf("leading zeros")
	}
	if a+b+c != strings.TrimSpace(in) {
		return fmt.Errorf("digits mismatch")
	}
	ai := new(big.Int)
	bi := new(big.Int)
	ci := new(big.Int)
	ai.SetString(a, 10)
	bi.SetString(b, 10)
	ci.SetString(c, 10)
	sum := new(big.Int).Add(ai, bi)
	if sum.Cmp(ci) != 0 {
		return fmt.Errorf("not equal")
	}
	return nil
}

func runCaseF(bin, in string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	if err := checkOutputF(in, buf.String()); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, _ := buildInputF(rng)
		if err := runCaseF(bin, in); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
