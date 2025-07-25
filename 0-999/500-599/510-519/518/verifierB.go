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

func expected(s, t string) string {
	var cntLow [26]int
	var cntUp [26]int
	for _, ch := range t {
		if ch >= 'a' && ch <= 'z' {
			cntLow[ch-'a']++
		} else if ch >= 'A' && ch <= 'Z' {
			cntUp[ch-'A']++
		}
	}
	yay := 0
	whoops := 0
	unmatched := make([]rune, 0, len(s))
	for _, ch := range s {
		if ch >= 'a' && ch <= 'z' {
			idx := ch - 'a'
			if cntLow[idx] > 0 {
				yay++
				cntLow[idx]--
			} else {
				unmatched = append(unmatched, ch)
			}
		} else {
			idx := ch - 'A'
			if cntUp[idx] > 0 {
				yay++
				cntUp[idx]--
			} else {
				unmatched = append(unmatched, ch)
			}
		}
	}
	for _, ch := range unmatched {
		if ch >= 'a' && ch <= 'z' {
			idx := ch - 'a'
			if cntUp[idx] > 0 {
				whoops++
				cntUp[idx]--
			}
		} else {
			idx := ch - 'A'
			if cntLow[idx] > 0 {
				whoops++
				cntLow[idx]--
			}
		}
	}
	return fmt.Sprintf("%d %d", yay, whoops)
}

func generateCase(rng *rand.Rand) (string, string, string) {
	n := rng.Intn(20) + 1
	m := n + rng.Intn(20)
	sbS := make([]byte, n)
	sbT := make([]byte, m)
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for i := 0; i < n; i++ {
		sbS[i] = letters[rng.Intn(len(letters))]
	}
	for i := 0; i < m; i++ {
		sbT[i] = letters[rng.Intn(len(letters))]
	}
	s := string(sbS)
	t := string(sbT)
	return s, t, expected(s, t)
}

func runCase(bin, s, t, exp string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	input := fmt.Sprintf("%s\n%s\n", s, t)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := []struct{ s, t, exp string }{}
	for i := 0; i < 100; i++ {
		s, t, e := generateCase(rng)
		cases = append(cases, struct{ s, t, exp string }{s, t, e})
	}

	for i, tc := range cases {
		if err := runCase(bin, tc.s, tc.t, tc.exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\n%s\n", i+1, err, tc.s, tc.t)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
