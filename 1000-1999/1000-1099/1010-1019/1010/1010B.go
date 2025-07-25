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

   var m int64
   var n int
   // Read bounds and pattern length
   if _, err := fmt.Fscan(reader, &m, &n); err != nil {
       return
   }
   // Determine lying pattern of length n by querying y=1
   pattern := make([]int, n)
   cur := 0
   for i := 0; i < n; i++ {
       // query 1
       fmt.Fprintln(writer, 1)
       writer.Flush()
       var r int
       if _, err := fmt.Fscan(reader, &r); err != nil {
           return
       }
       // correct answer zero => found x
       if r == 0 {
           return
       }
       // record truth or lie: r == 1 means truth, r == -1 means lie
       if r == 1 {
           pattern[cur] = 1
       } else {
           pattern[cur] = 0
       }
       cur = (cur + 1) % n
   }
   // Binary search using known pattern
   left := int64(1)
   right := m
   for left <= right {
       mid := (left + right) / 2
       fmt.Fprintln(writer, mid)
       writer.Flush()
       var r int
       if _, err := fmt.Fscan(reader, &r); err != nil {
           return
       }
       if r == 0 {
           return
       }
       // correct the response according to pattern
       t := r
       if pattern[cur] == 0 {
           t = -t
       }
       if t > 0 {
           left = mid + 1
       } else {
           right = mid - 1
       }
       cur = (cur + 1) % n
   }
}
