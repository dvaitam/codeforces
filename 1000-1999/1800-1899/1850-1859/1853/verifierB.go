package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var fib []int

func initFib() []int {
	f := []int{0, 1}
	for f[len(f)-1] <= 200000 {
		f = append(f, f[len(f)-1]+f[len(f)-2])
	}
	return f
}

func extendedEuclid(a, b int) (int, int, int) {
	if b == 0 {
		return 1, 0, a
	}
	x1, y1, g := extendedEuclid(b, a%b)
	x := y1
	y := x1 - (a/b)*y1
	return x, y, g
}

func modInverse(a, mod int) int {
	x, _, g := extendedEuclid(a, mod)
	if g != 1 {
		return 0
	}
	if x < 0 {
		x += mod
	}
	return x % mod
}

func expectedB(n, k int) int {
	if k >= len(fib) || fib[k-1] > n {
		return 0
	}
	fk2 := fib[k-2]
	fk1 := fib[k-1]
	fk := fib[k]
	inv := modInverse(fk2, fk1)
	a0 := (n % fk1) * inv % fk1
	b0 := (n - fk2*a0) / fk1
	tMax1 := b0 / fk2
	tMax2 := (b0 - a0) / fk
	if tMax1 < tMax2 {
		tMax2 = tMax1
	}
	if tMax2 < 0 {
		return 0
	}
	return tMax2 + 1
}

func genCase(rng *rand.Rand) (int, int) {
	n := rng.Intn(200000) + 1
	// len(fib) may be about 29
	if rng.Intn(5) == 0 {
		k := len(fib) + rng.Intn(10) + 1
		return n, k
	}
	k := rng.Intn(len(fib)-2) + 3
	return n, k
}

func runBinary(bin, input string) (string, string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var outBuf bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	return outBuf.String(), errBuf.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	fib = initFib()
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	const tests = 100
	ns := make([]int, tests)
	ks := make([]int, tests)
	expected := make([]int, tests)
	for i := 0; i < tests; i++ {
		n, k := genCase(rng)
		ns[i] = n
		ks[i] = k
		expected[i] = expectedB(n, k)
	}

	var input strings.Builder
	fmt.Fprintf(&input, "%d\n", tests)
	for i := 0; i < tests; i++ {
		fmt.Fprintf(&input, "%d %d\n", ns[i], ks[i])
	}

	out, errOut, err := runBinary(bin, input.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\n%s", err, errOut)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(strings.NewReader(out))
	scanner.Split(bufio.ScanWords)
	for i := 0; i < tests; i++ {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "missing output for test %d\n", i+1)
			os.Exit(1)
		}
		val, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid integer on test %d\n", i+1)
			os.Exit(1)
		}
		if val != expected[i] {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d got %d\n", i+1, expected[i], val)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
