package main

import (
	"bytes"
	"context"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func countDivisors(x int) int {
	count := 1
	for p := 2; p*p <= x; p++ {
		if x%p == 0 {
			e := 0
			for x%p == 0 {
				x /= p
				e++
			}
			count *= (e + 1)
		}
	}
	if x > 1 {
		count *= 2
	}
	return count
}

var validNums []int

func initValid() {
	for x := 1; x <= 1000; x++ {
		if countDivisors(x) <= 7 {
			validNums = append(validNums, x)
		}
	}
}

func bruteForce(arr []int) int {
	n := len(arr)
	best := -1
	for mask := 1; mask < (1 << n); mask++ {
		cnt := bits.OnesCount(uint(mask))
		if best != -1 && cnt >= best {
			continue
		}
		exponents := make(map[int]int)
		for i := 0; i < n; i++ {
			if mask&(1<<i) == 0 {
				continue
			}
			x := arr[i]
			for p := 2; p*p <= x; p++ {
				for x%p == 0 {
					exponents[p]++
					x /= p
				}
			}
			if x > 1 {
				exponents[x]++
			}
		}
		perfect := true
		for _, e := range exponents {
			if e%2 != 0 {
				perfect = false
				break
			}
		}
		if perfect {
			best = cnt
		}
	}
	return best
}

func run(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("time limit exceeded")
		}
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	initValid()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 40; i++ {
		n := rng.Intn(10) + 1
		arr := make([]int, n)
		for j := range arr {
			arr[j] = validNums[rng.Intn(len(validNums))]
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j, v := range arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expected := strconv.Itoa(bruteForce(arr))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
