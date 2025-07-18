#include <bits/stdc++.h>
using namespace std;
typedef unsigned long long ul;
typedef signed long long ll;

mt19937 mt(time(nullptr));
// uniform_int_distribution< int > rand_int(0, 10) // rand[0, 10] inclusive

void longest(int idx, vector< vector< int > > & path, vector< int > & p, vector< bool > & v) {
  if (v[idx]) return;
  v[idx] = true;
  for (auto c : path[idx]) {
    if (v[c]) continue;
    vector< int > pc;
    longest(c, path, pc, v);
    pc.push_back(c);
    if (p.size() < pc.size()) {
      p.swap(pc);
    }
  }
}

void check(int idx, vector< vector< int > > & path, vector< bool > & v, vector< int > & cl, set< int > & depl, int depth) {
  if (v[idx]) return;
  v[idx] = true;
  if (path[idx].size()==1 && depl.count(depth) == 0) {
    cl.push_back(idx);
    depl.insert(depth);
  }
  for (auto c : path[idx]) {
    check(c, path, v, cl, depl, depth+1);
  }
}

bool solve(int idx, vector< vector< int > > & path, int depth, vector< int > & dd, vector< bool > & v) {
  if (v[idx]) return true;
  v[idx] = true;
  if (dd[depth] == -1) {
    dd[depth] = path[idx].size();
  }
  else if (dd[depth] != (int)(path[idx].size())) {
    return false;
  }
  for (auto c : path[idx]) {
    if (!solve(c, path, depth+1, dd, v)) return false;
  }
  return true;
}

int main(void)
{
  cin.tie(0);
  ios::sync_with_stdio(false);
  cout << fixed;
  int n;
  cin >> n;
  if (n==1) {
    cout << 1 << endl;
    return 0;
  }
  vector< vector< int > > path(n);
  for (int i=0; i<n-1; ++i) {
    int v1, v2;
    cin >> v1 >> v2;
    v1--; v2--;
    path[v1].push_back(v2);
    path[v2].push_back(v1);
  }
  int mn = 0;
  for (int i=0; i<n; ++i) {
    int pis = path[i].size();
    mn = max(mn, pis);
  }
  if (mn<=2) {
    for (int i=0; i<n; ++i) {
      if (path[i].size()==1) {
        cout << i+1 << endl;
        return 0;
      }
    }
  }
  vector< int > p;
  vector< bool > v(n, false);
  longest(0, path, p, v);
  /*
  for (int i=0; i<p.size(); ++i) cout << p[i]+1 << " ";
  cout << endl;
  */
  vector< int > p2;
  vector< bool > v2(n, false);
  longest(p[0], path, p2, v2);
  p2.push_back(p[0]);
  /*
  for (int i=0; i<p2.size(); ++i) cout << p2[i]+1 << " ";
  cout << endl;
  */
  vector< int > dd(n, -1);
  vector< bool > v3(n, false);
  if (solve(p2[0], path, 0, dd, v3)) {
    cout << p2[0]+1 << endl;
    return 0;
  }
  vector< int > dd2(n, -1);
  vector< bool > v4(n, false);
  if (solve(p2[p2.size()-1], path, 0, dd2, v4)) {
    cout << p2[p2.size()-1]+1 << endl;
    return 0;
  }
  vector< int > dd3(n, -1);
  vector< bool > v5(n, false);
  if (solve(p2[p2.size()/2], path, 0, dd3, v5)) {
    cout << p2[p2.size()/2]+1 << endl;
    return 0;
  }
  vector< bool > v6(n, false);
  for (int i=0; i<p2.size(); ++i) {
    if (i==p2.size()/2) continue;
    v6[p2[i]] = true;
  }
  vector< int > cl;
  set< int > depl;
  check(p2[p2.size()/2], path, v6, cl, depl, 0);
  for (int i=0; i<cl.size(); ++i) {
    vector< int > ddx(n, -1);
    vector< bool > vx(n, false);
    if (solve(cl[i], path, 0, ddx, vx)) {
      cout << cl[i]+1 << endl;
      return 0;
    }
  }
  cout << -1 << endl;
  return 0;
}