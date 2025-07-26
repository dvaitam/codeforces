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

const numTestsF = 100

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: verifierF.go <binary>")
		os.Exit(1)
	}
	binPath, cleanup, err := prepareBinary(os.Args[1])
	if err != nil {
		fmt.Println("compile error:", err)
		os.Exit(1)
	}
	if cleanup != nil {
		defer cleanup()
	}
	r := rand.New(rand.NewSource(1))
	for t := 1; t <= numTestsF; t++ {
		n := r.Intn(20) + 1
		a := make([]int, n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			a[i] = r.Intn(30) + 1
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", a[i]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expected := solveF(a)
		out, err := run(binPath, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", t, err)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if out != expected {
			fmt.Printf("test %d failed\ninput:%sexpected:%s got:%s\n", t, input, expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("OK")
}

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp := filepath.Join(os.TempDir(), "verify_binF")
		cmd := exec.Command("go", "build", "-o", tmp, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, string(out))
		}
		return tmp, func() { os.Remove(tmp) }, nil
	}
	return path, nil, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return buf.String(), err
}

func solveF(a []int) string {
	const MAX = 200005
	cnt := make([]int, MAX)
	for _, x := range a {
		if x >= 0 && x < MAX {
			cnt[x]++
		}
	}
	sumcnt := make([]int, MAX)
	for i := 1; i < MAX; i++ {
		sumcnt[i] = sumcnt[i-1] + cnt[i]
	}
	ans := -1
	ansL, ansR := 0, 0
	r := 0
	for l := 1; l < MAX; l++ {
		if cnt[l] == 0 {
			continue
		}
		if r < l {
			r = l
		}
		for r+1 < MAX && cnt[r+1] >= 1 && (cnt[r] >= 2 || l == r) {
			r++
		}
		tans := sumcnt[r] - sumcnt[l-1]
		if tans > ans {
			ans = tans
			ansL = l
			ansR = r
		}
	}
	var out strings.Builder
	out.WriteString(fmt.Sprintf("%d\n", ans))
	for i := ansL; i <= ansR; i++ {
		if cnt[i] > 0 {
			out.WriteString(fmt.Sprintf("%d ", i))
			cnt[i]--
		}
	}
	for i := ansR; i >= ansL; i-- {
		for cnt[i] > 0 {
			out.WriteString(fmt.Sprintf("%d ", i))
			cnt[i]--
		}
	}
	s := strings.TrimSpace(out.String())
	return s
}
