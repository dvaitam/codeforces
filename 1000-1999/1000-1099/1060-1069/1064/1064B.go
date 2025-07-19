package main

import (
    "bufio"
    "fmt"
    "math/bits"
    "os"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    var t int
    fmt.Fscan(reader, &t)
    for i := 0; i < t; i++ {
        var x uint64
        fmt.Fscan(reader, &x)
        cnt := bits.OnesCount64(x)
        ans := uint64(1) << cnt
        fmt.Println(ans)
    }
}
