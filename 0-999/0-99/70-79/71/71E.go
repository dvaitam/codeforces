package main

import (
   "fmt"
)

var elements = []string{
   "H", "He", "Li", "Be", "B", "C", "N", "O", "F", "Ne",
   "Na", "Mg", "Al", "Si", "P", "S", "Cl", "Ar", "K", "Ca",
   "Sc", "Ti", "V", "Cr", "Mn", "Fe", "Co", "Ni", "Cu", "Zn",
   "Ga", "Ge", "As", "Se", "Br", "Kr", "Rb", "Sr", "Y", "Zr",
   "Nb", "Mo", "Tc", "Ru", "Rh", "Pd", "Ag", "Cd", "In", "Sn",
   "Sb", "Te", "I", "Xe", "Cs", "Ba", "La", "Ce", "Pr", "Nd",
   "Pm", "Sm", "Eu", "Gd", "Tb", "Dy", "Ho", "Er", "Tm", "Yb",
   "Lu", "Hf", "Ta", "W", "Re", "Os", "Ir", "Pt", "Au", "Hg",
   "Tl", "Pb", "Bi", "Po", "At", "Rn", "Fr", "Ra", "Ac", "Th",
   "Pa", "U", "Np", "Pu", "Am", "Cm", "Bk", "Cf", "Es", "Fm",
}

func main() {
   var n, m int
   if _, err := fmt.Scan(&n, &m); err != nil {
       return
   }
   seq := make([]int, n)
   tar := make([]int, m)
   elementMap := make(map[string]int, len(elements))
   for i, s := range elements {
       elementMap[s] = i + 1
   }
   var s string
   for i := 0; i < n; i++ {
       fmt.Scan(&s)
       seq[i] = elementMap[s]
   }
   for i := 0; i < m; i++ {
       fmt.Scan(&s)
       tar[i] = elementMap[s]
   }
   lim := (1 << n) - 1
   type data struct{ id, val int }
   f := make([]data, lim+1)
   link := make([]int, lim+1)
   for i := range link {
       link[i] = -1
   }
   for mask := 0; mask <= lim; mask++ {
       cur := f[mask]
       for j := 0; j < n; j++ {
           bit := 1 << j
           if mask&bit != 0 {
               continue
           }
           p := seq[j]
           id0, val0 := cur.id, cur.val
           rem := tar[id0] - val0
           if rem < p {
               continue
           }
           var nd data
           if rem == p {
               nd = data{id: id0 + 1, val: 0}
           } else {
               nd = data{id: id0, val: val0 + p}
           }
           nxt := mask | bit
           prev := f[nxt]
           if prev.id < nd.id || (prev.id == nd.id && prev.val < nd.val) {
               f[nxt] = nd
               link[nxt] = mask
           }
       }
   }
   end := f[lim]
   if end.id == m && end.val == 0 {
       fmt.Println("YES")
       now := lim
       v := 0
       for now > 0 {
           prev := link[now]
           diff := now ^ prev
           for j := 0; j < n; j++ {
               if diff&(1<<j) != 0 {
                   if v > 0 {
                       fmt.Printf("+%s", elements[seq[j]-1])
                   } else {
                       fmt.Printf("%s", elements[seq[j]-1])
                   }
                   v++
               }
           }
           if f[prev].val == 0 {
               v = 0
               fmt.Printf("->%s\n", elements[tar[f[prev].id]-1])
           }
           now = prev
       }
   } else {
       fmt.Println("NO")
   }
}
