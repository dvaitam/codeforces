package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   var p, k int64
   fmt.Fscan(reader, &n, &p, &k)
   freq := make(map[int64]int64, n)
   for i := 0; i < n; i++ {
       var a int64
       fmt.Fscan(reader, &a)
       a %= p
       a2 := (a * a) % p
       a4 := (a2 * a2) % p
       key := (a4 - (k*a)%p + p) % p
       freq[key]++
   }
   var ans int64
   for _, cnt := range freq {
       if cnt > 1 {
           ans += cnt * (cnt - 1) / 2
       }
   }
   fmt.Println(ans)
}
