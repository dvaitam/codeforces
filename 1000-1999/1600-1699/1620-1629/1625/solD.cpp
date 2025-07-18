/**

 *  author: ivanzuki   

 *  created: Fri Jun 17 2022

**/

#include <bits/stdc++.h>



using namespace std;



const int N = 30 * 300010;



int tne[N][2];

pair<int, int> tdp[N];

int tsz = 1;



void Add(int x, int y, int i) {

  int v = 0;

  for (int bit = 30; bit >= 0; bit--) {

    while (tsz == N) {

    }

    int z = (x >> bit) & 1;

    if (tne[v][z] == -1) {

      tne[v][z] = tsz++;

    }

    v = tne[v][z];

    while (v == -1) {}

    tdp[v] = max(tdp[v], make_pair(y, i));

  }

}



pair<int, int> Get(int x, int k) {

  pair<int, int> ret(0, 0);

  int v = 0;

  for (int bit = 30; bit >= 0; bit--) {

    while (tsz == N) {}

    int zk = (k >> bit) & 1;

    int zx = (x >> bit) & 1;

    if (zk == 0) {

      if (tne[v][zx ^ 1] != -1) {

        ret = max(ret, tdp[tne[v][zx ^ 1]]);

      }

      v = tne[v][zx];

    } else {

      v = tne[v][zx ^ 1];

    }

    if (v == -1) {

      break;

    }

  }

  if (v != -1) {

    ret = max(ret, tdp[v]);

  }

  return ret;

}



int main() {

  ios_base::sync_with_stdio(false);

  cin.tie(0);

  int n, k;

  cin >> n >> k;

  memset(tne, -1, sizeof(tne));

  vector<pair<int, int>> a(n);

  for (int i = 0; i < n; i++) {

    cin >> a[i].first;

    a[i].second = i;

  }

  if (k == 0) {

    cout << n << '\n';

    for (int i = 0; i < n; i++) {

      cout << i + 1 << " \n"[i == n - 1];

    }

    return 0;

  }

  sort(a.begin(), a.end());

  pair<int, int> ans;

  vector<int> parent(n + 1);

  for (int i = 0; i < n; i++) {

    pair<int, int> mx = Get(a[i].first, k);

    parent[i + 1] = mx.second;

    ans = max(ans, make_pair(mx.first + 1, i + 1));

    Add(a[i].first, mx.first + 1, i + 1);

  }

  if (ans.first <= 1) {

    cout << -1 << '\n';

  } else {

    cout << ans.first << '\n';

    vector<int> seq;

    int j = ans.second;

    while (j > 0) {

      seq.push_back(j);

      while (j < 1 || j > n) {}

      j = parent[j];

    }

    // assert((int) seq.size() == ans.first);

    for (int i = 0; i < (int) seq.size(); i++) {

      while (seq[i] - 1 < 0 || seq[i] - 1 >= n) {}

      cout << a[seq[i] - 1].second + 1 << ' ';

    }

    cout << '\n';

  }

  return 0;

}