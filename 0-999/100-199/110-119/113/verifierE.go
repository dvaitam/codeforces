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

func digits(x int64) int64 {
	if x < 0 {
		x = -x
	}
	cnt := int64(0)
	for {
		cnt++
		x /= 10
		if x == 0 {
			break
		}
	}
	return cnt
}

func change(N int64, a, b int64) int64 {
	var cnt int64
	for i := int64(0); i < N; i++ {
		if a%10 != b%10 {
			cnt++
		}
		a /= 10
		b /= 10
	}
	return cnt
}

func solveE(H, M, K int64, x1, y1, x2, y2 int64) int64 {
	N1 := digits(H - 1)
	N2 := digits(M - 1)
	if K > N1+N2 {
		return 0
	}
	maxP := N1 + N2
	if maxP < K {
		maxP = K
	}
	if maxP > 18 {
		maxP = 18
	}
	TEN := make([]int64, maxP+1)
	TEN[0] = 1
	for i := int64(1); i <= maxP; i++ {
		TEN[i] = TEN[i-1] * 10
	}
	C2 := change(N2, M-1, 0)
	C1 := change(N1, H-1, 0)
	F := func(x, y int64) int64 {
		if K <= N2 {
			d := TEN[K-1]
			part := ((M - 1) / d) * x
			part += y / d
			need := K - C2
			if need <= 0 {
				part += x
			} else {
				d2 := TEN[need-1]
				part += x / d2
			}
			return part
		}
		need := K - C2
		d2 := TEN[need-1]
		return x / d2
	}
	a := F(x1, y1)
	b := F(x2, y2)
	var ans int64
	if x1 < x2 || (x1 == x2 && y1 <= y2) {
		ans = b - a
	} else {
		total := F(H-1, M-1)
		ans = total - (a - b)
		if C1+C2 >= K {
			ans++
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	H := int64(rng.Intn(1000) + 2)
	M := int64(rng.Intn(1000) + 2)
	K := int64(rng.Intn(6) + 1)
	x1 := int64(rng.Intn(int(H)))
	y1 := int64(rng.Intn(int(M)))
	x2 := int64(rng.Intn(int(H)))
	y2 := int64(rng.Intn(int(M)))
	input := fmt.Sprintf("%d %d %d\n%d %d\n%d %d\n", H, M, K, x1, y1, x2, y2)
	exp := fmt.Sprintf("%d", solveE(H, M, K, x1, y1, x2, y2))
	return input, exp
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
