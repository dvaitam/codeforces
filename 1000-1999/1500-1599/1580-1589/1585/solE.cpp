#include "bits/stdc++.h"



using namespace std;



template<class T>

inline void setmax(T &a, T b) {

  if (a < b) {

    a = b;

  }

}



constexpr int maxn = 1e6 + 9;



int a[maxn], res[maxn];

int cnt[maxn], f[maxn];

int p[maxn], q[maxn];

vector<int> g[maxn];

vector<tuple<int, int, int>> queries[maxn];



inline void swap_pos(int a, int b) {

  swap(p[a], p[b]);

  swap(q[p[a]], q[p[b]]);

}



inline void insert(int x) {

  int pos = q[x];

  int newpos = ++f[++cnt[x]];

  swap_pos(pos, newpos);

}



inline void erase(int x) {

  int pos = q[x];

  int newpos = f[cnt[x]--]--;

  swap_pos(pos, newpos);

}



inline int query(int l, int k) {

  if (f[l] < k) {

    return -1;

  }

  return p[f[l] - k + 1];

}



void dfs(int u) {

  insert(a[u]);

  for (auto q : queries[u]) {

    int i, l, k;

    tie(i, l, k) = q;

    res[i] = query(l, k);

  }

  for (int &v : g[u]) {

    dfs(v);

  }

  erase(a[u]);

}



void solution() {

  int n, que;

  cin >> n >> que;

  for (int i = 1; i <= n; i++) {

    cin >> a[i];

    g[i].clear();

    queries[i].clear();

  }

  for (int v = 2; v <= n; v++) {

    int u;

    cin >> u;

    g[u].emplace_back(v);

  }

  iota(p, p + n + 1, 0);

  iota(q, q + n + 1, 0);

  for (int i = 0; i < que; i++) {

    int u, l, k;

    cin >> u >> l >> k;

    queries[u].emplace_back(i, l, k);

  }

  dfs(1);

  for (int i = 0; i < que; i++) {

    cout << res[i] << " ";

  }

  cout << "\n";

}



signed main() {

  ios::sync_with_stdio(0);

  cin.tie(0);

#ifdef DEBUG

  freopen("Input.txt", "r", stdin);

  auto start = chrono::high_resolution_clock::now();

#endif

  signed tc = 1;

  cin >> tc;



  for (int tt = 0; tt < tc; tt++) {

    // cout << "Case #" << tt << ": ";

    solution();

  }



#ifdef DEBUG

	auto end = chrono::high_resolution_clock::now();

	auto duration = chrono::duration_cast<chrono::milliseconds>(end - start).count();

	cout << endl << "_ Execution time: [ " << duration << "ms ] _" << endl;

#endif



  return 0;

}