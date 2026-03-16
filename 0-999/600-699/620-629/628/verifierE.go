package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func compileRef() (string, error) {
	exe := filepath.Join(os.TempDir(), "ref628E")
	cmd := exec.Command("go", "build", "-o", exe, "628E.go")
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return exe, nil
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
	bin := os.Args[1]
	ref, err := compileRef()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	const testcasesERaw = `100
5 5
z..zz
zzzz.
..z.z
zz.zz
z..z.
1 5
.zzzz
5 1
.
z
z
.
z
5 3
zz.
..z
.zz
...
..z
3 5
.zz.z
zz.z.
.zzz.
4 3
.z.
..z
zz.
..z
4 2
.z
z.
z.
zz
2 2
zz
..
1 5
zzz.z
2 1
z
z
5 5
....z
..zzz
.zzzz
z..z.
....z
5 1
.
.
z
z
z
4 2
z.
.z
.z
z.
3 1
.
z
z
3 5
...zz
.z.zz
..zz.
3 2
zz
z.
.z
1 3
.z.
1 1
.
2 5
zzz..
..zzz
5 5
z.zzz
zzzzz
z.z..
..zzz
.z...
4 4
z.zz
zzzz
..z.
.z..
2 3
z..
z..
1 3
z.z
4 4
z...
zzz.
z..z
....
2 5
..zzz
..zz.
3 3
z..
zzz
...
2 3
...
.zz
5 4
zzz.
.z.z
.z.z
..z.
zzzz
3 3
.zz
...
z..
1 1
z
1 3
z.z
5 3
zz.
zzz
z.z
z..
z.z
5 4
zzz.
z.z.
.z..
z.z.
z...
4 4
...z
zz..
z...
zzz.
5 2
z.
.z
z.
..
..
1 2
zz
4 2
z.
.z
z.
z.
3 1
z
.
.
4 5
.zzz.
...z.
..z..
.z...
5 5
.....
z.z..
..zz.
z.z.z
z.zzz
5 1
z
z
z
.
.
2 3
...
zz.
2 1
z
.
3 3
..z
...
zz.
2 4
.z.z
.zzz
3 1
.
.
z
1 1
z
2 5
z.zzz
zzzzz
2 1
.
.
1 1
.
4 5
zzz..
zzzzz
.....
zzz..
3 3
.zz
zz.
zzz
1 4
.z..
2 1
z
.
3 2
z.
z.
.z
1 5
z..z.
4 5
.z.z.
z..z.
z.z..
z.zz.
3 5
zz.zz
..z.z
zz.z.
2 5
zzz..
...zz
4 3
z.z
.zz
...
zz.
4 4
z.zz
z.zz
z.z.
....
4 4
z...
...z
.z..
.zz.
5 1
.
z
z
.
.
5 3
...
zz.
z..
zzz
z.z
4 3
.z.
...
zz.
..z
3 2
zz
z.
zz
1 4
.zzz
4 5
z.z..
z..z.
zzz..
.zz..
4 2
z.
..
..
.z
4 2
.z
.z
z.
..
3 4
.zz.
z.z.
zz..
2 1
.
z
3 3
zz.
...
.zz
3 2
zz
.z
z.
1 1
.
1 4
..zz
1 1
z
1 4
z.zz
2 3
.z.
..z
4 3
zzz
z..
z..
..z
2 3
z.z
zzz
4 4
z.zz
.zzz
....
..zz
3 2
z.
.z
zz
4 5
....z
..zz.
zzzzz
zz..z
1 1
z
2 3
...
z.z
5 2
z.
z.
zz
zz
..
3 4
....
....
....
1 2
..
4 4
zz.z
zzz.
zz.z
.z.z
1 1
z
1 1
z
2 5
z..zz
..z.z
5 2
..
zz
.z
zz
z.
3 3
z..
zz.
zzz
3 4
z..z
.z..
.z.z
3 3
zz.
.zz
z..
1 4
...z
4 1
.
z
z
.
3 3
z..
z.z
.z.`

	scan := bufio.NewReader(strings.NewReader(testcasesERaw))
	var t int
	fmt.Fscan(scan, &t)
	for caseNum := 1; caseNum <= t; caseNum++ {
		var n, m int
		if _, err := fmt.Fscan(scan, &n, &m); err != nil {
			fmt.Printf("bad test file at case %d\n", caseNum)
			os.Exit(1)
		}
		rows := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(scan, &rows[i])
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < n; i++ {
			sb.WriteString(rows[i])
			sb.WriteByte('\n')
		}
		input := sb.String()
		expect, err := run(ref, input)
		if err != nil {
			fmt.Printf("reference failed on case %d: %v\n", caseNum, err)
			os.Exit(1)
		}
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", caseNum, err, input)
			os.Exit(1)
		}
		if out != expect {
			fmt.Printf("case %d failed: expected %s got %s\n", caseNum, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
