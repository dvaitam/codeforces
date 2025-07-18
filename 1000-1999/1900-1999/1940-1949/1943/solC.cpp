#include <iostream>
#include <vector>

using namespace std;

const int MAX_N = 2e3 + 5;

vector<int> adj [MAX_N];
int par [MAX_N];
int lvl [MAX_N];

void dfs (int u, int p, int &out) {
  lvl[u] = lvl[p] + 1;
  par[u] = p;

  if (lvl[u] > lvl[out])
    out = u;
  
  for (int nxt : adj[u]) {
    if (nxt == p)
      continue;

    dfs(nxt, u, out);
  }
}

void solve () {
  int n;
  cin >> n;

  for (int i = 1; i <= n; i++)
    adj[i].clear();

  for (int i = 0; i < n - 1; i++) {
    int u, v;
    cin >> u >> v;

    adj[u].push_back(v);
    adj[v].push_back(u);
  }

  int root = 1;
  dfs(1, 0, root);

  int farth = root;
  dfs(root, 0, farth);

  vector<int> diam;
  for (int u = farth; u != root; u = par[u])
    diam.push_back(u);

  diam.push_back(root);

  int m = diam.size();
  vector<pair<int, int>> ops;
  if (m % 2 == 1) {
    int c = diam[m / 2];
    for (int i = 0; i <= m / 2; i++)
      ops.emplace_back(c, i);
  } else {
    int c1 = diam[m / 2];
    int c2 = diam[m / 2 - 1];

    for (int i = m / 2 - 1; i >= 0; i -= 2) {
      ops.emplace_back(c1, i);
      ops.emplace_back(c2, i);
    }
  }

  cout << (int) ops.size() << '\n';
  for (auto pr : ops) {
    cout << pr.first << " " << pr.second << '\n';
  }
}

int main () {
  ios::sync_with_stdio(false);
  cin.tie(0);

  int testc;
  cin >> testc;

  for (int i = 0; i < testc; i++) {
    solve();
  }
}