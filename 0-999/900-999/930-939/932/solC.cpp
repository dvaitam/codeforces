#include<bits/stdc++.h>
#define RG register
using namespace std;
inline int gi() {
    int x=0,o=1;char ch=getchar();
    while((ch<'0'||ch>'9')&&ch!='-') ch=getchar();
    if(ch=='-') o=-1,ch=getchar();
    while(ch>='0'&&ch<='9') x=x*10+ch-'0',ch=getchar();
    return x*o;
}
void put(int x)  
{  
    int num = 0; char c[15];
    while(x) c[++num] = (x%10)+48, x /= 10;
    while(num) putchar(c[num--]);
    putchar(' '); 
}
int n,a,b,x=-1,y=-1;
int main() {
    cin>>n>>a>>b;
    for(int i=0;i<=n/a;i++) {
        x=i;
        double ans=1.0*(n-a*x)/b;
        if(ans!=(int)ans||ans<0) continue;
        y=ans;break;
    }
    if(y==-1) { puts("-1");return 0; }
    //int gg=exgcd(a,b,x,y);
    //if(n%gg!=0) { puts("-1");return 0; }
    //printf("%d %d\n",x,y);
    RG int i,cnt;
    for(i=1,cnt=0;cnt<x;i+=a,cnt++) {
        RG int l=i,r=i+a-1;
        for(RG int j=l+1;j<=r;j++) put(j);
        put(l);
    }
    for(cnt=0;cnt<y;i+=b,cnt++) {
        RG int l=i,r=i+b-1;
        for(RG int j=l+1;j<=r;j++) put(j);
        put(l);
    }    
    return 0;
}