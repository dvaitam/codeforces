package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   a []int64
   k  int64
)

func solve(l, r int) int64 {
   if l >= r {
       return 0
   }
   mid := (l + r) >> 1
   var ans int64
   ans += solve(l, mid)
   ans += solve(mid+1, r)
   // build left sums and maxes
   leftLen := mid - l + 1
   rightLen := r - mid
   leftMax := make([]int64, leftLen)
   leftSumMod := make([]int, leftLen)
   rightMax := make([]int64, rightLen)
   rightSumMod := make([]int, rightLen)
   // left side from mid to l
   var sum int64
   var mx int64
   for i := 0; i < leftLen; i++ {
       v := a[mid-i]
       sum += v
       if v > mx {
           mx = v
       }
       leftMax[i] = mx
       leftSumMod[i] = int(sum % k)
   }
   // right side from mid+1 to r
   sum = 0
   mx = 0
   for j := 0; j < rightLen; j++ {
       v := a[mid+1+j]
       sum += v
       if v > mx {
           mx = v
       }
       rightMax[j] = mx
       rightSumMod[j] = int(sum % k)
   }
   // case A: left max >= right max, max is leftMax
   cntR := make(map[int]int)
   pR := 0
   for i := 0; i < leftLen; i++ {
       lm := leftMax[i] % k
       // include rights with rightMax <= leftMax[i]
       for pR < rightLen && rightMax[pR] <= leftMax[i] {
           cntR[rightSumMod[pR]]++
           pR++
       }
       // target right sum mod to match
       target := int((lm - int64(leftSumMod[i])%k + k) % k)
       ans += int64(cntR[target])
   }
   // case B: right max > left max, max is rightMax
   cntL := make(map[int]int)
   pL := 0
   for j := 0; j < rightLen; j++ {
       rm := rightMax[j] % k
       // include lefts with leftMax < rightMax[j]
       for pL < leftLen && leftMax[pL] < rightMax[j] {
           cntL[leftSumMod[pL]]++
           pL++
       }
       // target left sum mod to match
       target := int((rm - int64(rightSumMod[j])%k + k) % k)
       ans += int64(cntL[target])
   }
   return ans
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n int
   fmt.Fscan(reader, &n, &k)
   a = make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   res := solve(0, n-1)
   fmt.Fprintln(writer, res)
}
