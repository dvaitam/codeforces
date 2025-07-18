#ifndef LOCAL
#pragma GCC optimize ("O3")
#endif

#include <bits/stdc++.h>

using namespace std;

#define sim template < class c
#define ris return * this
#define dor > debug & operator <<
#define eni(x) sim > typename \
enable_if<sizeof dud<c>(0) x 1, debug&>::type operator<<(c i) {
sim > struct rge { c b, e; };
sim > rge<c> range(c i, c j) { return {i, j}; }
sim > auto dud(c* x) -> decltype(cerr << *x, 0);
sim > char dud(...);
struct debug {
#ifdef LOCAL
~debug() { cerr << endl; }
eni(!=) cerr << boolalpha << i; ris; }
eni(==) ris << range(begin(i), end(i)); }
sim, class b dor(pair < b, c > d) {
  ris << "(" << d.first << ", " << d.second << ")";
}
sim dor(rge<c> d) {
  *this << "[";
  for (c it = d.b; it != d.e; ++it)
    *this << ", " + 2 * (it == d.b) << *it;
  ris << "]";
}
#else
sim dor(const c&) { ris; }
#endif
};
#define imie(x...) " [" #x ": " << (x) << "] "

using ld = long double;
using ll = long long;

constexpr int mod = 1000 * 1000 * 1000 + 7;
constexpr int odw2 = (mod + 1) / 2;

void OdejmijOd(int& a, int b) { a -= b; if (a < 0) a += mod; }
int Odejmij(int a, int b) { OdejmijOd(a, b); return a; }
void DodajDo(int& a, int b) { a += b; if (a >= mod) a -= mod; }
int Dodaj(int a, int b) { DodajDo(a, b); return a; }
int Mnoz(int a, int b) { return (ll) a * b % mod; }
void MnozDo(int& a, int b) { a = Mnoz(a, b); }
int Pot(int a, int b) { int res = 1; while (b) { if (b % 2 == 1) MnozDo(res, a); a = Mnoz(a, a); b /= 2; } return res; }
int Odw(int a) { return Pot(a, mod - 2); }
void PodzielDo(int& a, int b) { MnozDo(a, Odw(b)); }
int Podziel(int a, int b) { return Mnoz(a, Odw(b)); }
int Moduluj(ll x) { x %= mod; if (x < 0) x += mod; return x; }

template <typename T> T Maxi(T& a, T b) { return a = max(a, b); }
template <typename T> T Mini(T& a, T b) { return a = min(a, b); }

constexpr int nax = 200 * 1000 + 105;

// Przedziały odpowiadające ścieżce z v do lca mają first>=second, zaś te dla
// ścieżki z lca do u mają first<=second, przedziały są po kolei, lca występuje
// tam dwa razy, najpierw jako second, a zaraz potem jako first.
vector<int> drz[nax];
int prel, roz[nax], jump[nax], pre[nax], post[nax], fad[nax];

void usun(int w, int o) {
  for (int i = 0; i < (int) drz[w].size(); i++) {
    if (drz[w][i] == o) {
      swap(drz[w][i], drz[w].back());
      i--;
      drz[w].pop_back();
      continue;
    }
    usun(drz[w][i], w);
  }
}

void dfs_roz(int v) {
  roz[v] = 1;
  for (int& i : drz[v]) {
    fad[i] = v;
    dfs_roz(i);
    roz[v] += roz[i];
    if (roz[i] > roz[drz[v][0]]) swap(i, drz[v][0]);
  }
}
void dfs_pre(int v) {
  if (!jump[v]) jump[v] = v;
  pre[v] = ++prel;
  if (!drz[v].empty()) jump[drz[v][0]] = jump[v];
  for (int i : drz[v]) dfs_pre(i);
  post[v] = prel;
}
int lca(int v, int u) {
  while (jump[v] != jump[u]) {
    if (pre[v] < pre[u]) swap(v, u);
    v = fad[jump[v]];
  }
  return (pre[v] < pre[u] ? v : u);
}
vector<pair<int, int>> path_up(int v, int u) {
  vector<pair<int, int>> ret;
  while (jump[v] != jump[u]) {
    ret.emplace_back(pre[jump[v]], pre[v]);
    v = fad[jump[v]];
  }
  ret.emplace_back(pre[u], pre[v]);
  return ret;
}
vector<pair<int, int>> get_path(int v, int u) {
  int w = lca(v, u);
  auto ret = path_up(v, w);
  auto pom = path_up(u, w);
  for (auto& i : ret) swap(i.first, i.second);
  while (!pom.empty()) {
    ret.push_back(pom.back());
    pom.pop_back();
  }
  return ret;
}
struct Sat {
  int n = 0, ile = 0;
  vector<vector<int>> imp;
  vector<bool> vis;
  vector<int> val, sort;

  void DfsMark(int x) {
    vis[x] = false;
    val[x] = (val[x ^ 1] == -1);
    for (int i : imp[x]) if (vis[i]) DfsMark(i);
  }
  void Dfs(int x) {
    vis[x] = true;
    for (int i : imp[x ^ 1]) if (!vis[i ^ 1]) Dfs(i ^ 1);
    sort[--ile] = x;
  }
  //Sat(int m) : n(m * 2), ile(n), imp(n), vis(n), val(n, -1), sort(n) {}
  int Add() {
    n += 2;
    ile += 2;
    imp.emplace_back();
    imp.emplace_back();
    vis.push_back(false);
    vis.push_back(false);
    val.push_back(-1);
    val.push_back(-1);
    sort.push_back(0);
    sort.push_back(0);
    return n - 2;
  }
  void Or(int a, int b) {
    imp[a ^ 1].push_back(b);
    imp[b ^ 1].push_back(a);
  }
  void Impl(int a, int b) {
    Or(a ^ 1, b);
  }
  void Nie(int a) {
    Or(a ^ 1, a ^ 1);
  }
  void Tak(int a) {
    Or(a, a);
  }
  void Nand(int a, int b) {
    Or(a ^ 1, b ^ 1);
  }
  bool Run() {
    for (int i = 0; i < n; i++) if (!vis[i]) Dfs(i);
    for (int i : sort) if (vis[i]) DfsMark(i);
    for (int i = 0; i < n; i++)
      if (val[i]) for (int x : imp[i]) if (!val[x]) return false;
    return true;
  }
};

//int main() {
//  Sat sat(3);
//  sat.Or(2 * 0 + 1, 2 * 1 + 0);  // !x_0 or  x_1
//  sat.Or(2 * 1 + 0, 2 * 2 + 1);  //  x_1 or !x_2
//  debug() << imie(sat.Run());
//  debug() << imie(sat.val); // [0, 1, 1, 0, 0, 1] = [x0, !x0, x1, !x1, x2, !x2]
//}

int n2;
vector<int> przedzial[nax * 8];
vector<int> nizej[nax * 8];
Sat sat;

void Wrzuc(int w, int p, int k, int a, int b, int zm) {
  if (k < a or b < p) return;
  if (a <= p and k <= b) {
    if (przedzial[w].empty() or przedzial[w].back() != zm) {
      przedzial[w].push_back(zm);
    }
    return;
  }
  if (nizej[w].empty() or nizej[w].back() != zm) {
    nizej[w].push_back(zm);
  }
  Wrzuc(w * 2, p, (p + k) / 2, a, b, zm);
  Wrzuc(w * 2 + 1, (p + k) / 2 + 1, k, a, b, zm);
}

void Dodaj(int zm, int a, int b) {
  const int l = pre[lca(a, b)];
  #ifdef LOCAL
  debug() << imie(zm) imie(a) imie(b) imie(l) imie(get_path(a, b));
  #endif
  for (auto p : get_path(a, b)) {
    if (p.first > p.second) {
      swap(p.first, p.second);
    }
    if (p.first == l) p.first++;
    if (p.second == l) p.second--;
    if (p.first <= p.second) {
      assert(0 <= p.first and p.first <= p.second and p.second < n2);
      debug() << imie(p);
      Wrzuc(1, 0, n2 - 1, p.first, p.second, zm);
    }
  }
}

void Ogarnij(const vector<int>& zm, const vector<int>& nizejx) {
  if (zm.empty()) return;
  debug() << imie(zm) imie(nizejx);
  int nor = sat.Add();
  for (int x : nizejx) {
    sat.Impl(x, nor);
  }
  int last = nor;
  for (int z : zm) {
    sat.Nand(z, last);
    const int nl = sat.Add();
    sat.Impl(last, nl);
    sat.Impl(z, nl);
    last = nl;
  }
}

int main() {
  ios_base::sync_with_stdio(0);
  cin.tie(0);
  int n;
  cin >> n;
  for (int i = 1; i < n; i++) {
    int a, b;
    cin >> a >> b;
    drz[a].push_back(b);
    drz[b].push_back(a);
  }

  usun(1, -1);
  dfs_roz(1);
  dfs_pre(1);

  n2 = 1;
  while (n2 <= n) n2 *= 2;

  int m;
  cin >> m;
  vector<int> zmienne;

  for (int i = 0; i < m; i++) {
    int a, b, c, d;
    cin >> a >> b >> c >> d;
    const int zm = sat.Add();
    zmienne.push_back(zm);
    Dodaj(zm, a, b);
    Dodaj(zm ^ 1, c, d);
  }

  #ifdef LOCAL
  for (int i = 1; i < n2 * 2; i++) {
    debug() << "przedzial[" << i << "] = " << przedzial[i];
    debug() << "nizej[" << i << "] = " << nizej[i];
  }
  #endif

  for (int i = 1; i < n2 * 2; i++) {
    Ogarnij(przedzial[i], nizej[i]);
  }

  if (!sat.Run()) {
    cout << "NO" << endl;
    return 0;
  }
  cout << "YES" << endl;
  for (int zm : zmienne) {
    if (sat.val[zm]) {
      cout << 1 << '\n';
    } else {
      cout << 2 << '\n';
    }
  }
  return 0;
}