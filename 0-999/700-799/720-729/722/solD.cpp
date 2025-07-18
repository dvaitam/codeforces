// Copyright (C) 2016 __debug.

// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License as
// published by the Free Software Foundation; version 3

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program; If not, see <http://www.gnu.org/licenses/>.


#include <cstdio>
#include <cstdlib>
#include <cstring>
#include <cmath>
#include <cctype>
#include <climits>
#include <cassert>
#include <ctime>
#include <iostream>
#include <fstream>
#include <algorithm>
#include <functional>
#include <string>

#define x first
#define y second
#define MP std::make_pair
#define VAL(x) #x " = " << x << " "
#define DEBUG(...) fprintf(stderr, __VA_ARGS__)
#ifdef __linux__
#define getchar getchar_unlocked
#define putchar putchar_unlocked
#endif

typedef long long LL;
typedef std::pair<int, int> Pii;

const int oo = 0x3f3f3f3f;

template<typename T> inline bool chkmax(T &a, T b) { return a < b ? a = b, true : false; }
template<typename T> inline bool chkmin(T &a, T b) { return b < a ? a = b, true : false; }
std::string procStatus()
{
    std::ifstream t("/proc/self/status");
    return std::string(std::istreambuf_iterator<char>(t), std::istreambuf_iterator<char>());
}
template<typename T> T read(T &x)
{
    int f = 1;
    char ch = getchar();
    for (; !isdigit(ch); ch = getchar())
        if (ch == '-')
            f = -1;
    for (x = 0; isdigit(ch); ch = getchar())
        x = 10 * x + ch - '0';
    return x *= f;
}
template<typename T> void write(T x)
{
    if (x == 0) {
        putchar('0');
        return;
    }
    if (x < 0) {
        putchar('-');
        x = -x;
    }
    static char s[20];
    int top = 0;
    for (; x; x /= 10)
        s[++top] = x % 10 + '0';
    while (top)
        putchar(s[top--]);
}


#include <queue>

const int MAXN = 5e4 + 5;

int N;
char Y[MAXN][11];

namespace Mitsuha
{

const int SIZE = MAXN * 31;

int size, ch[SIZE][2];
bool isword[SIZE];
int totheap;
std::priority_queue<int> heap[SIZE];
int id[SIZE];

void insert(int s[], int sz)
{
    int u = 0;
    for (int i = sz; i > 0; --i) {
        int x = s[i];
        // printf("%d", x);
        if (!ch[u][x])
            ch[u][x] = ++size;
        u = ch[u][x];
    }
    // printf("\n");
    isword[u] = true;
}

void dfs(int u, int val)
{
    // printf("dfs(%d, %d)\n", u, val);
    if (ch[u][0])
        dfs(ch[u][0], val << 1);
    if (ch[u][1])
        dfs(ch[u][1], val << 1 | 1);

    int x = id[ch[u][0]], y = id[ch[u][1]];
    if (ch[u][0] && !ch[u][1])
        id[u] = x;
    else if (!ch[u][0] && ch[u][1])
        id[u] = y;
    else if (ch[u][0] && ch[u][1]) {
        // printf("u = %d\n", u);
        if (x == 0)
            id[u] = y;
        else if (y == 0)
            id[u] = x;
        else {
            if (heap[x].size() > heap[y].size())
                std::swap(x, y);
            while (!heap[x].empty()) {
                heap[y].push(heap[x].top());
                // printf("%d: %d -> %d\n", heap[x].top(), x, y);
                heap[x].pop();
            }
            id[u] = y;
        }
    } else {
        // assert(isword[u]);
        id[u] = ++totheap;
    }

    if (isword[u]) {
        heap[id[u]].push(val);
    } else if (u && id[u]) {
        heap[id[u]].pop();
        heap[id[u]].push(val);
    }
}

}

void input()
{
    read(N);
    for (int i = 1; i <= N; ++i) {
        int x;
        read(x);

        static int s[35];
        int sz = 0;
        for (; x; x >>= 1) {
            s[++sz] = x & 1;
        }
        Mitsuha::insert(s, sz);
    }
}

void solve()
{
    Mitsuha::dfs(0, 0);

    static int ans[MAXN];

    // assert((int)Mitsuha::heap[Mitsuha::id[0]].size() == N);

    for (int i = 1; i <= N; ++i) {
        ans[i] = Mitsuha::heap[Mitsuha::id[0]].top();
        Mitsuha::heap[Mitsuha::id[0]].pop();
    }

    for (int i = 1; i <= N; ++i) {
        write(ans[i]); putchar(' ');
    }
    putchar('\n');
}

int main()
{
#ifndef ONLINE_JUDGE
    freopen("D.in", "r", stdin);
    freopen("D.out", "w", stdout);
#endif

    input();
    solve();

    return 0;
}