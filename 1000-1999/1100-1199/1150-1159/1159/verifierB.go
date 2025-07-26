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

func expected(a []int) int {
	n := len(a)
	const INF int64 = 1 << 60
	ans := INF
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			diff := int64(j - i)
			m := int64(a[i])
			if a[j] < a[i] {
				m = int64(a[j])
			}
			v := m / diff
			if v < ans {
				ans = v
				if ans == 0 {
					return 0
				}
			}
		}
	}
	if ans == INF {
		return 0
	}
	return int(ans)
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 2
	a := make([]int, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(1000)
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", a[i])
	}
	sb.WriteString("\n")
	exp := fmt.Sprintf("%d", expected(a))
	return sb.String(), exp
}

func runCase(bin, input, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
