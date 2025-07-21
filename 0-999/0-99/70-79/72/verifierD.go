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

func reverseString(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func substrVal(s string, a, b, step int) string {
	if a < 1 {
		a = 1
	}
	if b > len(s) {
		b = len(s)
	}
	if step <= 1 {
		if a > b {
			return ""
		}
		return s[a-1 : b]
	}
	var sb strings.Builder
	for i := a; i <= b && i <= len(s); i += step {
		sb.WriteByte(s[i-1])
	}
	return sb.String()
}

func genExpr(rng *rand.Rand, depth int) (string, string) {
	if depth <= 0 || rng.Intn(4) == 0 {
		l := rng.Intn(5) + 1
		b := make([]byte, l)
		for i := range b {
			b[i] = byte('a' + rng.Intn(3))
		}
		str := string(b)
		return fmt.Sprintf("\"%s\"", str), str
	}
	switch rng.Intn(3) {
	case 0:
		e1, v1 := genExpr(rng, depth-1)
		e2, v2 := genExpr(rng, depth-1)
		return fmt.Sprintf("concat(%s,%s)", e1, e2), v1 + v2
	case 1:
		e, v := genExpr(rng, depth-1)
		return fmt.Sprintf("reverse(%s)", e), reverseString(v)
	default:
		e, v := genExpr(rng, depth-1)
		n := len(v)
		if n == 0 {
			return fmt.Sprintf("substr(%s,1,0)", e), ""
		}
		a := rng.Intn(n) + 1
		b := rng.Intn(n) + 1
		if a > b {
			a, b = b, a
		}
		if rng.Intn(2) == 0 {
			return fmt.Sprintf("substr(%s,%d,%d)", e, a, b), substrVal(v, a, b, 1)
		}
		step := rng.Intn(n) + 1
		return fmt.Sprintf("substr(%s,%d,%d,%d)", e, a, b, step), substrVal(v, a, b, step)
	}
}

func generateCase(rng *rand.Rand) (string, string) {
	expr, val := genExpr(rng, 3)
	return expr + "\n", val
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
	if got != expected {
		return fmt.Errorf("expected '%s' got '%s'", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
