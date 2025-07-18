#include <bits/stdc++.h>

using namespace std;

using ll = long long;

const int maxn = 1e5 + 10;
const int inf = 1e9 + 5;
const ll llinf = 1e18 + 5;
const int mod = 998244353;

vector<pair<int, int>> pos[maxn];
int dp[maxn];

void add(int& x, int y) {
  x += y;
  if (x >= mod) {
    x -= mod;
  }
}

void sub(int& x, int y) {
  x -= y;
  if (x < 0) {
    x += mod;
  }
}

int total = 0;

const int szb = 330;
const int cntb = 330;

int cnt[maxn];
int freq[cntb][maxn];
int global[cntb];
int k;

void incseg(int l, int r) {
//  cout << "incseg" << ' ' << l << ' ' << r << ' ';
//  cout << total << endl;
  int L = l / szb;
  for (int i = l; i <= r; ++i) {
    sub(freq[L][cnt[i]], dp[i - 1]);
    if (cnt[i] + global[L] == k) {
      sub(total, dp[i - 1]);
    }
    cnt[i]++;
    add(freq[L][cnt[i]], dp[i - 1]);
  }
//  cout << total << endl;
}

int getondep(int i, int j) {
  if (j < 0) return 0;
  return freq[i][j];
}

void inc(int l, int r) {
  int L = l / szb;
  int R = r / szb;
  if (L == R) {
    incseg(l, r);
  } else {
    incseg(l, szb * (L + 1) - 1);
    incseg(szb * R, r);
    for (int i = L + 1; i < R; ++i) {
      sub(total, getondep(i, k - global[i]));
      ++global[i];
    }
  }
}

void decseg(int l, int r) {
//  cout << "decseg" << ' ' << l << ' ' << r << ' ';
//  cout << total << endl;
  int L = l / szb;
  for (int i = l; i <= r; ++i) {
    sub(freq[L][cnt[i]], dp[i - 1]);
    cnt[i]--;
    add(freq[L][cnt[i]], dp[i - 1]);
    if (cnt[i] + global[L] == k) {
      add(total, dp[i - 1]);
    }
  }
}


void dec(int l, int r) {
  int L = l / szb;
  int R = r / szb;
  if (L == R) {
    decseg(l, r);
  } else {
    decseg(l, szb * (L + 1) - 1);
    decseg(szb * R, r);
    for (int i = L + 1; i < R; ++i) {
      --global[i];
      add(total, getondep(i, k - global[i]));
    }
  }
}

int main() {
//  freopen("input.txt", "r", stdin);
//  freopen("output.txt", "w", stdout);
  std::ios_base::sync_with_stdio(false);

  int n;
  cin >> n >> k;
  dp[0] = 1;
  for (int i = 1; i <= n; ++i) {

    add(freq[i / szb][0], dp[i - 1]);
    add(total, dp[i - 1]);
    int x;
    cin >> x;

    if (pos[x].empty()) {
      inc(1, i);
//      for (int j = 1; j <= i; ++j) {
//        ++cnt[j];
//      }
      pos[x].emplace_back(1, i);
    } else {
      auto last = pos[x].back();
      int q = last.second;
      dec(last.first, last.second);
      inc(last.second + 1, i);
//      for (int j = last.first; j <= last.second; ++j) {
//        cnt[j]--;
//      }
//      for (int j = last.second + 1; j <= i; ++j) {
//        ++cnt[j];
//      }
      pos[x].emplace_back(last.second + 1, i);
    }

//    for (int j = 1; j <= i; ++j) {
//      if (cnt[j] <= k) {
//        add(dp[i], dp[j - 1]);
//      }
//    }
    dp[i] = total;
//    cout << i << ' ' << dp[i] << endl;

  }
  cout << dp[n] << endl;
  return 0;
}