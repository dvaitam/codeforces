#define _CRT_SECURE_NO_WARNINGS
#include <limits.h>
#include <math.h>
#include <numeric>
#include <cstring>
#include <fstream>
#include <map>
#include <iostream>
#include <utility>
#include <set>
#include <algorithm>
#include <bitset>
#include <queue>
#include <functional>
#include <assert.h>
#include <ctime>
#include <utility>
#include <stack>
#include <stdint.h>
#define SI(i) scanf("%d\n", &i)
#define SII(i,j) scanf("%d%d", &i, &j)
#define SFF(i,j) scanf("%.9lf%.9lf", &i, &j)
#define SIII(i,j, k) scanf("%d%d%d", &i, &j, &k)
#define SF(i) scanf("%.9lf", &i)
#define SFFF(i,j, k) scanf("%.9lf%.9lf%.9lf", &i, &j, &k)
#define SL(i) scanf("%I64", &i)
#define SLL(i,j) scanf("%I64%I64", &i, &j)
#define SLLL(i,j, k) scanf("%I64 %I64 %I64", &i, &j, &k)
#define SS(i) scanf("%s\n", i)
#define SSS(i,j) scanf("%s %s", i, j)
#define SC(i) scanf("%c ", &i)
#define SCC(i,j) scanf("%c %c\n", &i, &j)

#define PF(i) printf("%.9lf", i)
#define PI(i) printf("%d ", i)
#define PI2(i) printf("%d", i)
#define PII(i,j) printf("%d %d ", i, j)
#define PL(i) printf("%I64 ", i)
#define PLS(i) printf("%I64", i)
#define PLL(i,j) printf("%I64 %I64 ", i, j)
#define PS(i) printf("%s", i)
#define PSS(i,j) printf("%s %s ", i, j)
#define PC(i) printf("%c", i)
#define PCC(i,j) printf("%c%c ", i, j)
#define PN printf("\n")
#define forin(i, k, n) for(int32_t i = (k); i < (n); ++i)
#define forin2(i, k, n) for(int32_t i = (k); i <= (n); ++i)
#define rforin(i, k, n) for(int32_t i = (k); i > (n); --i)
#pragma comment(linker, "/STACK:100000000000")
using namespace std;
int n, m;
char grid[103][103];
bool was[103][103];
struct a {
	int x, y, s;
};
vector<a> ans;
int main() {
#ifdef _DEBUG
	freopen("1.in", "rt", stdin);
	freopen("1.out", "wt", stdout);
#endif
	SII(n, m);
	forin(i, 0, n) {
		SS(grid[i + 1] + 1);
	}
	forin2(i, 1, n) {
		forin2(j, 1, m) {
			if (grid[i][j] != '*') {
				continue;
			}
			int ii = 1;
			for(; ii < 100; ++ii){
				if (grid[i + ii][j] == '*' && grid[i - ii][j] == '*' && grid[i][j + ii] == '*' && grid[i][j - ii] == '*') {
					was[i][j] = true;
					was[i + ii][j] = true;
					was[i - ii][j] = true;
					was[i][j + ii] = true;
					was[i][j - ii] = true;
				}
				else {
					break;
				}
			}
			if (ii > 1) {
				--ii;
				ans.push_back({ i, j, ii });
			}
		}
	}
	forin2(i, 1, n) {
		forin2(j, 1, m) {
			if (grid[i][j] == '*' && !was[i][j]) {
				PI(-1);
				return 0;
			}
		}
	}
	PI(ans.size());
	PN;
	forin(i, 0, ans.size()) {
		PII(ans[i].x, ans[i].y);
		PI(ans[i].s);
		PN;
	}
}