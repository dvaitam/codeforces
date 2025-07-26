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

func expectedAnswerA(n, m int, games, bills []int) int {
	count := 0
	j := 0
	for i := 0; i < n && j < m; i++ {
		if games[i] <= bills[j] {
			count++
			j++
		}
	}
	return count
}

func genCaseA(rng *rand.Rand) (int, int, []int, []int) {
	n := rng.Intn(20) + 1
	m := rng.Intn(20) + 1
	games := make([]int, n)
	bills := make([]int, m)
	for i := range games {
		games[i] = rng.Intn(1000) + 1
	}
	for i := range bills {
		bills[i] = rng.Intn(1000) + 1
	}
	return n, m, games, bills
}

func runCaseA(bin string, n, m int, games, bills []int) error {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	sb.WriteString(strconv.Itoa(m))
	sb.WriteByte('\n')
	for i, v := range games {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, v := range bills {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	input := sb.String()
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	got := strings.TrimSpace(out.String())
	expected := strconv.Itoa(expectedAnswerA(n, m, games, bills))
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, m, g, b := genCaseA(rng)
		if err := runCaseA(bin, n, m, g, b); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%v %v\n%v\n%v\n", i+1, err, n, m, g, b)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
