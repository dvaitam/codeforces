#include <bits/stdc++.h>
#define maxn 200005
using namespace std;

int s, n, m, px, py, a[605][605], b[605][605], c[maxn], p[605][605];
struct temp
{
    int x, y, c;
    temp(int x1 = 0, int y1 = 0, int c1 = 0) { x = x1; y = y1; c = c1; }
    bool operator <(const temp& B) const { return p[x - px + n][y - py + m] < p[B.x - px + n][B.y - py + m]; }
} q[maxn];
vector<temp> v[maxn];
map<int, int> g;

void thu()
{
    cin >> n >> m;
    g[0] = 1;
    s = 1;
    for (int i = 1; i <= n; i++)
        for (int j = 1; j <= m; j++)
        {
            int k;
            cin >> k;
            a[i][j] = (g[k] ? g[k] : g[k] = ++s);
            c[a[i][j]]++;
        }
    for (int i = 1; i <= n; i++)
        for (int j = 1; j <= m; j++)
        {
            int k;
            cin >> k;
            if (k > -1)
            {
                b[i][j] = (g[k] ? g[k] : g[k] = ++s);
                v[a[i][j]].push_back(temp(i, j, b[i][j]));
            }
        }
    cin >> px >> py;
    p[n][m] = 1;
    for (int i = 2, x = n - 1, y = m - 1, k = 1; i <= 600; i += 2, x--, y--)
    {
        for (int j = 0; j < i; j++) p[x][++y] = ++k;
        for (int j = 0; j < i; j++) p[++x][y] = ++k;
        for (int j = 0; j < i; j++) p[x][--y] = ++k;
        for (int j = 0; j < i; j++) p[--x][y] = ++k;
    }
}

void lam()
{
    long long ans = 0;
    q[1] = temp(px, py, b[px][py]);
    for (int l = 1, r = 1, k = a[px][py], u; l <= r; l++, k = u)
    {
        if (k > 1 && (u = q[l].c) != k && c[k])
        {
            c[u] += c[k];
            ans += c[k];
            c[k] = 0;
            px = q[l].x;
            py = q[l].y;
            sort(v[k].begin(), v[k].end());
            for (auto x : v[k]) q[++r] = x;
            v[k].clear();
        }
    }
    cout << ans << endl;
}

int main()
{
    ios_base::sync_with_stdio(0); cin.tie(0); cout.tie(0);
    thu();
    lam();
    return 0;
}