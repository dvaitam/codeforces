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

func can(a []int, k int, D int) bool {
	mNeeded := (D + 1) / 2
	cnt := 0
	for _, v := range a {
		if v < mNeeded {
			cnt++
		}
	}
	if cnt > k {
		return false
	}
	rem := k - cnt
	l := 0
	bad := 0
	for r, v := range a {
		if v >= mNeeded && v < D {
			bad++
		}
		for bad > rem {
			y := a[l]
			if y >= mNeeded && y < D {
				bad--
			}
			l++
		}
		if r-l+1 >= 2 {
			return true
		}
	}
	return false
}

func solve(a []int, k int) int {
	lo, hi := 1, 1000000000
	for lo < hi {
		mid := (lo + hi + 1) / 2
		if can(a, k, mid) {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	return lo
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 2
	k := rng.Intn(n) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(50) + 1
	}
	ans := solve(arr, k)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteString("\n")
	return sb.String(), fmt.Sprint(ans)
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
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
