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

func runBinary(bin, input string) (string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solveCase(n int, d int64, a []int64) int {
	pref := make([]int64, n+1)
	for i := 0; i < n; i++ {
		pref[i+1] = pref[i] + a[i]
	}
	maxAfter := make([]int64, n+1)
	maxAfter[n] = pref[n]
	for i := n - 1; i >= 0; i-- {
		if pref[i] > maxAfter[i+1] {
			maxAfter[i] = pref[i]
		} else {
			maxAfter[i] = maxAfter[i+1]
		}
	}
	shift := int64(0)
	count := 0
	for i := 1; i <= n; i++ {
		cur := pref[i] + shift
		if cur > d {
			return -1
		}
		if a[i-1] == 0 && cur < 0 {
			limit := d - maxAfter[i]
			if limit < -pref[i] {
				return -1
			}
			if shift < limit {
				shift = limit
				count++
			}
			cur = pref[i] + shift
			if cur < 0 || cur > d {
				return -1
			}
		}
	}
	return count
}

func generateCase(r *rand.Rand) (string, string) {
	n := r.Intn(20) + 1
	d := int64(r.Intn(50) + 1)
	a := make([]int64, n)
	for i := range a {
		a[i] = int64(r.Intn(21) - 10)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, d))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(a[i]))
	}
	sb.WriteByte('\n')
	expect := fmt.Sprint(solveCase(n, d, a))
	return sb.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(r)
		out, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected: %s\nGot: %s\n", i+1, in, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
