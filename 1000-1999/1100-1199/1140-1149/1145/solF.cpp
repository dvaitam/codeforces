#include <bits/stdc++.h>
using namespace std;

typedef long long  ll;

const int maxn = 2e5 + 5;


int main() {
  ios::sync_with_stdio(0), cin.tie(0), cout.tie(0);

  string str = "AEFHIKLMNTVWXYZ";
  string s; cin >> s;
  if (str.find(s[0]) == -1) {
    for (auto& c : s) {
      if (str.find(c) != -1) {
        cout << "NO\n";
        return 0;
      }
    }
  } else {
    for (auto& c : s) {
      if (str.find(c) == -1) {
        cout << "NO\n";
        return 0;
      }
    }
  }
  cout << "YES\n";
}