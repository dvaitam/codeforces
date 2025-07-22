package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCaseE struct {
	n        int
	words    []string
	expected string
}

func solveCase(tc testCaseE) string {
	n := tc.n
	masks := make([]uint32, n)
	used := make([]bool, 24)
	for i, w := range tc.words {
		var m uint32
		for _, ch := range w {
			idx := ch - 'a'
			if idx < 0 || idx >= 24 {
				continue
			}
			used[idx] = true
			m |= 1 << idx
		}
		masks[i] = m
	}
	cntUsed := 0
	for _, u := range used {
		if u {
			cntUsed++
		}
	}
	if cntUsed < 24 {
		return "0"
	}
	size := 1 << 24
	f := make([]uint16, size)
	for _, m := range masks {
		f[m]++
	}
	for bit := 0; bit < 24; bit++ {
		step := 1 << bit
		for mask := 0; mask < size; mask++ {
			if mask&step != 0 {
				f[mask] += f[mask^step]
			}
		}
	}
	var res uint32
	for mask := 0; mask < size; mask++ {
		v := uint32(f[mask])
		res ^= v * v
	}
	return strconv.FormatUint(uint64(res), 10)
}

func run(bin string, input string) (string, error) {
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

func randWord(rng *rand.Rand) string {
	b := []byte{byte(rng.Intn(24) + 'a'), byte(rng.Intn(24) + 'a'), byte(rng.Intn(24) + 'a')}
	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(6) + 1
		words := make([]string, n)
		for j := 0; j < n; j++ {
			words[j] = randWord(rng)
		}
		tc := testCaseE{n: n, words: words}
		expected := solveCase(tc)
		var sb strings.Builder
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for _, w := range words {
			sb.WriteString(w)
			sb.WriteByte('\n')
		}
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
