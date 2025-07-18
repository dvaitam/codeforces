//#pragma comment(linker, "/STACK:134217728")

#include <iostream>
#include <sstream>
#include <cstdio>
#include <cstdlib>
#include <cmath>
#include <memory>
#include <cctype>
#include <cstring>
#include <vector>
#include <list>
#include <queue>
#include <deque>
#include <stack>
#include <map>
#include <set>
#include <algorithm>
#include <numeric>
using namespace std;

typedef long long Int;
typedef pair<int,int> PII;
typedef vector<int> VInt;

#define FOR(i, a, b) for(i = (a); i < (b); ++i)
#define RFOR(i, a, b) for(i = (a) - 1; i >= (b); --i)
#define CLEAR(a, b) memset(a, b, sizeof(a))
#define SIZE(a) int((a).size())
#define ALL(a) (a).begin(),(a).end()
#define PB push_back
#define MP make_pair

const int DR[] = { -1, 1, 0, 0 };
const int DC[] = { 0, 0, -1, 1 };

int R[64];
int C[64];
int TR[64];
int TC[64];
int A[64][64];
vector<pair<PII, PII>> res;

void move(int id, int r, int c)
{
	if (R[id] != r && C[id] != c) throw 0;
	int dir = -1;
	if (R[id] > r) dir = 0;
	if (R[id] < r) dir = 1;
	if (C[id] > c) dir = 2;
	if (C[id] < c) dir = 3;
	if (dir == -1) return;

	while(R[id] != r || C[id] != c)
	{
		PII p1(R[id], C[id]);
		if (A[R[id]][C[id]] != id) throw 0;
		A[R[id]][C[id]] = -1;
		R[id] += DR[dir];
		C[id] += DC[dir];
		if (A[R[id]][C[id]] != -1) throw 0;
		A[R[id]][C[id]] = id;
		PII p2(R[id], C[id]);
		res.push_back(MP(p1, p2));
	}
}

int solve(int n, int m)
{
	res.clear();
	if (n == 1) return 0;
	int i;
	vector<pair<PII, int>> v(m);
	FOR(i, 0, m)
		v[i] = MP(PII(R[i], R[i] == 0 ? C[i] : -C[i]), i);

	sort(ALL(v));
	FOR(i, 0, m)
	{
		int tr = 0;
		int tc = i;
		int id = v[i].second;

		if (C[id] < tc)	move(id, R[id], tc);
		move(id, tr, C[id]);
		move(id, tr, tc);
	}

	FOR(i, 0, m)
		v[i] = MP(PII(TR[i], TC[i]), i);

	sort(ALL(v));
	RFOR(i, m, 0)
		if (v[i].first.first > 1)
		{
			int id = v[i].second;
			move(id, 1, C[id]);
			move(id, R[id], TC[id]);
			move(id, TR[id], TC[id]);
		}

	v.clear();
	FOR(i, 0, m)
		if (TR[i] <= 1)
			v.push_back(MP(PII(TR[i], TR[i] == 0 ? TC[i] : -TC[i]), i));


	sort(ALL(v));
	int am = SIZE(v);
	FOR(i, 0, am)
	{
		int id = v[i].second;
		if (C[id] < i) throw 0;
		while (C[id] != i)
		{
			if (A[0][C[id] - 1] == -1) move(id, 0, C[id] - 1);
			else
			{
				int id2 = A[0][C[id] - 1];
				move(id2, 1, C[id2]);
				move(id, 0, C[id] - 1);
				move(id2, 1, C[id2] + 1);
				move(id2, 0, C[id2]);
			}
		}
	}

	RFOR(i, am, 0)
	{
		int id = v[i].second;
		if (TR[id] != 0)
		{
			move(id, 0, n - 1);
			move(id, 1, n - 1);
		}
		move(id, TR[id], TC[id]);
	}

	return SIZE(res);
}

int main()
{
	int n, m;
	n = m = 4;

	int i;

	scanf("%d%d", &n, &m);
	CLEAR(A, -1);
	FOR(i, 0, m)
	{
		int r, c;
		scanf("%d%d", &r, &c);
		--r;
		--c;
		R[i] = r;
		C[i] = c;
		A[r][c] = i;
	}
	FOR(i, 0, m)
	{
		int r, c;
		scanf("%d%d", &r, &c);
		--r;
		--c;
		TR[i] = r;
		TC[i] = c;
	}

	solve(n, m);


	printf("%d\n", SIZE(res));
	FOR(i, 0, SIZE(res))
		printf("%d %d %d %d\n", res[i].first.first + 1, res[i].first.second + 1, res[i].second.first + 1, res[i].second.second + 1);

	return 0;
};