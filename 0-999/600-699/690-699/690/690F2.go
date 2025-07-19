package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "strconv"
   "strings"
)

var reader = bufio.NewReader(os.Stdin)
var writer = bufio.NewWriter(os.Stdout)

// key for slice of ints
func keyOf(c []int) string {
   if len(c) == 0 {
       return ""
   }
   sb := strings.Builder{}
   sb.WriteString(strconv.Itoa(c[0]))
   for i := 1; i < len(c); i++ {
       sb.WriteByte(',')
       sb.WriteString(strconv.Itoa(c[i]))
   }
   return sb.String()
}

// Global map for canonical labels
var M map[string]int

// mapT: map slice c to unique int
func mapT(c []int) int {
   k := keyOf(c)
   if v, ok := M[k]; ok {
       return v
   }
   id := len(M)
   M[k] = id
   return id
}

// canonizeIndices: remap node labels to 0..n-1 by first appearance
func canonizeIndices(e [][2]int) [][2]int {
   m2 := make(map[int]int)
   idx := 0
   for i := range e {
       u := e[i][0]
       v := e[i][1]
       if _, ok := m2[u]; !ok {
           m2[u] = idx
           idx++
       }
       if _, ok := m2[v]; !ok {
           m2[v] = idx
           idx++
       }
       e[i][0] = m2[u]
       e[i][1] = m2[v]
   }
   return e
}

// encodeSubtree: return canonical label of subtree rooted at x
func encodeSubtree(x, dad int, E [][]int) int {
   var ch []int
   for _, y := range E[x] {
       if y != dad {
           ch = append(ch, encodeSubtree(y, x, E))
       }
   }
   sort.Ints(ch)
   return mapT(ch)
}

// dfsC: compute subtree sizes and balance factors
func dfsC(x, dad int, E [][]int, sz, bal []int) {
   sz[x] = 1
   bal[x] = 0
   for _, y := range E[x] {
       if y != dad {
           dfsC(y, x, E, sz, bal)
           if sz[y] > bal[x] {
               bal[x] = sz[y]
           }
           sz[x] += sz[y]
       }
   }
}

// canonizeTree: return canonical label of tree e
func canonizeTree(e [][2]int) int {
   e = canonizeIndices(e)
   n := 0
   for _, p := range e {
       if p[0]+1 > n {
           n = p[0] + 1
       }
       if p[1]+1 > n {
           n = p[1] + 1
       }
   }
   E := make([][]int, n)
   for _, p := range e {
       E[p[0]] = append(E[p[0]], p[1])
       E[p[1]] = append(E[p[1]], p[0])
   }
   sz := make([]int, n)
   bal := make([]int, n)
   dfsC(0, -1, E, sz, bal)
   ans := -1
   for x := 0; x < n; x++ {
       rem := n - sz[x]
       if rem > bal[x] {
           bal[x] = rem
       }
       if 2*bal[x] <= n {
           h := encodeSubtree(x, -1, E)
           if ans == -1 || h < ans {
               ans = h
           }
       }
   }
   return ans
}

// dfs: mark component via t
func dfsMark(x, t int, E [][]int, bio []int) {
   bio[x] = t
   for _, y := range E[x] {
       if bio[y] == 0 {
           dfsMark(y, t, E, bio)
       }
   }
}

// canonizeForest: return canonical label of forest e
func canonizeForest(e [][2]int) int {
   e = canonizeIndices(e)
   n := 0
   for _, p := range e {
       if p[0]+1 > n {
           n = p[0] + 1
       }
       if p[1]+1 > n {
           n = p[1] + 1
       }
   }
   E := make([][]int, n)
   for _, p := range e {
       E[p[0]] = append(E[p[0]], p[1])
       E[p[1]] = append(E[p[1]], p[0])
   }
   bio := make([]int, n)
   var all []int
   for i := 0; i < n; i++ {
       if bio[i] == 0 {
           dfsMark(i, i+1, E, bio)
           var ne [][2]int
           for _, p := range e {
               if bio[p[0]] == i+1 && bio[p[1]] == i+1 {
                   ne = append(ne, p)
               }
           }
           all = append(all, canonizeTree(ne))
       }
   }
   sort.Ints(all)
   return mapT(all)
}

// solve one test case
func solve() [][2]int {
   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return nil
   }
   // clear map
   M = make(map[string]int)
   // read graphs
   v := make([][][2]int, n)
   for i := 0; i < n; i++ {
       var m int
       fmt.Fscan(reader, &m)
       w := make([][2]int, m)
       for j := 0; j < m; j++ {
           fmt.Fscan(reader, &w[j][0], &w[j][1])
           w[j][0]--
           w[j][1]--
       }
       v[i] = canonizeIndices(w)
   }
   p := 0
   for p < n && len(v[p]) != n-2 {
       p++
   }
   if p == n {
       return nil
   }
   allHashes := make([]int, n)
   for i := 0; i < n; i++ {
       allHashes[i] = canonizeForest(v[i])
   }
   sort.Ints(allHashes)
   g := make([][2]int, len(v[p]))
   copy(g, v[p])
   // try adding edge
   for i := 0; i < n-1; i++ {
       g = append(g, [2]int{i, n - 1})
       // multiset
       freq := make(map[int]int)
       for _, h := range allHashes {
           freq[h]++
       }
       ok := true
       for j := 0; j < n; j++ {
           var ng [][2]int
           for _, p := range g {
               if p[0] != j && p[1] != j {
                   ng = append(ng, p)
               }
           }
           h := canonizeForest(ng)
           if freq[h] > 0 {
               freq[h]--
           } else {
               ok = false
               break
           }
       }
       if ok {
           return g
       }
       g = g[:len(g)-1]
   }
   return nil
}

func main() {
   defer writer.Flush()
   var t int
   fmt.Fscan(reader, &t)
   for ; t > 0; t-- {
       ans := solve()
       if ans == nil || len(ans) == 0 {
           fmt.Fprintln(writer, "NO")
       } else {
           fmt.Fprintln(writer, "YES")
           for _, p := range ans {
               fmt.Fprintf(writer, "%d %d\n", p[0]+1, p[1]+1)
           }
       }
   }
}
