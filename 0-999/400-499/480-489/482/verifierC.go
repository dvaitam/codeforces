package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func run(bin, input string) (string, error) {
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

func expected(stringsList []string) float64 {
	n := len(stringsList)
	m := len(stringsList[0])
	size := 1 << m
	c := make([]uint64, size)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			mask := 0
			for w := 0; w < m; w++ {
				if stringsList[i][w] == stringsList[j][w] {
					mask |= 1 << w
				}
			}
			bit := (uint64(1) << uint(i)) | (uint64(1) << uint(j))
			c[mask] |= bit
		}
	}
	for mask := size - 1; mask >= 0; mask-- {
		for b := 0; b < m; b++ {
			if (mask & (1 << b)) == 0 {
				c[mask] |= c[mask|(1<<b)]
			}
		}
		if mask == 0 {
			break
		}
	}
	d := make([][]float64, m+1)
	for i := 0; i <= m; i++ {
		d[i] = make([]float64, i+1)
		d[i][0] = 1
		d[i][i] = 1
		for j := 1; j < i; j++ {
			d[i][j] = d[i-1][j] + d[i-1][j-1]
		}
	}
	ans := 0.0
	for mask := 0; mask < size; mask++ {
		t := bits.OnesCount(uint(mask))
		cnt := bits.OnesCount64(c[mask])
		ans += float64(cnt) / float64(n) * (1.0 / d[m][t])
	}
	return ans
}

func generateCase(rng *rand.Rand) []string {
	m := rng.Intn(5) + 1
	n := rng.Intn(4) + 1
	used := make(map[string]bool)
	res := make([]string, 0, n)
	for len(res) < n {
		b := make([]byte, m)
		for i := 0; i < m; i++ {
			b[i] = byte('a' + rng.Intn(26))
		}
		s := string(b)
		if !used[s] {
			used[s] = true
			res = append(res, s)
		}
	}
	return res
}

func check(stringsList []string, out string) error {
	got, err := strconv.ParseFloat(strings.TrimSpace(out), 64)
	if err != nil {
		return fmt.Errorf("invalid float output")
	}
	want := expected(stringsList)
	if diff := got - want; diff < -1e-6 || diff > 1e-6 {
		return fmt.Errorf("expected %.6f got %.6f", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		strs := generateCase(rng)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", len(strs)))
		for _, s := range strs {
			sb.WriteString(s)
			sb.WriteByte('\n')
		}
		out, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := check(strs, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
