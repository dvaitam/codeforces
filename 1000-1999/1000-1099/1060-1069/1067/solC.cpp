#include <bits/stdc++.h>
#pragma GCC optimize("O3")
#ifdef ONLINE_JUDGE
#pragma GCC target("sse,sse2,sse3,ssse3,sse4,popcnt,abm,mmx,avx,tune=native")
#endif
#pragma GCC optimize("unroll-loops")
#pragma comment(linker, "/stack:200000000")
#include <bits/stdc++.h>
#include <deque>
#include <type_traits>
#include <string.h>
#include <stdlib.h>
#include <stdio.h>
#include <time.h>
using namespace std;
#include <ext/pb_ds/assoc_container.hpp>
#include <ext/pb_ds/tree_policy.hpp>
#include <ext/pb_ds/detail/standard_policies.hpp>
template<typename Key=int,typename Mapped=__gnu_pbds::null_type,typename Cmp_Fn=std::less<Key>,typename Tag=__gnu_pbds::rb_tree_tag,template<typename Const_Node_Iterator,typename Node_Iterator,typename Cmp_Fn_,typename Allocator_>class Node_Update=__gnu_pbds::tree_order_statistics_node_update,typename Allocator=std::allocator<char>>class ordered_set_t:__gnu_pbds::tree<Key,Mapped,Cmp_Fn,Tag,Node_Update,Allocator>{};const int GRANDOM=std::chrono::high_resolution_clock::now().time_since_epoch().count();struct ghash{int operator()(int x){return std::hash<int>{}(x^GRANDOM);}};template<typename KeyType>class hash_table_t:__gnu_pbds::gp_hash_table<KeyType,int,ghash>{};namespace vh{template<class T>void _R(T&x){std::cin>>x;}inline void _R(int&x){scanf("%d",&x);}inline void _R(int64_t&x){scanf("%" SCNd64,&x);}inline void _R(double&x){scanf("%lf",&x);}inline void _R(char&x){scanf(" %c",&x);}inline void _R(char*x){scanf("%s",x);}inline void R(){}template<class T,class... U>inline void R(T&head,U&... tail){_R(head);R(tail...);}template<class T>void _W(const T&x){cout<<x;}inline void _W(const int&x){printf("%d",x);}inline void _W(const int64_t&x){printf("%" PRId64,x);}inline void _W(const double&x){printf("%.16f",x);}inline void _W(const char&x){putchar(x);}inline void _W(const char*x){printf("%s",x);}template<class T>inline void _W(const vector<T>&x){for(auto i=x.begin();i!=x.end();_W(*i++))if(i!=x.cbegin())putchar(' ');}inline void W(){}template<class T,class... U>inline void W(const T&head,const U&... tail){_W(head);if(sizeof...(tail))putchar(' '),W(tail...);}template<class T,class... U>inline void WL(const T&head,const U&... tail){_W(head);putchar(sizeof...(tail)? ' ':'\n');W(tail...);}};
#define prec setprecision
namespace vh{typedef long long ll;typedef unsigned long long llu;template<typename T>T gcd(T m,T n){while(n){T t=m%n;m=n;n=t;};return m;}template<typename T>T exgcd(T a,T b,T&sa,T&ta){T q,r,sb=0,tb=1,sc,tc;sa=1,ta=0;if(b)do q=a/b,r=a-q*b,a=b,b=r,sc=sa-q*sb,sa=sb,sb=sc,tc=ta-q*tb,ta=tb,tb=tc;while(b);return a;}template<typename T>T mul_inv(T a,T b){T t1=a,t2=b,t3;T v1=1,v2=0,v3;T x;while(t2!=1)x=t1/t2,t3=t1-x*t2,v3=v1-x*v2,t1=t2,t2=t3,v1=v2,v2=v3;return(v2+b)%b;}template<typename T>T powmod(T a,T b,T MOD){if(b<0)return 0;T rv=1;while(b)(b%2)&&(rv=(rv*a)%MOD),a=a*a%MOD,b/=2;return rv;}template<typename T>inline T isqrt(T k){T r=sqrt((double)k)+1;while(r*r>k)r--;return r;}template<typename T>inline T icbrt(T k){T r=cbrt((double)k)+1;while(r*r*r>k)r--;return r;}template<typename T>bool mul_overflow(T&r,T a,T b){return __builtin_mul_overflow(a,b,&r);}template<ll n>struct BitSize{enum{Size=BitSize<n/2>::Size+1};};template<>struct BitSize<0>{enum{Size=1};};template<>struct BitSize<1>{enum{Size=1};};
#define BITSIZE(n) (BitSize<n>::Size)
#define BITMAX(n) (BitSize<n>::Size - 1)
#define DEBUG !defined(ONLINE_JUDGE)
#define SZ(x) ((int)(x).size())
#define ALL(x) begin(x), end(x)
template<typename TH>void _dbg(const char*sdbg,TH h){cerr<<sdbg<<"="<<h<<"\n";}template<typename TH,typename... TA>void _dbg(const char*sdbg,TH h,TA... t){while(*sdbg!=',')cerr<<*sdbg++;cerr<<"="<<h<<",";_dbg(sdbg+1,t...);}
#ifdef DEBUG
#define debug(...) _dbg(#__VA_ARGS__, __VA_ARGS__)
#define debugv(x) {{cerr <<#x <<" = "; for(auto itt: x) cerr <<*itt <<", "; cerr <<"\n"; }}
#else
#define debug(...) (__VA_ARGS__)
#define debugv(x)
#define cerr if(0)cout
#endif
}namespace vh{template<typename T>int min_rotation(vector<T>s){int a=0,N=s.size();for(int b=0;b<N;b++)for(int i=0;i<N;i++){if(a+i==b||s[(a+i)%N]<s[(b+i)%N]){b+=max(0,i-1);break;}if(s[(a+i)%N]>s[(b+i)%N]){a=b;break;}}return a;}}namespace vh{};
namespace vh{
  int main(){
    srand(GRANDOM);
    int n;cin>>n;
    int x=0;
    for(int i=0;i<n;i++){
      if(i%3==0)x++,cout<<x<<' '<<0<<'\n';
      else if(i%3==1)x++,cout<<x<<' '<<0<<'\n';
      else if(i%3==2)cout<<x<<' '<<3<<'\n';
    }
    return 0;
  }
};
int main(int argc,char*argv[]){
  std::cin.sync_with_stdio(false);std::cin.tie(nullptr);
  return vh::main();
}