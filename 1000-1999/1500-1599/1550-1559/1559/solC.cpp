#include<bits/stdc++.h>

using namespace std;

#define ll long long

#define ld long double

#define pb push_back

#define sz(x) (int)(x).size()

#define all(x) (x).begin(),(x).end()

#define rall(x) (x).rbegin(),(x).rend()

#define PI 2*acos(0)

#define torad(x) (x)*((2*acos(0))/180.0)

#define eps 1e-15

#define PQ priority_queue

#define ff first

#define ss second

const int mod=1e9+7;

const ll INF=1e14;

const int MAXN=2e5+5;

void solve()

{

    int n,x=0;

    cin>>n;

    for(int i=1;i<=n;i++){

        int a;

        cin>>a;

        if(!a)

            x=i;

    }

    if(x==0)

        cout<<n+1<<" ";

    for(int i=1;i<=n;i++){

        cout<<i<<" ";

        if(i==x)

            cout<<n+1<<" ";

    }

    cout<<endl;

}

int main()

{

    #ifndef ONLINE_JUDGE

    freopen("input.txt","r",stdin);

    freopen("output.txt","w",stdout);

    #endif

    ios::sync_with_stdio(false);

    cin.tie(nullptr);

    int t;

    cin>>t;

    while(t--){

        solve();

    }

    return 0;

}