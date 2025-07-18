#include <bits/stdc++.h>

using namespace std;

#define _ ios_base::sync_with_stdio(0); cin.tie(0); cout.tie(0);

#define ll long long

#define pb push_back

#define sz(x) (int)x.size()

#define all(x) x.begin(),x.end()

#define f first

#define s second

#define L(x) (x<<1)

#define R(x) ((x<<1)+1)

#define lsb(x) ((x)&(-x))

#define inf (int)1e9

#define linf (ll)1e17

typedef pair<int,int> ii;

typedef vector<int> vi;

const ll mod = 1e9 + 7;



int n, tam[200005], ans[200005];

ll totSum;

vi v[200005], o;



int findSize(int a,int p){ // Acha o tamanho dos sizes das Subtrees

    tam[a]=1; o.pb(a);

    

    for(auto x : v[a]){

        if(x==p) continue;

        findSize(x,a);

        totSum+=min(tam[x],n-tam[x]);

        tam[a]+=tam[x];

    }



    return tam[a];

}





int main(){_

    cin>>n;



    for(int i=1;i<n;i++){

        int a, b; cin>>a>>b;

        v[a].pb(b); v[b].pb(a);

    }

    

    findSize(1,1);



    for(int i=0;i<n;i++){

        int a=o[i], b=o[(i+(n>>1))%n]; ans[b]=a;

    }



    cout<<2*totSum<<'\n';

    for(int i=1;i<=n;i++) cout<<ans[i]<<' ';

    cout<<'\n';



    return 0;

}