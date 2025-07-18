#include <bits/stdc++.h>

#define sys ios_base::sync_with_stdio(0);cin.tie(0);

#define mod 1000000007

using namespace std;

//#pragma comment(linker, "/STACK:268435456");

#define count_setbits(n) __builtin_popcount(n)

#define count_bits(n)  ((ll)log2(n))+1

#define no_of_digits(n) ((ll)log10(n))+1 

#define str string 

#define c(itretor) cout<<itretor

#define cp(itretor) cout<<setprecision(itretor)

#define cys  cout<<"YES"<<endl

#define cno  cout<<"NO"<<endl

#define endl "\n"

#define imx INT_MAX

#define imn INT_MIN

#define lmx LLONG_MAX

#define lmn LLONG_MIN

#define ll long long 

#define f(i,l,r) for(long long i=l;i<r;++i) 

#define fr(i,r,l) for(long long i=r-1;i>=l;--i) 

#define vi vector<int> 

#define vs vector<string> 

#define vll vector <long long> 

#define mii map<int,int> 

#define mll map<long long,long long> 

#define tsolve long long t;cin>>t; while(t--) solve();

#define inp(x) for(auto &i:x) cin>>i;

#define all(x) x.begin(),x.end()

#define print(x) for(auto i:x) cout<<i<<" ";

#define pprint(x) for(auto [i,j]:x) cout<<i<<" "<<j

ll modpower(ll a, ll b, ll m = mod)

{ ll ans = 1; while (b) { if (b & 1) ans = (ans * a) % m; a = (a * a) % m; b >>= 1; } return ans; }

ll mod_inverse(ll x,ll y) { return modpower(x,y-2,y); }

ll max_ele(ll * arr ,ll n) { ll max=*max_element(arr,arr+n); return max;}

ll max_i(ll * arr,ll n){ return max_element(arr,arr+n)-arr+1; }

void solve()

{

    ll n1,n2,n3; cin>>n1>>n2>>n3;

    str s="",t="1";

    bool c=0;

    swap(n1,n3);

    if(n2==0){

        if(n1){

            s+='1';

            while(n1--) s+='1';

        } else{

            s+='0';

            while(n3--) s+='0';

        }

    }

    else if(n2&1){

    while(n2--){

        c?(t+='1',c=0):(t+='0',c=1);

    }

    while(n1--) s+='1';

    s+=t;

    while(n3--) s+='0';

    }else{

        --n2;

        while(n2--) c?(t+='1',c=0):(t+='0',c=1);

        while(n3--) t+='0';

        t+='1';

        while(n1--) s+='1';

        s+=t;

    }

    c(s)<<endl;

}

int main()

{

sys;

#ifndef ONLINE_JUDGE

freopen("input.txt","r",stdin);

freopen("output.txt","w",stdout);

#endif

tsolve;

return 0;

}