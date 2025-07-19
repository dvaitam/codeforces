package main

import (
   "bufio"
   "fmt"
   "os"
)

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}

// calc computes the sum as in the original C++ code
func calc(n, m int, a []int) int64 {
   var s int64
   for i := 1; i <= n; i++ {
       var xSum int64
       ySum := m - a[i]
       l, r := i-1, i+1
       for j := 1; j <= n; j++ {
           if j == i {
               continue
           }
           var weight int
           if a[j] != 0 {
               weight = m - 1
           } else {
               weight = m
           }
           xSum += int64(weight * abs(i-j))
       }
       var weightI int
       if a[i] != 0 {
           weightI = m - 1
       } else {
           weightI = m
       }
       s += xSum * int64(weightI)
       if a[i] != 0 {
           for ; l >= 0 && a[l] > a[l+1]; l-- {
               ySum += m - a[l]
           }
           for ; r <= n+1 && r < len(a) && a[r] > a[r-1]; r++ {
               ySum += m - a[r]
           }
           // a[i] >= 1 here
           s += 4 * int64(a[i]-1) * int64(ySum)
       }
   }
   return s
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var N, M int
   if _, err := fmt.Fscan(reader, &N, &M); err != nil {
       return
   }
   total := N * M
   x := make([]int, M+2)
   y := make([]int, N+2)
   S := total
   for i := 1; i <= N; i++ {
       var line string
       fmt.Fscan(reader, &line)
       for j := 1; j <= M && j <= len(line); j++ {
           if line[j-1] == 'X' {
               y[i] = j
               x[j] = i
               S--
           }
       }
   }
   sum1 := calc(N, M, y)
   sum2 := calc(M, N, x)
   denom := float64(S) * float64(S)
   result := (float64(sum1) + float64(sum2)) / denom
   fmt.Printf("%.9f\n", result)
}
