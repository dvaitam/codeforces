package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

var k int
var strs []string

// dfs processes strs[l:r] at depth d, returns dp slice where dp[x]=max representativity picking x strings, and total available strings count
func dfs(l, r, d int) ([]int64, int) {
   total := 0
   // initial dp for zero picks
   dp := make([]int64, 1)
   dp[0] = 0
   // count strings ending here (length == d)
   i := l
   for i < r && len(strs[i]) == d {
       i++
   }
   leafCount := i - l
   if leafCount > 0 {
       newTotal := leafCount
       if newTotal > k {
           newTotal = k
       }
       newDp := make([]int64, newTotal+1)
       // merge dp (size total=0) with leafCount copies (dpLeaf[x]=0)
       for t1 := 0; t1 <= leafCount && t1 <= k; t1++ {
           // dp[0] + 0
           newDp[t1] = 0
       }
       dp = newDp
       total = newTotal
   }
   // process each group by next character
   for j := i; j < r; {
       ch := strs[j][d]
       m := j + 1
       for m < r && len(strs[m]) > d && strs[m][d] == ch {
           m++
       }
       // dfs on group [j:m]
       childDp, childCnt := dfs(j, m, d+1)
       // merge dp and childDp
       newTotal := total + childCnt
       if newTotal > k {
           newTotal = k
       }
       newDp := make([]int64, newTotal+1)
       // initialize as zero
       for t0 := 0; t0 <= total; t0++ {
           for t1 := 0; t1 <= childCnt && t0+t1 <= k; t1++ {
               // childDp[t1] is defined for t1<=childCnt
               val := dp[t0] + childDp[t1]
               if val > newDp[t0+t1] {
                   newDp[t0+t1] = val
               }
           }
       }
       dp = newDp
       total = newTotal
       j = m
   }
   // add representativity at this depth (prefix length d > 0)
   if d > 0 {
       for x := 0; x <= total; x++ {
           dp[x] += int64(x * (x - 1) / 2)
       }
   }
   return dp, total
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   strs = make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &strs[i])
   }
   sort.Strings(strs)
   dp, _ := dfs(0, n, 0)
   if k < len(dp) {
       fmt.Fprintln(writer, dp[k])
   } else {
       fmt.Fprintln(writer, 0)
   }
}
