// https://codeforces.com/problemset/problem/1907/A

#include <bits/stdc++.h>
using namespace std;
#define SET_IO() ios_base::sync_with_stdio(false), cin.tie(nullptr);
using ll = long long;

const string columns = "abcdefgh";

void solve() {
  char c;
  int row;
  cin >> c >> row;

  // print squares in same row
  for(int i=0; i<8; ++i) {
    if(columns[i] == c) continue;  
    cout << columns[i] << row << '\n';
  }

  // print squares in same column
  for(int i=1; i<=8; ++i) {
    if(i == row) continue;  
    cout << c << i << '\n';
  }
}

int main() {
  SET_IO();

  int t;
  cin >> t;
  while (t--) {
    solve();
  }

}