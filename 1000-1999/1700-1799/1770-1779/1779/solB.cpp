/**

 *    Author:   Shihab

 *    Created:  2023-02-02 20:50:16

 *

 *    Problem:  https://codeforces.com/contest/1779/problem/B

 **/

// GCC Optimization Pragmas



/*#pragma GCC optimize("O3,unroll-loops")

#pragma GCC target("avx2,bmi,bmi2,lzcnt,popcnt")*/



#include <bits/stdc++.h>



using namespace std;



int main() {

  ios::sync_with_stdio(0);

  cin.tie(0);



  int tt;

  cin >> tt;

  while (tt--) {

    int n;

    cin >> n;

    if (n & 1) {

      if (n != 3) {

        int p = (n / 2) - 1;

        int q = n / 2;

        cout << "YES\n";

        for (int i = 0; i < n; i++) {

          cout << (i & 1 ? -q : p) << " \n"[i == n - 1];

        }

      } else

        cout << "NO\n";

    } else {

      cout << "YES\n";

      for (int i = 0; i < n; i++) {

        cout << (i & 1 ? -1 : 1) << " \n"[i == n - 1];

      }

    }

  }

  return 0;

}



// https://www.youtube.com/watch?v=2L326KCBqug