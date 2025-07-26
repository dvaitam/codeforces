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

// Precompute prefix sums up to exceed 1e9 digits
var digitsPrefix []int
var blocksPrefix []int

func initPrefix() {
	const limit = 1000000000
	digitsPrefix = []int{0}
	blocksPrefix = []int{0}
	sumDigits := 0
	sumBlocks := 0
	for i := 1; sumBlocks < limit; i++ {
		l := 0
		for x := i; x > 0; x /= 10 {
			l++
		}
		sumDigits += l
		sumBlocks += sumDigits
		digitsPrefix = append(digitsPrefix, sumDigits)
		blocksPrefix = append(blocksPrefix, sumBlocks)
	}
}

func digitInBlock(n, idx int) byte {
	for d := 1; ; d++ {
		start := 1
		for i := 1; i < d; i++ {
			start *= 10
		}
		end := start*10 - 1
		if end > n {
			end = n
		}
		if end < start {
			continue
		}
		cnt := end - start + 1
		total := cnt * d
		if idx <= total {
			number := start + (idx-1)/d
			digitIdx := (idx - 1) % d
			digits := make([]byte, d)
			for i := d - 1; i >= 0; i-- {
				digits[i] = byte('0' + number%10)
				number /= 10
			}
			return digits[digitIdx]
		}
		idx -= total
	}
}

func query(k int) byte {
	l, r := 1, len(blocksPrefix)-1
	for l < r {
		m := (l + r) / 2
		if blocksPrefix[m] < k {
			l = m + 1
		} else {
			r = m
		}
	}
	block := l
	prev := blocksPrefix[block-1]
	idx := k - prev
	return digitInBlock(block, idx)
}

func genCase(rng *rand.Rand) (string, []int) {
	q := rng.Intn(5) + 1
	ks := make([]int, q)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for i := 0; i < q; i++ {
		ks[i] = rng.Intn(1000000000) + 1
		sb.WriteString(fmt.Sprintf("%d\n", ks[i]))
	}
	return sb.String(), ks
}

func expected(ks []int) string {
	var sb strings.Builder
	for _, k := range ks {
		b := query(k)
		sb.WriteByte(b)
		sb.WriteByte('\n')
	}
	return strings.TrimSpace(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE1.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	initPrefix()
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
