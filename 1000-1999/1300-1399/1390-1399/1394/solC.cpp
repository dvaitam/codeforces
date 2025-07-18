#include <bits/stdc++.h>

using namespace std;



#define FOR(x,n) for(int x=0;x<n;x++)

using ll = long long;

using ii = pair<int,int>;



int main() {

  ios_base::sync_with_stdio(false);

  cin.tie(NULL);

  // mt19937 rng(chrono::steady_clock::now().time_since_epoch().count());

  

  int n;

  cin >> n;

  

  int xmin = INT_MAX, ymin = INT_MAX;

  int xmax = -1, ymax = -1;

  int a = INT_MIN, b = INT_MIN;

  FOR(i,n) {

    string s;

    cin >> s;

    int ctb = 0, ctn = 0; 

    for(char c : s) {

      if(c == 'B') ctb++;

      else ctn++;

    }

    xmin = min(xmin, ctb); ymin = min(ymin, ctn); 

    xmax = max(xmax, ctb); ymax = max(ymax, ctn);

    a = max(a, -ctb + ctn);

    b = max(b, ctb - ctn);

  }

  

  int best = INT_MAX, ansb = 0, ansn = 0;

  auto check = [&](int t) {

    FOR(xpr,xmax+1) {

      if(abs(xpr - xmin) > t) continue;

      if(abs(xpr - xmax) > t) continue;

      int ylow = ymax - t;

      int yhi = ymin + t;

      ylow = max(ylow, xpr + a - t);

      yhi = min(yhi, t - b + xpr);

      if(ylow <= yhi && (xpr > 0 || yhi > 0)) {

        if(t < best) {

          best = t;

          ansb = xpr;

          ansn = yhi;

        }

        return true;

      }

    }

    return false;

  };

  

  int lo = 0, hi = 1e6; 

  while(lo <= hi) {

    int mid = (lo + hi)/2;

    if(check(mid)) hi = mid - 1;

    else lo = mid + 1;

  }

  

  cout << best << '\n';

  cout << string(ansb, 'B') + string(ansn, 'N') << '\n';

}