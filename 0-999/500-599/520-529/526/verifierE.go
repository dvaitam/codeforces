package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveECase(a []int64, b int64) int {
	n := len(a)
	arr := make([]int64, 2*n)
	copy(arr, a)
	for i := 0; i < n; i++ {
		arr[n+i] = a[i]
	}
	r := make([]int, 2*n)
	var sum int64
	j := 0
	for i := 0; i < 2*n; i++ {
		for j < 2*n && sum+arr[j] <= b {
			sum += arr[j]
			j++
		}
		r[i] = j
		sum -= arr[i]
	}
	dq := make([]int, 0, n+2)
	head := 0
	cur := 0
	dq = append(dq, cur)
	for cur < n {
		cur = r[cur]
		dq = append(dq, cur)
	}
	ans := len(dq) - head - 1
	for s := 1; s < n; s++ {
		for head < len(dq) && dq[head] < s {
			head++
		}
		last := dq[len(dq)-1]
		for last < s+n {
			last = r[last]
			dq = append(dq, last)
		}
		cnt := len(dq) - head - 1
		if cnt < ans {
			ans = cnt
		}
	}
	return ans
}

func solveE(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, q int
	fmt.Fscan(reader, &n, &q)
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}
	var sb strings.Builder
	for qi := 0; qi < q; qi++ {
		var b int64
		fmt.Fscan(reader, &b)
		res := solveECase(arr, b)
		if qi > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(fmt.Sprint(res))
	}
	return sb.String()
}

func genTestE(rng *rand.Rand) string {
	n := rng.Intn(8) + 1
	q := rng.Intn(5) + 1
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		arr[i] = int64(rng.Intn(20) + 1)
	}
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%d %d\n", n, q)
	for i, v := range arr {
		if i > 0 {
			buf.WriteByte(' ')
		}
		fmt.Fprintf(&buf, "%d", v)
	}
	buf.WriteByte('\n')
	for i := 0; i < q; i++ {
		fmt.Fprintf(&buf, "%d\n", rng.Intn(40)+1)
	}
	return buf.String()
}

func run(bin, input string) (string, error) {
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	for i := 1; i <= 100; i++ {
		in := genTestE(rng)
		expect := solveE(in)
		got, err := run(exe, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i, in, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
