#include <bits/stdc++.h>
using namespace std;

int solve() {
  int x;
  cin >> x;

  int n = 1 << x;
  vector g(n, vector(0, pair<int, int>{}));

  for (int i = 1, x, y; i < n; ++i) {
    cin >> x >> y;
    x -= 1, y -= 1;
    g[x].emplace_back(y, i - 1);
    g[y].emplace_back(x, i - 1);
  }

  vector<int> p1(n), p2(n - 1);
  p1.front() = n;

  function<void(int, int, bool)> dfs = [&](int u, int fa, bool left) {
    for (auto [v, i] : g[u]) {
      if (v != fa) {
        // x ^ n <-> x
        p1[v] = v ^ (left ? n : 0);
        p2[i] = v ^ (left ? 0 : n);
        dfs(v, u, !left);
      }
    }
  };

  dfs(0, -1, false);

  cout << "1\n";
  ranges::copy(p1, ostream_iterator<int>(cout, " ")), cout << "\n";
  ranges::copy(p2, ostream_iterator<int>(cout, " ")), cout << "\n";

  return {};
}

int main() {
  cin.tie(0), ios::sync_with_stdio(0);

  int tests;
  cin >> tests;
  while (tests--) {
    solve();
  }

  return 0 ^ 0;
}