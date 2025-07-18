//#pragma GCC optimize("O3,unroll-loops")

//#pragma GCC target("avx2,bmi,bmi2,lzcnt,popcnt")

#include <bits/stdc++.h>

//#include <ext/pb_ds/assoc_container.hpp>

//using namespace __gnu_pbds;

using namespace std;



#define int long long 

//typedef tree<int, null_type, less<int>, rb_tree_tag, tree_order_statistics_node_update> indexed_set;

//indexed_set s;

//#define pii pair<int,int>

//#define ll long long  

//#define int unsigned int

//#define pb push_back

//#define mp make_pair

#define ff first

#define ss second

#define countbit(x) __builtin_popcountll( x )

//#define lb lower_bound

//#define ub upper_bound

//#define bs binary_search

#define mod 1000000007 // mod+2

#define double long double

#define all(x) x.begin(), x.end()

#define debug(x) cout << #x << " = " << x << "\n"

const int inf=1e18;

const int N=100005;

//int d;

vector <int> fac(N);

vector <int> adj[N];  

//vector <int> vis(N);

//vector <vector <int>> v(N);



void google(int t){

    cout<<"Case #"<<t<<": ";

}

//mt19937 rng(chrono::steady_clock::now().time_since_epoch().count());

int gcd(int a, int b) {if (b > a) {return gcd(b, a);} if (b == 0) {return a;} return gcd(b, a % b);}



vector<int> dx = {1, -1, 0, 0}, dy = {0, 0, 1, -1};

/*

void dfs(int s){

    if(vis[s])return;

    vis[s]=1;

    for(auto u:adj[s]){

        dfs(u);

    }

}

int ncr(int n, int r){

    int x=(fac[r]*fac[n-r]);

    x=inv(fac[n],x);

    return x;

}

*/

int fun(int a,int b){

    if(b==0)return 1;

    else if(a==1)return 1;

    int temp=fun(a,b/2)%mod;

    temp=(temp*temp)%mod;

     if(b%2==1)temp=(temp*a)%mod;

    return temp%mod;

}

 

int inv(int a,int b){

    return (a*fun(b,mod-2))%mod;

}



vector <int> p(N,0), done(N,0);

int cnt=0;



void dfs(int s, int par){

    int c=-1;

    for(auto u:adj[s]){

        if(u==par)continue;

        c=u;

        dfs(u,s);

    }

    if(!done[s]&&par){

        swap(p[s],p[par]);

        done[s]=done[par]=1;

        cnt+=2;

    }

    if(!done[s]){

        swap(p[s],p[c]);

        done[s]=done[c]=1;

        cnt+=2;

    }

}



void solve(){

    int i,j,k=0,n,m,l=0;

    

    cin>>n;

    /*

    for(i=0;i<=n;i++)vis[i]=0;

    for(i=0;i<=n;i++){

         adj[i].clear();

    }

    */

    /*string s;

    cin>>s;

    n=s.length();*/

    for(i=0;i<n-1;i++){

        int x,y;

        cin>>x>>y;

        adj[x].push_back(y);

        adj[y].push_back(x);

    }

    

    for(i=0;i<=n;i++)p[i]=i;

    

    dfs(1,0);

    cout<<cnt<<"\n";

    for(i=1;i<=n;i++)cout<<p[i]<<" ";cout<<"\n";

    //memset(a,0,sizeof(a));

    //cout<<fixed<<setprecision(10)<<ans<<"\n";

       

}



 

signed main()

{   /*  

   #ifndef ONLINE_JUDGE 

   freopen("inputhsr.txt", "r", stdin);  

   freopen("outputhsr.txt", "w", stdout); 

   #endif 

   */

    ios_base::sync_with_stdio(false);

    cin.tie(NULL);

    cout.tie(NULL);

    

    /*

    fac[0]=1;

    for(i=1;i<N;i++){

        fac[i]=(i*fac[i-1])%mod;

    }

    */ 

    

    int t=1;

    int T=1;

    //cin>>T;

    for(t=1;t<=T;t++){

        

       

        //google(t);

        solve();

        

        

    }

    

    return 0;

}