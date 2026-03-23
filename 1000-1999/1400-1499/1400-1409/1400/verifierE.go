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

const oracleSource = `package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReaderSize(os.Stdin, 1<<20)

	var n int
	fmt.Fscan(in, &n)

	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	var solve func(int, int, int64) int64
	solve = func(l, r int, base int64) int64 {
		if l > r {
			return 0
		}
		width := int64(r - l + 1)
		mn := a[l]
		for i := l + 1; i <= r; i++ {
			if a[i] < mn {
				mn = a[i]
			}
		}
		if mn-base >= width {
			return width
		}
		cost := mn - base
		i := l
		for i <= r {
			if a[i] == mn {
				i++
				continue
			}
			j := i
			for j <= r && a[j] > mn {
				j++
			}
			cost += solve(i, j-1, mn)
			if cost >= width {
				return width
			}
			i = j
		}
		if cost < width {
			return cost
		}
		return width
	}

	var ans int64
	i := 0
	for i < n {
		for i < n && a[i] == 0 {
			i++
		}
		if i == n {
			break
		}
		j := i
		for j < n && a[j] > 0 {
			j++
		}
		ans += solve(i, j-1, 0)
		i = j
	}

	fmt.Println(ans)
}
`

func buildOracle() (string, error) {
	dir := os.TempDir()
	src := filepath.Join(dir, fmt.Sprintf("oracle1400E_%d.go", time.Now().UnixNano()))
	if err := os.WriteFile(src, []byte(oracleSource), 0644); err != nil {
		return "", fmt.Errorf("write oracle source: %v", err)
	}
	defer os.Remove(src)
	oracle := src[:len(src)-3]
	cmd := exec.Command("go", "build", "-o", oracle, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(r *rand.Rand) string {
	n := r.Intn(20) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", r.Intn(10)))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input := genCase(rng)
		expect, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", i, err)
			os.Exit(1)
		}
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
