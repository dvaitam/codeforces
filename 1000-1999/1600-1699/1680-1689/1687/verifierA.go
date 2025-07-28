package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func runCandidate(bin, input string) (string, error) {
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

func expected(n, k int, a []int64) int64 {
	var sum int64
	for _, v := range a {
		sum += v
	}
	if k < n {
		pref := make([]int64, n+1)
		for i := 0; i < n; i++ {
			pref[i+1] = pref[i] + a[i]
		}
		best := pref[k] - pref[0]
		for i := 1; i+k <= n; i++ {
			s := pref[i+k] - pref[i]
			if s > best {
				best = s
			}
		}
		k64 := int64(k)
		return best + k64*(k64-1)/2
	}
	k64 := int64(k)
	n64 := int64(n)
	add := n64*(n64-1)/2 + (k64-n64)*n64
	return sum + add
}

func genTest(rng *rand.Rand) (int, int, []int64) {
	n := rng.Intn(10) + 1
	k := rng.Intn(15)
	if k == 0 {
		k = 1
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		a[i] = int64(rng.Intn(20))
	}
	return n, k, a
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < 100; i++ {
		n, k, a := genTest(rng)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for j, v := range a {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		expect := fmt.Sprint(expected(n, k, a))
		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
