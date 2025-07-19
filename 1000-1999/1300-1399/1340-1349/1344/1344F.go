package main

import (
   "bufio"
   "fmt"
   "math/big"
   "os"
)

const NN = 2001

// basis represents a linear basis over GF(2) for bitsets of size NN
type basis struct {
   f [NN]*big.Int
}

// insert tries to insert vector x into the basis. Returns false if inconsistent.
func (b *basis) insert(x *big.Int) bool {
   // make a copy of x to work on
   v := new(big.Int).Set(x)
   for i := NN - 1; i >= 1; i-- {
       if v.Bit(i) == 1 {
           if b.f[i] != nil && b.f[i].Bit(i) == 1 {
               v.Xor(v, b.f[i])
           } else {
               // set basis vector
               b.f[i] = v
               return true
           }
       }
   }
   // check constant term
   if v.Bit(0) == 1 {
       return false
   }
   return true
}

// getans returns one solution vector satisfying the basis
func (b *basis) getans() *big.Int {
   ans := new(big.Int)
   for i := 1; i < NN; i++ {
       if b.f[i] != nil && b.f[i].Bit(i) == 1 {
           tmp := b.f[i].Bit(0)
           for j := 1; j < i; j++ {
               if ans.Bit(j) == 1 && b.f[i].Bit(j) == 1 {
                   tmp ^= 1
               }
           }
           if tmp == 1 {
               ans.SetBit(ans, i, 1)
           }
       }
   }
   return ans
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   // a[i]: 2x2 matrix
   type mat2 [2][2]bool
   a := make([]mat2, n+1)
   for i := 1; i <= n; i++ {
       a[i][0][0] = true
       a[i][1][1] = true
   }
   b := &basis{}
   for k := 0; k < m; k++ {
       var op string
       var ln int
       fmt.Fscan(reader, &op, &ln)
       switch op {
       case "mix":
           A := new(big.Int)
           B := new(big.Int)
           for j := 0; j < ln; j++ {
               var x int
               fmt.Fscan(reader, &x)
               if a[x][0][0] {
                   A.SetBit(A, x*2-1, 1)
               }
               if a[x][0][1] {
                   A.SetBit(A, x*2, 1)
               }
               if a[x][1][0] {
                   B.SetBit(B, x*2-1, 1)
               }
               if a[x][1][1] {
                   B.SetBit(B, x*2, 1)
               }
           }
           var c byte
           fmt.Fscan(reader, &c)
           // constant term at bit 0
           if c == 'Y' || c == 'B' {
               A.SetBit(A, 0, 1)
           }
           if c == 'R' || c == 'B' {
               B.SetBit(B, 0, 1)
           }
           if !b.insert(A) || !b.insert(B) {
               fmt.Fprintln(writer, "NO")
               return
           }
       case "RY":
           for j := 0; j < ln; j++ {
               var x int
               fmt.Fscan(reader, &x)
               // swap rows
               a[x][0], a[x][1] = a[x][1], a[x][0]
           }
       case "RB":
           for j := 0; j < ln; j++ {
               var x int
               fmt.Fscan(reader, &x)
               a[x][0][0] = a[x][0][0] != a[x][1][0]
               a[x][0][1] = a[x][0][1] != a[x][1][1]
           }
       case "YB":
           for j := 0; j < ln; j++ {
               var x int
               fmt.Fscan(reader, &x)
               a[x][1][0] = a[x][1][0] != a[x][0][0]
               a[x][1][1] = a[x][1][1] != a[x][0][1]
           }
       }
   }
   fmt.Fprintln(writer, "YES")
   ans := b.getans()
   // output colors
   for i := 1; i <= n; i++ {
       bit1 := ans.Bit(i*2 - 1)
       bit2 := ans.Bit(i * 2)
       idx := bit1<<1 | bit2
       // mapping: 0->'.',1->'R',2->'Y',3->'B'
       var ch byte
       switch idx {
       case 0:
           ch = '.'
       case 1:
           ch = 'R'
       case 2:
           ch = 'Y'
       case 3:
           ch = 'B'
       }
       writer.WriteByte(ch)
   }
   writer.WriteByte('\n')
}
