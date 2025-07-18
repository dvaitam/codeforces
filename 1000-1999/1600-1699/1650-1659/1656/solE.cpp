#include <bits/stdc++.h>

using namespace std;



int solve() {

  int n;

  cin >> n;

  vector G(n, vector<int>{});

  vector<int> deg(n);

  vector<bool> color(n);



  for (int i = 1, u, v; i < n; ++i) {

    cin >> u >> v;

    deg[--u] += 1, deg[--v] += 1;

    G[u].push_back(v);

    G[v].push_back(u);

  }



  function<void(int, int)> dfs = [&](int u, int fa) {

    for (int v : G[u]) {

      if (v != fa) {

        color[v] = color[u] ^ 1;

        dfs(v, u);

      }

    }

  };



  dfs(0, -1);



  for (int i = 0; i < n; ++i) {

    cout << (color[i] ? 1 : -1) * deg[i] << ' ';

  }



  cout << "\n";

  

  return {};

}



int main() {

  cin.tie(0), ios::sync_with_stdio(0);



  int tests;

  cin >> tests;

  while (tests --) {

    solve();

  }



  return 0 ^ 0;

}