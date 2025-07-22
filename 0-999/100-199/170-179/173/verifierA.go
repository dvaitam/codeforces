package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func solveA(n int64, A, B string) (int64, int64) {
	m := int64(len(A))
	k := int64(len(B))
	g := gcd(m, k)
	l := m / g * k
	var winA, winB int64
	for i := int64(0); i < l; i++ {
		a := A[i%m]
		b := B[i%k]
		if (a == 'R' && b == 'S') || (a == 'S' && b == 'P') || (a == 'P' && b == 'R') {
			winA++
		} else if (b == 'R' && a == 'S') || (b == 'S' && a == 'P') || (b == 'P' && a == 'R') {
			winB++
		}
	}
	full := n / l
	rem := n % l
	totalA := winA * full
	totalB := winB * full
	for i := int64(0); i < rem; i++ {
		a := A[i%m]
		b := B[i%k]
		if (a == 'R' && b == 'S') || (a == 'S' && b == 'P') || (a == 'P' && b == 'R') {
			totalA++
		} else if (b == 'R' && a == 'S') || (b == 'S' && a == 'P') || (b == 'P' && a == 'R') {
			totalB++
		}
	}
	return totalB, totalA
}

func randSeq(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	letters := []byte{'R', 'S', 'P'}
	for i := range b {
		b[i] = letters[rng.Intn(3)]
	}
	return string(b)
}

func runCase(bin string, n int64, A, B string) error {
	var input bytes.Buffer
	fmt.Fprintf(&input, "%d\n%s\n%s\n", n, A, B)
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	reader := bufio.NewReader(&out)
	var x, y int64
	if _, err := fmt.Fscan(reader, &x, &y); err != nil {
		return fmt.Errorf("cannot parse output: %v", err)
	}
	expX, expY := solveA(n, A, B)
	if x != expX || y != expY {
		return fmt.Errorf("expected %d %d got %d %d", expX, expY, x, y)
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

	tests := 100
	for i := 0; i < tests; i++ {
		n := rng.Int63n(1_000_000_000) + 1
		lenA := rng.Intn(50) + 1
		lenB := rng.Intn(50) + 1
		A := randSeq(rng, lenA)
		B := randSeq(rng, lenB)
		if err := runCase(bin, n, A, B); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: n=%d A=%s B=%s\n", i+1, err, n, A, B)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
