package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

const M = 200010

var vis [M]bool
type pair struct{ first, second int }
var p [M]pair

func initVis() {
   for i := 0; i < 310; i++ {
       for j := 0; j < 310; j++ {
           s := i*i + j*j
           if s < M {
               vis[s] = true
               p[s] = pair{i, j}
           }
       }
   }
}

// getSign assigns signs to arr to minimize angle ordering variance
func getSign(arr []int) {
   type pi struct{ val, idx int }
   n := len(arr)
   arrp := make([]pi, n)
   for i := 0; i < n; i++ {
       arrp[i] = pi{arr[i], i}
   }
   sort.Slice(arrp, func(i, j int) bool { return arrp[i].val > arrp[j].val })
   d := 0
   for i := 0; i+1 < n; i += 2 {
       if arrp[i].val == arrp[i+1].val {
           idx := arrp[i].idx
           arr[idx] = -arr[idx]
       } else {
           diff := arrp[i].val - arrp[i+1].val
           if d <= 0 {
               idx := arrp[i+1].idx
               arr[idx] = -arr[idx]
               d += diff
           } else {
               idx := arrp[i].idx
               arr[idx] = -arr[idx]
               d -= diff
           }
       }
   }
   if n%2 == 1 {
       idx := arrp[n-1].idx
       arr[idx] = -arr[idx]
   }
}

func main() {
   initVis()
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   x := make([]int, n)
   y := make([]int, n)
   cnt := 0
   sum := 0
   i := 1
   // select n-1 smallest sums of two squares
   for cnt < n-1 {
       if i < M && vis[i] {
           x[cnt] = p[i].first
           y[cnt] = p[i].second
           sum += i
           cnt++
       }
       i++
   }
   // find two candidates for last element
   var a, b, tcount int
   for j := i; tcount < 2; j++ {
       if j < M && vis[j] {
           if tcount == 0 {
               a = j
           } else {
               b = j
           }
           tcount++
       }
   }
   ok := false
   // choose last based on parity
   if (sum+a)%2 == 0 {
       x[cnt] = p[a].first
       y[cnt] = p[a].second
   } else if (sum+b)%2 == 0 {
       x[cnt] = p[b].first
       y[cnt] = p[b].second
   } else {
       x[cnt] = p[a].first
       y[cnt] = p[a].second
       if sum%2 != 0 {
           x[0] = p[b].first
           y[0] = p[b].second
       } else {
           if n >= 4 {
               x[3] = p[b].first
               y[3] = p[b].second
           } else {
               x[0] = p[b].first
               y[0] = p[b].second
           }
       }
       ok = true
   }
   cnt++
   // adjust parity if needed
   sum = 0
   for k := 0; k < n; k++ {
       sum += x[k]
   }
   if sum%2 != 0 {
       if !ok {
           x[0], y[0] = y[0], x[0]
       } else {
           if n >= 4 {
               x[3], y[3] = y[3], x[3]
           } else {
               x[0], y[0] = y[0], x[0]
           }
       }
   }
   // assign signs
   getSign(x)
   getSign(y)
   // compute angles and order
   ang := make([]float64, n)
   for k := 0; k < n; k++ {
       ang[k] = math.Atan2(float64(y[k]), float64(x[k]))
   }
   ids := make([]int, n)
   for k := 0; k < n; k++ {
       ids[k] = k
   }
   sort.Slice(ids, func(i, j int) bool {
       return ang[ids[i]] < ang[ids[j]]
   })
   // output
   fmt.Fprintln(writer, "YES")
   tx, ty := 0, 0
   for _, idx := range ids {
       tx += x[idx]
       ty += y[idx]
       fmt.Fprintln(writer, tx, ty)
   }
}
