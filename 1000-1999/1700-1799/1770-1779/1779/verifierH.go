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

type Player struct {
	val int64
	idx int
}

func solveH(vals []int64) string {
	n := len(vals)
	players := make([]Player, n)
	var total int64
	for i := 0; i < n; i++ {
		players[i] = Player{vals[i], i}
		total += vals[i]
	}
	sort.Slice(players, func(i, j int) bool { return players[i].val > players[j].val })
	var pref int64
	k := 0
	half := total / 2
	for k < n {
		pref += players[k].val
		k++
		if pref >= half {
			break
		}
	}
	res := make([]byte, n)
	for i := 0; i < n; i++ {
		res[i] = '0'
	}
	for i := 0; i < k; i++ {
		res[players[i].idx] = '1'
	}
	return string(res)
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

func genCase(rng *rand.Rand) (string, string) {
	nPow := rng.Intn(3) + 2 // 2..4 => n from 4 to 16
	n := 1 << nPow
	vals := make([]int64, n)
	for i := 0; i < n; i++ {
		vals[i] = rng.Int63n(100) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range vals {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteString("\n")
	input := sb.String()
	expected := solveH(vals)
	return input, expected
}

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append([]string{os.Args[0]}, os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := genCase(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
