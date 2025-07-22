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

func run(bin string, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func countFriday13(dates []string) int {
	cnt := 0
	for _, s := range dates {
		if len(s) >= 10 && s[8:10] == "13" {
			t, err := time.Parse("2006-01-02", s)
			if err == nil && t.Weekday() == time.Friday {
				cnt++
			}
		}
	}
	return cnt
}

func runCase(bin string, dates []string) error {
	input := fmt.Sprintf("%d\n%s\n", len(dates), strings.Join(dates, "\n"))
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	got, err := strconv.Atoi(strings.TrimSpace(out))
	if err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	exp := countFriday13(dates)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func randomDate(rng *rand.Rand) string {
	year := 1974 + rng.Intn(57) // up to 2030
	month := rng.Intn(12) + 1
	day := rng.Intn(28) + 1
	return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(10) + 1
		dates := make([]string, n)
		for j := 0; j < n; j++ {
			dates[j] = randomDate(rng)
		}
		if err := runCase(bin, dates); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
