package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Fast input reader
type FastReader struct {
   r *bufio.Reader
}

func NewReader() *FastReader {
   return &FastReader{r: bufio.NewReader(os.Stdin)}
}

func (fr *FastReader) ReadInt() (int, error) {
   var x int
   var c byte
   var err error
   // skip non-digit
   for {
       c, err = fr.r.ReadByte()
       if err != nil {
           return 0, err
       }
       if c == '-' || (c >= '0' && c <= '9') {
           break
       }
   }
   neg := false
   if c == '-' {
       neg = true
       c, _ = fr.r.ReadByte()
   }
   for ; err == nil && c >= '0' && c <= '9'; c, err = fr.r.ReadByte() {
       x = x*10 + int(c-'0')
   }
   if neg {
       x = -x
   }
   return x, nil
}

func (fr *FastReader) ReadToken() (byte, error) {
   var c byte
   var err error
   for {
       c, err = fr.r.ReadByte()
       if err != nil {
           return 0, err
       }
       if c >= 'a' && c <= 'z' {
           return c, nil
       }
   }
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   fr := NewReader()
   n, _ := fr.ReadInt()
   parent := make([]int, n+1)
   parent[0] = -1
   cArr := make([]byte, n+1)
   for i := 1; i <= n; i++ {
       p, _ := fr.ReadInt()
       c, _ := fr.ReadToken()
       parent[i] = p
       cArr[i] = c
   }
   k, _ := fr.ReadInt()
   a := make([]int, k)
   mark := make([]bool, n+1)
   for i := 0; i < k; i++ {
       ai, _ := fr.ReadInt()
       a[i] = ai
       mark[ai] = true
   }
   // build adjacency flat arrays
   deg := make([]int, n+1)
   for i := 1; i <= n; i++ {
       deg[parent[i]]++
   }
   pos := make([]int, n+1)
   for u := 1; u <= n; u++ {
       pos[u] = pos[u-1] + deg[u-1]
   }
   childTo := make([]int, n)
   childC := make([]byte, n)
   cur := make([]int, n+1)
   copy(cur, pos)
   for i := 1; i <= n; i++ {
       p := parent[i]
       idx := cur[p]
       childTo[idx] = i
       childC[idx] = cArr[i]
       cur[p]++
   }
   // sort children by char
   for u := 0; u <= n; u++ {
       st := pos[u]
       cnt := deg[u]
       if cnt > 1 {
           sort.Slice(childTo[st:st+cnt], func(i, j int) bool {
               return childC[st+i] < childC[st+j]
           })
       }
   }
   // first DFS: compute lex order and intervals
   const INF = int(1e9)
   l := make([]int, n+1)
   r := make([]int, n+1)
   for i := range l {
       l[i] = INF
   }
   order := make([]int, n+1)
   cntOrder := 0
   type Frame struct{ u, idx int }
   stack := make([]Frame, 0, n+1)
   stack = append(stack, Frame{u: 0, idx: 0})
   for len(stack) > 0 {
       fr := &stack[len(stack)-1]
       u := fr.u
       if fr.idx == 0 {
           if mark[u] {
               cntOrder++
               order[u] = cntOrder
               l[u], r[u] = cntOrder, cntOrder
           }
       }
       if fr.idx < deg[u] {
           v := childTo[pos[u]+fr.idx]
           fr.idx++
           stack = append(stack, Frame{u: v, idx: 0})
       } else {
           // post
           if u != 0 {
               p := parent[u]
               if l[u] <= r[u] {
                   if l[p] > l[u] {
                       l[p] = l[u]
                   }
                   if r[p] < r[u] {
                       r[p] = r[u]
                   }
               }
           }
           stack = stack[:len(stack)-1]
       }
   }
   // dp with second DFS
   dp := make([]int, n+1)
   type Frame2 struct{ u, idx int; act bool }
   stack2 := make([]Frame2, 0, n+1)
   valStack := make([]int, 0, n+1)
   minStack := make([]int, 0, n+1)
   // init root
   dp[0] = 0
   stack2 = append(stack2, Frame2{u: 0, idx: 0, act: true})
   // push root val
   rootVal := dp[0] - l[0] + 1
   valStack = append(valStack, rootVal)
   minStack = append(minStack, rootVal)
   for len(stack2) > 0 {
       fr2 := &stack2[len(stack2)-1]
       u := fr2.u
       if fr2.idx == 0 && u != 0 {
           p := parent[u]
           dp[u] = dp[p] + 1
           if mark[u] {
               // autocomplete
               best := minStack[len(minStack)-1]
               cand := order[u] + best
               if cand < dp[u] {
                   dp[u] = cand
               }
           }
           if l[u] <= r[u] {
               val := dp[u] - l[u] + 1
               stack2[len(stack2)-1].act = true
               valStack = append(valStack, val)
               // update minStack
               m := val
               if last := minStack[len(minStack)-1]; last < m {
                   m = last
               }
               minStack = append(minStack, m)
           }
       }
       if fr2.idx < deg[u] {
           v := childTo[pos[u]+fr2.idx]
           stack2[len(stack2)-1].idx++
           stack2 = append(stack2, Frame2{u: v, idx: 0, act: false})
           continue
       }
       // exit
       if fr2.act {
           valStack = valStack[:len(valStack)-1]
           minStack = minStack[:len(minStack)-1]
       }
       stack2 = stack2[:len(stack2)-1]
   }
   // output
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   for i, ai := range a {
       if i > 0 {
           w.WriteByte(' ')
       }
       fmt.Fprintf(w, "%d", dp[ai])
   }
}
