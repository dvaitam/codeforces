#include<bits/stdc++.h>
using namespace std;

/*-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-*/

#include <ext/pb_ds/assoc_container.hpp> 
#include <ext/pb_ds/tree_policy.hpp>
using namespace __gnu_pbds;

template<class T> using ordered_set = tree<int, null_type,less_equal<int>, rb_tree_tag,tree_order_statistics_node_update>;
/*
*p.find_by_order(3)
p.order_of_key(6)P
*/

/*-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-*/
#define ll long long 
#define int long long
#define rep(i,a,b) for(int i=a;i<b;i++)
#define red(i,a) for(int i=a;i>=0;i--)
#define pb push_back
#define pii pair<ll,ll> 
#define lb lower_bound 
//#define up upper_bound
#define endl "\n"
#define as(a) sort(a.begin(),a.end())
#define des(a) sort(a.begin(),a.end(),greater<ll>());
#define vi vector<int> 
/*-----------------------------------------*/

ll N=(int)998244353;

/*------------------------------------------*/
#define get(arr,n) rep(i,0,n) cin>>arr[i]
bool isprime(int a){for(int i=2;i*i<=a;i++)if(a%i==0)return false;return true;}
ll power(ll a,ll b){if(b==0) return 1; ll ans=1; while(b){ if(b & 1) ans=ans*a;a=a*a; b=b/2;}return ans;}
#define fast ios::sync_with_stdio(0); cin.tie(0); cout.tie(0)

#define inf (int)1e18+2
#define cont continue 
#ifndef ONLINE_JUDGE
#define debug(x) cerr<<#x<<" "; _print(x); cerr<<endl;
#else
#define debug(x)
#endif
void _print(int a){cerr<<a<<" ";}
void _print(float a){cerr<<a<<" ";}
void _print(char a){cerr<<a<<endl;}
void _print(string a){cerr<<a<<" ";}
void _print(bool a){cerr<<a<<" ";}
void _print(double a){cerr<<a<<" ";}
template<class T> void _print(vector<T> v){cerr<<"[ ";for(auto i:v){ _print(i);}cerr<<" ]";}
template<class T,class V> void _print(map<T,V> m){for(auto u:m){cerr<<u.first<<" "<<u.second<<endl;}}
template<class T,class V> void _print(multimap<T,V> m){for(auto u:m) {cerr<<u.first<<" "<<u.second<<endl;}}
template<class T> void _print(set<T> s){cerr<<"[";for(auto i:s)cerr<<i<<" ";cerr<<"]";}
template<class T> void _print(multiset<T> s){cerr<<"["; for(auto i:s)cerr<<i<<" ";cerr<<"]";}
#define gcd __gcd
#define float double 


int add(int a,int b){return (a%N+b%N)%N;}
/*------------------------------------------------------------------*/

bool sorted(vi& v){for(int i=0;i<(v.size()-1);i++){if(v[i+1]<v[i])return false;}return true;}
ll binary(ll a,ll b){ll res=1;a=a;while(b){if(b&1)res=((res%N)*(a%N))%N;a=((a%N)*(a%N))%N;res=res%N;b=b>>1;}return (int)res;}
int modinverse(int a,int b){return binary(a,b-2);}
int log_a_to_base_b(int a, int b){return log(a) / log(b);}
double pi = 2 * acos(0.0);

/*--------------------------------------------------------------*/

struct custom_hash
{
    static uint64_t splitmix64(uint64_t x) 
    {
        // http://xorshift.di.unimi.it/splitmix64.c
        x += 0x9e3779b97f4a7c15;
        x = (x ^ (x >> 30)) * 0xbf58476d1ce4e5b9;
        x = (x ^ (x >> 27)) * 0x94d049bb133111eb;
        return x ^ (x >> 31);
    }
 
    size_t operator()(uint64_t x) const {
        static const uint64_t FIXED_RANDOM = chrono::steady_clock::now().time_since_epoch().count();
        return splitmix64(x + FIXED_RANDOM);
    }
};
#define all(x) (x).begin(), (x).end()
#define mii unordered_map<int,int,custom_hash>
/*------------------------------------------------------------------------*/
/*----------------------- PRAJWAL TS -------------------------------------*/
bool left(string a,string b,int bit);
vector<int> v;

bool right(string a,string b,int bit)
{
  if(a==b)
     return true;

   int i=a.size()-1;
   while((i>=0) and (a[i]!='1'))
     i--;

   if(i==-1)
    {
       if(bit==0)
         return false;

       return left(a,b,0);
    }

   for(int j=i-1;j>=0;j--)
   {
     if(a[j]!=b[j])
     {
       int k=(i-j);
       v.push_back(k);
       string x="";
       // for(int r=0;r<k;r++)
       // {
       //   x+='0';
       // }

       for(int r=k;r<a.size();r++)
       {
         // if(x.size()==a.size()) 
         //   break;

         x+=a[r];
       }

       while(x.size()<a.size())
         x+='0';


       for(int r=0;r<a.size();r++)
       {
          if(a[r]=='0')
          {
               if(x[r]=='1')
                 a[r]='1';
          }
          else if(a[r]=='1')
          {
                    if(x[r]=='1')
                       a[r]='0';
          }
       }
     }
   }


   if(bit)
     return left(a,b,0);

   return (a==b);
}
bool left(string a,string b,int bit)
{
  if(a==b)
     return true;

   int i=0;
   while((i<a.size()) and a[i]!='1')
     i++;

   if(i==a.size())
   {
           if(bit)
             return right(a,b,0);

           return false;
    }

   for(int j=i+1;j<a.size();j++)
   {
     if(a[j]!=b[j])
     {
       int k=(j-i);
       v.push_back(-1*k);
       string x="";
       for(int r=0;r<k;r++)
       {
         x+='0';
       }

       for(int r=0;r<a.size();r++)
       {
         if(x.size()==a.size()) 
           break;

         x+=a[r];
       }

       
       for(int r=0;r<a.size();r++)
       {
          if(a[r]=='0')
          {
               if(x[r]=='1')
                 a[r]='1';
          }
          else if(a[r]=='1')
          {
                    if(x[r]=='1')
                       a[r]='0';
          }
       }
     }
   }


   if(bit)
     return right(a,b,0);

   return (a==b);
}


void solve()
{
   int n;
   cin>>n;

   string a,b;
   cin>>a>>b;

   v.clear();

   bool a1=left(a,b,1);

 //  cout<<a1<<endl;

   if(a1)
   {
     cout<<v.size()<<endl;
     for(auto x:v)
       cout<<x<<" ";

      if(v.size())
     cout<<endl;
     return;
   }

   v.clear();

    a1=right(a,b,1);

   // cout<<a1<<endl;

   if(a1)
   {
     cout<<v.size()<<endl;
     for(auto x:v)
       cout<<x<<" ";

  if(v.size())
     cout<<endl;
     return;
   }

   cout<<"-1"<<endl;
   return;


}
 
int32_t main()
{
  fast;
 #ifndef ONLINE_JUDGE
  freopen("errorf.in","w",stderr);
 #endif 

   int t=1;
   cin>>t;
   rep(i,1,t+1)
   {
     solve();
   }
  return 0;
}