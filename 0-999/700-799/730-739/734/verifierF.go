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

func expected(b, c []int64) string {
	n := len(b)
	var sum int64
	for i := 0; i < n; i++ {
		sum += b[i] + c[i]
	}
	denom := int64(2 * n)
	if sum%denom != 0 {
		return "-1"
	}
	S := sum / denom
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		tmp := b[i] + c[i] - S
		if tmp%int64(n) != 0 {
			return "-1"
		}
		ai := tmp / int64(n)
		if ai < 0 {
			return "-1"
		}
		a[i] = ai
	}
	const maxbit = 31
	cnt := make([]int64, maxbit+1)
	for _, ai := range a {
		for k := 0; k <= maxbit; k++ {
			if (ai>>k)&1 == 1 {
				cnt[k]++
			}
		}
	}
	for i, ai := range a {
		var b2, c2 int64
		for k := 0; k <= maxbit; k++ {
			bit := int64(1) << k
			if (ai>>k)&1 == 1 {
				b2 += cnt[k] * bit
				c2 += int64(n) * bit
			} else {
				c2 += cnt[k] * bit
			}
		}
		if b2 != b[i] || c2 != c[i] {
			return "-1"
		}
	}
	var sb strings.Builder
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	return sb.String()
}

func generateValidCase(rng *rand.Rand) (int, []int64, []int64, string) {
	n := rng.Intn(5) + 1
	a := make([]int64, n)
	for i := range a {
		a[i] = int64(rng.Intn(10))
	}
	const maxbit = 31
	cnt := make([]int64, maxbit+1)
	for _, ai := range a {
		for k := 0; k <= maxbit; k++ {
			if (ai>>k)&1 == 1 {
				cnt[k]++
			}
		}
	}
	b := make([]int64, n)
	c := make([]int64, n)
	for i, ai := range a {
		var bi, ci int64
		for k := 0; k <= maxbit; k++ {
			bit := int64(1) << k
			if (ai>>k)&1 == 1 {
				bi += cnt[k] * bit
				ci += int64(n) * bit
			} else {
				ci += cnt[k] * bit
			}
		}
		b[i] = bi
		c[i] = ci
	}
	var sb strings.Builder
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	expect := sb.String()
	return n, b, c, expect
}

func generateCase(rng *rand.Rand) (string, string) {
	if rng.Float64() < 0.6 {
		n, b, c, expect := generateValidCase(rng)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		for i, v := range c {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		input := sb.String()
		return input, expect
	}
	n := rng.Intn(5) + 1
	b := make([]int64, n)
	c := make([]int64, n)
	for i := 0; i < n; i++ {
		b[i] = int64(rng.Intn(50))
		c[i] = int64(rng.Intn(50))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	for i, v := range c {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	input := sb.String()
	exp := expected(b, c)
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
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
