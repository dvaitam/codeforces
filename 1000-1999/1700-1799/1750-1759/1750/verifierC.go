package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type op struct{ l, r int }

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveC(r *bufio.Reader) string {
	var T int
	if _, err := fmt.Fscan(r, &T); err != nil {
		return ""
	}
	var out strings.Builder
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(r, &n)
		var s1, s2 string
		fmt.Fscan(r, &s1)
		fmt.Fscan(r, &s2)
		fl1, fl2 := false, false
		for i := 0; i < n; i++ {
			if s1[i] == s2[i] {
				fl1 = true
			} else {
				fl2 = true
			}
		}
		if fl1 && fl2 {
			out.WriteString("NO\n")
			continue
		}
		var ops []op
		k := 0
		for i := 0; i < n; i++ {
			if s1[i] == '1' {
				ops = append(ops, op{i + 1, i + 1})
				if i >= 1 {
					k++
				}
			}
		}
		if ((int(s2[0]-'0') + k) & 1) != 0 {
			ops = append(ops, op{1, n - 1})
			ops = append(ops, op{n, n})
			ops = append(ops, op{1, n})
		}
		out.WriteString("YES\n")
		out.WriteString(fmt.Sprintf("%d\n", len(ops)))
		for _, e := range ops {
			out.WriteString(fmt.Sprintf("%d %d\n", e.l, e.r))
		}
	}
	return strings.TrimSpace(out.String())
}

func generateCaseC(rng *rand.Rand) string {
	n := rng.Intn(4) + 2 // 2..5
	var s1, s2 strings.Builder
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			s1.WriteByte('0')
		} else {
			s1.WriteByte('1')
		}
	}
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			s2.WriteByte('0')
		} else {
			s2.WriteByte('1')
		}
	}
	return fmt.Sprintf("1\n%d\n%s\n%s\n", n, s1.String(), s2.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseC(rng)
		expect := solveC(bufio.NewReader(strings.NewReader(tc)))
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
