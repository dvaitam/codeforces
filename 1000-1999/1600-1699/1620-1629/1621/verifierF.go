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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func solveCaseF(n int, a, b, c int64, s string) int64 {
	cnt00 := 0
	cnt11 := 0
	for i := 0; i+1 < n; i++ {
		if s[i] == '0' && s[i+1] == '0' {
			cnt00++
		}
		if s[i] == '1' && s[i+1] == '1' {
			cnt11++
		}
	}
	singles := 0
	i := 0
	for i < n {
		if s[i] == '0' {
			j := i
			for j < n && s[j] == '0' {
				j++
			}
			if j-i == 1 && i > 0 && j < n {
				singles++
			}
			i = j
		} else {
			i++
		}
	}
	ans := int64(min(cnt00, cnt11)) * (a + b)
	if cnt00 > cnt11 {
		ans += a
	} else if cnt11 > cnt00 {
		ans += b
	}
	if b > c {
		ans += (b - c) * int64(singles)
	}
	return ans
}

func runCandidate(bin, input string) (string, error) {
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

func generateCaseF(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 1
	a := int64(rng.Intn(5) + 1)
	b := int64(rng.Intn(5) + 1)
	c := int64(rng.Intn(5) + 1)
	s := make([]byte, n)
	for i := range s {
		if rng.Intn(2) == 0 {
			s[i] = '0'
		} else {
			s[i] = '1'
		}
	}
	input := fmt.Sprintf("1\n%d %d %d %d\n%s\n", n, a, b, c, string(s))
	exp := fmt.Sprintf("%d", solveCaseF(n, a, b, c, string(s)))
	return input, exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseF(rng)
		got, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
