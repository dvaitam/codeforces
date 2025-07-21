package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var nStr string
   var m int
   fmt.Fscan(reader, &nStr, &m)
   // extract digits and sort
   digits := make([]int, len(nStr))
   for i, ch := range nStr {
       digits[i] = int(ch - '0')
   }
   sort.Ints(digits)
   // count frequencies
   type pair struct{ d, cnt int }
   var vp []pair
   for i := 0; i < len(digits); {
       j := i + 1
       for j < len(digits) && digits[j] == digits[i] {
           j++
       }
       vp = append(vp, pair{digits[i], j - i})
       i = j
   }
   L := len(vp)
   // compute base multipliers
   totalStates := 1
   bases := make([]int, L)
   for i := 0; i < L; i++ {
       bases[i] = totalStates
       totalStates *= (vp[i].cnt + 1)
   }
   totalDigits := len(digits)
   // precompute counts and sums for states
   countsState := make([][]int, totalStates)
   sumState := make([]int, totalStates)
   for idx := 0; idx < totalStates; idx++ {
       cnts := make([]int, L)
       rem := idx
       sum := 0
       for i := 0; i < L; i++ {
           c := rem / bases[i] % (vp[i].cnt + 1)
           cnts[i] = c
           sum += c
           rem -= c * bases[i]
       }
       countsState[idx] = cnts
       sumState[idx] = sum
   }
   // DP[state*m + rem]
   dp := make([]int64, totalStates*m)
   dp[0*m+0] = 1
   // fill DP
   for idx := 0; idx < totalStates; idx++ {
       s := sumState[idx]
       if s >= totalDigits {
           continue
       }
       for rem0 := 0; rem0 < m; rem0++ {
           cur := dp[idx*m+rem0]
           if cur == 0 {
               continue
           }
           // place next digit
           for i := 0; i < L; i++ {
               if countsState[idx][i] < vp[i].cnt {
                   // no leading zero
                   if s == 0 && vp[i].d == 0 {
                       continue
                   }
                   newIdx := idx + bases[i]
                   newRem := (rem0*10 + vp[i].d) % m
                   dp[newIdx*m+newRem] += cur
               }
           }
       }
   }
   // answer at full state
   lastIdx := totalStates - 1
   fmt.Println(dp[lastIdx*m+0])
}
