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

func buildRef() (string, error) {
	refBin := "./1762F_ref"
	if err := exec.Command("go", "build", "-o", refBin, "1762F.go").Run(); err != nil {
		return "", err
	}
	return refBin, nil
}

func run(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return out.String() + errBuf.String(), err
	}
	return strings.TrimSpace(out.String()), nil
}

type testF struct {
	n   int
	k   int
	arr []int
}

func genTestF(rng *rand.Rand) testF {
	n := rng.Intn(20) + 1
	k := rng.Intn(10)
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(100000) + 1
	}
	return testF{n: n, k: k, arr: arr}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	userBin := os.Args[1]
	refBin, err := buildRef()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		return
	}
	defer os.Remove(refBin)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		tc := genTestF(rng)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		input := sb.String()

		expect, err := run(refBin, input)
		if err != nil {
			fmt.Printf("reference failed on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		got, err := run(userBin, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n%s", t+1, err, got)
			os.Exit(1)
		}
		if strings.TrimSpace(expect) != strings.TrimSpace(got) {
			fmt.Printf("test %d failed: expected %s got %s\n", t+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", 100)
}
