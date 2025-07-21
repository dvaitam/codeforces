package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const MOD = 1000000007

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   n := len(s)
   // build match
   match := make([]int, n)
   posStack := make([]int, 0, n)
   for i := 0; i < n; i++ {
       if s[i] == '(' {
           posStack = append(posStack, i)
       } else {
           l := posStack[len(posStack)-1]
           posStack = posStack[:len(posStack)-1]
           match[i] = l
           match[l] = i
       }
   }
   // build tree of nodes using child stacks
   type node struct { l, r int; children []int }
   var nodes []node
   // childStack holds slices of child indices for each open bracket
   childStack := make([][]int, 0, n)
   var roots []int
   for i := 0; i < n; i++ {
       if s[i] == '(' {
           // start new children list for this '(' level
           childStack = append(childStack, nil)
       } else {
           l := match[i]
           // pop children of this node
           last := len(childStack) - 1
           kids := childStack[last]
           childStack = childStack[:last]
           // create node
           idx := len(nodes)
           nodes = append(nodes, node{l: l, r: i, children: kids})
           // append to parent children list or roots
           if len(childStack) > 0 {
               childStack[len(childStack)-1] = append(childStack[len(childStack)-1], idx)
           } else {
               roots = append(roots, idx)
           }
       }
   }
   // sort children by l
   for i := range nodes {
       if len(nodes[i].children) > 1 {
           // simple insertion sort small
           ch := nodes[i].children
           for a := 1; a < len(ch); a++ {
               v := ch[a]
               j := a
               for j > 0 && nodes[ch[j-1]].l > nodes[v].l {
                   ch[j] = ch[j-1]
                   j--
               }
               ch[j] = v
           }
           nodes[i].children = ch
       }
   }
   // sort root nodes by l
   if len(roots) > 1 {
       sort.Slice(roots, func(i, j int) bool {
           return nodes[roots[i]].l < nodes[roots[j]].l
       })
   }
   // DP f per node
   f := make([][3][3]int, len(nodes))
   var dfs func(int)
   dfs = func(u int) {
       // compute for children first
       for _, v := range nodes[u].children {
           dfs(v)
       }
       // children sequence DP
       // children_dp[a][c] = ways for nested children with bracket.l color a and last child.r color c
       var children_dp [3][3]int
       for a := 0; a < 3; a++ {
           // cur[x] = ways for prev child's r color = x
           cur := [3]int{0, 0, 0}
           cur[a] = 1
           for _, v := range nodes[u].children {
               next := [3]int{0, 0, 0}
               // combine cur with f[v]
               for x := 0; x < 3; x++ {
                   if cur[x] == 0 {
                       continue
                   }
                   for cl := 0; cl < 3; cl++ {
                       // adjacency between prev r color x and child.l color cl
                       if x > 0 && cl > 0 && x == cl {
                           continue
                       }
                       for cr := 0; cr < 3; cr++ {
                           ways := f[v][cl][cr]
                           if ways == 0 {
                               continue
                           }
                           next[cr] = (next[cr] + cur[x]*ways) % MOD
                       }
                   }
               }
               cur = next
           }
           for c := 0; c < 3; c++ {
               children_dp[a][c] = cur[c]
           }
       }
       // compute f[u][a][b]
       for a := 0; a < 3; a++ {
           for b := 0; b < 3; b++ {
               // exactly one of a,b colored
               if (a > 0) == (b > 0) {
                   f[u][a][b] = 0
                   continue
               }
               total := 0
               for c := 0; c < 3; c++ {
                   ways := children_dp[a][c]
                   if ways == 0 {
                       continue
                   }
                   // adjacency between last child.r (c) and bracket.r color b
                   if c > 0 && b > 0 && c == b {
                       continue
                   }
                   total = (total + ways) % MOD
               }
               f[u][a][b] = total
           }
       }
   }
   for _, u := range roots {
       dfs(u)
   }
   // combine roots as siblings
   cur := [3]int{0, 0, 0}
   cur[0] = 1
   for _, u := range roots {
       next := [3]int{0, 0, 0}
       for x := 0; x < 3; x++ {
           if cur[x] == 0 {
               continue
           }
           for cl := 0; cl < 3; cl++ {
               // adjacency between prev r color x and sibling.l color cl
               if x > 0 && cl > 0 && x == cl {
                   continue
               }
               for cr := 0; cr < 3; cr++ {
                   ways := f[u][cl][cr]
                   if ways == 0 {
                       continue
                   }
                   next[cr] = (next[cr] + cur[x]*ways) % MOD
               }
           }
       }
       cur = next
   }
   // answer is sum of cur[*]
   ans := 0
   for c := 0; c < 3; c++ {
       ans = (ans + cur[c]) % MOD
   }
   fmt.Println(ans)
}
