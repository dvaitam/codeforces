package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   for i := 1; i <= n; i++ {
       a[i] ^= a[i-1]
   }
   const maxm = 1 << 20
   odd := make([]int, maxm)
   even := make([]int, maxm)
   var ans int64
   for i := 1; i <= n; i += 2 {
       x := a[i]
       ans += int64(odd[x])
       odd[x]++
   }
   for i := 0; i <= n; i += 2 {
       x := a[i]
       ans += int64(even[x])
       even[x]++
   }
   fmt.Println(ans)
}

