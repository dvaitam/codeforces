package main

import (
   "bufio"
   "fmt"
   "math/rand"
   "os"
   "strconv"
   "time"
)

// Treap implementation for ordered set of (key, idx) pairs
type node struct {
   key, idx, priority int
   left, right        *node
}

type Treap struct {
   root *node
}

// cmp returns true if (k1,i1) < (k2,i2)
func cmp(k1, i1, k2, i2 int) bool {
   if k1 != k2 {
       return k1 < k2
   }
   return i1 < i2
}

// split splits t into (< (key,idx)) and (>= (key,idx))
func split(t *node, key, idx int) (l, r *node) {
   if t == nil {
       return nil, nil
   }
   if cmp(t.key, t.idx, key, idx) {
       ll, rr := split(t.right, key, idx)
       t.right = ll
       return t, rr
   }
   ll, rr := split(t.left, key, idx)
   t.left = rr
   return ll, t
}

// merge merges two treaps a and b
func merge(a, b *node) *node {
   if a == nil {
       return b
   }
   if b == nil {
       return a
   }
   if a.priority < b.priority {
       a.right = merge(a.right, b)
       return a
   }
   b.left = merge(a, b.left)
   return b
}

// Insert adds a (key, idx) pair
func (tr *Treap) Insert(key, idx int) {
   n := &node{key: key, idx: idx, priority: rand.Int()}
   l, r := split(tr.root, key, idx)
   tr.root = merge(merge(l, n), r)
}

// delete removes the node with (key, idx)
func (tr *Treap) delete(t *node, key, idx int) *node {
   if t == nil {
       return nil
   }
   if t.key == key && t.idx == idx {
       return merge(t.left, t.right)
   }
   if cmp(t.key, t.idx, key, idx) {
       t.right = tr.delete(t.right, key, idx)
   } else {
       t.left = tr.delete(t.left, key, idx)
   }
   return t
}

// Delete removes (key, idx) from treap
func (tr *Treap) Delete(key, idx int) {
   tr.root = tr.delete(tr.root, key, idx)
}

// Clear empties the treap
func (tr *Treap) Clear() {
   tr.root = nil
}

// Empty checks if treap is empty
func (tr *Treap) Empty() bool {
   return tr.root == nil
}

// Min returns the minimum (key, idx)
func (tr *Treap) Min() (key, idx int, ok bool) {
   t := tr.root
   if t == nil {
       return 0, 0, false
   }
   for t.left != nil {
       t = t.left
   }
   return t.key, t.idx, true
}

// Max returns the maximum (key, idx)
func (tr *Treap) Max() (key, idx int, ok bool) {
   t := tr.root
   if t == nil {
       return 0, 0, false
   }
   for t.right != nil {
       t = t.right
   }
   return t.key, t.idx, true
}

// RemoveMinWhile removes all nodes with key < minKey
func (tr *Treap) RemoveMinWhile(minKey int) {
   for {
       k, idx, ok := tr.Min()
       if !ok || k >= minKey {
           break
       }
       tr.Delete(k, idx)
   }
}

// RemoveMaxWhile removes all nodes with key > maxKey
func (tr *Treap) RemoveMaxWhile(maxKey int) {
   for {
       k, idx, ok := tr.Max()
       if !ok || k <= maxKey {
           break
       }
       tr.Delete(k, idx)
   }
}

var (
   reader *bufio.Reader
   writer *bufio.Writer
)

// readInt reads next integer from stdin
func readInt() (int, error) {
   var x int
   var c byte
   var err error
   // skip non-digits
   for {
       c, err = reader.ReadByte()
       if err != nil {
           return 0, err
       }
       if c >= '0' && c <= '9' {
           break
       }
   }
   x = int(c - '0')
   for {
       c, err = reader.ReadByte()
       if err != nil || c < '0' || c > '9' {
           break
       }
       x = x*10 + int(c-'0')
   }
   return x, nil
}

func main() {
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   rand.Seed(time.Now().UnixNano())

   n, _ := readInt()
   m, _ := readInt()
   _ = m // m is unused
   k := make([]int, n+1)
   xl := make([]int, n+1)
   xr := make([]int, n+1)
   yl := make([]int, n+1)
   yr := make([]int, n+1)
   for i := 1; i <= n; i++ {
       k[i], _ = readInt()
       xl[i], _ = readInt()
       xr[i], _ = readInt()
       yl[i], _ = readInt()
       yr[i], _ = readInt()
   }
   f0 := make([]int, n+1)
   f1 := make([]int, n+1)
   w := make([]int, n+1)
   A := &Treap{}
   B := &Treap{}
   A.Insert(0, 0)
   B.Insert(0, 0)
   for i := 1; i <= n; i++ {
       if k[i] < xl[i] || k[i] > xr[i] {
           A.Clear()
       }
       if k[i] < yl[i] || k[i] > yr[i] {
           B.Clear()
       }
       A.RemoveMinWhile(yl[i])
       A.RemoveMaxWhile(yr[i])
       B.RemoveMinWhile(xl[i])
       B.RemoveMaxWhile(xr[i])
       if A.Empty() && B.Empty() {
           fmt.Fprintln(writer, "No")
           return
       }
       if _, idx, ok := A.Min(); ok {
           f0[i] = idx
       } else {
           f0[i] = -1
       }
       if _, idx, ok := B.Min(); ok {
           f1[i] = idx
       } else {
           f1[i] = -1
       }
       if f0[i] != -1 {
           B.Insert(k[i], i)
       }
       if f1[i] != -1 {
           A.Insert(k[i], i)
       }
   }
   o := 0
   if f0[n] == -1 {
       o = 1
   }
   x := 0
   if o == 0 {
       x = f0[n]
   } else {
       x = f1[n]
   }
   cur := o
   for i := n; i >= 1; i-- {
       w[i] = cur
       if i == x {
           cur ^= 1
           if cur == 0 {
               x = f0[x]
           } else {
               x = f1[x]
           }
       }
   }
   fmt.Fprintln(writer, "Yes")
   for i := 1; i <= n; i++ {
       if i > 1 {
           writer.WriteByte(' ')
       }
       writer.WriteString(strconv.Itoa(w[i]))
   }
   writer.WriteByte('\n')
}
