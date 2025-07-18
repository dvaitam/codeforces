#include<bits/stdc++.h>
using namespace std;
inline void read(int &x){
    register char c=getchar();
    x=0;
    for(;c<'0'||c>'9';c=getchar());
    for(;c>='0'&&c<='9';c=getchar())
        x=(x<<1)+(x<<3)+c-'0';
}
int main(){
    int n,m,i,j,k=1e9;
    read(n);read(m);
    while(m--){
        read(i);read(j);
        k=min(k,j-i);
    }
    ++k;
    printf("%d\n",k);
    while(n--)
        printf("%d ",n%k);
}