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

func runProg(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		tmp := fmt.Sprintf("/tmp/tmpbin-%d", rand.Int())
		build := exec.Command("go", "build", "-o", tmp, bin)
		if bErr := build.Run(); bErr != nil {
			return "", fmt.Errorf("build failed: %v", bErr)
		}
		defer os.Remove(tmp)
		cmd = exec.Command(tmp)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("%v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func solveA(s string) string {
	n := len(s)
	type pair struct{ bal, idx int }
	arr := make([]pair, n)
	bal := 0
	for i := 0; i < n; i++ {
		arr[i] = pair{bal, i}
		if s[i] == '(' {
			bal++
		} else {
			bal--
		}
	}
	sort.Slice(arr, func(i, j int) bool {
		if arr[i].bal == arr[j].bal {
			return arr[i].idx > arr[j].idx
		}
		return arr[i].bal < arr[j].bal
	})
	res := make([]byte, n)
	for i := 0; i < n; i++ {
		res[i] = s[arr[i].idx]
	}
	return string(res)
}

func randomBalanced(rng *rand.Rand, n int) string {
	if n%2 == 1 {
		n++
	}
	open := n / 2
	close := n / 2
	bal := 0
	var sb strings.Builder
	for open+close > 0 {
		if open == close || bal == 0 {
			sb.WriteByte('(')
			open--
			bal++
			continue
		}
		if open == 0 {
			sb.WriteByte(')')
			close--
			bal--
			continue
		}
		if rng.Intn(2) == 0 {
			sb.WriteByte('(')
			open--
			bal++
		} else {
			sb.WriteByte(')')
			close--
			bal--
		}
	}
	return sb.String()
}

func runCase(bin string, s string) error {
	expect := solveA(s)
	out, err := runProg(bin, s+"\n")
	if err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out)
	}
	if out != expect {
		return fmt.Errorf("expected %s got %s", expect, out)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(40) + 2
		if n%2 == 1 {
			n++
		}
		s := randomBalanced(rng, n)
		if err := runCase(bin, s); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\n", i+1, err, s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
