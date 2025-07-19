package main

import (
   "bufio"
   "fmt"
   "os"
)

// Pair represents a command segment [first, second]
type Pair struct{ first, second int }

var n int
var a []int
var path []Pair

func dfs(remain int) bool {
   ok := true
   for i := 1; i <= n; i++ {
       if a[i] != i {
           ok = false
           break
       }
   }
   if ok {
       return true
   }
   if remain == 0 {
       return false
   }
   // generate candidate segments
   vc := make([]Pair, 0)
   vc = append(vc, Pair{1, n})
   pos := make([]int, n+1)
   for i := 1; i <= n; i++ {
       pos[a[i]] = i
   }
   // consecutive runs
   for i := 1; i <= n; {
       j := i
       for j+1 <= n && a[j+1]-a[j] == a[i+1]-a[i] && abs(a[i+1]-a[i]) == 1 {
           j++
       }
       vc = append(vc, Pair{i, j})
       i = j + 1
   }
   // adjacent value positions
   for i := 1; i <= n; i++ {
       if a[i] != 1 {
           j := pos[a[i]-1]
           if j+1 < i {
               vc = append(vc, Pair{j + 1, i})
           }
           if i+1 < j {
               vc = append(vc, Pair{i, j - 1})
           }
       }
       if a[i] != n {
           j := pos[a[i]+1]
           if j+1 < i {
               vc = append(vc, Pair{j + 1, i})
           }
           if i+1 < j {
               vc = append(vc, Pair{i, j - 1})
           }
       }
   }
   // try each candidate
   for _, p := range vc {
       L, R := p.first, p.second
       path = append(path, p)
       reverse(L, R)
       if dfs(remain - 1) {
           return true
       }
       reverse(L, R)
       path = path[:len(path)-1]
   }
   return false
}

// reverse reverses segment [l, r] in a
func reverse(l, r int) {
   for l < r {
       a[l], a[r] = a[r], a[l]
       l++
       r--
   }
}

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fscan(reader, &n)
   a = make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   dfs(3)
   k := len(path)
   fmt.Fprintln(writer, k)
   for i := k - 1; i >= 0; i-- {
       fmt.Fprintln(writer, path[i].first, path[i].second)
   }
}
