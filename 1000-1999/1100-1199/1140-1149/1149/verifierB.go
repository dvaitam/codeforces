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

func possible(S string, A, B, C []byte) bool {
	n := len(S)
	la, lb, lc := len(A), len(B), len(C)
	nxt := make([][26]int, n+2)
	for c := 0; c < 26; c++ {
		nxt[n][c] = n
		nxt[n+1][c] = n
	}
	for i := n - 1; i >= 0; i-- {
		for c := 0; c < 26; c++ {
			nxt[i][c] = nxt[i+1][c]
		}
		nxt[i][S[i]-'a'] = i + 1
	}
	dp := make([][][]int, la+1)
	for i := range dp {
		dp[i] = make([][]int, lb+1)
		for j := range dp[i] {
			dp[i][j] = make([]int, lc+1)
			for k := range dp[i][j] {
				dp[i][j][k] = n + 1
			}
		}
	}
	dp[0][0][0] = 0
	for i := 0; i <= la; i++ {
		for j := 0; j <= lb; j++ {
			for k := 0; k <= lc; k++ {
				pos := dp[i][j][k]
				if pos > n {
					continue
				}
				if i < la {
					p := nxt[pos][A[i]-'a']
					if p <= n && p < dp[i+1][j][k] {
						dp[i+1][j][k] = p
					}
				}
				if j < lb {
					p := nxt[pos][B[j]-'a']
					if p <= n && p < dp[i][j+1][k] {
						dp[i][j+1][k] = p
					}
				}
				if k < lc {
					p := nxt[pos][C[k]-'a']
					if p <= n && p < dp[i][j][k+1] {
						dp[i][j][k+1] = p
					}
				}
			}
		}
	}
	return dp[la][lb][lc] <= n
}

func genCaseB(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	q := rng.Intn(6) + 1
	letters := []byte("abc")
	var sb strings.Builder
	Sbytes := make([]byte, n)
	for i := range Sbytes {
		Sbytes[i] = letters[rng.Intn(len(letters))]
	}
	fmt.Fprintf(&sb, "%d %d\n", n, q)
	sb.WriteString(string(Sbytes))
	sb.WriteByte('\n')
	var A, B, C []byte
	var out strings.Builder
	for t := 0; t < q; t++ {
		if rng.Intn(2) == 0 {
			idx := rng.Intn(3) + 1
			ch := letters[rng.Intn(len(letters))]
			fmt.Fprintf(&sb, "+ %d %c\n", idx, ch)
			if idx == 1 {
				A = append(A, ch)
			} else if idx == 2 {
				B = append(B, ch)
			} else {
				C = append(C, ch)
			}
		} else {
			// ensure chosen string not empty
			idx := rng.Intn(3) + 1
			for (idx == 1 && len(A) == 0) || (idx == 2 && len(B) == 0) || (idx == 3 && len(C) == 0) {
				idx = rng.Intn(3) + 1
			}
			fmt.Fprintf(&sb, "- %d\n", idx)
			if idx == 1 {
				A = A[:len(A)-1]
			} else if idx == 2 {
				B = B[:len(B)-1]
			} else {
				C = C[:len(C)-1]
			}
		}
		if possible(string(Sbytes), A, B, C) {
			out.WriteString("YES\n")
		} else {
			out.WriteString("NO\n")
		}
	}
	return sb.String(), out.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		in, expect := genCaseB(rand.New(rand.NewSource(time.Now().UnixNano() + int64(i))))
		got, err := run(bin, in)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\ninput:\n%soutput:\n%s", i+1, err, in, got)
			return
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%sbut got:\n%s", i+1, in, expect, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
