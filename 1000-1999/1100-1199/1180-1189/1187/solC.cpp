//...author:MAUNIE...//
#include<bits/stdc++.h>

using namespace std;

int main() {
  ios::sync_with_stdio(false);
  cin.tie(0);
  int n, m;
  cin >> n >> m;
  vector<int> a(n);
  vector<bool> s(n);
  vector<int> t(m), l(m), r(m);
  for (int i = 0; i < n; i++) {
      a[i] = 0;
  }
  for (int i = 0; i < m; i++) {
      cin >> t[i] >> l[i] >> r[i];
      l[i]--;
      r[i]--;
      if (t[i]) {
          for (int j = l[i] + 1; j <= r[i]; j++) {
              s[j] = true;
          }
      }
  }
  for (int i = 0; i < n; i++) {
      if (s[i]) {
          a[i] = a[i - 1];
      }
      else {
          a[i] = a[i - 1] - 1;
      }
  }
  int k = 1e9;
  for (int i = 0; i < n; i++) {
      k = min(k, a[i]);
  }
  if (k <= 0) {
      for (int i = 0; i < n; i++) {
          a[i] += abs(k);
          a[i]++;
      }
  }
  for (int i = 0; i < m; i++) {
      bool sorted = true;
      for (int j = l[i] + 1; j <= r[i]; j++) {
          if (a[j] < a[j - 1]) {
              sorted = false;
          }
      }
      if (sorted != t[i]) {
          cout << "NO";
          return 0;
      }
  }
  cout << "YES" << '\n';
  for (int i = 0; i < n; i++) {
      cout << a[i] << " ";
  }
  return 0;
}