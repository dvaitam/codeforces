package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]int, n)
   freq := make([]int, n+2)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
       if a[i] <= n+1 {
           freq[a[i]]++
       }
   }
   // compute mex of full array
   mex := 1
   for ; mex <= n+1; mex++ {
       if freq[mex] == 0 {
           break
       }
   }
   // check if array is constant
   constVal := a[0]
   allConst := true
   for i := 1; i < n; i++ {
       if a[i] != constVal {
           allConst = false
           break
       }
   }
   ans := mex + 1
   if allConst {
       // minimal k such that freq[k]==n is constVal
       if constVal < ans {
           ans = constVal
       }
   }
   fmt.Println(ans)
}
