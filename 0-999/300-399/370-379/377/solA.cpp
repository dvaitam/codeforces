/**
 * Copyright (c) 2013 Authors. All rights reserved.
 * 
 * FileName: A.cpp
 * Author: Beiyu Li <sysulby@gmail.com>
 * Date: 2013-12-29
 */
#include <cstdio>
#include <cstdlib>
#include <cstring>
#include <cmath>
#include <cctype>
#include <climits>
#include <deque>
#include <list>
#include <map>
#include <queue>
#include <set>
#include <stack>
#include <vector>
#include <iostream>
#include <iomanip>
#include <sstream>
#include <algorithm>
#include <bitset>
#include <complex>
#include <string>
#include <utility>

using namespace std;

// #pragma comment(linker, "/STACK:102400000,102400000")

#define REP(i,n) for (int i = 0; i < n; ++i)
#define FOR(i,l,r) for (int i = l; i <= r; ++i)
#define FOREACH(i,c) for (__typeof(c.begin()) i = c.begin(); i != c.end(); ++i)
#define PB push_back
#define MP make_pair

typedef long long LL;
typedef pair<int, int> Pii;
typedef pair<LL, LL> Pll;
#define FS first
#define SC second
typedef complex<double> Point;
typedef complex<double> Vector;
#define X real()
#define Y imag()

const int inf = 0x3f3f3f3f;
const LL infLL = 0x3f3f3f3f3f3f3f3fLL;
const int hash_mod = 1000037;
const double eps = 1e-10;
const double pi = acos(-1);

const int maxn = 500 + 5;

int n, m, k;
char g[maxn][maxn];
int dx[] = {-1, 1, 0, 0}, dy[] = {0, 0, -1, 1};

bool inside(int x, int y)
{
        return x >= 0 && x < n && y >= 0 && y < m;
}

void floor_fill(int x, int y, int s)
{
        queue<Pii> Q;
        g[x][y] = '.'; --s;
        if (s == 0) return;
        Q.push(MP(x, y));
        while (!Q.empty()) {
                x = Q.front().FS, y = Q.front().SC; Q.pop();
                for (int i = 0; i < 4; ++i) {
                        int nx = x + dx[i];
                        int ny = y + dy[i];
                        if (!inside(nx, ny) || g[nx][ny] != 'X') continue;
                        g[nx][ny] = '.'; --s;
                        if (s == 0) return;
                        Q.push(MP(nx, ny));
                }
        }
}

int main()
{
        scanf("%d%d%d", &n, &m, &k);
        for (int i = 0; i < n; ++i)
                scanf("%s", g[i]);
        int s = 0, xx = 0, yy = 0;
        for (int i = 0; i < n; ++i)
                for (int j = 0; j < m; ++j) if (g[i][j] == '.') {
                        ++s;
                        g[i][j] = 'X';
                        xx = i;
                        yy = j;
                }
        if (g[xx][yy] == 'X') floor_fill(xx, yy, s - k);
        for (int i = 0; i < n; ++i)
                puts(g[i]);

        return 0;
}