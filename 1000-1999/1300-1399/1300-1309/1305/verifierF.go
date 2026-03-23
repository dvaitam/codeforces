package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Embedded oracle solver for 1305F
func solveOracle(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var n int
	fmt.Fscan(reader, &n)

	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(n, func(i, j int) {
		a[i], a[j] = a[j], a[i]
	})

	var ans int64 = 0
	for _, x := range a {
		if x%2 != 0 {
			ans++
		}
	}

	tested := make(map[int64]bool)
	tested[2] = true

	cost := func(p int64, currentMin int64) int64 {
		var c int64 = 0
		for _, x := range a {
			if x < p {
				c += p - x
			} else {
				rem := x % p
				if rem < p-rem {
					c += rem
				} else {
					c += p - rem
				}
			}
			if c >= currentMin {
				return currentMin
			}
		}
		return c
	}

	limit := 40
	if limit > n {
		limit = n
	}

	for i := 0; i < limit; i++ {
		for d := int64(-1); d <= 1; d++ {
			val := a[i] + d
			if val <= 1 {
				continue
			}
			for val%2 == 0 {
				val /= 2
			}
			for div := int64(3); div*div <= val; div += 2 {
				if val%div == 0 {
					if !tested[div] {
						tested[div] = true
						ans = oracleMin(ans, cost(div, ans))
					}
					for val%div == 0 {
						val /= div
					}
				}
			}
			if val > 1 && !tested[val] {
				tested[val] = true
				ans = oracleMin(ans, cost(val, ans))
			}
		}
	}

	return strconv.FormatInt(ans, 10)
}

func oracleMin(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(6) + 2
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", rng.Intn(1000)+1))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(bin, input string) (string, error) {
	expected := strings.TrimSpace(solveOracle(input))

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr: %s", err, errb.String())
	}
	result := strings.TrimSpace(out.String())
	if result != expected {
		return "", fmt.Errorf("expected %q got %q", expected, result)
	}
	return result, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		if _, err := runCase(bin, input); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
