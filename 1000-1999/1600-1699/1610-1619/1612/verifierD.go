package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func run(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func isMagic(a, b, x int64) bool {
	g := gcd(a, b)
	if x%g != 0 {
		return false
	}
	a /= g
	b /= g
	x /= g
	for a != 0 && b != 0 && max64(a, b) >= x {
		if a == x || b == x {
			return true
		}
		if a < b {
			a, b = b, a
		}
		if (a-x)%b == 0 {
			return true
		}
		a %= b
	}
	return a == x || b == x
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func check(a, b, x int64, output string) error {
	got := strings.TrimSpace(output)
	want := "NO"
	if isMagic(a, b, x) {
		want = "YES"
	}
	if got != want {
		return fmt.Errorf("expected %s got %s", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		a := rand.Int63n(1000000) + 1
		b := rand.Int63n(1000000) + 1
		x := rand.Int63n(1000000) + 1
		input := fmt.Sprintf("1\n%d %d %d\n", a, b, x)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := check(a, b, x, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: a=%d b=%d x=%d\n", i+1, err, a, b, x)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
