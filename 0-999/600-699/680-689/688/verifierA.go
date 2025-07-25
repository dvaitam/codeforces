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

type testCaseA struct {
	days []string
}

func generateCase(rng *rand.Rand) (string, testCaseA) {
	n := rng.Intn(10) + 1
	d := rng.Intn(10) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, d))
	days := make([]string, d)
	for i := 0; i < d; i++ {
		var row strings.Builder
		for j := 0; j < n; j++ {
			if rng.Intn(2) == 0 {
				row.WriteByte('0')
			} else {
				row.WriteByte('1')
			}
		}
		s := row.String()
		days[i] = s
		sb.WriteString(s)
		sb.WriteByte('\n')
	}
	return sb.String(), testCaseA{days: days}
}

func expected(tc testCaseA) int {
	maxSeq := 0
	curSeq := 0
	for _, s := range tc.days {
		ok := false
		for i := 0; i < len(s); i++ {
			if s[i] == '0' {
				ok = true
				break
			}
		}
		if ok {
			curSeq++
			if curSeq > maxSeq {
				maxSeq = curSeq
			}
		} else {
			curSeq = 0
		}
	}
	return maxSeq
}

func runCase(bin, input string, tc testCaseA) error {
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
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("invalid output: %v", err)
	}
	want := expected(tc)
	if got != want {
		return fmt.Errorf("expected %d got %d", want, got)
	}
	return nil
}

func main() {
	var bin string
	switch len(os.Args) {
	case 2:
		bin = os.Args[1]
	case 3:
		if os.Args[1] == "--" {
			bin = os.Args[2]
		}
	}
	if bin == "" {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go [--] /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, tc := generateCase(rng)
		if err := runCase(bin, input, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
