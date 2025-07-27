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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for tc := 0; tc < t; tc++ {
       var n, k int
       fmt.Fscan(reader, &n, &k)
       var aStr, bStr string
       fmt.Fscan(reader, &aStr)
       fmt.Fscan(reader, &bStr)
       a := []byte(aStr)
       b := []byte(bStr)
       if solve(n, k, a, b) {
           fmt.Fprintln(writer, "Yes")
       } else {
           fmt.Fprintln(writer, "No")
       }
   }
}

func solve(n, k int, a, b []byte) bool {
   var cntA [26]int
   var cntB [26]int
   for i := 0; i < n; i++ {
       cntA[a[i]-'a']++
       cntB[b[i]-'a']++
   }
   carry := 0
   for i := 0; i < 26; i++ {
       available := cntA[i] + carry
       if available < cntB[i] {
           return false
       }
       diff := available - cntB[i]
       if diff % k != 0 {
           return false
       }
       carry = diff
   }
   // After processing 'z', no surplus should remain
   if carry != 0 {
       return false
   }
   return true
}
