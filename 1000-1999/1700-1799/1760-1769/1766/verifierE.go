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

const stateCount = 16

var states = [stateCount]int{0, 1, 2, 3, 5, 6, 7, 9, 10, 11, 25, 26, 27, 37, 38, 39}
var indexMap = map[int]int{}
var trans [stateCount][4]int
var inc [stateCount][4]int

func encode(m1, m2, m3 int) int { return m1 | (m2 << 2) | (m3 << 4) }

func appendState(state, x int) (int, int) {
	m1 := state & 3
	m2 := (state >> 2) & 3
	m3 := (state >> 4) & 3
	if x == 0 {
		return state, 1
	}
	if m1 != 0 && (m1&x) != 0 {
		m1 = x
		return encode(m1, m2, m3), 0
	}
	if m2 != 0 && (m2&x) != 0 {
		m2 = x
		return encode(m1, m2, m3), 0
	}
	if m3 != 0 && (m3&x) != 0 {
		m3 = x
		return encode(m1, m2, m3), 0
	}
	if m1 == 0 {
		m1 = x
		return encode(m1, m2, m3), 1
	}
	if m2 == 0 {
		m2 = x
		return encode(m1, m2, m3), 1
	}
	if m3 == 0 {
		m3 = x
		return encode(m1, m2, m3), 1
	}
	return state, 1
}

func init() {
	for i, s := range states {
		indexMap[s] = i
	}
	for i, s := range states {
		for x := 0; x < 4; x++ {
			ns, d := appendState(s, x)
			trans[i][x] = indexMap[ns]
			inc[i][x] = d
		}
	}
}

func solveE(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return ""
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}
	var cnt [stateCount]int64
	var sum [stateCount]int64
	var ans int64
	for _, v := range arr {
		var nextCnt [stateCount]int64
		var nextSum [stateCount]int64
		ns := trans[0][v]
		d := inc[0][v]
		nextCnt[ns] += 1
		nextSum[ns] += int64(d)
		for i := 0; i < stateCount; i++ {
			if cnt[i] == 0 {
				continue
			}
			ns := trans[i][v]
			dd := inc[i][v]
			nextCnt[ns] += cnt[i]
			nextSum[ns] += sum[i] + int64(dd)*cnt[i]
		}
		var partial int64
		for i := 0; i < stateCount; i++ {
			partial += nextSum[i]
		}
		ans += partial
		cnt = nextCnt
		sum = nextSum
	}
	return fmt.Sprint(ans)
}

func genTestE(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	var buf strings.Builder
	fmt.Fprintf(&buf, "%d\n", n)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&buf, "%d ", rng.Intn(4))
	}
	buf.WriteByte('\n')
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
