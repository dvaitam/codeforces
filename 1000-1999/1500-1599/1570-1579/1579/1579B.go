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
   fmt.Fscan(reader, &t)
   for ; t > 0; t-- {
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       // check if already sorted
       sorted := true
       for i := 0; i+1 < n; i++ {
           if a[i] > a[i+1] {
               sorted = false
               break
           }
       }
       if sorted {
           writer.WriteString("0\n")
           continue
       }
       // perform shifting sort
       ops := make([][2]int, 0, n)
       for i := n - 1; i >= 0; i-- {
           maxI := 0
           for j := 1; j <= i; j++ {
               if a[j] > a[maxI] {
                   maxI = j
               }
           }
           if maxI == i {
               continue
           }
           // record operation: shift element at maxI to position i
           ops = append(ops, [2]int{maxI + 1, i + 1})
           tmp := a[maxI]
           for j := maxI; j < i; j++ {
               a[j] = a[j+1]
           }
           a[i] = tmp
       }
       // output operations
       writer.WriteString(fmt.Sprintf("%d\n", len(ops)))
       for _, op := range ops {
           writer.WriteString(fmt.Sprintf("%d %d 1\n", op[0], op[1]))
       }
   }
}
