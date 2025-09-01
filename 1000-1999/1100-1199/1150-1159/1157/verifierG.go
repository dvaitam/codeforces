package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const numTestsG = 100

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: verifierG.go <binary>")
		os.Exit(1)
	}
	binPath, cleanup, err := prepareBinary1157G(os.Args[1])
	if err != nil {
		fmt.Println("compile error:", err)
		os.Exit(1)
	}
	if cleanup != nil {
		defer cleanup()
	}
	r := rand.New(rand.NewSource(1))
	for t := 1; t <= numTestsG; t++ {
		n := r.Intn(4) + 1
		m := r.Intn(4) + 1
		a := make([][]int, n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < n; i++ {
			a[i] = make([]int, m)
			for j := 0; j < m; j++ {
				a[i][j] = r.Intn(2)
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(fmt.Sprintf("%d", a[i][j]))
			}
			sb.WriteByte('\n')
		}
		input := sb.String()
		out, err := run1157G(binPath, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", t, err)
			os.Exit(1)
		}
		first := firstNonEmptyLine(out)
		if strings.ToUpper(first) == "YES" {
			if !validateGOutput(n, m, a, out) {
				ok, rr, cc := existsSolution1157G(a)
				expected := "NO"
				if ok {
					expected = fmt.Sprintf("YES\n%s\n%s", rr, cc)
				}
				fmt.Printf("test %d failed\ninput:%sexpected:%s got:%s\n", t, input, expected, strings.TrimSpace(out))
				os.Exit(1)
			}
			continue
		}
		// Treat anything else as NO
		ok, rr, cc := existsSolution1157G(a)
		if ok {
			expected := fmt.Sprintf("YES\n%s\n%s", rr, cc)
			fmt.Printf("test %d failed\ninput:%sexpected:%s got:%s\n", t, input, expected, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Println("OK")
}

func prepareBinary1157G(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp := filepath.Join(os.TempDir(), "verify_binG")
		cmd := exec.Command("go", "build", "-o", tmp, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, string(out))
		}
		return tmp, func() { os.Remove(tmp) }, nil
	}
	return path, nil, nil
}

func run1157G(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return buf.String(), err
}

func solveG(a [][]int) string {
	n := len(a)
	m := len(a[0])
	F := make([][2][2]bool, m-1)
	for i := 0; i < n; i++ {
		for j := 0; j+1 < m; j++ {
			x := a[i][j]
			y := a[i][j+1]
			u := 1 - x
			v := y
			F[j][u][v] = true
		}
	}
	dp := make([][2]bool, m)
	parent := make([][2]int, m)
	dp[0][0] = true
	dp[0][1] = true
	for j := 0; j+1 < m; j++ {
		for b := 0; b < 2; b++ {
			if !dp[j][b] {
				continue
			}
			for nb := 0; nb < 2; nb++ {
				if !F[j][b][nb] {
					if !dp[j+1][nb] {
						dp[j+1][nb] = true
						parent[j+1][nb] = b
					}
				}
			}
		}
	}
	c := make([]int, m)
	if m > 0 {
		if dp[m-1][0] {
			c[m-1] = 0
		} else if dp[m-1][1] {
			c[m-1] = 1
		} else {
			return "NO"
		}
		for j := m - 1; j > 0; j-- {
			c[j-1] = parent[j][c[j]]
		}
	}
	rarr := make([]int, n)
	if n > 0 {
		rarr[0] = 0
		for i := 0; i+1 < n; i++ {
			d := a[i][m-1] ^ c[m-1] ^ a[i+1][0] ^ c[0]
			if d == 0 {
				rarr[i+1] = rarr[i]
			} else {
				rarr[i+1] = 0
			}
		}
	}
	var sb strings.Builder
	sb.WriteString("YES\n")
	for i := 0; i < n; i++ {
		sb.WriteByte(byte('0' + rarr[i]))
	}
	sb.WriteByte('\n')
	for j := 0; j < m; j++ {
		sb.WriteByte(byte('0' + c[j]))
	}
	return strings.TrimSpace(sb.String())
}

func validateGOutput(n, m int, a [][]int, out string) bool {
	// Expect format:
	// YES\n
	// row bits of length n (0/1) possibly with spaces\n
	// col bits of length m (0/1) possibly with spaces
	lines := strings.Split(out, "\n")
	pick := func(k int) (string, int) {
		for k < len(lines) {
			ln := strings.TrimSpace(lines[k])
			k++
			if ln != "" {
				return ln, k
			}
		}
		return "", k
	}
	ln, idx := pick(0)
	if strings.ToUpper(strings.TrimSpace(ln)) != "YES" {
		return false
	}
	rowLine, idx := pick(idx)
	colLine, _ := pick(idx)
	rStr := strings.ReplaceAll(strings.TrimSpace(rowLine), " ", "")
	cStr := strings.ReplaceAll(strings.TrimSpace(colLine), " ", "")
	if len(rStr) != n || len(cStr) != m {
		return false
	}
	r := make([]int, n)
	c := make([]int, m)
	for i := 0; i < n; i++ {
		if rStr[i] != '0' && rStr[i] != '1' {
			return false
		}
		r[i] = int(rStr[i] - '0')
	}
	for j := 0; j < m; j++ {
		if cStr[j] != '0' && cStr[j] != '1' {
			return false
		}
		c[j] = int(cStr[j] - '0')
	}
	// Build flattened sequence after flips and ensure nondecreasing
	prev := 0
	first := true
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			v := a[i][j] ^ r[i] ^ c[j]
			if first {
				prev = v
				first = false
				continue
			}
			if prev == 1 && v == 0 {
				return false
			}
			prev = v
		}
	}
	return true
}

// firstNonEmptyLine returns the first non-empty trimmed line (or "").
func firstNonEmptyLine(s string) string {
	for _, ln := range strings.Split(s, "\n") {
		l := strings.TrimSpace(ln)
		if l != "" {
			return l
		}
	}
	return ""
}

// existsSolution1157G brute-forces r and c to find a valid sorted configuration.
// Returns true along with bit strings r and c if found.
func existsSolution1157G(a [][]int) (bool, string, string) {
	n := len(a)
	if n == 0 {
		return true, "", ""
	}
	m := len(a[0])
	isSorted := func(rmask, cmask int) bool {
		first := true
		prev := 0
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				v := a[i][j] ^ ((rmask >> i) & 1) ^ ((cmask >> j) & 1)
				if first {
					prev = v
					first = false
					continue
				}
				if prev == 1 && v == 0 {
					return false
				}
				prev = v
			}
		}
		return true
	}
	for rmask := 0; rmask < (1 << n); rmask++ {
		for cmask := 0; cmask < (1 << m); cmask++ {
			if isSorted(rmask, cmask) {
				r := make([]byte, n)
				for i := 0; i < n; i++ {
					if (rmask>>i)&1 == 1 {
						r[i] = '1'
					} else {
						r[i] = '0'
					}
				}
				c := make([]byte, m)
				for j := 0; j < m; j++ {
					if (cmask>>j)&1 == 1 {
						c[j] = '1'
					} else {
						c[j] = '0'
					}
				}
				return true, string(r), string(c)
			}
		}
	}
	return false, "", ""
}
