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

func solveE(a []int) int {
	const INF = int(^uint(0) >> 1)
	n := len(a)
	ans := INF
	for i := n - 1; i >= 0; i-- {
		u := a[i] - 1
		x := 1
		if a[i] >= ans {
			break
		}
		for j := i + 1; j < n; j++ {
			v1 := u / a[j]
			x += v1
			u -= v1 * a[j]
			v2 := a[i] + a[j] - u - 1
			if v2 >= ans {
				continue
			}
			y := 0
			vv := v2
			k := 0
			for vv != 0 && y <= x {
				y += vv / a[k]
				vv %= a[k]
				k++
			}
			if y > x {
				ans = v2
			}
		}
	}
	if ans == INF {
		return -1
	}
	return ans
}

func generateCaseE(rng *rand.Rand) (string, int) {
	n := rng.Intn(5) + 2
	a := make([]int, n)
	a[0] = 1
	for i := 1; i < n; i++ {
		a[i] = a[i-1] + rng.Intn(5) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range a {
		sb.WriteString(fmt.Sprintf("%d", v))
		if i+1 < n {
			sb.WriteByte(' ')
		}
	}
	sb.WriteByte('\n')
	return sb.String(), solveE(a)
}

func runCaseE(bin, input string, expected int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Sscan(strings.TrimSpace(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, exp := generateCaseE(rng)
		if err := runCaseE(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
