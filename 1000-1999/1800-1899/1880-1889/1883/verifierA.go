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

func solveCase(pin string) string {
	pos := 1
	ans := 0
	for _, ch := range pin {
		d := int(ch - '0')
		if d > pos {
			ans += d - pos
		} else {
			ans += pos - d
		}
		ans++
		pos = d
	}
	return fmt.Sprint(ans)
}

func runCase(bin, pin string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(fmt.Sprintf("1\n%s\n", pin))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := solveCase(pin)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func randomCase(rng *rand.Rand) string {
	b := make([]byte, 4)
	for i := range b {
		b[i] = byte('0' + rng.Intn(10))
	}
	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []string{"0000", "9999", "1234", "4321"}
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCase(rng))
	}
	for idx, pin := range cases {
		if err := runCase(bin, pin); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: %s\n", idx+1, err, pin)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
