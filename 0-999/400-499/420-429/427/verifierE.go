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

func expected(n, m int, a []int64) int64 {
	sL := make([]int64, n)
	for i := 0; i < n; i++ {
		if i >= m {
			sL[i] = a[i] + sL[i-m]
		} else {
			sL[i] = a[i]
		}
	}
	sR := make([]int64, n)
	for i := n - 1; i >= 0; i-- {
		if i+m < n {
			sR[i] = a[i] + sR[i+m]
		} else {
			sR[i] = a[i]
		}
	}
	var best int64 = -1
	for i := 0; i < n; {
		j := i + 1
		for j < n && a[j] == a[i] {
			j++
		}
		S := a[i]
		L := int64(i)
		var sumL int64
		if i > 0 {
			sumL = sL[i-1]
		}
		TL := (L + int64(m) - 1) / int64(m)
		distL := TL*S - sumL
		Rn := int64(n - j)
		var sumRsel int64
		TR := (Rn + int64(m) - 1) / int64(m)
		if Rn > 0 {
			idx0 := int64(n-1) - (TR-1)*int64(m)
			if idx0 < int64(j) {
				modR := int64(n-1) % int64(m)
				rem := int64(j) % int64(m)
				delta := (modR - rem + int64(m)) % int64(m)
				idx0 = int64(j) + delta
			}
			if idx0 < int64(n) {
				sumRsel = sR[idx0]
			}
		}
		distR := sumRsel - TR*S
		total := 2 * (distL + distR)
		if best < 0 || total < best {
			best = total
		}
		i = j
	}
	return best
}

func generateCase(rng *rand.Rand) (string, int64) {
	n := rng.Intn(8) + 1
	m := rng.Intn(4) + 1
	a := make([]int64, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		a[i] = int64(rng.Intn(41) - 20)
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", a[i]))
	}
	sb.WriteByte('\n')
	return sb.String(), expected(n, m, a)
}

func runCase(exe, input string, exp int64) error {
	cmd := exec.Command(exe)
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
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
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
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
