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

func expected(n int, langs []int, m int, audio []int, sub []int) string {
	freq := make(map[int]int)
	for _, v := range langs {
		freq[v]++
	}
	bestIdx := 0
	bestAudio := -1
	bestSub := -1
	for i := 0; i < m; i++ {
		a := freq[audio[i]]
		s := freq[sub[i]]
		if a > bestAudio || (a == bestAudio && s > bestSub) {
			bestAudio = a
			bestSub = s
			bestIdx = i + 1
		}
	}
	return fmt.Sprintf("%d", bestIdx)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	langs := make([]int, n)
	for i := range langs {
		langs[i] = rng.Intn(50) + 1
	}
	m := rng.Intn(20) + 1
	audio := make([]int, m)
	sub := make([]int, m)
	for i := 0; i < m; i++ {
		audio[i] = rng.Intn(50) + 1
		for {
			sub[i] = rng.Intn(50) + 1
			if sub[i] != audio[i] {
				break
			}
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range langs {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", m))
	for i, v := range audio {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, v := range sub {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	inp := sb.String()
	exp := expected(n, langs, m, audio, sub)
	return inp, exp
}

func runProg(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("%v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := generateCase(rng)
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\noutput:\n%s", i+1, err, got)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed:\nexpected: %s\n got: %s\n", i+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
