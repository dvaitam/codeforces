#include <bits/stdc++.h>

#include <ext/hash_map>

#include <ext/numeric>

using namespace std;

using namespace __gnu_cxx;



#define EPS      1e-9

#define F        first

#define S        second

#define pi       acos(-1)

#define ll       long long

#define oo       0x3f3f3f3f

#define sz(x)    (int)x.size()

#define sc(x)    scanf("%d",&x)

#define all(x)   x.begin(),x.end()

#define rall(x)  x.rbegin(),x.rend()



int n, m;



int dist(pair<int, int> &a, pair<int, int> &b) {

  return ((a.F - b.F) * (a.F - b.F) + (a.S - b.S) * (a.S - b.S));

}



vector<pair<int, int> > lol;

vector<pair<int, int> > best, picks;

int mx;



void solve(int i) {

  if (i == sz(lol)) {

    if (sz(picks) != 4)

      return;

    vector<pair<int, int> > tmp = picks;

    sort(all(tmp));

    do {

      int d = dist(tmp[0], tmp[1]) + dist(tmp[1], tmp[2])

          + dist(tmp[2], tmp[3]);

      if (d > mx)

        mx = d, best = tmp;

    } while (next_permutation(all(tmp)));

    return;

  }

  picks.push_back(lol[i]);

  solve(i + 1);

  picks.pop_back();

  solve(i + 1);

}



int main() {

#ifndef ONLINE_JUDGE

  freopen("input.txt", "r", stdin);

//freopen("output.txt", "w", stdout);

#endif

  sc(n), sc(m);

  if (n == 0) {

    cout << 0 << " " << 1 << endl;

    cout << 0 << " " << m << endl;

    cout << 0 << " " << 0 << endl;

    cout << 0 << " " << m - 1 << endl;

    return 0;

  }

  if (m == 0) {

    cout << 1 << " " << 0 << endl;

    cout << n << " " << 0 << endl;

    cout << 0 << " " << 0 << endl;

    cout << n - 1 << " " << 0 << endl;

    return 0;

  }

  lol.push_back(make_pair(0, 0));

  lol.push_back(make_pair(n, 0));

  lol.push_back(make_pair(0, m));

  lol.push_back(make_pair(n, m));

  lol.push_back(make_pair(1, 0));

  lol.push_back(make_pair(0, 1));

  lol.push_back(make_pair(n - 1, 0));

  lol.push_back(make_pair(0, m - 1));

  lol.push_back(make_pair(n, 1));

  lol.push_back(make_pair(1, m));

  lol.push_back(make_pair(n - 1, m));

  lol.push_back(make_pair(n, m - 1));

  sort(all(lol));

  lol.resize(unique(all(lol)) - lol.begin());

  solve(0);

  for (int i = 0; i < 4; i++)

    cout << best[i].F << " " << best[i].S << endl;

}