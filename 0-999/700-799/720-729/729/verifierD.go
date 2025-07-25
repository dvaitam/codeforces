package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
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

func expected(n, a, b, k int, s string) (string, []int) {
	var candidates []int
	for i := 0; i < n; {
		if s[i] == '1' {
			i++
			continue
		}
		j := i
		for j < n && s[j] == '0' {
			j++
		}
		L := j - i
		cnt := L / b
		for t := 1; t <= cnt; t++ {
			pos := i + t*b
			candidates = append(candidates, pos)
		}
		i = j
	}
	R := len(candidates) - a + 1
	if R < 0 {
		R = 0
	}
	if R > len(candidates) {
		R = len(candidates)
	}
	return fmt.Sprintf("%d", R), candidates[:R]
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	r := rand.New(rand.NewSource(4))
	for tc := 1; tc <= 100; tc++ {
		n := r.Intn(20) + 1
		a := r.Intn(n) + 1
		b := r.Intn(n) + 1
		if b == 0 {
			b = 1
		}
		k := r.Intn(n)
		onesPos := make([]int, 0, k)
		used := make(map[int]bool)
		for len(onesPos) < k {
			p := r.Intn(n)
			if !used[p] {
				used[p] = true
				onesPos = append(onesPos, p)
			}
		}
		str := make([]byte, n)
		for i := 0; i < n; i++ {
			if used[i] {
				str[i] = '1'
			} else {
				str[i] = '0'
			}
		}
		s := string(str)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d %d\n", n, a, b, k)
		sb.WriteString(s)
		sb.WriteByte('\n')
		input := sb.String()
		expectCount, expectPositions := expected(n, a, b, k, s)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", tc, err)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		fields := strings.Fields(out)
		if len(fields) < 1 {
			fmt.Printf("test %d: output too short\n", tc)
			os.Exit(1)
		}
		if fields[0] != expectCount {
			fmt.Printf("test %d failed\ninput:\n%sexpected count: %s got: %s\n", tc, input, expectCount, fields[0])
			os.Exit(1)
		}
		if len(fields)-1 != len(expectPositions) {
			fmt.Printf("test %d failed\nexpected %d positions got %d\n", tc, len(expectPositions), len(fields)-1)
			os.Exit(1)
		}
		for i, pos := range expectPositions {
			if fields[i+1] != fmt.Sprintf("%d", pos) {
				fmt.Printf("test %d failed\nposition %d expected %d got %s\n", tc, i+1, pos, fields[i+1])
				os.Exit(1)
			}
		}
	}
	fmt.Println("All 100 tests passed")
}
