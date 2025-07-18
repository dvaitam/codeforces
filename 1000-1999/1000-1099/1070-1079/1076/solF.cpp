/*********************** 
 c++ file 
***********************/

#include<iostream> 
#include<cstdio> 
#include<vector> 
#include<map> 
#include<algorithm> 
#include<string> 
#include<cstring> 
#include<set> 
#include<cmath> 
#include<queue> 
#include<fstream> 
#include<ctime> 
#include<iomanip> 
#include<bitset> 
#include<cstdlib> 
#include<deque> 
#include<list> 
#include<stack> 
#include<utility> 
 
#define ll long long 
#define inf 0x3f3f3f3f 
#define lb(x) (x&(-x)) 
#define ls (rt<<1) 
#define rs (rt<<1|1) 
#define ms(x) memset(x,0,sizeof(x)) 
#define msinf(x) memset(x,0x3f,sizeof(x)) 
#define IO ios::sync_with_stdio(false),cin.tie(0),cout.tie(0) 
#define forn(i,a,b) for(int i=a;i<b;++i) 
 
using namespace std; 
 
inline char get(void) 
{ 
    static char buf[100000],*S=buf,*T=buf; 
    if(S==T) 
    { 
        T= (S=buf) + fread(buf,1,100000,stdin); 
        if(T==S) return EOF; 
    } 
    return *S++; 
} 
inline void read(int & x) 
{ 
    static char c; x=0; 
    for(c=get();c<'0'||c>'9';c=get()); 
    for(;c>='0'&&c<='9';c=get()) 
        x = x*10+c-'0'; 
} 
const int M = 1000*1000+10; 
int x[M];
int y[M];
int n,k;
 
bool cal(int & u,int & v,int a,int b)
{
    if(1ll*(b+1)*k-u<a) return false;
    if(1ll*(a+1)*k-v<b) return false;
        if(1ll*b*k-u<a)
            u = a-b*k+u;
        else
            u=0;
        if(1ll*a*k-v<b)
            v = b-a*k+v;
        else
            v=0;
        return true;
}
bool solve()
{
    int u=0, v =0;
    int a,b;
    forn(i,0,n)
    {
        a=x[i];
        b=y[i];
        if(!cal(u,v,a,b))
            return false;
    }
    return true;
}
int main() 
{ 
    read(n),read(k);
    forn(i,0,n)
        read(x[i]);
   forn(i,0,n) read(y[i]);
    if(solve())
        cout<<"YES"<<endl;
    else
        cout<<"NO"<<endl;
    return 0; 
}