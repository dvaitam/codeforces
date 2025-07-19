package main

import (
   "bufio"
   "fmt"
   "os"
)

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, T, k int
   if _, err := fmt.Fscan(reader, &n, &T, &k); err != nil {
       return
   }
   t := make([]int, T+2)
   for i := 1; i <= T; i++ {
       fmt.Fscan(reader, &t[i])
   }
   capArr := make([]int, T+2)
   lowp := make([]int, T+2)
   num := make([]int, T+2)
   sum := 0
   maxSum := 0
   // compute capacities and lower bounds
   for i := 1; i <= T; i++ {
       capArr[i] = t[i] - 1
       lowp[i] = max(0, t[i]-t[i+1])
       sum += lowp[i]
       maxSum += capArr[i]
   }
   // adjust for last level
   maxSum -= capArr[T]
   sum -= lowp[T]

   need := k - t[T]
   if sum > need || maxSum < need {
       fmt.Fprintln(writer, -1)
       return
   }
   // distribute additional
   rem := need - sum
   for i := 1; i <= T && rem > 0; i++ {
       add := min(rem, capArr[i]-lowp[i])
       num[i] = lowp[i] + add
       rem -= add
   }
   // construct tree
   node := make([]int, n+2)
   node[1] = 1
   cnt := 1
   l, r := 1, 1
   // output number of nodes
   fmt.Fprintln(writer, n)
   // for each level
   for i := 1; i <= T; i++ {
       temp := cnt
       templ := l
       // connect t[i] new nodes
       for j := temp + 1; j <= temp+t[i]; j++ {
           parent := node[l]
           fmt.Fprintf(writer, "%d %d\n", j, parent)
           l++
           if l > r {
               l = templ
           }
           cnt++
           node[cnt] = j
       }
       // update window for next
       r = cnt
       // number of children at this level is t[i], but we keep only num[i] for branching next
       l = cnt - (t[i] - num[i]) + 1
   }
}
