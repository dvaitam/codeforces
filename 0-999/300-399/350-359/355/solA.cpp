#include <bits/stdc++.h>
#include <iostream>
#include <cstdio>
#include <vector>
#include <algorithm>
#include <ctime>
#include <cstdlib>
#include <stack>
#include <string>
#include <queue>
#include <cmath>
#include <map>
#define ll int64_t
#define ld long double
#define frp(i, a, b) for(int i=a, i<=b; i++)
#define frm(i, a, b) for(int i=a; i>=b; i--)
#define f first
#define s second
#define mp make_pair
#define pb push_back
#define MAXN ((int)(1e9)+7)
#define len length()
#define sz size()
#define a() (a.begin(), a.end())
using namespace std;
bool used[1001][10];

int main() {
	long long k, d;
	cin >> k >> d;
	ll sum = 0;
		
	if(k >= 2 && d == 0)  {
		cout << "No solution";
		return 0;
	}
	if(d == 0) {
		cout << 0;
		return 0;
	}
	if(k == 1) {
		cout << d;
		return 0;
	}
	int u = d - 1;
	int v = 1;
	cout << max(u, v);
	for(int i = 1; i <= k - 2; ++i)
		cout << 0;
	cout << min(u, v);
}