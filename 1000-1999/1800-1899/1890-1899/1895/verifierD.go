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

func countOnes(n, bit int) int {
	cycle := 1 << (bit + 1)
	full := n / cycle
	cnt := full * (1 << bit)
	rem := n % cycle
	if rem > (1 << bit) {
		cnt += rem - (1 << bit)
	}
	return cnt
}

func solveD(n int, arr []int) string {
	pref := make([]int, n)
	cur := 0
	for i := 1; i < n; i++ {
		cur ^= arr[i-1]
		pref[i] = cur
	}
	bits := 20
	cntPref := make([]int, bits)
	for _, v := range pref {
		for j := 0; j < bits; j++ {
			if (v>>j)&1 == 1 {
				cntPref[j]++
			}
		}
	}
	cntRange := make([]int, bits)
	for j := 0; j < bits; j++ {
		cntRange[j] = countOnes(n, j)
	}
	key := 0
	for j := 0; j < bits; j++ {
		if cntPref[j] != cntRange[j] {
			key |= 1 << j
		}
	}
	var sb strings.Builder
	for i, v := range pref {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v ^ key))
	}
	return sb.String()
}

func genCase(rng *rand.Rand) (int, []int) {
	n := rng.Intn(10) + 2
	arr := make([]int, n-1)
	for i := range arr {
		arr[i] = rng.Intn(512)
	}
	return n, arr
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, arr := genCase(rng)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j, v := range arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		expect := solveD(n, arr)
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
