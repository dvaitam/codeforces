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

type TestCase struct {
	ticket string
}

func genCase(rng *rand.Rand) (string, string) {
	var sb strings.Builder
	for i := 0; i < 6; i++ {
		sb.WriteByte(byte('0' + rng.Intn(10)))
	}
	s := sb.String()
	sum1 := int(s[0]-'0') + int(s[1]-'0') + int(s[2]-'0')
	sum2 := int(s[3]-'0') + int(s[4]-'0') + int(s[5]-'0')
	ans := "NO"
	if sum1 == sum2 {
		ans = "YES"
	}
	input := fmt.Sprintf("1\n%s\n", s)
	expected := ans
	return input, expected
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
