package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func runBinary(bin, input string) (string, error) {
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
	err := cmd.Run()
	if err != nil {
		return out.String(), fmt.Errorf("%v: %s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func lcs(a, b string) int {
	n := len(a)
	m := len(b)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else if dp[i-1][j] > dp[i][j-1] {
				dp[i][j] = dp[i-1][j]
			} else {
				dp[i][j] = dp[i][j-1]
			}
		}
	}
	return dp[n][m]
}

func expected(base []string, idx []int, s string) int {
	var b strings.Builder
	for _, id := range idx {
		b.WriteString(base[id])
	}
	t := b.String()
	return lcs(s, t)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	if bin == "--" && len(os.Args) >= 3 {
		bin = os.Args[2]
	}
	rand.Seed(3)
	letters := []rune("abc")
	for tcase := 0; tcase < 100; tcase++ {
		n := rand.Intn(3) + 1
		base := make([]string, n)
		for i := 0; i < n; i++ {
			l := rand.Intn(3) + 1
			var sb strings.Builder
			for j := 0; j < l; j++ {
				sb.WriteRune(letters[rand.Intn(len(letters))])
			}
			base[i] = sb.String()
		}
		m := rand.Intn(3) + 1
		idx := make([]int, m)
		for i := 0; i < m; i++ {
			idx[i] = rand.Intn(n)
		}
		slen := rand.Intn(5)
		var sb strings.Builder
		for i := 0; i < slen; i++ {
			sb.WriteRune(letters[rand.Intn(len(letters))])
		}
		s := sb.String()
		// build input
		var in strings.Builder
		fmt.Fprintf(&in, "%d\n", n)
		for i := 0; i < n; i++ {
			fmt.Fprintf(&in, "%s\n", base[i])
		}
		fmt.Fprintf(&in, "%d\n", m)
		for i := 0; i < m; i++ {
			fmt.Fprintf(&in, "%d ", idx[i]+1)
		}
		fmt.Fprint(&in, "\n")
		fmt.Fprintf(&in, "%s\n", s)
		input := in.String()
		exp := expected(base, idx, s)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", tcase+1, err)
			return
		}
		var got int
		fmt.Sscan(out, &got)
		if got != exp {
			fmt.Printf("test %d failed: expected %d got %d\ninput:\n%s", tcase+1, exp, got, input)
			return
		}
	}
	fmt.Println("All tests passed.")
}
