package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

var (
   x []int64
   y []int64
)

func sharp(a, b, c int) bool {
   // dot product (x[a]-x[b],y[a]-y[b]) . (x[c]-x[b],y[c]-y[b]) > 0
   dx1 := x[a] - x[b]
   dy1 := y[a] - y[b]
   dx2 := x[c] - x[b]
   dy2 := y[c] - y[b]
   return dx1*dx2+dy1*dy2 > 0
}

func distribute(axis []int, take int) (q1, q2 []int) {
   for i, v := range axis {
       if i < take {
           q1 = append(q1, v)
       } else {
           q2 = append(q2, v)
       }
   }
   return
}

func appendPair(line *[]int, q1, q2 []int) {
   for i := 0; i < len(q1); i++ {
       *line = append(*line, q1[i], q2[i])
   }
}

func appendPairRev(line *[]int, q1, q2 []int) {
   for i := len(q1) - 1; i >= 0; i-- {
       *line = append(*line, q1[i], q2[i])
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   fmt.Fscan(in, &n)
   x = make([]int64, n+2)
   y = make([]int64, n+2)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &x[i], &y[i])
   }
   odd := n%2 == 1
   if odd {
       n--
   }
   // order indices by x and y
   xOrder := make([]int, n)
   yOrder := make([]int, n)
   for i := 0; i < n; i++ {
       xOrder[i] = i + 1
       yOrder[i] = i + 1
   }
   sort.Slice(xOrder, func(i, j int) bool { return x[xOrder[i]] < x[xOrder[j]] })
   sort.Slice(yOrder, func(i, j int) bool { return y[yOrder[i]] < y[yOrder[j]] })
   xMedian := x[xOrder[n/2]]
   yMedian := y[yOrder[n/2]]

   // place array not needed explicitly
   quarter := make([][]int, 4)
   axis := make([][]int, 4)
   center := make([]int, 0)
   for i := 0; i < 4; i++ {
       quarter[i] = make([]int, 0)
       axis[i] = make([]int, 0)
   }
   for i := 1; i <= n; i++ {
       p := 3*(1+boolToInt(y[i]<yMedian)-boolToInt(y[i]>yMedian)) +
           (1 + boolToInt(x[i]<xMedian) - boolToInt(x[i]>xMedian))
       switch p {
       case 0:
           quarter[0] = append(quarter[0], i)
       case 1:
           axis[1] = append(axis[1], i)
       case 2:
           quarter[1] = append(quarter[1], i)
       case 3:
           axis[0] = append(axis[0], i)
       case 4:
           center = append(center, i)
       case 5:
           axis[2] = append(axis[2], i)
       case 6:
           quarter[3] = append(quarter[3], i)
       case 7:
           axis[3] = append(axis[3], i)
       case 8:
           quarter[2] = append(quarter[2], i)
       }
   }
   upCount := len(quarter[0]) + len(quarter[1]) + len(axis[1])
   downCount := len(quarter[2]) + len(quarter[3]) + len(axis[3])
   rightCount := len(quarter[0]) + len(quarter[3]) + len(axis[0])
   leftCount := len(quarter[1]) + len(quarter[2]) + len(axis[2])
   for len(center) > 0 {
       v := center[len(center)-1]
       center = center[:len(center)-1]
       if len(axis[0])+len(axis[2]) < abs(downCount-upCount) {
           if leftCount > rightCount {
               axis[0] = append(axis[0], v)
               rightCount++
           } else {
               axis[2] = append(axis[2], v)
               leftCount++
           }
       } else {
           if downCount > upCount {
               axis[1] = append(axis[1], v)
               upCount++
           } else {
               axis[3] = append(axis[3], v)
               downCount++
           }
       }
   }
   balance := downCount - upCount + len(axis[0]) - len(axis[2])
   take := [4]int{}
   take[0] = max(0, balance/2)
   take[2] = take[0] - balance/2
   diff := len(quarter[0]) + take[0] - len(quarter[2]) - take[2] + len(axis[1]) - len(axis[3])
   take[1] = max(0, diff)
   take[3] = take[1] - diff
   // distribute axis to quarters
   q0a, q3a := distribute(axis[0], take[0])
   quarter[0] = append(quarter[0], q0a...)
   quarter[3] = append(quarter[3], q3a...)
   q1a, q0b := distribute(axis[1], take[1])
   quarter[1] = append(quarter[1], q1a...)
   quarter[0] = append(quarter[0], q0b...)
   q2a, q1b := distribute(axis[2], take[2])
   quarter[2] = append(quarter[2], q2a...)
   quarter[1] = append(quarter[1], q1b...)
   q3b, q2b := distribute(axis[3], take[3])
   quarter[3] = append(quarter[3], q3b...)
   quarter[2] = append(quarter[2], q2b...)
   // build line
   total := make([]int, 0, n+1)
   if len(quarter[0]) == 0 {
       appendPair(&total, quarter[1], quarter[3])
   } else if len(quarter[1]) == 0 {
       appendPair(&total, quarter[0], quarter[2])
   } else {
       order := []int{0, 2, 1, 3}
       if !sharp(quarter[order[0]][len(quarter[order[0]])-1], quarter[order[1]][len(quarter[order[1]])-1], quarter[order[2]][len(quarter[order[2]])-1]) ||
           !sharp(quarter[order[1]][len(quarter[order[1]])-1], quarter[order[2]][len(quarter[order[2]])-1], quarter[order[3]][len(quarter[order[3]])-1]) {
           order[0], order[1] = order[1], order[0]
       }
       if !sharp(quarter[order[0]][len(quarter[order[0]])-1], quarter[order[1]][len(quarter[order[1]])-1], quarter[order[2]][len(quarter[order[2]])-1]) ||
           !sharp(quarter[order[1]][len(quarter[order[1]])-1], quarter[order[2]][len(quarter[order[2]])-1], quarter[order[3]][len(quarter[order[3]])-1]) {
           order[2], order[3] = order[3], order[2]
       }
       if !sharp(quarter[order[0]][len(quarter[order[0]])-1], quarter[order[1]][len(quarter[order[1]])-1], quarter[order[2]][len(quarter[order[2]])-1]) ||
           !sharp(quarter[order[1]][len(quarter[order[1]])-1], quarter[order[2]][len(quarter[order[2]])-1], quarter[order[3]][len(quarter[order[3]])-1]) {
           order[0], order[1] = order[1], order[0]
       }
       appendPair(&total, quarter[order[0]], quarter[order[1]])
       appendPairRev(&total, quarter[order[2]], quarter[order[3]])
   }
   if odd {
       idx := n + 1
       total = append(total, idx)
       i := len(total) - 1
       for i > 1 && !sharp(total[i-2], total[i-1], idx) {
           total[i] = total[i-1]
           i--
       }
       if i == 1 && !sharp(total[0], idx, total[1]) {
           total[1] = total[0]
           i--
       }
       total[i] = idx
   }
   // output
   for _, v := range total {
       fmt.Fprintf(out, "%d ", v)
   }
   fmt.Fprintln(out)
}

func boolToInt(b bool) int {
   if b {
       return 1
   }
   return 0
}

func abs(a int) int {
   if a < 0 {
       return -a
   }
   return a
}

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

// helper to compare y[i] < yMedian conveniently
