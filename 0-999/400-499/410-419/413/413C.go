package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   isAuction := make([]bool, n)
   for i := 0; i < m; i++ {
       var bi int
       fmt.Fscan(in, &bi)
       if bi >= 1 && bi <= n {
           isAuction[bi-1] = true
       }
   }
   // Separate auction and regular questions
   var regSum int64
   var A []int64
   for i := 0; i < n; i++ {
       if isAuction[i] {
           A = append(A, a[i])
       } else {
           regSum += a[i]
       }
   }
   sort.Slice(A, func(i, j int) bool { return A[i] < A[j] })
   totalSum := regSum
   for _, v := range A {
       totalSum += v
   }
   mA := len(A)
   // prefix sums of auction values
   PS := make([]int64, mA+1)
   for i := 1; i <= mA; i++ {
       PS[i] = PS[i-1] + A[i-1]
   }
   // best result: no doubling
   best := totalSum
   // try using d auctions for doubling: the first d smallest auctions
   for d := 1; d <= mA; d++ {
       sumD := PS[d]
       S0 := totalSum - sumD
       if S0 <= 0 {
           continue
       }
       ok := true
       // check feasibility: before j-th doubling, score S0 * 2^(j-1) > A[j-1]
       cur := S0
       for j := 1; j <= d; j++ {
           if cur <= A[j-1] {
               ok = false
               break
           }
           // double for next
           cur <<= 1
       }
       if !ok {
           continue
       }
       // final score after d doublings: S0 * 2^d
       endScore := S0 << d
       if endScore > best {
           best = endScore
       }
   }
   // output
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprint(writer, best)
}
