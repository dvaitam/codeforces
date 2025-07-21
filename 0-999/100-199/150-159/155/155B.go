package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   zeroList := make([]int, 0, n)
   var sumA1 int64
   var extra int64
   for i := 0; i < n; i++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       if b == 0 {
           zeroList = append(zeroList, a)
       } else {
           sumA1 += int64(a)
           extra += int64(b - 1)
       }
   }
   sort.Slice(zeroList, func(i, j int) bool {
       return zeroList[i] > zeroList[j]
   })
   // can take up to extra+1 zero-cards
   k0 := int(extra + 1)
   if k0 > len(zeroList) {
       k0 = len(zeroList)
   }
   var sumA0 int64
   for i := 0; i < k0; i++ {
       sumA0 += int64(zeroList[i])
   }
   fmt.Println(sumA1 + sumA0)
}
