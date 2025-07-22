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

func isLucky(x int) bool {
	if x <= 0 {
		return false
	}
	for x > 0 {
		d := x % 10
		if d != 4 && d != 7 {
			return false
		}
		x /= 10
	}
	return true
}

func noCommonLucky(a1, a2 []int) bool {
	seen := make(map[int]bool)
	for _, v := range a1 {
		if isLucky(v) {
			seen[v] = true
		}
	}
	for _, v := range a2 {
		if isLucky(v) && seen[v] {
			return false
		}
	}
	return true
}

func solve(arr []int) int64 {
	n := len(arr)
	var ans int64
	for l1 := 0; l1 < n; l1++ {
		for r1 := l1; r1 < n; r1++ {
			for l2 := r1 + 1; l2 < n; l2++ {
				for r2 := l2; r2 < n; r2++ {
					if noCommonLucky(arr[l1:r1+1], arr[l2:r2+1]) {
						ans++
					}
				}
			}
		}
	}
	return ans
}

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

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 2
	arr := make([]int, n)
	luckyVals := []int{4, 7, 44, 47, 74, 77}
	for i := range arr {
		if rng.Intn(3) == 0 {
			arr[i] = luckyVals[rng.Intn(len(luckyVals))]
		} else {
			arr[i] = rng.Intn(100) + 1
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	exp := fmt.Sprintf("%d", solve(arr))
	return sb.String(), exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
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
