#include <bits/stdc++.h>
#include <ext/pb_ds/assoc_container.hpp>
#define ll                  long long
#define pb                  push_back
#define pii                 pair<ll,ll>
#define vi                  vector<ll>
#define loop(i,a,b)         for(ll i=(a);i<=(b);i++)
#define looprev(i,a,b)      for(ll i=(a);i>=(b);i--)
#define all(v)              v.begin(),v.end()
#define ff                  first 
#define ss                  second
#define boost               ios_base::sync_with_stdio(false);cin.tie(NULL);cout.tie(NULL);
#define log(args...)        {string _s=#args;replace(_s.begin(),_s.end(),',',' ');stringstream _ss(_s);istream_iterator<string> _it(_ss);err(_it,args);}
#define logarr(arr,a,b)     for(int z=(a);z<=(b);z++) cerr<<(arr[z])<<" ";cerr<<endl;
using namespace std;
using namespace __gnu_pbds;
void err(istream_iterator<string> it){}
template<typename T,typename... Args> 
void err(istream_iterator<string> it,T a,Args... args){
    cerr<<*it<<" = "<<a<<endl;
    err(++it,args...);
}
template<class T> void debug(set<T> S){cerr<<"[ ";for(auto i:S) cerr<<i<<" ";cerr<<"] "<<endl;}
template<class T> void debug(multiset<T> S){cerr<<"[ ";for(auto i:S) cerr<<i<<" ";cerr<<"] "<<endl;}
template<class T> void debug(unordered_set<T> S){cerr<<"[ ";for(auto i:S) cerr<<i<<" ";cerr<<"] "<<endl;}
template<class T,class X> void debug(T *arr,X s,X e){cerr<<"[ ";loop(i,s,e) cerr<<arr[i]<<" ";cerr<<"] "<<endl;}
template<class T,class V> void debug(pair<T,V> p){cerr<<"{";cerr<<p.ff<<" "<<p.ss<<"}"<<endl;}
template<class T,class V> void debug(vector<pair<T,V>> v){cerr<<"[ "<<endl;for(auto i:v) debug(i);cerr<<"]"<<endl;}
template<class T> void debug(vector<T> v){cerr<<"[ ";for(auto i:v) cerr<<i<<" ";cerr<<"] "<<endl;}
template<class T> void debug(vector<vector<T>> v){cerr<<"[ "<<endl;for(auto i:v) debug(i);cerr<<"] "<<endl;}
typedef tree<int, null_type, less<int>,
            rb_tree_tag, tree_order_statistics_node_update> ordered_set;
vi primes(){
    int N=1e7;
    vi isPrime(N+5,1);
    isPrime[0]=0;
    isPrime[1]=0;
    loop(i,2,N){
        if(isPrime[i]==0) continue;
        int z=2*i;
        while(z<N){
            isPrime[z]=0;
            z+=i;
        }
    }
    return isPrime;
}
vi *getfactors(){
    int N=1e6;
    vi *factors=new vi[N+5];
    for(int i=2;i<=N;i++){
        for(int j=i;j<=N;j+=i){
            factors[j].pb(i);
        }
    }
    // time complexity O(nlogn)
    return factors;
}
const int mod=1000000007;
ll *getfactorial(){
    int m=2e5;
    ll *fact=new ll[m+1];
    fact[0]=1;
    loop(i,1,m){
        fact[i]=(fact[i-1]*i)%mod;
    }
    // time complexity O(n)
    return fact;
}
ll *fac;
ll power(ll a,ll b,ll m=mod){
    if(b==0) return 1;
    if(b==1) return a;
    ll ans=power(a,b/2,m);
    ans=(1ll*ans*ans)%m;
    if(b%2==1){
        ans=(1ll*ans*a)%m;
    }
    return ans;
}
ll nCr(ll n,ll r,int m=mod){
    if(n<r) return 0; 
    ll ans=fac[n];
    ans=(ans*power(fac[n-r],m-2))%m;
    ans=(ans*power(fac[r],m-2))%m;
    return ans;
}
ll gcd(ll a,ll b){
    if(b==0) return a;
    if(a==0) return b;
    return gcd(b,a%b);
}
struct dsu{
    vi par,sz;
    dsu(int n){
        par.resize(n+1);
        sz.resize(n+1);
        loop(i,0,n) par[i]=i;
        loop(i,0,n) sz[i]=0;
    }
    int find(int x){
        if(x==par[x]) return x;
        return par[x]=find(par[x]);
    }
    int comb(int a,int b){
        if(sz[a]<sz[b]) swap(a,b);
        int x=find(a);
        int y=find(b);
        if(x==y) return 0;
        sz[x]+=sz[y];
        par[y]=x;
        return 1;
    }
};

void test_case(){    
    int n,k,x;
    cin>>n>>k>>x;
    if(x!=1){
        cout<<"YES"<<endl;
        cout<<n<<endl;
        loop(i,1,n) cout<<1<<" ";
        cout<<endl;
    }
    else{
        if(k==1){
            cout<<"NO"<<endl;
        }
        else{
            if(n%2==0){
                cout<<"YES"<<endl;
                cout<<n/2<<endl;
                loop(i,1,n/2){
                    cout<<2<<" ";
                }
                cout<<endl;
            }
            else{
                if(k==2){
                    cout<<"NO"<<endl;
                }
                else{
                    if(n==1){
                        cout<<"NO"<<endl;
                        return ;
                    }
                    cout<<"YES"<<endl;
                    cout<<n/2<<endl;
                    loop(i,1,n/2-1){
                        cout<<2<<" ";
                    }
                    cout<<3<<endl;
                }
            }
        }
    }
}  
int main(){
    boost
    int t=1;
    cin>>t;
    while(t--){
        test_case();
    }  

}