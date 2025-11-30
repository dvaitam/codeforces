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

const testcasesF64 = "MTAwCjcgMwoxIDUgOSA4IDcgNSA4CjEgMgoxIDMKMyA0CjIgNQozIDYKMiA3CjIgMgo5IDEwCjEgMgo1IDAKMiA2IDggOSAyCjEgMgoyIDMKMiA0CjIgNQo4IDMKOSA1IDEgOSAxIDIgNyAxCjEgMgoyIDMKMSA0CjMgNQoxIDYKMiA3CjUgOAo0IDEKMyA5IDggMgoxIDIKMiAzCjMgNAo4IDAKNSA5IDUgMiA5IDYgOSA0CjEgMgoyIDMKMSA0CjQgNQozIDYKNSA3CjIgOAo1IDEKNCAzIDEgMTAgNQoxIDIKMSAzCjEgNAoyIDUKMyAwCjIgOSA3CjEgMgoxIDMKNCAzCjEwIDUgOCA4CjEgMgoxIDMKMiA0CjIgMwoxMCA2CjEgMgo0IDAKNSAyIDQgNgoxIDIKMiAzCjIgNAoxIDAKMwo0IDAKMTAgOSAxMCAyCjEgMgoxIDMKMyA0CjQgMAo3IDIgNiAyCjEgMgoxIDMKMSA0CjMgMAo4IDQgMQoxIDIKMiAzCjIgMgoyIDQKMSAyCjUgMgo3IDMgMSA5IDgKMSAyCjEgMwozIDQKNCA1CjQgMgo2IDggMTAgMwoxIDIKMSAzCjMgNAozIDEKNiA5IDUKMSAyCjIgMwozIDAKOCA3IDEwCjEgMgoyIDMKNyAyCjMgOSAxIDggMiA2IDEKMSAyCjEgMwoxIDQKNCA1CjMgNgo1IDcKNSAyCjEwIDEwIDMgNSA3CjEgMgoxIDMKMSA0CjIgNQo2IDEKNCA0IDggNyAxMCA3CjEgMgoyIDMKMyA0CjQgNQoxIDYKMyAzCjIgNSAzCjEgMgoyIDMKMSAwCjgKNiAyCjggMSA3IDQgOSAyCjEgMgoxIDMKMiA0CjQgNQozIDYKMSAxCjEKMSAwCjQKMiAxCjUgNQoxIDIKMiAzCjcgMgoxIDIKNSAzCjIgNSAzIDkgNgoxIDIKMSAzCjIgNAoxIDUKMSAwCjQKNSAyCjYgMTAgMSAxMCA4CjEgMgoyIDMKMiA0CjIgNQo0IDMKMTAgNSAxIDMKMSAyCjIgMwoyIDQKNiAyCjIgNiAxMCAxIDEgNQoxIDIKMSAzCjMgNAozIDUKMyA2CjcgMQo1IDIgOCA0IDEgNSAzCjEgMgoyIDMKMiA0CjMgNQozIDYKNCA3CjIgMAo5IDgKMSAyCjYgMgoyIDggMiA4IDcgMQoxIDIKMiAzCjMgNAoyIDUKMiA2CjcgMAoyIDIgNCA0IDEgNyAxCjEgMgoyIDMKMyA0CjMgNQo0IDYKNCA3CjQgMwoyIDYgNCA1CjEgMgoyIDMKMSA0CjYgMAoyIDEgOSA4IDQgMgoxIDIKMiAzCjIgNAoyIDUKMSA2CjQgMQoyIDQgOCA3CjEgMgoxIDMKMSA0CjggMQoxMCA3IDcgOSA4IDYgOCA4CjEgMgoxIDMKMSA0CjMgNQozIDYKMyA3CjEgOAozIDIKMTAgMyA3CjEgMgoyIDMKMiAwCjkgMQoxIDIKNCAxCjEgNSAxIDgKMSAyCjEgMwoxIDQKOCAyCjkgNyA5IDkgMSAxMCAyIDkKMSAyCjIgMwoxIDQKMyA1CjUgNgo1IDcKNCA4CjggMwoxMCAxMCA0IDEgMSAzIDUgOQoxIDIKMiAzCjEgNAo0IDUKMyA2CjMgNwo3IDgKNyAzCjcgMSAzIDMgNCA1IDYKMSAyCjEgMwoyIDQKNCA1CjIgNgo0IDcKMiAxCjYgNwoxIDIKOCAzCjggMSAyIDggMyAxIDEgMTAKMSAyCjIgMwoxIDQKMyA1CjIgNgo0IDcKNyA4CjggMAoxIDEwIDggMTAgNiAyIDEwIDUKMSAyCjIgMwoyIDQKMSA1CjUgNgoyIDcKMSA4CjcgMwo2IDQgOCA2IDIgMSAxCjEgMgoyIDMKMSA0CjIgNQoyIDYKMSA3CjcgMgoyIDMgNyAxMCA3IDIgMgoxIDIKMSAzCjEgNAo0IDUKMiA2CjYgNwoxIDMKNwo3IDAKOCA2IDUgMiA2IDIgMgoxIDIKMSAzCjIgNAozIDUKMiA2CjEgNwo0IDIKMiAxMCAzIDQKMSAyCjEgMwozIDQKMiAwCjUgNgoxIDIKNCAxCjMgOCAyIDgKMSAyCjIgMwoxIDQKMSAxCjYKNiAzCjUgNSA5IDYgMyAxMAoxIDIKMSAzCjMgNAozIDUKMiA2CjcgMQozIDQgNiA5IDQgNCAzCjEgMgoyIDMKMiA0CjEgNQoyIDYKNSA3CjEgMwoyCjIgMQo3IDUKMSAyCjMgMwo1IDYgMgoxIDIKMiAzCjYgMAo3IDcgMSA3IDYgOAoxIDIKMiAzCjIgNAo0IDUKMSA2CjMgMAo1IDIgOQoxIDIKMiAzCjcgMQo3IDcgMyA0IDggNiA5CjEgMgoyIDMKMiA0CjEgNQo0IDYKMiA3CjUgMAo4IDEwIDggMSA0CjEgMgoxIDMKMyA0CjMgNQozIDMKOCAyIDgKMSAyCjIgMwo1IDAKMiA1IDEgMSA1CjEgMgoyIDMKMiA0CjEgNQo1IDIKNSA0IDEwIDIgMQoxIDIKMiAzCjIgNAozIDUKMiAxCjMgMgoxIDIKNSAyCjkgMyAxMCA5IDQKMSAyCjIgMwozIDQKNCA1CjUgMgo4IDYgMTAgMyAzCjEgMgoxIDMKMiA0CjQgNQo4IDEKOSA1IDYgOCA3IDQgOCA4CjEgMgoyIDMKMyA0CjEgNQo0IDYKMyA3CjIgOAo4IDAKMTAgNCAxIDYgOCA3IDEgOQoxIDIKMSAzCjMgNAo0IDUKMSA2CjMgNwoxIDgKMiAwCjUgNQoxIDIKMyAyCjQgMiA3CjEgMgoyIDMKNyAxCjYgNyA3IDMgOCAzIDkKMSAyCjEgMwoxIDQKMiA1CjQgNgozIDcKNyAzCjggNyA0IDQgOCA0IDEwCjEgMgoyIDMKMSA0CjIgNQoxIDYKMiA3CjYgMAozIDQgMTAgNSAxMCAyCjEgMgoyIDMKMiA0CjQgNQoxIDYKNyAzCjggNSA4IDQgNiA1IDEKMSAyCjEgMwoxIDQKMyA1CjEgNgozIDcKMSAxCjIKNyAxCjEwIDcgOSA0IDggNCA2CjEgMgoxIDMKMiA0CjMgNQo1IDYKNCA3CjYgMgoxIDkgMSA0IDYgMgoxIDIKMiAzCjEgNAoyIDUKMyA2CjUgMgo5IDcgNSA4IDYKMSAyCjEgMwoyIDQKMSA1CjEgMwo4CjggMAo3IDggOCA4IDIgMiAyIDQKMSAyCjEgMwoyIDQKMiA1CjQgNgo1IDcKMSA4CjcgMwoxIDMgNCA4IDQgMyA1CjEgMgoyIDMKMiA0CjEgNQo1IDYKMyA3CjQgMgo4IDkgMTAgOAoxIDIKMiAzCjEgNAoxIDAKMTAKMiAxCjcgNAoxIDIKNSAwCjkgOSA3IDEgMgoxIDIKMiAzCjEgNAozIDUKNCAyCjQgNCAyIDkKMSAyCjIgMwoxIDQKNiAzCjUgMTAgMyAzIDEgOQoxIDIKMiAzCjMgNAoxIDUKMiA2CjcgMQozIDkgMiAzIDQgOCAxMAoxIDIKMSAzCjMgNAoyIDUKMiA2CjQgNwo2IDEKOCAyIDEwIDEgOSAxMAoxIDIKMiAzCjIgNAozIDUKMSA2CjQgMQo4IDggOSA2CjEgMgoyIDMKMSA0CjcgMQo2IDUgNyAxIDQgMSA2CjEgMgoyIDMKMiA0CjIgNQozIDYKMyA3Cg=="

type testCase struct {
    n, k int
    a    []int
    edges [][2]int
}

func buildIfGo(path string) (string, func(), error) {
    if strings.HasSuffix(path, ".go") {
        tmp, err := os.CreateTemp("", "solbin*")
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
    var errb bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &errb
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
    }
    return strings.TrimSpace(out.String()), nil
}

func parseTestcases() ([]testCase, error) {
    raw, err := base64.StdEncoding.DecodeString(testcasesF64)
    if err != nil {
        return nil, err
    }
    sc := bufio.NewScanner(bytes.NewReader(raw))
    sc.Split(bufio.ScanWords)
    if !sc.Scan() {
        return nil, fmt.Errorf("empty data")
    }
    t, err := strconv.Atoi(sc.Text())
    if err != nil {
        return nil, fmt.Errorf("parse t: %v", err)
    }
    cases := make([]testCase, 0, t)
    for i := 0; i < t; i++ {
        var n, k int
        for _, dst := range []*int{&n, &k} {
            if !sc.Scan() {
                return nil, fmt.Errorf("truncated at case %d", i+1)
            }
            v, err := strconv.Atoi(sc.Text())
            if err != nil {
                return nil, fmt.Errorf("case %d: %v", i+1, err)
            }
            *dst = v
        }
        a := make([]int, n)
        for j := 0; j < n; j++ {
            if !sc.Scan() {
                return nil, fmt.Errorf("case %d truncated a", i+1)
            }
            val, err := strconv.Atoi(sc.Text())
            if err != nil {
                return nil, fmt.Errorf("case %d a[%d]: %v", i+1, j, err)
            }
            a[j] = val
        }
        edges := make([][2]int, n-1)
        for j := 0; j < n-1; j++ {
            if !sc.Scan() {
                return nil, fmt.Errorf("case %d truncated edge u", i+1)
            }
            u, err := strconv.Atoi(sc.Text())
            if err != nil {
                return nil, fmt.Errorf("case %d edge u%d: %v", i+1, j, err)
            }
            if !sc.Scan() {
                return nil, fmt.Errorf("case %d truncated edge v", i+1)
            }
            v, err := strconv.Atoi(sc.Text())
            if err != nil {
                return nil, fmt.Errorf("case %d edge v%d: %v", i+1, j, err)
            }
            edges[j] = [2]int{u, v}
        }
        cases = append(cases, testCase{n: n, k: k, a: a, edges: edges})
    }
    return cases, nil
}

func gcd(a, b int) int {
    for b != 0 {
        a, b = b, a%b
    }
    return a
}

// Embedded solver logic from 1778F.go (current implementation): uses gcd of values and optionally multiplies.
func solve(tc testCase) int {
    g := tc.a[0]
    for i := 1; i < tc.n; i++ {
        g = gcd(g, tc.a[i])
    }
    ans := tc.a[0]
    if tc.k > 0 {
        ans *= g
    }
    return ans
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("usage: go run verifierF.go /path/to/binary")
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
        fmt.Fprintf(&input, "1\n%d %d\n", tc.n, tc.k)
        for i, v := range tc.a {
            if i > 0 {
                input.WriteByte(' ')
            }
            input.WriteString(strconv.Itoa(v))
        }
        input.WriteByte('\n')
        for _, e := range tc.edges {
            fmt.Fprintf(&input, "%d %d\n", e[0], e[1])
        }

        want := solve(tc)
        gotStr, err := runCandidate(bin, input.String())
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
            os.Exit(1)
        }
        got, err := strconv.Atoi(strings.TrimSpace(gotStr))
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d parse output: %v\n", idx+1, err)
            os.Exit(1)
        }
        if got != want {
            fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx+1, want, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", len(cases))
}
