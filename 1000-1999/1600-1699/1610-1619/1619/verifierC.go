package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("time limit")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

func solveCase(aStr, sStr string) string {
	i := len(aStr) - 1
	j := len(sStr) - 1
	ans := make([]byte, 0, len(sStr))
	for i >= 0 {
		if j < 0 {
			return "-1"
		}
		aDigit := aStr[i] - '0'
		sDigit := sStr[j] - '0'
		if sDigit >= aDigit {
			ans = append(ans, '0'+(sDigit-aDigit))
			j--
		} else {
			if j == 0 {
				return "-1"
			}
			twoDigit := (sStr[j-1]-'0')*10 + sDigit
			diff := twoDigit - aDigit
			if diff < 0 || diff > 9 {
				return "-1"
			}
			ans = append(ans, '0'+diff)
			j -= 2
		}
		i--
	}
	for j >= 0 {
		ans = append(ans, sStr[j])
		j--
	}
	for len(ans) > 1 && ans[len(ans)-1] == '0' {
		ans = ans[:len(ans)-1]
	}
	for l, r := 0, len(ans)-1; l < r; l, r = l+1, r-1 {
		ans[l], ans[r] = ans[r], ans[l]
	}
	return string(ans)
}

func generateCase(rng *rand.Rand) (string, string) {
	t := rng.Intn(5) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	var out strings.Builder
	for i := 0; i < t; i++ {
		a := rng.Int63n(1_000_000_000_000_000_000) + 1
		s := a + rng.Int63n(1_000_000_000_000_000_000-a) + 1
		fmt.Fprintf(&sb, "%d %d\n", a, s)
		fmt.Fprintf(&out, "%s\n", solveCase(fmt.Sprint(a), fmt.Sprint(s)))
	}
	return sb.String(), strings.TrimSpace(out.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runBinary(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
