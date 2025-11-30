package main

import (
    "bufio"
    "bytes"
    "encoding/base64"
    "fmt"
    "os"
    "os/exec"
    "strconv"
    "strings"
)

const maxA = 1000000

// Base64-encoded contents of testcasesG1.txt.
const testcasesG1B64 = "MTAwCjQKMTYgOSAyIDEKNQoxOSAxNiAxMiAxMSAxCjcKMTYgNyAxNCAxOCAxOCA0IDcKNwoyMCAzIDE0IDExIDMgMTIgMTQKNwoxNSA0IDcgMTAgNCAyIDE5CjYKMTIgMTYgNyAxNyAxOSAxNwozCjEyIDggMjAKOQoxMCAxMiAxOSA0IDMgMTcgMTcgNyA0CjcKMTAgNyAxMyAxNiA4IDUgMjAKNgoxNyAxIDcgNiAxIDExCjcKMTIgMTMgMTcgMTMgMTAgNSAxNgozCjYgMTQgMjAKOQo0IDE1IDggMyAyMCAxNSAxNSAxMyAzCjkKMTYgMTAgMTQgMyA3IDkgMTUgMTYgNgozCjEgMTggNAo3CjE5IDEyIDcgOSAxNyAxNSAxMQo3CjE0IDE0IDIwIDE2IDkgMjAgMTYKMTAKMTYgNSAxMyAxNiAxMCAxNSAxMSAxMiA2IDIwCjkKMjAgOSAxMSAxMyAxNiA2IDEwIDE4IDEKMTAKMiA2IDEgMjAgMTkgNCAxMiAxMiAxNiAxOQozCjcgNSA5CjMKMTQgMTcgMTYKNAoxNiA4IDQgMTIKOAo1IDggMjAgMTEgNSAyIDQgNAozCjE2IDE1IDMKMwoxOCA1IDIwCjQKNiAyMCA5IDE1CjEwCjEgNSA3IDE1IDE2IDE1IDE1IDE4IDkgMTYKNAoxMCAxMiAxMCAxMgozCjMgMTcgOQo3CjkgMTYgMTQgMTUgMTIgMiAzCjYKMjAgOCAxOCAyIDYgMTEKOAoyIDIgNyAxNCAxMSAzIDYgNAo5CjE4IDUgMTUgMTYgMyA1IDMgMTIgMTcKMwo2IDE3IDUKNQoxMCA1IDE3IDE5IDEzCjYKOCAxOSA1IDE0IDYgMTEKNgoxOCAxOSAxIDE5IDEzIDcKOAoxMyAxNyAxNSAxIDEzIDE3IDE4IDUKNgoxOCAxMiA5IDYgNSA2CjUKMTMgMyAyIDEgMgoxMAoyIDQgMTggNCAxMiAyMCAyIDIgMTggOAo5CjggMTkgMTAgMyAxNiAzIDE5IDIgNgoxMAo0IDQgNiAyIDEgMTQgMTAgOSAxNSA4CjYKMTggNyA5IDE5IDEgOQo4CjE3IDMgMTkgMjAgMiAxMiAxMiA0CjYKMSAxMiAxMSAxNCAyIDEyCjUKMSAxIDExIDMgMTYKMwoxNSA5IDMKNwo4IDUgNSAxNCAxOCA2IDgKOQoyMCAxNCAxOSAxNyAxNiA5IDkgNSAxMAozCjcgMTUgNAo1CjE5IDQgMSAxNyA2CjQKNSAxOSAxMyAxNgoxMAoxNyAxOCAxNiA3IDEzIDUgMTUgMSAxMSAxNAo3CjExIDUgMTAgMTIgMTkgMTIgOAozCjEwIDYgMwo5CjIwIDIwIDggMyA3IDEzIDE0IDUgNQozCjE3IDE3IDE2CjgKMTEgMTkgMTMgOSAzIDkgNSAyCjQKMTkgMjAgMjAgNQo0CjMgMyA0IDcKOAoxMSA2IDggNiA2IDggNSAxNgoxMAoyMCA0IDIwIDUgMTUgNyAxOSAxMiA1IDExCjcKMyA5IDEgMiAxIDE2IDE1CjgKMiAxIDkgNCAxMyAyMCAxMSA0CjgKNSA1IDEyIDE2IDEzIDEgMTggNQo0CjE1IDEyIDEgNQoxMAo5IDMgMTYgNCAxMyA3IDE4IDYgMTUgMTEKNwoxOCAxMiAxMCAxOCA5IDE1IDE4CjMKMTcgOCAyMAozCjE3IDE2IDE4CjcKNiAyIDExIDEgNyAxOSAxMAo2CjUgOSAyIDkgMTAgMTIKNgo2IDIwIDEgMyAxMiAxOQo1CjEyIDYgMTYgOCAxNgo5CjE4IDIwIDEyIDQgMyAxNCAyIDE5IDIwCjYKMjAgMTUgOSA4IDIgMTYKNwoxNyA5IDEyIDEwIDE0IDEyIDQKOQo2IDkgMiAxNyAxMSAxNyAyIDMgNAozCjE1IDE0IDE4CjQKNiA5IDEyIDEKNAo5IDEwIDEyIDEzCjkKMTkgMyA0IDEwIDE2IDE1IDkgMTAgOAozCjIgMTcgMTgKOQo1IDggMTAgMSA0IDkgNSAxOCA0CjMKMTMgOSAxNgo2CjMgOSA2IDUgNSA0CjMKNyAyMCAxMQo5CjE2IDIwIDE5IDE2IDE1IDExIDIwIDE1IDcKNQoxOCA0IDggNSAxMwo2CjE5IDEyIDE1IDE5IDkgMTkKNgoxOSA4IDEgMTQgOCAxMwo0CjIgMTEgMTIgMTYKNwoxNCAxMyAxMCAxMiAxNyAxOSAxMgozCjE0IDE5IDYKNwoxMCAxOSAxNiA4IDQgMTUgMTkKNQoxMiAxOSAxNCAxNSAxMQo0CjE1IDE0IDE2IDEKNAoyMCAzIDE2IDE0Cg=="

type testCase struct {
    n int
    arr []int
}

// Embedded solver logic from 1822G1.go.
func solve(tc testCase) string {
    freq := make(map[int]int)
    used := make([]int, 0, len(tc.arr))
    for _, v := range tc.arr {
        if freq[v] == 0 {
            used = append(used, v)
        }
        freq[v]++
    }

    divisors := func(x int) []int {
        res := []int{}
        for d := 1; d*d <= x; d++ {
            if x%d == 0 {
                res = append(res, d)
                if d != x/d {
                    res = append(res, x/d)
                }
            }
        }
        return res
    }

    var ans int64
    for _, x := range used {
        cx := freq[x]
        divs := divisors(x)
        for _, b := range divs {
            if x*b > maxA {
                continue
            }
            y := x / b
            z := x * b
            cy := freq[y]
            cz := freq[z]
            if cy == 0 || cz == 0 {
                continue
            }
            if y == x && z == x {
                if cx >= 3 {
                    ans += int64(cx) * int64(cx-1) * int64(cx-2)
                }
            } else if y == x {
                if cx >= 2 {
                    ans += int64(cx) * int64(cx-1) * int64(cz)
                }
            } else if z == x {
                if cx >= 2 {
                    ans += int64(cy) * int64(cx) * int64(cx-1)
                }
            } else {
                ans += int64(cx) * int64(cy) * int64(cz)
            }
        }
    }
    return strconv.FormatInt(ans, 10)
}

func parseTestcases() ([]testCase, error) {
    raw, err := base64.StdEncoding.DecodeString(testcasesG1B64)
    if err != nil {
        return nil, err
    }
    sc := bufio.NewScanner(bytes.NewReader(raw))
    sc.Split(bufio.ScanWords)
    if !sc.Scan() {
        return nil, fmt.Errorf("invalid test data")
    }
    t, err := strconv.Atoi(sc.Text())
    if err != nil {
        return nil, fmt.Errorf("parse t: %v", err)
    }
    cases := make([]testCase, 0, t)
    for i := 0; i < t; i++ {
        if !sc.Scan() {
            return nil, fmt.Errorf("case %d missing n", i+1)
        }
        n, err := strconv.Atoi(sc.Text())
        if err != nil {
            return nil, fmt.Errorf("case %d n: %v", i+1, err)
        }
        arr := make([]int, n)
        for j := 0; j < n; j++ {
            if !sc.Scan() {
                return nil, fmt.Errorf("case %d truncated array", i+1)
            }
            val, err := strconv.Atoi(sc.Text())
            if err != nil {
                return nil, fmt.Errorf("case %d a[%d]: %v", i+1, j, err)
            }
            arr[j] = val
        }
        cases = append(cases, testCase{n: n, arr: arr})
    }
    return cases, nil
}

func buildIfGo(path string) (string, func(), error) {
    if strings.HasSuffix(path, ".go") {
        tmp, err := os.CreateTemp("", "bin*")
        if err != nil {
            return "", nil, err
        }
        tmp.Close()
        if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
            os.Remove(tmp.Name())
            return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
        }
        return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
    }
    return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    var errBuf bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &errBuf
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
    }
    return strings.TrimSpace(out.String()), nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("usage: go run verifierG1.go /path/to/binary")
        os.Exit(1)
    }

    cases, err := parseTestcases()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    bin, cleanup, err := buildIfGo(os.Args[1])
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    defer cleanup()

    for idx, tc := range cases {
        var input strings.Builder
        input.WriteString("1\n")
        input.WriteString(fmt.Sprintf("%d\n", tc.n))
        for i, v := range tc.arr {
            if i > 0 {
                input.WriteByte(' ')
            }
            input.WriteString(strconv.Itoa(v))
        }
        input.WriteByte('\n')

        want := solve(tc)
        got, err := runCandidate(bin, input.String())
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
            os.Exit(1)
        }
        if got != want {
            fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, input.String(), want, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", len(cases))
}

