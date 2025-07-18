#include <bits/stdc++.h>

using namespace std;



const int N = 2e5+5;

vector<int> ans_four[N], ans_three[N];

int a[N], freq[N], BIT[N], next_greater[N], next_less[N];



void add(int i, int x) {

  for (; i < N; i += i&-i) BIT[i] += x;

}



int query(int i) {

  int sum = 0;

  for (; i; i -= i&-i) sum += BIT[i];

  return sum;

}



int get_index(int target) {

  if (!target) return 0;

  int sum = 0;

  int z = 0;

  for (int i = 1<<__lg(N-1); i; i >>= 1) {

    if (z+i >= N) continue;

    if (sum+BIT[z+i] < target) {

      sum += BIT[z+i];

      z += i;

    }

  }

  return z+1;

}



int main () {

  ios_base::sync_with_stdio(0); cin.tie(0);

  int n, q;

  cin >> n >> q;

  for (int i = 1; i <= n; i++) cin >> a[i];

  ans_four[n+1].resize(4, n+1);

  ans_three[n+1].resize(3, n+1);

  vector<int> maxstk, minstk;

  int r = n+1;

  add(r, 1);

  next_greater[n+1] = n+1;

  next_less[n+1] = n+1;

  for (int i = n; i >= 1; i--) {

    ans_four[i] = ans_four[i+1];

    ans_three[i] = ans_three[i+1];

    while (!maxstk.empty() && a[maxstk.back()] < a[i]) {

      if (!--freq[maxstk.back()]) add(maxstk.back(), 1);

      maxstk.pop_back();

    }

    if (maxstk.empty()) next_greater[i] = n+1;

    else if (a[maxstk.back()] > a[i]) next_greater[i] = maxstk.back();

    else next_greater[i] = next_greater[maxstk.back()];



    while (!minstk.empty() && a[minstk.back()] > a[i]) {

      if (!--freq[minstk.back()]) add(minstk.back(), 1);

      minstk.pop_back();

    }

    if (minstk.empty()) next_less[i] = n+1;

    else if (a[minstk.back()] < a[i]) next_less[i] = minstk.back();

    else next_less[i] = next_less[minstk.back()];



    maxstk.push_back(i);

    minstk.push_back(i);

    freq[i] = 2;

    do {

      int rp = get_index(query(r)-1);

      if (!rp) break;

      assert(rp > i);

      int u = *lower_bound(maxstk.begin(), maxstk.end(), rp, greater<>());

      int v = *lower_bound(minstk.begin(), minstk.end(), rp, greater<>());

      if (a[u] > a[i] && a[u] > a[rp] && a[v] < a[i] && a[v] < a[rp]) {

        r = rp;

        if (u > v) swap(u, v);

        ans_four[i] = {i, u, v, r};

      }

      else break;

    } while (1);



    if (next_less[next_greater[i]] < ans_three[i].back()) {

      ans_three[i] = {i, next_greater[i], next_less[next_greater[i]]};

    }

    if (next_greater[next_less[i]] < ans_three[i].back()) {

      ans_three[i] = {i, next_less[i], next_greater[next_less[i]]};

    }

  }



  while (q--) {

    int l, r;

    cin >> l >> r;

    if (ans_four[l].back() <= r) {

      cout << "4\n";

      for (int x: ans_four[l]) cout << x << ' ';

      cout << '\n';

    }

    else if (ans_three[l].back() <= r) {

      cout << "3\n";

      for (int x: ans_three[l]) cout << x << ' ';

      cout << '\n';

    }

    else cout << "0\n";

  }

}