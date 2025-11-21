package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var t int
    if _, err := fmt.Fscan(in, &t); err != nil {
        return
    }
    for ; t > 0; t-- {
        var n, k int
        fmt.Fscan(in, &n, &k)
        full := n / k
        rem := n % k
        var result []byte
        for i := 0; i < full; i++ {
            for j := 0; j < k; j++ {
                result = append(result, byte('a'+j))
            }
        }
        for j := 0; j < rem; j++ {
            result = append(result, byte('a'+j))
        }
        fmt.Fprintln(out, string(result))
    }
}
