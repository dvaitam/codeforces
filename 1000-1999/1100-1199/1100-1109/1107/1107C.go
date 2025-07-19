package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   var s string
   fmt.Fscan(reader, &s)

   var ans int64
   for i := 0; i < n; {
       j := i
       for j+1 < n && s[j+1] == s[i] {
           j++
       }
       // sort block a[i:j+1] descending
       block := a[i : j+1]
       sort.Ints(block)
       // reverse
       for l, r := 0, len(block)-1; l < r; l, r = l+1, r-1 {
           block[l], block[r] = block[r], block[l]
       }
       // sum top m or until end
       limit := m
       if len(block) < m {
           limit = len(block)
       }
       for k := 0; k < limit; k++ {
           ans += int64(block[k])
       }
       i = j + 1
   }
   fmt.Println(ans)
}
