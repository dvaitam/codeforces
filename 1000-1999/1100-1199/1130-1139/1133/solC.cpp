#include <bits/stdc++.h>
using namespace std;
using namespace __gnu_cxx;


#define IOC std::ios::sync_with_stdio(false)
typedef long long LL;
typedef unsigned long long uLL;
typedef unsigned int uit;
typedef long double ld;
typedef pair<int,int> pii;
const int INF = 0x3f3f3f3f;
const double enf = 1e-6;
const double pi = acos(-1.0);
inline LL read()
{
    LL X=0,w=0; char ch=0;
    while(!isdigit(ch)) {w|=ch=='-';ch=getchar();}
    while(isdigit(ch)) X=(X<<3)+(X<<1)+(ch^48),ch=getchar();
    return w?-X:X;
}
inline void write(LL x)
{
     if(x<0) putchar('-'),x=-x;
     if(x>9) write(x/10);
     putchar(x%10+'0');
}
//-------head------------


const int maxn = 2e6 + 10;
LL dict[maxn];
int main(){

    int n = read();
    for(int i = 0; i < n; ++i){
        dict[i] = read();
    }
    sort(dict,dict + n);
    LL who = 0,out = 1;
    LL mm = dict[who];
    for(int i = 1; i < n; ++i){
        while(dict[i] > mm + 5){
            mm = dict[++who];
        }
        out = max(out,i - who + 1);
    }
    cout << out << endl;
    return 0;
}