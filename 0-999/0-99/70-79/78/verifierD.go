package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func isqrt(n int64) int64 {
	if n <= 0 {
		return 0
	}
	x := int64(math.Sqrt(float64(n)))
	for (x+1)*(x+1) <= n {
		x++
	}
	for x*x > n {
		x--
	}
	return x
}

func rowCount(j, n int64) int64 {
	jj := j * j
	if jj+j <= n {
		d := 4*n - 3*jj - 2*j + 1
		s := isqrt(d)
		r := (s - j - 1) / 2
		return 2*r + j + 1
	}
	d := 4*n - 3*jj - 4*j
	if d < 0 {
		return 0
	}
	s := isqrt(d)
	t := (j - s + 1) / 2
	return j - 2*t + 1
}

func solve(k int64) int64 {
	n := (k*k - 1) / 3
	var ans int64
	for j := int64(0); ; j++ {
		jj := j * j
		fmin := jj + j - jj/4
		if fmin > n {
			break
		}
		c := rowCount(j, n)
		if j == 0 {
			ans += c
		} else {
			ans += 2 * c
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, int64) {
	k := int64(rng.Intn(50) + 1)
	input := fmt.Sprintf("%d\n", k)
	return input, solve(k)
}

func runCase(bin, input string, expected int64) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
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
