#include <bits/stdc++.h>

#define sys ios_base::sync_with_stdio(0);cin.tie(0);

#define mod 1000000007

using namespace std;

//#pragma comment(linker, "/STACK:268435456");

#define count_setbits(n) __builtin_popcount(n)

#define fixed cout<<fixed<<setprecision(16)

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

#define pb push_back

#define ff first

#define ss second

#define print(x) for(auto i:x) cout<<i<<" ";

#define pprint(x) for(auto [i,j]:x) cout<<i<<" "<<j

#define dbg1(x) cout << #x << "= " << x << endl;

#define dbg2(x,y) cout << #x << "= " << x << "	" << #y << "= " << y <<endl;

#define dbg3(x,y,z) cout << #x << "= " << x << "	" << #y << "= " << y << "	" << #z << "= " << z << endl;

#define dbg4(x,y,z,w) cout << #x << "= " << x << "	" << #y << "= " << y << "	" << #z << "= " << z << "	" << #w << "= " << w << endl;

//const ll N=1e5;

//vector<bool> isprime(N+1,1);

//inline void sieve(){ isprime[0]=isprime[1]=1; for(int i=2;i<=N;i++) if(isprime[i]) for(int j=i*i;j<=N;j+=i) isprime[j]=false;}

inline vector<ll> get_factors(ll n) { vector<ll> factors; if(n==0) return factors; for(ll i=1;i*i<=n;++i){ if(n%i==0){factors.push_back(i); 

if(i!=n/i) factors.push_back(n/i);}}  return factors;}

inline ll modpower(ll a, ll b, ll m = mod)

{ ll ans = 1; while (b) { if (b & 1) ans = (ans * a) % m; a = (a * a) % m; b >>= 1; } return ans; }

inline ll mod_inverse(ll x,ll y) { return modpower(x,y-2,y); }

inline ll max_ele(ll * arr ,ll n) { ll max=*max_element(arr,arr+n); return max;}

inline ll max_i(ll * arr,ll n){ return max_element(arr,arr+n)-arr+1; }



inline void solve()

{

    ll n ; cin>>n;

    str s; cin>>s;

    str k="abacaba";

    ssize_t p=s.find(k);

    if(p!=str::npos){

        p=s.find(k,p+1);

        if(p!=str::npos){

            cno;

            return;

        }else{

                for(auto &i:s) if(i=='?') i='z';

                cys;

                c(s)<<endl;

                return;

        }

    }

    str opt=s;

    bool check=1;

    f(i,0,n-k.size() + 1){

        s=opt;

        str op=s.substr(i,k.size());

        check=1;

        f(j,0,k.size()){

            if(op[j]!='?'){

                if(k[j]!=op[j]) {

                    check=0;

                    goto there;

                }

            }

        }

        there:{}

        if(check){

            f(j,0,k.size()){

                s[i+j]=k[j];

            }

            if((p=s.find(k))!=str::npos and s.find(k,p+1)==str::npos){

                goto end;

            }

        }

    }

    end:{};

    f(i,0,n) if(s[i]=='?') s[i]='z';

    

    p=s.find(k);

    if(p!=str::npos){

        p=s.find(k,p+1);

        if(p!=str::npos){

            cno;

            return;

        }else{

            cys;

            c(s)<<endl;

            return;

        }

        }else{

        cno;

        return;

    }

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