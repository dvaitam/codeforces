package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   fmt.Fscan(in, &n)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   var m int
   fmt.Fscan(in, &m)
   b := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &b[i])
   }
   // collect values for compression
   value := make([]int, 0, n+m)
   for _, v := range a {
       if v != -1 {
           value = append(value, v)
       }
   }
   for _, v := range b {
       value = append(value, v)
   }
   sort.Ints(value)
   uniq := make([]int, 0, len(value))
   prev := 0
   first := true
   for _, v := range value {
       if first || v != prev {
           uniq = append(uniq, v)
           prev = v
           first = false
       }
   }
   k := len(uniq)
   // build 1-based value map
   uniqVal := make([]int, k+2)
   for i, v := range uniq {
       uniqVal[i+1] = v
   }
   inf := k + 1
   // compress arrays
   aComp := make([]int, n)
   for i, v := range a {
       if v == -1 {
           aComp[i] = -1
       } else {
           idx := sort.SearchInts(uniq, v) + 1
           aComp[i] = idx
       }
   }
   bComp := make([]int, m)
   for i, v := range b {
       bComp[i] = sort.SearchInts(uniq, v) + 1
   }
   sort.Ints(bComp)
   // nxt[i]: first index in bComp with value > i
   nxt := make([]int, k+2)
   for i := 0; i <= k; i++ {
       j := sort.Search(m, func(j int) bool { return bComp[j] > i })
       nxt[i] = j
   }
   // DP for LIS with gaps
   c := make([]int, n+2)
   for i := range c {
       c[i] = inf
   }
   c[0] = 0
   pos := make([]int, n+2)
   pre := make([]int, n)
   for i := range pre {
       pre[i] = -1
   }
   pos[0] = -1
   length := 0
   for i := 0; i < n; i++ {
       ai := aComp[i]
       if ai != -1 {
           // find largest l in [0..length] with c[l] < ai
           l, r := 0, length
           for l < r {
               mid := (l + r + 1) >> 1
               if c[mid] < ai {
                   l = mid
               } else {
                   r = mid - 1
               }
           }
           if c[l+1] > ai {
               c[l+1] = ai
               pos[l+1] = i
           }
           pre[i] = pos[l]
           length = max(length, l+1)
       } else {
           for j := length; j >= 0; j-- {
               x := c[j]
               p := nxt[x]
               if p >= m {
                   continue
               }
               nextVal := bComp[p]
               if nextVal < c[j+1] {
                   c[j+1] = nextVal
                   pos[j+1] = pos[j]
                   length = max(length, j+1)
               }
           }
       }
   }
   // reconstruct LIS
   mark := make([]bool, n)
   valComp := make([]int, 0, length)
   curr := pos[length]
   for curr != -1 {
       if aComp[curr] != -1 {
           valComp = append(valComp, aComp[curr])
           mark[curr] = true
       }
       curr = pre[curr]
   }
   // reverse
   for i, j := 0, len(valComp)-1; i < j; i, j = i+1, j-1 {
       valComp[i], valComp[j] = valComp[j], valComp[i]
   }
   // add sentinel
   valComp = append(valComp, inf)
   // fill gaps greedily preserving LIS
   nowCount := 0
   last := 0
   nodeB := 0
   used := make([]bool, m)
   for i := 0; i < n; i++ {
       if mark[i] {
           nowCount++
       }
       if aComp[i] == -1 {
           // next limit
           var nextLimit int
           if nowCount < len(valComp) {
               nextLimit = valComp[nowCount]
           } else {
               nextLimit = inf
           }
           for nodeB < m && bComp[nodeB] <= last {
               nodeB++
           }
           if nodeB >= m {
               break
           }
           if bComp[nodeB] < nextLimit {
               aComp[i] = bComp[nodeB]
               used[nodeB] = true
               last = max(last, aComp[i])
               nodeB++
           }
       }
   }
   // fill remaining gaps
   nodeB2 := 0
   for i := 0; i < n; i++ {
       if aComp[i] == -1 {
           for nodeB2 < m && used[nodeB2] {
               nodeB2++
           }
           aComp[i] = bComp[nodeB2]
           nodeB2++
       }
   }
   // output
   for i := 0; i < n; i++ {
       out.WriteString(fmt.Sprint(uniqVal[aComp[i]]))
       if i+1 < n {
           out.WriteByte(' ')
       }
   }
   out.WriteByte('\n')
}
