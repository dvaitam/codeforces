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

// Embedded correct solver for 1401F
var solverData []byte
var solverIdx int
var seg []int64
var solverMask int

func solverNextInt() int {
	n := len(solverData)
	for solverIdx < n && (solverData[solverIdx] < '0' || solverData[solverIdx] > '9') && solverData[solverIdx] != '-' {
		solverIdx++
	}
	sign := 1
	if solverIdx < n && solverData[solverIdx] == '-' {
		sign = -1
		solverIdx++
	}
	val := 0
	for solverIdx < n && solverData[solverIdx] >= '0' && solverData[solverIdx] <= '9' {
		val = val*10 + int(solverData[solverIdx]-'0')
		solverIdx++
	}
	return val * sign
}

func build(node, level, start int, arr []int64) {
	if level == 0 {
		seg[node] = arr[start]
		return
	}
	half := 1 << (level - 1)
	build(node<<1, level-1, start, arr)
	build(node<<1|1, level-1, start+half, arr)
	seg[node] = seg[node<<1] + seg[node<<1|1]
}

func update(node, level, pos int, val int64) {
	if level == 0 {
		seg[node] = val
		return
	}
	half := 1 << (level - 1)
	leftNode, rightNode := node<<1, node<<1|1
	if (solverMask>>uint(level-1))&1 == 1 {
		leftNode, rightNode = rightNode, leftNode
	}
	if pos < half {
		update(leftNode, level-1, pos, val)
	} else {
		update(rightNode, level-1, pos-half, val)
	}
	seg[node] = seg[node<<1] + seg[node<<1|1]
}

func query(node, level, l, r int) int64 {
	if l == 0 && r == (1<<level)-1 {
		return seg[node]
	}
	if level == 0 {
		return seg[node]
	}
	half := 1 << (level - 1)
	leftNode, rightNode := node<<1, node<<1|1
	if (solverMask>>uint(level-1))&1 == 1 {
		leftNode, rightNode = rightNode, leftNode
	}
	if r < half {
		return query(leftNode, level-1, l, r)
	}
	if l >= half {
		return query(rightNode, level-1, l-half, r-half)
	}
	return query(leftNode, level-1, l, half-1) + query(rightNode, level-1, 0, r-half)
}

func solve(input string) string {
	solverData = []byte(input)
	solverIdx = 0
	solverMask = 0

	n := solverNextInt()
	q := solverNextInt()
	size := 1 << n
	arr := make([]int64, size)
	for i := 0; i < size; i++ {
		arr[i] = int64(solverNextInt())
	}
	seg = make([]int64, size*4)
	build(1, n, 0, arr)

	var out strings.Builder
	out.Grow(q * 20)

	for i := 0; i < q; i++ {
		t := solverNextInt()
		switch t {
		case 1:
			x := solverNextInt() - 1
			k := int64(solverNextInt())
			update(1, n, x, k)
		case 2:
			k := solverNextInt()
			solverMask ^= (1 << k) - 1
		case 3:
			k := solverNextInt()
			solverMask ^= 1 << k
		case 4:
			l := solverNextInt() - 1
			r := solverNextInt() - 1
			out.WriteString(strconv.FormatInt(query(1, n, l, r), 10))
			out.WriteByte('\n')
		}
	}

	return strings.TrimSpace(out.String())
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(5)
	N := 1 << n
	q := rng.Intn(15) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	for i := 0; i < N; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(rng.Intn(1000)))
	}
	sb.WriteByte('\n')
	sumCnt := 0
	for i := 0; i < q; i++ {
		typ := rng.Intn(4) + 1
		if i == q-1 && sumCnt == 0 {
			typ = 4
		}
		switch typ {
		case 1:
			x := rng.Intn(N) + 1
			k := rng.Intn(1000)
			sb.WriteString(fmt.Sprintf("1 %d %d\n", x, k))
		case 2:
			k := 0
			if n > 0 {
				k = rng.Intn(n + 1)
			}
			sb.WriteString(fmt.Sprintf("2 %d\n", k))
		case 3:
			k := 0
			if n > 0 {
				k = rng.Intn(n)
			}
			sb.WriteString(fmt.Sprintf("3 %d\n", k))
		case 4:
			l := rng.Intn(N) + 1
			r := rng.Intn(N-l+1) + l
			sb.WriteString(fmt.Sprintf("4 %d %d\n", l, r))
			sumCnt++
		}
	}
	return sb.String()
}

func runCase(bin, input string) error {
	exp := solve(input)
	got, err := run(bin, input)
	if err != nil {
		return err
	}
	if strings.TrimSpace(got) != strings.TrimSpace(exp) {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		if err := runCase(bin, input); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
