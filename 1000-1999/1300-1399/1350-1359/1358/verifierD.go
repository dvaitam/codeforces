package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func tri(x int64) int64 { return x * (x + 1) / 2 }

func expected(d []int64, x int64) string {
	n := len(d)
	arr := make([]int64, 2*n)
	for i := 0; i < 2*n; i++ {
		arr[i] = d[i%n]
	}
	prefDays := make([]int64, 2*n+1)
	prefHugs := make([]int64, 2*n+1)
	for i := 0; i < 2*n; i++ {
		prefDays[i+1] = prefDays[i] + arr[i]
		prefHugs[i+1] = prefHugs[i] + tri(arr[i])
	}
	var best int64
	for r := 1; r <= 2*n; r++ {
		if prefDays[r] < x {
			continue
		}
		need := prefDays[r] - x
		idx := sort.Search(len(prefDays), func(i int) bool { return prefDays[i] > need }) - 1
		if idx < 0 {
			idx = 0
		}
		leftover := need - prefDays[idx]
		partial := tri(arr[idx]) - tri(leftover)
		val := partial + prefHugs[r] - prefHugs[idx+1]
		if val > best {
			best = val
		}
	}
	return fmt.Sprintf("%d", best)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 1; t <= 100; t++ {
		n := rng.Intn(10) + 1
		x := int64(rng.Intn(50) + 1)
		d := make([]int64, n)
		for i := range d {
			d[i] = int64(rng.Intn(10) + 1)
		}
		// ensure x <= sum(d)
		sum := int64(0)
		for _, v := range d {
			sum += v
		}
		if x > sum {
			x = sum
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, x))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", d[i]))
		}
		sb.WriteByte('\n')
		want := expected(d, x)
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", t, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
