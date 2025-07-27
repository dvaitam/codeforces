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

func solveCase(a []int) string {
	n := len(a)
	pref := make([]int, n+1)
	for i := 0; i < n; i++ {
		pref[i+1] = pref[i] + a[i]
	}
	cnt := 0
	for l := 0; l < n; l++ {
		for r := l + 2; r < n; r++ {
			sum := pref[r] - pref[l+1]
			if (a[l] ^ a[r]) == sum {
				cnt++
			}
		}
	}
	return fmt.Sprint(cnt)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(40) + 3
	var in strings.Builder
	in.WriteString(fmt.Sprintf("%d\n", n))
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(1<<15) + 1
		in.WriteString(fmt.Sprintf("%d", arr[i]))
		if i+1 < n {
			in.WriteByte(' ')
		}
	}
	in.WriteByte('\n')
	exp := solveCase(arr)
	return in.String(), exp
}

func runCase(bin, input, expected string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
