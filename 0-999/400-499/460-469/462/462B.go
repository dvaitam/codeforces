package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   var s string
   fmt.Fscan(reader, &s)
   cnts := make([]int64, 26)
   for i := 0; i < n; i++ {
       cnts[s[i]-'A']++
   }
   sort.Slice(cnts, func(i, j int) bool { return cnts[i] > cnts[j] })
   want := int64(k)
   var ans int64
   for _, c := range cnts {
       if want <= 0 {
           break
       }
       var t int64
       if c < want {
           t = c
       } else {
           t = want
       }
       ans += t * t
       want -= t
   }
   fmt.Println(ans)
}
