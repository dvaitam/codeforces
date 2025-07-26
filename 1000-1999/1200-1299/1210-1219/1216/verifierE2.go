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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func sumLen(n int64) int64 {
	var res int64
	pow10 := int64(1)
	d := int64(1)
	for pow10 <= n {
		nxt := pow10*10 - 1
		if nxt > n {
			nxt = n
		}
		res += (nxt - pow10 + 1) * d
		pow10 *= 10
		d++
	}
	return res
}

func sumXLen(n int64) int64 {
	var res int64
	pow10 := int64(1)
	d := int64(1)
	for pow10 <= n {
		nxt := pow10*10 - 1
		if nxt > n {
			nxt = n
		}
		cnt := nxt - pow10 + 1
		res += d * (pow10 + nxt) * cnt / 2
		pow10 *= 10
		d++
	}
	return res
}

func prefix(n int64) int64 {
	if n <= 0 {
		return 0
	}
	return (n+1)*sumLen(n) - sumXLen(n)
}

func digitInRange(n, k int64) int {
	pow10 := int64(1)
	d := int64(1)
	for {
		nxt := pow10*10 - 1
		if nxt > n {
			nxt = n
		}
		cnt := nxt - pow10 + 1
		total := cnt * d
		if k > total {
			k -= total
		} else {
			idx := (k - 1) / d
			num := pow10 + idx
			pos := (k - 1) % d
			str := fmt.Sprintf("%d", num)
			return int(str[pos] - '0')
		}
		if nxt == n {
			break
		}
		pow10 *= 10
		d++
	}
	return 0
}

func query(k int64) int {
	l, r := int64(1), int64(1e9)
	for l < r {
		m := (l + r) / 2
		if prefix(m) >= k {
			r = m
		} else {
			l = m + 1
		}
	}
	n := l
	k -= prefix(n - 1)
	return digitInRange(n, k)
}

func genCase(rng *rand.Rand) (string, []int64) {
	q := rng.Intn(5) + 1
	ks := make([]int64, q)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for i := 0; i < q; i++ {
		ks[i] = rng.Int63n(1_000_000_000_000) + 1
		sb.WriteString(fmt.Sprintf("%d\n", ks[i]))
	}
	return sb.String(), ks
}

func expected(ks []int64) string {
	var sb strings.Builder
	for _, k := range ks {
		sb.WriteString(fmt.Sprintf("%d\n", query(k)))
	}
	return strings.TrimSpace(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE2.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, ks := genCase(rng)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		exp := expected(ks)
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected:\n%s\ngot:\n%s\ninput:\n%s", i+1, exp, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
