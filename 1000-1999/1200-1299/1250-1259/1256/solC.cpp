#include<bits/stdc++.h>
using namespace std;

#define MOD           1000000007/*998244353*/
#define pi            acos(-1)
#define int           long long 
#define D             double
#define S             second
#define F             first
#define pb            push_back
#define ff            fflush(stdout)
#define ppb           pop_back
#define B(c)          (c).begin()
#define E(c)          (c).end()
#define all(c)        (c).begin(),(c).end()
#define rall(c)       (c).rbegin(),(c).rend() 
#define lb            lower_bound
#define ub            upper_bound
#define si(c)         (int)((c).size())
#define L(c)           c[si(c)-1]
#define gcd(a,b)      __gcd(a,b)
#define lcm(a,b)      (a*(b/gcd(a,b)))  
#define accuracy      cout << fixed << setprecision(18);
#define inf           (int)2e18
#define pow(i,n)      (int)pow((int)i,n)
#define err           cerr<<"move"<<'\n';
#define print         cout<<"move"<<'\n';
#define en            '\n'

typedef vector<int>                           vi;
typedef pair<int,int>                         pii;
typedef vector<pii>                           vpi;
typedef pair<int,pii>                         pipii;
typedef vector<vector<int> >                  vvi;
typedef map<int,int>                          mp;
typedef map<string,int>                       msi;
typedef priority_queue<pii, vector<pii>, greater<pii> > pq;


void solve(){
     int n,d,m;
     cin>>n>>m>>d;
     int c[m],sum=0;
     for(int i=0;i<m;i++){
          cin>>c[i];
          sum+=c[i];
     }
     d--;

      if(sum+(m+1)*d>=n){
          cout<<"YES"<<endl;
          int k=(d>0?(n-sum)/d:0),res=(d>0?(n-sum)%d:0),cnt=0,id=0;
          //cout<<k<< " "<<res<<endl;
          bool flag=1;
          for(int i=1;i<=n;i++){
               int j;
               if(cnt==k&&res){
                    k++,d=res;
                    res=0;
               }
               if(cnt<k&&flag){
                    for(j=i;j<i+d&&j<=n;j++){
                         cout<<0<<" ";
                    }
                    cnt++;
               }
               else{
                    for(j=i;j<i+c[id];j++){
                         cout<<id+1<<" ";
                    }
                    id++;
               }
               flag=!flag;
               i=j-1;
          }
          cout<<endl;
     }
     else cout<<"NO"<<endl;

}

 
signed main(){
    ios_base::sync_with_stdio(0);cin.tie(0);cout.tie(0);    
    int t=1;
    //accuracy;
    #ifndef ONLINE_JUDGE
    freopen("input.txt", "r", stdin); 
    freopen("output.txt", "w", stdout);
    freopen("error.txt","w", stderr);
    #endif
    //cin>>t;

    while(t--){
        solve();
    }
}