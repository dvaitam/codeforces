// Problem: E. Cells Arrangement
// Contest: Codeforces - Codeforces Round 943 (Div. 3)
// URL: https://codeforces.com/contest/1968/problem/E
// Author : Setsuna
// Memory Limit: 256 MB
// Time Limit: 2000 ms
// 
// Powered by CP Editor (https://cpeditor.org)

#include<bits/stdc++.h>
using namespace std;
using ll = long long;

void solve() {
	int n; cin >> n;
	if(n <= 3) {
		for(int i = 1; i < n; i++) cout << 1 << " " << i << '\n';
		cout << n << " " << n << '\n';
	} else {
		cout << 1 << " " << 1 << '\n';
		for(int i = 3; i < n; i++) {
			cout << 1 << " " << i << '\n';
		}
		cout << n << " " << n-1 << '\n';
		cout << n << " " << n << '\n';
	}
	cout << '\n';
}

int main() {
	ios::sync_with_stdio(false); cin.tie(nullptr);
	int t; cin >> t; while(t--)
	solve();
	return 0;
}