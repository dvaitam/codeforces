#include <iostream>
#include <fstream>
#include <cstdio>
#include <sstream>
#include <string>
#include <cmath>
#include <stack>
#include <list>
#include <queue>
#include <deque>
#include <set>
#include <map>
#include <vector>
#include <algorithm>
#include <numeric>
#include <utility>
#include <functional>
#include <iomanip>
#include <cstring>
using namespace std;

typedef long long ll;
typedef unsigned long long ull;
#define SZ size()
#define PB push_back
#define SORT(a) sort((a).begin(), (a).end())
#define REV(a) reverse((a).begin(), (a).end())
#define FOR(i, a, b) for(int i = (a); i < (b); i++)
#define TR(i, a) for(typeof(a.begin()) i = a.begin(); i != a.end(); i++)

#define DEBUG(a) cout << #a << ": " << (a) << endl;
#define DEBUG1D(a, x1, x2) { cout << #a << ":"; for(int _i = (x1); _i < (x2); _i++){ cout << " " << a[_i]; } cout << endl; }
#define DEBUG2D(a, x1, x2, y1, y2) { cout << #a << ":" << endl; for(int _i = (x1); _i < (x2); _i++){ for(int _j = (y1); _j < (y2); _j++){ cout << (_j > y1 ? " " : "") << a[_i][_j]; } cout << endl; } }

int n;
bool a[111][111];

void solve()
{
    bool badr = 0, badc = 0;

    FOR(i, 0, n)
    {
        int cnt = 0;
        FOR(j, 0, n) cnt += a[i][j];
        badr |= (cnt == n);
    }

    FOR(i, 0, n)
    {
        int cnt = 0;
        FOR(j, 0, n) cnt += a[j][i];
        badc |= (cnt == n);
    }

    if(badc and badr)
    {
        cout << -1 << endl;
        return;
    }

    if(badc)
    {
        FOR(i, 0, n)
        {
            FOR(j, 0, n)
            {
                if(!a[i][j])
                {
                    cout << i + 1 << " " << j + 1 << endl;
                    break;
                }
            }
        }
    }

    else
    {
        FOR(i, 0, n)
        {
            FOR(j, 0, n)
            {
                if(!a[j][i])
                {
                    cout << j + 1 << " " << i + 1 << endl;
                    break;
                }
            }
        }
    }
}

int main()
{
    //freopen("A.in", "r", stdin);
    while(cin >> n)
    {
        FOR(i, 0, n) FOR(j, 0, n)
        {
            char c;
            cin >> c;
            a[i][j] = (c == '.' ? 0 : 1);
        }
        solve();
    }
}