package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func expectedD(a []int) string {
	n := len(a)
	if n < 3 {
		sorted := true
		for i := 1; i < n; i++ {
			if a[i] < a[i-1] {
				sorted = false
				break
			}
		}
		if sorted {
			return "YES"
		}
		return "NO"
	}
	freq := make(map[int]int)
	for _, v := range a {
		freq[v]++
	}
	for _, v := range freq {
		if v > 1 {
			return "YES"
		}
	}
	visited := make([]bool, n+1)
	cycles := 0
	for i := 1; i <= n; i++ {
		if !visited[i] {
			cycles++
			for j := i; !visited[j]; j = a[j-1] {
				visited[j] = true
			}
		}
	}
	if (n-cycles)%2 == 0 {
		return "YES"
	}
	return "NO"
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 1
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(n) + 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(arr[i]))
	}
	sb.WriteByte('\n')
	return sb.String(), expectedD(arr)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := generateCase(rng)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got := strings.TrimSpace(out)
		if strings.ToUpper(got) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
