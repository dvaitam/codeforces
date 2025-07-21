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
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   var sizes []int
   currLen := 0
   negCnt := 0
   for _, v := range a {
       if v < 0 {
           if negCnt == 2 {
               // start new folder before this day
               sizes = append(sizes, currLen)
               currLen = 1
               negCnt = 1
           } else {
               currLen++
               negCnt++
           }
       } else {
           currLen++
       }
   }
   if currLen > 0 {
       sizes = append(sizes, currLen)
   }
   // output result
   fmt.Println(len(sizes))
   for i, sz := range sizes {
       if i > 0 {
           fmt.Print(" ")
       }
       fmt.Print(sz)
   }
   fmt.Println()
}
