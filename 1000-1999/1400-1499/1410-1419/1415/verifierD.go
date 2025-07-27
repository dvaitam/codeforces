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

func run(bin string, input string) (string, error) {
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveD(a []int) int {
	n := len(a)
	if n > 130 {
		return 1
	}
	px := make([]int, n+1)
	for i := 1; i <= n; i++ {
		px[i] = px[i-1] ^ a[i-1]
	}
	const INF = int(1e9)
	ans := INF
	for l := 1; l <= n; l++ {
		for r := l + 1; r <= n; r++ {
			for i := l; i < r; i++ {
				left := px[i] ^ px[l-1]
				right := px[r] ^ px[i]
				if left > right {
					ops := r - l - 1
					if ops < ans {
						ans = ops
					}
				}
			}
		}
	}
	if ans == INF {
		return -1
	}
	return ans
}

func runCase(bin string, arr []int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(arr)))
	for i := 0; i < len(arr); i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", arr[i]))
	}
	sb.WriteByte('\n')
	expect := fmt.Sprintf("%d", solveD(arr))
	out, err := run(bin, sb.String())
	if err != nil {
		return err
	}
	if out != expect {
		return fmt.Errorf("expected %s got %s", expect, out)
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
	total := 0
	// edge cases
	if err := runCase(bin, []int{1, 2}); err != nil {
		fmt.Fprintf(os.Stderr, "case %d failed: %v\n", total+1, err)
		os.Exit(1)
	}
	total++
	if err := runCase(bin, []int{1, 1, 1}); err != nil {
		fmt.Fprintf(os.Stderr, "case %d failed: %v\n", total+1, err)
		os.Exit(1)
	}
	total++
	for total < 100 {
		n := rng.Intn(20) + 2
		arr := make([]int, n)
		for i := range arr {
			arr[i] = rng.Intn(1000)
		}
		if err := runCase(bin, arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", total+1, err)
			os.Exit(1)
		}
		total++
	}
	fmt.Printf("All %d tests passed\n", total)
}
