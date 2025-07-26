package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func compileRef() (string, error) {
	out := filepath.Join(os.TempDir(), "1170D_ref")
	cmd := exec.Command("go", "build", "-o", out, "1170D.go")
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return out, nil
}

func runBin(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func encode(seqs [][]int) []int {
	// append -1 to each sequence and interleave by rounds
	for i := range seqs {
		seqs[i] = append(seqs[i], -1)
	}
	var res []int
	for idx := 0; ; idx++ {
		added := false
		for i := range seqs {
			if idx < len(seqs[i]) {
				res = append(res, seqs[i][idx])
				added = true
			}
		}
		if !added {
			break
		}
	}
	return res
}

func genTest() string {
	n := rand.Intn(5) + 1
	seqs := make([][]int, n)
	for i := 0; i < n; i++ {
		l := rand.Intn(5)
		seq := make([]int, l)
		for j := 0; j < l; j++ {
			seq[j] = rand.Intn(20)
		}
		seqs[i] = seq
	}
	enc := encode(seqs)
	nums := make([]string, len(enc))
	for i, v := range enc {
		nums[i] = fmt.Sprintf("%d", v)
	}
	return fmt.Sprintf("%d\n%s\n", len(enc), strings.Join(nums, " "))
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref, err := compileRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to compile reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rand.Seed(1)
	for t := 0; t < 100; t++ {
		test := genTest()
		exp, err := runBin(ref, test)
		if err != nil {
			fmt.Fprintln(os.Stderr, "reference run failed:", err)
			os.Exit(1)
		}
		got, err := runBin(cand, test)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate run failed on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("mismatch on test %d\ninput:\n%sexpected:%s got:%s\n", t+1, test, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
