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

func expectedAnswerA(k int, s string) string {
	freq := make([]int, 26)
	for i := 0; i < len(s); i++ {
		freq[s[i]-'a']++
	}
	for j := 0; j < 26; j++ {
		if freq[j]%k != 0 {
			return "-1"
		}
	}
	var part []byte
	for j := 0; j < 26; j++ {
		cnt := freq[j] / k
		for i := 0; i < cnt; i++ {
			part = append(part, byte('a'+j))
		}
	}
	return strings.Repeat(string(part), k)
}

func generateCaseA(rng *rand.Rand) (int, string) {
	k := rng.Intn(5) + 1
	n := rng.Intn(10) + 1
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('a' + rng.Intn(26))
	}
	return k, string(b)
}

func runCaseA(bin string, k int, s string) error {
	input := fmt.Sprintf("%d\n%s\n", k, s)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expected := expectedAnswerA(k, s)
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

	// simple deterministic case
	if err := runCaseA(bin, 1, "a"); err != nil {
		fmt.Fprintln(os.Stderr, "deterministic case failed:", err)
		os.Exit(1)
	}

	for i := 0; i < 100; i++ {
		k, s := generateCaseA(rng)
		if err := runCaseA(bin, k, s); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%d\n%s\n", i+1, err, k, s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
