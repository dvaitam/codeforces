#include <stdio.h>
#include <stack>
#include <map>
#include <string.h>
#include <string>
#include <iostream>
#include <algorithm>
#include <iomanip>
#include <math.h>
#include <vector>
#include <set>
#include <queue>
#include <climits>
#include <unordered_map>
#include <iterator> 
#include <bitset>
#include <complex>
#include <random>
#include <chrono>
using namespace std;
#define ll long long
#define ull unsigned long long
#define ui unsigned int
#define mp make_pair
#define inf32 INT_MAX
#define inf64 LLONG_MAX
#define PI acos(-1)
#define cos45 cos(PI/4)
#define ld long double
#define inf 1000000
#define pii pair<int, int>
#define pll pair<ll, ll>
#define pli pair<ll, int>
#define pil pair<int, ll>
#pragma GCC optimize ("O3")
//#define x first
//#define y second
const int mod = (1e9) + 7, mod2 = 998244353;
const double eps = 1e-10;
const int siz = 1e7 + 5, siz2 = 21, lg = 21, block = 317, block2 = 1000, mxv = 1e6;
int k, ans[2005];
int main()
{
	scanf("%d", &k);
	ans[0] = -1;
	ans[1] = 2000;
	int lft = k;
	for (int i = 2; i < 2000; i++) {
		ans[i] = min(lft, 1000000);
		lft -= ans[i];
	}
	printf("2000\n");
	for (int i = 0; i < 2000; i++) {
		printf("%d ", ans[i]);
	}
	printf("\n");
	return 0;
}