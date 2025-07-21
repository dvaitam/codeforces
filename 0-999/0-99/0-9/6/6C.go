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
   times := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &times[i])
   }
   left, right := 0, n-1
   var timeA, timeB int64
   var countA, countB int
   for left <= right {
       if timeA <= timeB {
           timeA += int64(times[left])
           left++
           countA++
       } else {
           timeB += int64(times[right])
           right--
           countB++
       }
   }
   fmt.Println(countA, countB)
