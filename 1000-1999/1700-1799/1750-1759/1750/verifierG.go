package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const refSource = `package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	var m_in int64
	if _, err := fmt.Fscan(reader, &n, &m_in); err != nil {
		return
	}
	m := m_in

	p := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &p[i])
	}

	C := make([][]int64, n+1)
	for i := 0; i <= n; i++ {
		C[i] = make([]int64, n+1)
		C[i][0] = 1
		for j := 1; j <= i; j++ {
			C[i][j] = (C[i-1][j-1] + C[i-1][j]) % m
		}
	}

	fact := make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = (fact[i-1] * int64(i)) % m
	}

	f := make([][]int64, n+1)
	for i := 0; i <= n; i++ {
		f[i] = make([]int64, n+1)
	}
	for p_val := 0; p_val <= n; p_val++ {
		f[0][p_val] = fact[p_val]
	}
	for d := 1; d <= n; d++ {
		for p_val := d; p_val <= n; p_val++ {
			f[d][p_val] = (f[d-1][p_val] - f[d-1][p_val-1] + m) % m
		}
	}

	ans := make([]int64, n+1)
	in_S := make([]bool, n+2)
	for i := 1; i <= n; i++ {
		in_S[i] = true
	}
	c := n - 1
	c_prefix := 0

	for i := 0; i < n; i++ {
		L := n - i
		cnt := [2][2]int64{}

		for v := 1; v < p[i]; v++ {
			if in_S[v] {
				is_match := 0
				if i > 0 && v == p[i-1]+1 {
					is_match = 1
				}
				has_prev := 0
				if in_S[v-1] {
					has_prev = 1
				}
				cnt[is_match][has_prev]++
			}
		}

		for is_match := 0; is_match < 2; is_match++ {
			for has_prev := 0; has_prev < 2; has_prev++ {
				if cnt[is_match][has_prev] > 0 {
					FL := c_prefix + is_match
					C_star := c - has_prev
					for j := 0; j <= C_star; j++ {
						ways := (C[C_star][j] * f[C_star-j][L-1-j]) % m
						ways = (ways * cnt[is_match][has_prev]) % m
						ans[FL+j] = (ans[FL+j] + ways) % m
					}
				}
			}
		}

		in_S[p[i]] = false
		if in_S[p[i]-1] {
			c--
		}
		if in_S[p[i]+1] {
			c--
		}
		if i > 0 && p[i] == p[i-1]+1 {
			c_prefix++
		}
	}

	for k := 1; k <= n; k++ {
		if k > 1 {
			fmt.Print(" ")
		}
		fmt.Print(ans[n-k])
	}
	fmt.Println()
}
`

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

func buildRef() (string, func()) {
	tmpDir, err := os.MkdirTemp("", "ref1750G")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	srcPath := filepath.Join(tmpDir, "ref.go")
	os.WriteFile(srcPath, []byte(refSource), 0644)
	binPath := filepath.Join(tmpDir, "ref")
	cmd := exec.Command("go", "build", "-o", binPath, srcPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build ref: %v\n%s\n", err, string(out))
		os.Exit(1)
	}
	return binPath, func() { os.RemoveAll(tmpDir) }
}

func generateCaseG(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	perm := rng.Perm(n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", perm[i]+1)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	refBin, cleanup := buildRef()
	defer cleanup()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseG(rng)
		expect, err := runBinary(refBin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
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
