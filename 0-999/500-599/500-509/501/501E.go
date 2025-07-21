package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // Collect mismatched mirror pairs
   var I, J []int
   for i := 1; i*2 <= n; i++ {
       j := n - i + 1
       if a[i] != a[j] {
           I = append(I, i)
           J = append(J, j)
       }
   }
   k := len(I)
   // If already palindrome, all segments are valid
   if k == 0 {
       // n*(n+1)/2
       nn := int64(n)
       fmt.Println(nn * (nn + 1) / 2)
       return
   }
   // Sort pairs by I
   idx := make([]int, k)
   for t := 0; t < k; t++ {
       idx[t] = t
   }
   sort.Slice(idx, func(x, y int) bool {
       return I[idx[x]] < I[idx[y]]
   })
   si := make([]int, k)
   sj := make([]int, k)
   for t, v := range idx {
       si[t] = I[v]
       sj[t] = J[v]
   }
   I = si; J = sj
   // Precompute prefix min and max for J
   preMinJ := make([]int, k)
   preMaxJ := make([]int, k)
   for t := 0; t < k; t++ {
       if t == 0 {
           preMinJ[t] = J[t]
           preMaxJ[t] = J[t]
       } else {
           if J[t] < preMinJ[t-1] {
               preMinJ[t] = J[t]
           } else {
               preMinJ[t] = preMinJ[t-1]
           }
           if J[t] > preMaxJ[t-1] {
               preMaxJ[t] = J[t]
           } else {
               preMaxJ[t] = preMaxJ[t-1]
           }
       }
   }
   // Precompute suffix min for J-1
   suffMinJm1 := make([]int, k)
   for t := k - 1; t >= 0; t-- {
       j1 := J[t] - 1
       if t == k-1 {
           suffMinJm1[t] = j1
       } else {
           if j1 < suffMinJm1[t+1] {
               suffMinJm1[t] = j1
           } else {
               suffMinJm1[t] = suffMinJm1[t+1]
           }
       }
   }
   // maxI is maximum i among all mismatches
   maxI := I[k-1]
   var ans int64
   pos := 0
   // iterate over l
   for l := 1; l <= n; l++ {
       // move edges with i < l to E2
       for pos < k && I[pos] < l {
           pos++
       }
       // E2 constraint: if any j < l, invalid
       if pos > 0 && preMinJ[pos-1] < l {
           continue
       }
       // compute r range
       // r_low = max(l, maxI (if E1 nonempty), preMaxJ (if E2 nonempty))
       rLow := l
       if maxI > rLow {
           rLow = maxI
       }
       if pos > 0 {
           if preMaxJ[pos-1] > rLow {
               rLow = preMaxJ[pos-1]
           }
       }
       // r_high
       var rHigh int
       if pos <= k-1 {
           // E1 nonempty
           rHigh = suffMinJm1[pos]
       } else {
           rHigh = n
       }
       if rHigh < rLow {
           continue
       }
       ans += int64(rHigh - rLow + 1)
   }
   fmt.Println(ans)
}
