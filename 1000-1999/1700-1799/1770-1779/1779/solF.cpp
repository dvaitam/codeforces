#include <bits/stdc++.h>

using namespace std;

#define ll long long

// #define int ll

#define ff first

#define ss second

#define pb push_back

#define mkp make_pair

typedef pair<int, int> pii;

const ll inf=2e18;

const int maxn=2e5+10;

const int maxai=35;



int n;

int a[maxn]; 

int pai[maxn];

vector<int> filho[maxn];



vector<int> et;



void dfs(int v)

{

    et.pb(v);

    for(int u: filho[v])

        dfs(u);

}



int deg[maxn];

int sz[maxn];

int xr[maxn];



int dp[maxn][maxai];



int f(int i, int x)

{

    if(dp[i][x]!=-1)

        return dp[i][x];

    if(i==n+1)

    {

        if(x)

            return dp[i][x]=0;

        else

            return dp[i][x]=1;

    }

    if(f(i+1, x))

        return dp[i][x]=1;

    if(sz[et[i-1]]%2==0 && f(i+sz[et[i-1]], x^xr[et[i-1]]))

        return dp[i][x]=1;

    return dp[i][x]=0;

}



vector<int> ans;



void g(int i, int x)

{

    if(i==n+1)

        return;

    if(f(i+1, x))

    {

        g(i+1, x);

        return;

    }

    ans.pb(et[i-1]);

    g(i+sz[et[i-1]], x^xr[et[i-1]]);

}



int32_t main()

{

    #ifndef ONLINE_JUDGE

    freopen("input.txt","r",stdin);

    freopen("output.txt","w",stdout); 

    #endif

    ios::sync_with_stdio(0);cin.tie(0);

 

    cin>> n;

    for(int i=1; i<=n; i++)

        cin>> a[i];

    pai[1]=1;

    for(int i=2; i<=n; i++)

    {

        cin>> pai[i];

        filho[pai[i]].pb(i);

        deg[pai[i]]++;

    }

    dfs(1);

    queue<int> q;

    for(int i=1; i<=n; i++)

    {

        if(!deg[i])

            q.push(i);

        sz[i]=1;

        xr[i]=a[i];

    }

    deg[1]++;

    while(q.size())

    {

        int v=q.front();

        q.pop();

        sz[pai[v]]+=sz[v];

        xr[pai[v]]^=xr[v];

        deg[pai[v]]--;

        if(!deg[pai[v]])

            q.push(pai[v]);

    }

    for(int i=1; i<=n+1; i++)

        for(int x=0; x<=31; x++)

            dp[i][x]=-1;

    bool ok=f(1, xr[1]);

    if(!ok)

    {

        cout<< "-1";

        return 0;

    }

    g(1, xr[1]);

    ans.pb(1);

    cout<< ans.size()<< "\n";

    for(int v: ans)

        cout<< v<< " ";

 

    return 0;

}