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

func solve(arr []int) (int, []int) {
	n := len(arr)
	dif := make([]int, n)
	if n > 0 {
		dif[0] = arr[0]
	}
	for i := 1; i < n; i++ {
		dif[i] = arr[i] - arr[i-1]
	}
	per := make([]bool, n+1)
	ans := 0
	for i := 1; i <= n; i++ {
		per[i] = true
		for x := 0; x < n; x++ {
			if dif[x] != dif[x%i] {
				per[i] = false
				break
			}
		}
		if per[i] {
			ans++
		}
	}
	res := make([]int, 0, ans)
	for i := 1; i <= n; i++ {
		if per[i] {
			res = append(res, i)
		}
	}
	return ans, res
}

func runCase(bin, input, want string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if strings.TrimSpace(want) != got {
		return fmt.Errorf("expected %q got %q", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	tests := [][]int{{1, 2, 3, 4}, {1, 3, 6, 7, 9}}
	for i := 0; i < 100; i++ {
		n := rng.Intn(20) + 1
		arr := make([]int, n)
		cur := rng.Intn(10)
		for j := 0; j < n; j++ {
			cur += rng.Intn(5)
			arr[j] = cur
		}
		tests = append(tests, arr)
	}

	for idx, arr := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", len(arr)))
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		cnt, vals := solve(arr)
		want := fmt.Sprintf("%d\n", cnt)
		for i, v := range vals {
			if i > 0 {
				want += " "
			}
			want += fmt.Sprintf("%d", v)
		}
		want += "\n"
		if err := runCase(bin, sb.String(), want); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, sb.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
