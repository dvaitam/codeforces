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

func solveCase(input string) string {
	reader := strings.NewReader(input)
	weights := make([]int64, 26)
	for i := 0; i < 26; i++ {
		fmt.Fscan(reader, &weights[i])
	}
	var s string
	fmt.Fscan(reader, &s)
	n := len(s)
	prefix := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		prefix[i] = prefix[i-1] + weights[s[i-1]-'a']
	}
	cnt := make([]map[int64]int64, 26)
	for i := range cnt {
		cnt[i] = make(map[int64]int64)
	}
	var ans int64
	for i := 1; i <= n; i++ {
		c := s[i-1] - 'a'
		target := prefix[i-1]
		if v, ok := cnt[c][target]; ok {
			ans += v
		}
		cnt[c][prefix[i]]++
	}
	return fmt.Sprint(ans)
}

func generateCase(rng *rand.Rand) string {
	weights := make([]int64, 26)
	for i := 0; i < 26; i++ {
		weights[i] = rng.Int63n(5) + 1
	}
	n := rng.Intn(20) + 1
	var sb strings.Builder
	for i := 0; i < 26; i++ {
		fmt.Fprintf(&sb, "%d ", weights[i])
	}
	sb.WriteByte('\n')
	letters := make([]byte, n)
	for i := range letters {
		letters[i] = byte('a' + rng.Intn(26))
	}
	sb.WriteString(string(letters))
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(bin, input, expected string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		inp := generateCase(rng)
		exp := solveCase(inp)
		if err := runCase(bin, inp, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, inp)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
