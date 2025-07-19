package main
import (
    "bufio"
    "fmt"
    "os"
)
func main() {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()
    var n, m int
    if _, err := fmt.Fscan(reader, &n, &m); err != nil {
        return
    }
    t := make([]int, m)
    l := make([]int, m)
    r := make([]int, m)
    s := make([]bool, n)
    for i := 0; i < m; i++ {
        fmt.Fscan(reader, &t[i], &l[i], &r[i])
        l[i]--
        r[i]--
        if t[i] == 1 {
            // mark non-decreasing constraint positions
            for j := l[i] + 1; j <= r[i]; j++ {
                s[j] = true
            }
        }
    }
    a := make([]int, n)
    // initialize a[0]=0
    if n > 0 {
        a[0] = 0
    }
    // build array based on constraints
    for i := 1; i < n; i++ {
        if s[i] {
            a[i] = a[i-1]
        } else {
            a[i] = a[i-1] - 1
        }
    }
    // find minimum value
    const INF = 1000000000
    k := INF
    for i := 0; i < n; i++ {
        if a[i] < k {
            k = a[i]
        }
    }
    // adjust to make all values >= 1
    if k <= 0 {
        absk := -k
        for i := 0; i < n; i++ {
            a[i] += absk
            a[i]++
        }
    }
    // verify all constraints
    for i := 0; i < m; i++ {
        sorted := true
        for j := l[i] + 1; j <= r[i]; j++ {
            if a[j] < a[j-1] {
                sorted = false
                break
            }
        }
        if (sorted && t[i] == 0) || (!sorted && t[i] == 1) {
            fmt.Fprint(writer, "NO")
            return
        }
    }
    // output result
    fmt.Fprintln(writer, "YES")
    for i := 0; i < n; i++ {
        if i > 0 {
            writer.WriteByte(' ')
        }
        fmt.Fprint(writer, a[i])
    }
    fmt.Fprintln(writer)
}
