package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const negInf = math.MinInt64 / 4

func solve(l int, a []int64, s string) int64 {
	dp1 := make([][]int64, l)
	for i := range dp1 {
		dp1[i] = make([]int64, l)
		for j := range dp1[i] {
			dp1[i][j] = negInf
		}
	}
	for i := 0; i < l; i++ {
		if a[1] >= 0 {
			dp1[i][i] = a[1]
		}
	}
	for length := 2; length <= l; length++ {
		for i := 0; i+length-1 < l; i++ {
			j := i + length - 1
			if s[i] == s[j] && a[length] >= 0 {
				if length == 2 {
					dp1[i][j] = max64(dp1[i][j], a[length])
				} else if dp1[i+1][j-1] > negInf {
					dp1[i][j] = max64(dp1[i][j], dp1[i+1][j-1]+a[length])
				}
			}
			for k := i; k < j; k++ {
				if dp1[i][k] > negInf && dp1[k+1][j] > negInf {
					dp1[i][j] = max64(dp1[i][j], dp1[i][k]+dp1[k+1][j])
				}
			}
		}
	}
	dp2 := make([]int64, l+1)
	for i := 1; i <= l; i++ {
		dp2[i] = dp2[i-1]
		for j := 0; j < i; j++ {
			if dp1[j][i-1] > negInf {
				if dp2[j]+dp1[j][i-1] > dp2[i] {
					dp2[i] = dp2[j] + dp1[j][i-1]
				}
			}
		}
	}
	return dp2[l]
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		l := rng.Intn(6) + 1
		aVals := make([]int64, l+1)
		for j := 1; j <= l; j++ {
			v := rng.Intn(7) - 1
			if v == 5 {
				aVals[j] = -1
			} else {
				aVals[j] = int64(v)
			}
		}
		sb := strings.Builder{}
		sb.WriteString(fmt.Sprintf("%d\n", l))
		for j := 1; j <= l; j++ {
			sb.WriteString(fmt.Sprintf("%d ", aVals[j]))
		}
		sb.WriteString("\n")
		var letters = []byte("abcde")
		str := make([]byte, l)
		for j := 0; j < l; j++ {
			str[j] = letters[rng.Intn(len(letters))]
		}
		s := string(str)
		sb.WriteString(s)
		sb.WriteString("\n")
		input := sb.String()
		expected := fmt.Sprintf("%d", solve(l, aVals, s))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(expected) != strings.TrimSpace(got) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
