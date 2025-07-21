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
   size := 1 << (n + 1)
   a := make([]int, size)
   for i := 2; i < size; i++ {
       if _, err := fmt.Fscan(reader, &a[i]); err != nil {
           return
       }
   }
   dp := make([]int, size)
   var res int64
   // process from bottom internal nodes up to root
   // internal nodes are 1..size/2-1
   for i := size/2 - 1; i >= 1; i-- {
       left, right := 2*i, 2*i+1
       sL := a[left] + dp[left]
       sR := a[right] + dp[right]
       if sL > sR {
           res += int64(sL - sR)
           sR = sL
       } else {
           res += int64(sR - sL)
           sL = sR
       }
       // set dp to equalized maximum path sum
       dp[i] = sL
   }
   fmt.Println(res)
}
