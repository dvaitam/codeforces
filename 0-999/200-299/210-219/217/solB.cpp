#include <stdio.h>
#include <iostream>
#include <algorithm>
#include <vector>
#include <string>
#include <queue>
#include <map>
#include <set>
#include <cmath>
#include <sstream>
#include <stack>
#include <cassert>

#define pb push_back
#define mp make_pair
#define PI 3.1415926535897932384626433832795
#define sqr(x) (x)*(x)
#define forn(i, n) for(int i = 0; i < n; ++i)
#define ALL(x) x.begin(), x.end()
#define sz(x) int((x).size())
#define X first
#define Y second
typedef long long ll;
typedef unsigned long long ull;
typedef long double ld;
using namespace std;
typedef pair<int,int> pii;
const int INF = 2147483647;
const ll LLINF = 9223372036854775807LL;
int n, r;
const int maxn = 1000010;
char s[maxn];
void doh() {
	printf("IMPOSSIBLE\n"), exit(0);
}
int main()
{
#ifndef ONLINE_JUDGE
	freopen("input.txt", "r", stdin);
	//freopen("output.txt", "w", stdout);
#endif
	scanf("%d%d", &n, &r);
	if (n == 1) {
		if (r == 1) printf("0\nT\n"), exit(0);
		else doh();
	}
	if (n == 2) {
		if (r == 2) printf("0\nTB\n"), exit(0);
		else doh();
	}
	
	int bi = -1;
	const int cinf = 1000000000;
	int bval = cinf;
	for (int i = 1; i+i <= r; ++i) {
		int x = i;
		int y = r-i;
		int ans = 0;
		int err = 0;
		while (x&&y) {
			if (x<y) swap(x,y);
			int dd = x/y;
			ans += dd;
			err += dd-1;
			x%=y;
		}
		x+=y;
		err--;
		if (x == 1 && ans+1 == n) {
			if (bval>err) {
				bval = err;
				bi = i;
			}
		}
	}
	if (bi == -1) doh();
	int len = 1;
	int u = bi;
	int d = r-bi;
	while (u>1||d>1) {
		if (u>d) s[len++] = 'T', u-=d;
		else s[len++] = 'B', d-=u;
	}
	reverse(s+1,s+len);
	s[0] = s[1]^('T'^'B');
	s[len] = s[len-1]^('T'^'B');
	s[++len] = 0;
	if (s[0] == 'B') for (int i = 0; i < len; ++i) s[i] ^= ('T'^'B');
	printf("%d\n%s\n", bval, s);
	return 0;
}