#include <bits/stdc++.h>

using namespace std;

#ifdef LOCAL
#include "algo/debug.h"
#else
#define debug(...) 42
#endif

int main() {
  ios::sync_with_stdio(false);
  cin.tie(0);
  int tt;
  cin >> tt;
  while (tt--) {
    int n, m;
    cin >> n >> m;
    vector<vector<int>> g(n);
    for (int i = 0; i < m; i++) {
      int x, y;
      cin >> x >> y;
      --x; --y;
      g[y].push_back(x);
    }
    vector<int> d(n, -1);
    vector<int> que(1, 0);
    d[0] = 1;
    for (int b = 0; b < (int) que.size(); b++) {
      for (int u : g[que[b]]) {
        if (d[u] == -1) {
          que.push_back(u);
          d[u] = d[que[b]] + 1;
        }
      }
    }
    if (*min_element(d.begin(), d.end()) == -1) {
      cout << "INFINITE" << '\n';
      continue;
    }
    cout << "FINITE" << '\n';
    vector<vector<int>> at(n + 1);
    for (int i = 0; i < n; i++) {
      at[d[i]].push_back(i);
    }
    vector<int> seq;
    for (int from = 1; from <= n; from++) {
      for (int val = n; val >= from; val--) {
        for (int x : at[val]) {
          seq.push_back(x);
        }
      }
    }
    cout << seq.size() << '\n';
    for (int i = 0; i < (int) seq.size(); i++) {
      cout << seq[i] + 1 << " \n"[i == (int) seq.size() - 1];
    }
  }
  return 0;
}