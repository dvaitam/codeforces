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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveA(n int, k int64, arr []int64) string {
	rem := arr[0] % k
	minA := arr[0]
	for i := 1; i < n; i++ {
		if arr[i]%k != rem {
			return "-1"
		}
		if arr[i] < minA {
			minA = arr[i]
		}
	}
	var ans int64
	for i := 0; i < n; i++ {
		ans += (arr[i] - minA) / k
	}
	return fmt.Sprint(ans)
}

func genCase(rng *rand.Rand) (int, int64, []int64) {
	n := rng.Intn(8) + 2
	k := int64(rng.Intn(20) + 1)
	arr := make([]int64, n)
	base := int64(rng.Intn(50))
	rem := rng.Int63n(k)
	for i := 0; i < n; i++ {
		delta := int64(rng.Intn(50))
		arr[i] = base + delta*k + rem
	}
	return n, k, arr
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, k, arr := genCase(rng)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, k)
		for j, v := range arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
		expect := solveA(n, k, arr)
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected %s got %s\n", i+1, sb.String(), expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
