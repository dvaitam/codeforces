package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var r, c int
   if _, err := fmt.Fscan(reader, &r, &c); err != nil {
       return
   }
   rowsFree := make([]bool, r)
   colsFree := make([]bool, c)
   // initially assume all rows and cols are free
   for i := 0; i < r; i++ {
       rowsFree[i] = true
   }
   for j := 0; j < c; j++ {
       colsFree[j] = true
   }
   // read grid and mark rows/cols with strawberries
   for i := 0; i < r; i++ {
       var s string
       fmt.Fscan(reader, &s)
       for j, ch := range s {
           if ch == 'S' {
               rowsFree[i] = false
               colsFree[j] = false
           }
       }
   }
   // count free rows and cols
   freeRows := 0
   for _, ok := range rowsFree {
       if ok {
           freeRows++
       }
   }
   freeCols := 0
   for _, ok := range colsFree {
       if ok {
           freeCols++
       }
   }
   // cells eaten: eat all free rows, then eat in free cols the remaining cells
   result := freeRows*c + freeCols*(r-freeRows)
   fmt.Println(result)
}
