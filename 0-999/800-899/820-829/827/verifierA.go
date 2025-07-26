package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

func buildReference() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "827A.go")
	bin := filepath.Join(os.TempDir(), "ref827A.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("compile reference: %v\n%s", err, out)
	}
	return bin, nil
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func randString(n int) string {
	letters := "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func genCase() string {
	L := rand.Intn(20) + 5
	final := make([]byte, L)
	for i := range final {
		final[i] = byte('a' + rand.Intn(26))
	}
	n := rand.Intn(4) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		l := rand.Intn(4) + 1
		if l > L {
			l = L
		}
		start := rand.Intn(L-l+1) + 1
		t := string(final[start-1 : start-1+l])
		var occ []int
		for j := 1; j <= L-l+1; j++ {
			if string(final[j-1:j-1+l]) == t {
				occ = append(occ, j)
			}
		}
		k := rand.Intn(len(occ)) + 1
		rand.Shuffle(len(occ), func(a, b int) { occ[a], occ[b] = occ[b], occ[a] })
		pos := occ[:k]
		sort.Ints(pos)
		sb.WriteString(t)
		sb.WriteString(fmt.Sprintf(" %d", k))
		for _, p := range pos {
			sb.WriteString(fmt.Sprintf(" %d", p))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	userBin := os.Args[1]
	rand.Seed(1)
	ref, err := buildReference()
	if err != nil {
		fmt.Println("reference compile failed:", err)
		return
	}
	defer os.Remove(ref)

	for i := 0; i < 100; i++ {
		input := genCase()
		want, err1 := runBinary(ref, input)
		if err1 != nil {
			fmt.Printf("reference solution failed on test %d: %v\n", i+1, err1)
			return
		}
		got, err2 := runBinary(userBin, input)
		if err2 != nil {
			fmt.Printf("test %d runtime error: %v\n", i+1, err2)
			fmt.Println("input:\n" + input)
			return
		}
		if want != got {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\nactual:\n%s\n", i+1, input, want, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
