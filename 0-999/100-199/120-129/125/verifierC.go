package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const solution125CSource = `package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // find maximum k such that k*(k-1)/2 <= n
   k := 1
   for (k*(k-1))/2 <= n {
       k++
   }
   k--
   // prepare guest lists for days 1..k
   guests := make([][]int, k)
   id := 1
   for i := 0; i < k; i++ {
       for j := i + 1; j < k; j++ {
           guests[i] = append(guests[i], id)
           guests[j] = append(guests[j], id)
           id++
       }
   }
   // output
   fmt.Fprintln(writer, k)
   for i := 0; i < k; i++ {
       // each guest list for day i+1
       for j, x := range guests[i] {
           if j > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, x)
       }
       fmt.Fprintln(writer)
   }
}
`

// Keep the embedded reference solution reachable.
var _ = solution125CSource

const testcasesRaw = `100
27
29
5
19
35
34
28
22
33
25
40
16
35
11
21
11
9
42
19
37
48
41
12
22
9
49
7
46
24
33
38
9
25
30
23
42
43
16
38
33
31
36
19
6
38
3
8
49
28
48
45
43
3
42
34
24
18
49
23
48
7
15
39
17
18
12
37
31
8
8
23
35
34
9
22
38
21
48
10
38
24
37
16
41
38
40
21
31
8
41
27
23
39
18
21
14
15
14
5
42
`

func parseTestcases() []int {
	fields := strings.Fields(testcasesRaw)
	if len(fields) == 0 {
		return nil
	}
	t, _ := strconv.Atoi(fields[0])
	cases := make([]int, 0, t)
	for i := 1; i <= t && i < len(fields); i++ {
		n, _ := strconv.Atoi(fields[i])
		cases = append(cases, n)
	}
	return cases
}

func solveCase(n int) string {
	k := 1
	for (k*(k-1))/2 <= n {
		k++
	}
	k--
	guests := make([][]int, k)
	id := 1
	for i := 0; i < k; i++ {
		for j := i + 1; j < k; j++ {
			guests[i] = append(guests[i], id)
			guests[j] = append(guests[j], id)
			id++
		}
	}
	var out strings.Builder
	fmt.Fprintln(&out, k)
	for i := 0; i < k; i++ {
		for j, x := range guests[i] {
			if j > 0 {
				out.WriteByte(' ')
			}
			out.WriteString(strconv.Itoa(x))
		}
		if i+1 < k {
			out.WriteByte('\n')
		}
	}
	return out.String()
}

func runCase(bin string, idx int, n int) error {
	input := fmt.Sprintf("%d\n", n)
	expect := solveCase(n)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("case %d failed: %v\nstderr: %s", idx, err, string(out))
	}
	got := strings.TrimSpace(string(out))
	if got != strings.TrimSpace(expect) {
		return fmt.Errorf("case %d failed:\nexpected:\n%s\n----\ngot:\n%s\ninput:\n%s", idx, expect, got, input)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	testcases := parseTestcases()
	for i, n := range testcases {
		if err := runCase(bin, i+1, n); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
