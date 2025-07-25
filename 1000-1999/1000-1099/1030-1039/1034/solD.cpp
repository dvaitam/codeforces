#include <iostream>
#include <iomanip>
#include <algorithm>
#include <map>
#include <vector>
#include <string>
#include <set>
#include <cmath>

using namespace std;
typedef long long ll;
const int MAXN = 300100;
const int INF = 2e9;

int N, K;
int lo[MAXN], hi[MAXN];
vector <pair <int, int> > event[MAXN];

map <int, int> mm;
int ccover[MAXN];

bool works (ll M)
{
    for (int i = 0; i <= N; i++)
        ccover[i] = 0;
    ccover[N] = INF;

    ll ctot = INF;
    int cstop = N;

    ll ncnt = 0;
    for (int i = N - 1; i >= 0; i--)
    {
        for (int j = 0; j < event[i].size(); j++)
        {
            int color = event[i][j].first, diff = event[i][j].second;
            ccover[color] += diff;
            if (color <= cstop)
                ctot += diff;
        }
        while (ctot - ccover[cstop] >= M)
        {
            ctot -= ccover[cstop];
            cstop--;
        }
        ncnt += N - cstop;
    }

    return ncnt >= K;
}

ll vrun (ll M)
{
    for (int i = 0; i <= N; i++)
        ccover[i] = 0;
    ccover[N] = INF;

    ll ctot = INF;
    int cstop = N;
    ll cc = 0;

    ll ncnt = 0;
    ll ans = 0;
    for (int i = N - 1; i >= 0; i--)
    {
        for (int j = 0; j < event[i].size(); j++)
        {
            int color = event[i][j].first, diff = event[i][j].second;
            ccover[color] += diff;
            if (color <= cstop)
                ctot += diff;
            else
            {
                //cout << "ADD " << color << " " << diff << "\n";
                cc += diff * (ll) (N - color);
            }
        }
        while (ctot - ccover[cstop] >= M)
        {
            ctot -= ccover[cstop];
            //cout << "PASS " << cstop << " " << ccover[cstop] << "\n";
            cc += (N - cstop) * (ll) ccover[cstop];
            cstop--;
        }

        int nv = N - cstop;
        //cout << i << " " << nv << " " << ctot << " " << cc << "\n";
        ans += (cc + nv * ctot);
        ncnt += nv;
    }

    return ans - (ncnt - K) * M;
}

int main()
{
    scanf ("%d %d", &N, &K);
    for (int i = 0; i < N; i++)
        scanf ("%d %d", lo + i, hi + i);

    mm[0] = N; // all covered
    mm[INF] = N; // bookkeeping
    for (int i = N - 1; i >= 0; i--)
    {
        int l = lo[i], h = hi[i];
        auto it = mm.upper_bound(l);
        auto nit = it;
        nit--;
        int nv = nit->second;
        event[i].push_back(make_pair(nit->second, l - min (h, it->first)));
        int nl = nv;
        while (it->first < h)
        {
            nl = it->second;
            nit = it;
            nit++;
            event[i].push_back(make_pair(it->second, it->first - min (h, nit->first)));
            mm.erase(it);
            it = nit;
        }
        if (it->first != h)
        {
            mm[h] = nl;
        }
        mm[l] = i;

        event[i].push_back(make_pair (i, h - l));
    }

    /*for (int i = N - 1; i >= 0; i--)
    {
        cout << "EVENT " << i << "\n";
        for (int j = 0; j < event[i].size(); j++)
            cout << event[i][j].first << " " << event[i][j].second << "\n";
    }*/

    int hlo = 0, hhi = 1e9;
    while (hlo < hhi)
    {
        int mid = (hlo + hhi + 1) / 2;
        if (works (mid))
            hlo = mid;
        else
            hhi = mid - 1;
    }
    cout << vrun (hlo) << "\n";
    //cout << hlo << "\n";
}