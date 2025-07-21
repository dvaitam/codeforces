package main

import (
   "bufio"
   "fmt"
   "os"
)

type Node struct {
   next [26]int32
   link int32
}

// Fenwick tree for range add, point query
type BIT struct {
   n    int
   tree []int32
}

func NewBIT(n int) *BIT {
   return &BIT{n: n, tree: make([]int32, n+2)}
}

// add v at position i
func (b *BIT) add(i int, v int32) {
   for x := i + 1; x <= b.n+1; x += x & -x {
       b.tree[x] += v
   }
}

// sum of [0..i]
func (b *BIT) sum(i int) int32 {
   var s int32
   for x := i + 1; x > 0; x -= x & -x {
       s += b.tree[x]
   }
   return s
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   fmt.Fscan(reader, &n, &k)
   patterns := make([]string, k+1)
   for i := 1; i <= k; i++ {
       fmt.Fscan(reader, &patterns[i])
   }
   // build trie
   nodes := make([]Node, 1)
   patternEnd := make([]int32, k+1)
   for i := 1; i <= k; i++ {
       s := patterns[i]
       v := int32(0)
       for _, ch := range s {
           c := ch - 'a'
           if nodes[v].next[c] == 0 {
               nodes = append(nodes, Node{})
               nodes[v].next[c] = int32(len(nodes) - 1)
           }
           v = nodes[v].next[c]
       }
       patternEnd[i] = v
   }
   // build AC automaton
   q := make([]int32, 0, len(nodes))
   // init root links
   for c := 0; c < 26; c++ {
       if nodes[0].next[c] != 0 {
           child := nodes[0].next[c]
           nodes[child].link = 0
           q = append(q, child)
       }
   }
   // fill missing from root
   for c := 0; c < 26; c++ {
       if nodes[0].next[c] == 0 {
           nodes[0].next[c] = 0
       }
   }
   // BFS
   for i := 0; i < len(q); i++ {
       v := q[i]
       for c := 0; c < 26; c++ {
           u := nodes[v].next[c]
           if u != 0 {
               nodes[u].link = nodes[nodes[v].link].next[c]
               q = append(q, u)
           } else {
               nodes[v].next[c] = nodes[nodes[v].link].next[c]
           }
       }
   }
   // build fail tree
   N := len(nodes)
   children := make([][]int32, N)
   for v := 1; v < N; v++ {
       p := nodes[v].link
       children[p] = append(children[p], int32(v))
   }
   // Euler tour
   tin := make([]int, N)
   tout := make([]int, N)
   time := 0
   // iterative DFS
   stack := make([]int32, 0, N)
   idx := make([]int, N)
   stack = append(stack, 0)
   for len(stack) > 0 {
       v := stack[len(stack)-1]
       if idx[v] == 0 {
           tin[v] = time
           time++
       }
       if idx[v] < len(children[v]) {
           u := children[v][idx[v]]
           idx[v]++
           stack = append(stack, u)
       } else {
           tout[v] = time - 1
           stack = stack[:len(stack)-1]
       }
   }
   // BIT
   bit := NewBIT(time + 2)
   active := make([]bool, k+1)
   for i := 1; i <= k; i++ {
       active[i] = true
       v := int(patternEnd[i])
       l := tin[v]
       r := tout[v]
       bit.add(l, 1)
       bit.add(r+1, -1)
   }
   // process queries
   for i := 0; i < n; i++ {
       var line string
       fmt.Fscan(reader, &line)
       if line[0] == '+' {
           var id int
           fmt.Sscanf(line[1:], "%d", &id)
           if !active[id] {
               active[id] = true
               v := int(patternEnd[id])
               l, r := tin[v], tout[v]
               bit.add(l, 1)
               bit.add(r+1, -1)
           }
       } else if line[0] == '-' {
           var id int
           fmt.Sscanf(line[1:], "%d", &id)
           if active[id] {
               active[id] = false
               v := int(patternEnd[id])
               l, r := tin[v], tout[v]
               bit.add(l, -1)
               bit.add(r+1, 1)
           }
       } else if line[0] == '?' {
           text := line[1:]
           v := int32(0)
           var ans int64
           for _, ch := range text {
               v = nodes[v].next[ch-'a']
               ans += int64(bit.sum(tin[int(v)]))
           }
           fmt.Fprintln(writer, ans)
       }
   }
