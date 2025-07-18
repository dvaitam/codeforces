#define _CRT_SECURE_NO_WARNINGS
#include <iostream>
#include <vector>
#include <string>
#include <stdio.h>
#include <algorithm>
#include <set>
#include <map>
#include <math.h>
#include <cmath>
#include <queue>
#include <iomanip>
#include <bitset>
#include <memory.h>
#include <cstring>
#include <stack>
#include <unordered_set>
#pragma comment (linker, "/STACK:667177216")
#define ll long long
#define ull unsigned long long
#define INF 1000000007;
#define pb push_back
#define all(x) (x).begin(),(x).end()
#define rall(x) (x).rbegin(),(x).rend()
#define mp make_pair
#define vI vector<int>
#define vvI vector<vector<int>>
#define vLL vector<LL>
#define vS vector<string>
#define fori(i, n) for(int (i)=0; (i)<n; (i)++)
#define forn(it,from,to) for(int (it)=from; (it)<to; (it)++)
#define forI(tmp) for(auto(it)=(tmp).begin();(it)!=(tmp).end();(it)++)
#define PI 3.14159265356
#define LD long double
#define sc(a) scanf("%d", &(a))
#define scLL(a) scanf("%I64d", &(a))
#define mems(a, val) memset(a, val, sizeof(a))
typedef long long LL;
using namespace std;
const LL MOD = 1000000000 + 7;
const LL MAXN = 2000000 + 10;

bool build(vector<vector<int> > &v) {
	vector<int> b = v.back();
	vector<int> a = v[v.size() - 2];
	b.insert(b.begin(), 0);
	int mult = 1;
	for (int i = 0; i < a.size(); ++i) {
		if (mult == -1 && ((a[i] == -1 && b[i] == 1) || (a[i] == 1 && b[i] == -1))) {
			return false;
		}

		if ((a[i] == 1 && b[i] == 1) || (a[i] == -1 && b[i] == -1)) {
			mult = -1;
		}
	}

	for (int i = 0; i < a.size(); ++i) {
		b[i] += mult * a[i];
	}

	v.push_back(b);
	return true;
}

int main() {
#ifdef _DEBUG
	freopen("input.txt", "r", stdin); freopen("output.txt", "w", stdout);
#else
	//freopen("input.txt", "r", stdin); freopen("output.txt", "w", stdout);
#endif
	int n;

	vector<vector<int> > a;
	a.push_back({ 0 });
	a.push_back({ 1 });
	a.push_back({ 0, 1 });
	a.push_back({ 1, 0, 1 });
	for (int i = 0; i < 150; ++i) {
		build(a);
	}
	
	cin >> n;
	cout << n << endl;
	for (int i = 0; i < a[n+1].size(); ++i) {
		cout << a[n+1][i] << " ";
	}

	cout << endl;
	cout << n-1 << endl;
	for (int i = 0; i < a[n].size(); ++i) {
		cout << a[n][i] << " ";
	}
	return 0;
}