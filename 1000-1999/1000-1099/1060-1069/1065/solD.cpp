#include <iostream>
#include <cstdio>
#include <utility>
#include <queue>
#include <algorithm>
using namespace std;
const int maxN = 11;
const int maxD = maxN * maxN * 2 + 1;
const int maxKind = 3;
typedef pair<int, int> TCost;
int n;
int a[maxN][maxN];
int posX[maxN * maxN], posY[maxN * maxN];
TCost d[maxN][maxN][maxKind];
int dx[8] = {2, -2, 2, -2, 1, -1, 1, -1};
int dy[8] = {1, 1, -1, -1, 2, 2, -2, -2};
void ReadInput()
{
    cin >> n;
    for (int i = 1; i <= n; ++i)
        for (int j = 1; j <= n; ++j)
        {
            int a;
            cin >> a;
            posX[a] = i;
            posY[a] = j;
        }
}

struct TPQRec
{
    int x, y, p;
    TCost Cost;
    inline bool operator < (const TPQRec& other) const
    {
        return Cost > other.Cost;
    }
    inline TCost GetD() const
    {
        return d[x][y][p];
    }
    inline void SetD() const
    {
        d[x][y][p] = Cost;
    }
    inline bool Valid() const
    {
        return Cost == GetD();
    }
};

priority_queue<TPQRec> PQ;

void KnightMove(const TPQRec& u)
{
    TPQRec v = u;
    v.Cost.first = u.Cost.first + 1;
    for (int d = 0; d < 8; ++d)
    {
        v.x = u.x + dx[d];
        v.y = u.y + dy[d];
        if (v.x < 1 || v.x > n || v.y < 1 || v.y > n) continue;
        if (v.Cost < v.GetD())
        {
            v.SetD();
            PQ.push(v);
        }
    }
}

void RockMove(const TPQRec& u)
{
    TPQRec v = u;
    v.Cost.first = u.Cost.first + 1;
    for (int i = 1; i <= n; ++i)
    {
        if (i == u.x) continue;
        v.x = i; v.y = u.y;
        if (v.Cost < v.GetD())
        {
            v.SetD();
            PQ.push(v);
        }
    }
    for (int j = 1; j <= n; ++j)
    {
        if (j == u.y) continue;
        v.x = u.x; v.y = j;
        if (v.Cost < v.GetD())
        {
            v.SetD();
            PQ.push(v);
        }
    }
}

void BishopMove(const TPQRec& u)
{
    TPQRec v = u;
    v.Cost.first = u.Cost.first + 1;
    int c = u.x + u.y, d = u.x - u.y;
    for (int i = 1; i <= n; ++i)
    {
        if (i == u.x) continue;
        int j;
        j = c - i;
        if (j >= 1 && j <= n)
        {
            v.x = i; v.y = j;
            if (v.Cost < v.GetD())
            {
                v.SetD();
                PQ.push(v);
            }
        }
        j = i - d;
        if (j >= 1 && j <= n)
        {
            v.x = i; v.y = j;
            if (v.Cost < v.GetD())
            {
                v.SetD();
                PQ.push(v);
            }
        }
    }
}

void Dijkstra(int sx, int sy, int tx, int ty, const TCost dis[maxKind])
{
    while (!PQ.empty()) PQ.pop();
    TCost infty = make_pair(maxD, 0);
    fill_n(&d[0][0][0], sizeof(d) / sizeof(d[0][0][0]), infty);
    for (int piece = 0; piece < maxKind; ++piece)
    {
        TPQRec u = {x:sx, y:sy, p:piece};
        u.Cost = dis[piece];
        u.SetD();
        PQ.push(u);
    }
    int Cnt = 0;
    while (!PQ.empty())
    {
        TPQRec u = PQ.top(), v;
        PQ.pop();
        if (!u.Valid()) continue;
        if (u.x == tx && u.y == ty)
        {
            ++Cnt;
            if (Cnt == maxKind) break;
        }
        v = u;
        for (int piece = 0; piece < maxKind; ++piece)
        {
            if (piece == u.p) continue;
            v.p = piece;
            v.Cost = make_pair(u.Cost.first + 1, u.Cost.second + 1);
            if (v.Cost < v.GetD())
            {
                v.SetD();
                PQ.push(v);
            }
        }
        switch(u.p)
        {
            case 0: KnightMove(u); break;
            case 1: RockMove(u); break;
            case 2: BishopMove(u); break;
        }
    }
}

void Solve()
{
    TCost dis[maxKind];
    for (int piece = 0; piece < maxKind; ++piece)
        dis[piece] = make_pair(0, 0);
    int Finish = n * n;
    int x1 = posX[1], y1 = posY[1];
    for (int dest = 2; dest <= Finish; ++dest)
    {
        int x2 = posX[dest], y2 = posY[dest];
        Dijkstra(x1, y1, x2, y2, dis);
        for (int piece = 0; piece < maxKind; ++piece)
            dis[piece] = d[x2][y2][piece];
        x1 = x2; y1 = y2;
    }
    TCost res = *min_element(dis, dis + maxKind);
    cout << res.first << ' ' << res.second;
}

int main()
{
    #ifndef hoang
    ios_base::sync_with_stdio(false);
    cin.tie(nullptr);
    #endif
    #ifdef hoang
    freopen("input.txt", "r", stdin);
    #endif // hoang
    ReadInput();
    Solve();
}