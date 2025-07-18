#define _CRT_SECURE_NO_DEPRECATE
#pragma comment(linker, "/STACK:167772160000")
#include <iostream>
#include <fstream>
#include <cstdio>
#include <stdio.h>
#include <cstdlib>
#include <stdlib.h>
#include <string>
#include <list>
#include <fstream>
#include <algorithm>
#include <cmath>
#include <map>
#include <vector>
#include <iomanip>
#include <queue>
#include <deque>
#include <set>
#include <unordered_set>
#include <stack>
#include <sstream>
#include <assert.h>
#include <unordered_map>
#include <functional>
#include <climits>
#include <cstring>
using namespace std;
typedef unsigned long long ull;
typedef long long ll;
typedef pair<ll, ll> pll;
typedef pair<ull, ull> pull;
typedef pair<int, int> pii;
typedef pair<double, double> pdd;
//typedef uint64_t ull;
//typedef std::pair<long double,long double> pdd;
#define fori(N)          for(int i = 0; i<(N); i++)
#define forj(N)         for(int j = 0; j<(N); j++)
#define fork(N)         for(int k = 0; k<(N); k++)
#define forl(N)         for(int l = 0; l<(N); l++)
#define ford(N)         for(int d = 0; d<(N); d++)
#define fori1(N)          for(int i = 1; i<=(N); i++)
#define forj1(N)         for(int j = 1; j<=(N); j++)
#define fork1(N)         for(int k = 1; k<=(N); k++)
#define ford1(N)         for(int d = 1; d<=(N); d++)
#define PI (2*asin(1))
#define read(n) scanf("%d", &n);
#define read2(n, m) scanf("%d%d", &n, &m);
#define readll(n) scanf("%I64d", &n);
#define mp make_pair
#define double long double


void redirectIO() {
	ios::sync_with_stdio(false); cin.tie(0);
#if defined(_DEBUG) || defined(_RELEASE)
	freopen("input.txt", "r", stdin);
	freopen("output.txt", "w", stdout);
#endif
}
ll been[110000];
int solve(int nn, int mm, int ss, ll tt) {
	fori(nn + 5)been[i] = 0;
	int n, m; n = nn; m = mm;
	int s; s = ss;
	s--;
	ll t; t = tt;
	while (t % n != 0) {
		int move = t%n;
		if (s < m)s += move;
		else s -= move;
		s = (s%n + n) % n;
		t--;
	}
	ll times = t / n;
	while (times > 0) {
		if (been[s]) {
			times %= (been[s] - times);
			break;
		}
		been[s] = times;
		for (int i = n - 1; i > 0; i--) {
			if (s < m)s += i;
			else s -= i;
			if (s < 0)s += n;
			else if (s >= n)s -= n;
		}
		times--;
	}
	if (times != 0) {
		int i = 0;
		while (been[i] == 0 || been[s] - been[i] != times)i++;
		s = i;
		times = 0;
	}
	return s + 1;
	
}
int main()
{
	redirectIO();
	int nn, mm, ss; ll  tt; cin >> nn >> mm >> ss >> tt;
	cout << solve(nn, mm, ss, tt) << endl;
	return 0;
}