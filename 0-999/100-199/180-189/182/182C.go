package main

import (
   "bufio"
   "fmt"
   "math/rand"
   "os"
   "time"
)

// Treap node
type Treap struct {
   key   int64
   prio  int32
   count int
   size  int
   sum   int64
   left  *Treap
   right *Treap
}

// update recalculates size and sum
func (t *Treap) update() {
   t.size = t.count
   t.sum = t.key * int64(t.count)
   if t.left != nil {
       t.size += t.left.size
       t.sum += t.left.sum
   }
   if t.right != nil {
       t.size += t.right.size
       t.sum += t.right.sum
   }
}

// rotateRight performs right rotation
func rotateRight(t *Treap) *Treap {
   l := t.left
   t.left = l.right
   l.right = t
   t.update()
   l.update()
   return l
}

// rotateLeft performs left rotation
func rotateLeft(t *Treap) *Treap {
   r := t.right
   t.right = r.left
   r.left = t
   t.update()
   r.update()
   return r
}

// insert key into treap
func treapInsert(t *Treap, key int64) *Treap {
   if t == nil {
       return &Treap{key: key, prio: rand.Int31(), count: 1, size: 1, sum: key}
   }
   if key == t.key {
       t.count++
   } else if key < t.key {
       t.left = treapInsert(t.left, key)
       if t.left.prio < t.prio {
           t = rotateRight(t)
       }
   } else {
       t.right = treapInsert(t.right, key)
       if t.right.prio < t.prio {
           t = rotateLeft(t)
       }
   }
   t.update()
   return t
}

// delete one instance of key from treap
func treapDelete(t *Treap, key int64) *Treap {
   if t == nil {
       return nil
   }
   if key < t.key {
       t.left = treapDelete(t.left, key)
   } else if key > t.key {
       t.right = treapDelete(t.right, key)
   } else {
       if t.count > 1 {
           t.count--
       } else {
           // remove node
           if t.left == nil {
               return t.right
           } else if t.right == nil {
               return t.left
           } else {
               // merge by priority
               if t.left.prio < t.right.prio {
                   t = rotateRight(t)
                   t.right = treapDelete(t.right, key)
               } else {
                   t = rotateLeft(t)
                   t.left = treapDelete(t.left, key)
               }
           }
       }
   }
   if t != nil {
       t.update()
   }
   return t
}

// sumTopK returns sum of largest k keys in treap
func sumTopK(t *Treap, k int) int64 {
   if t == nil || k <= 0 {
       return 0
   }
   if k >= t.size {
       return t.sum
   }
   // consider right subtree first (larger keys)
   if t.right != nil {
       if t.right.size >= k {
           return sumTopK(t.right, k)
       }
       // include all right subtree
       res := t.right.sum
       rem := k - t.right.size
       if t.count >= rem {
           return res + t.key*int64(rem)
       }
       // include entire node
       res += t.key * int64(t.count)
       rem -= t.count
       // then go to left subtree
       return res + sumTopK(t.left, rem)
   }
   // no right subtree
   rem := k
   if t.count >= rem {
       return t.key * int64(rem)
   }
   // include node and then left
   res := t.key * int64(t.count)
   rem -= t.count
   return res + sumTopK(t.left, rem)
}

// TopK data structure
type TopK struct {
   root *Treap
}

// Insert a key
func (tk *TopK) Insert(key int64) {
   tk.root = treapInsert(tk.root, key)
}

// Delete a key
func (tk *TopK) Delete(key int64) {
   tk.root = treapDelete(tk.root, key)
}

// SumTopK returns sum of largest k keys
func (tk *TopK) SumTopK(k int) int64 {
   return sumTopK(tk.root, k)
}

func main() {
   rand.Seed(time.Now().UnixNano())
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, length, k int
   fmt.Fscan(reader, &n, &length)
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   fmt.Fscan(reader, &k)

   var S int64
   posDS := &TopK{}
   negDS := &TopK{}
   // initial window [0..length-1]
   for i := 0; i < length; i++ {
       v := a[i]
       S += v
       if v > 0 {
           posDS.Insert(v)
       } else if v < 0 {
           negDS.Insert(-v)
       }
   }
   ans := abs64(S)
   // slide windows
   for i := 1; i+length <= n; i++ {
       // remove a[i-1]
       old := a[i-1]
       S -= old
       if old > 0 {
           posDS.Delete(old)
       } else if old < 0 {
           negDS.Delete(-old)
       }
       // add a[i+length-1]
       v := a[i+length-1]
       S += v
       if v > 0 {
           posDS.Insert(v)
       } else if v < 0 {
           negDS.Insert(-v)
       }
       // compute best for this window
       gain := int64(2) * negDS.SumTopK(k)
       reduce := int64(2) * posDS.SumTopK(k)
       best1 := S + gain
       best2 := -S + reduce
       cur := best1
       if best2 > cur {
           cur = best2
       }
       if cur > ans {
           ans = cur
       }
   }
   fmt.Fprintln(writer, ans)
}

func abs64(x int64) int64 {
   if x < 0 {
       return -x
   }
   return x
}
