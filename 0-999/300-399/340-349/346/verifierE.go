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

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
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

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func expectedAnswer(instances [][4]int64) []string {
	res := make([]string, len(instances))
	for idx, ins := range instances {
		a, n, p, h := ins[0], ins[1], ins[2], ins[3]
		if n <= 2000000 {
			heights := make([]int64, n)
			for i := int64(1); i <= n; i++ {
				heights[i-1] = (a * i) % p
			}
			sort.Slice(heights, func(i, j int) bool { return heights[i] < heights[j] })
			ok := true
			if heights[0] > h {
				ok = false
			}
			for i := 1; ok && i < int(n); i++ {
				if heights[i]-heights[i-1] > h {
					ok = false
				}
			}
			if ok {
				res[idx] = "YES"
			} else {
				res[idx] = "NO"
			}
		} else {
			if (n+1)*h >= p {
				res[idx] = "YES"
			} else {
				res[idx] = "NO"
			}
		}
	}
	return res
}

func generateInstance(rng *rand.Rand) [4]int64 {
	var a, p int64
	for {
		a = int64(rng.Intn(20) + 1)
		p = int64(rng.Intn(50) + 30)
		if gcd(a, p) == 1 {
			break
		}
	}
	n := int64(rng.Intn(10) + 1)
	h := int64(rng.Intn(20))
	if n >= p {
		n = p - 1
	}
	return [4]int64{a, n, p, h}
}

func generateCase(rng *rand.Rand) string {
	t := rng.Intn(3) + 1
	inst := make([][4]int64, t)
	for i := 0; i < t; i++ {
		inst[i] = generateInstance(rng)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		v := inst[i]
		fmt.Fprintf(&sb, "%d %d %d %d\n", v[0], v[1], v[2], v[3])
	}
	return sb.String()
}

func parseCase(input string) [][4]int64 {
	r := strings.NewReader(input)
	var t int
	fmt.Fscan(r, &t)
	inst := make([][4]int64, t)
	for i := 0; i < t; i++ {
		fmt.Fscan(r, &inst[i][0], &inst[i][1], &inst[i][2], &inst[i][3])
	}
	return inst
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []string{}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}
	for idx, tc := range cases {
		inst := parseCase(tc)
		expect := expectedAnswer(inst)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, tc)
			os.Exit(1)
		}
		outLines := strings.Split(strings.TrimSpace(got), "\n")
		if len(outLines) != len(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d lines got %d\ninput:\n%s", idx+1, len(expect), len(outLines), tc)
			os.Exit(1)
		}
		for i := range expect {
			if strings.TrimSpace(outLines[i]) != expect[i] {
				fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s (line %d)\ninput:\n%s", idx+1, expect[i], outLines[i], i+1, tc)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
