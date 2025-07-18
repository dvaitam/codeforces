#include <bits/stdc++.h>

using namespace std;

int main() {
  int t;
  cin >> t;
  while(t--) {
  int x;
  cin >> x;
  int ans = 1;
  for(int i = 0; i < __builtin_popcount(x); i++) {
    ans *= 2;
  }
  cout << ans << "\n";
  }
  return 0;
}