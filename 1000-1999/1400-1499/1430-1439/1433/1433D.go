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
   for t > 0 {
       t--
       var n int
       fmt.Fscan(reader, &n)
       nums := make([]int, n)
       mx, mn := 0, 1000000001
       idmx, idmn := -1, -1
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &nums[i])
           if nums[i] > mx {
               mx = nums[i]
               idmx = i
           }
           if nums[i] < mn {
               mn = nums[i]
               idmn = i
           }
       }
       if mx == mn {
           fmt.Fprintln(writer, "NO")
       } else {
           fmt.Fprintln(writer, "YES")
           chk := false
           for i := 0; i < n; i++ {
               if !chk && i == idmx {
                   chk = true
                   continue
               }
               if chk && nums[i] == mx {
                   fmt.Fprintln(writer, i+1, idmn+1)
               } else {
                   fmt.Fprintln(writer, idmx+1, i+1)
               }
           }
       }
   }
}
