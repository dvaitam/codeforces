#include <bits/stdc++.h>

#define Nmax 5005
#define ll long long
#define md 1000000007

using namespace std;

int N, M;
int P[Nmax];
int cnt[Nmax][Nmax];
vector <int> G[Nmax];
int C[Nmax];
bool used[Nmax];
bool found[Nmax];
int D;
int K[Nmax];
int ans[Nmax];
int L = -1;
int Ri[Nmax], Le[Nmax];

bool pairUp(int node)
{
    if(used[node])
        return false;
    used[node] = true;
    for(auto it : G[node])
    {
        if(Le[it] == -1)
        {
            Le[it] = node;
            Ri[node] = it;
            return true;
        }
    }

    for(auto it : G[node])
        if(pairUp(Le[it]))
        {
            Le[it] = node;
            Ri[node] = it;
            return true;
        }
    return false;
}

int main()
{
    cin >> N >> M;
    for(int i = 1; i <= N; i++)
        cin >> P[i];
    for(int i = 1; i <= N; i++)
        cin >> C[i];
    cin >> D;
    for(int i = 1; i <= D; i++)
    {
        cin >> K[i];
        found[K[i]] = true;
    }
    int newD = D;
    for(int i = 1; i <= N; i++)
        if(!found[i])
            K[++newD] = i;
    for(int i = 0; i <= 5000; i++)
        Le[i] = -1;
    for(int i = N; i >= 1; i--)
    {
        if(i <= D)
            ans[i] = L + 1;
        if(cnt[C[K[i]]][P[K[i]]] == 0)
            G[P[K[i]]].push_back(C[K[i]]);
        for(bool ok = true; ok;)
        {
            ok = false;
            memset(used, false, sizeof(used));
            if(pairUp(L + 1))
                L++, ok = true;
        }
        cnt[C[K[i]]][P[K[i]]] = 1;
    }
    for(int i = 1; i <= D; i++)
        cout << ans[i] << "\n";
    return 0;
}