#include<bits/stdc++.h>
using namespace std;
typedef long long ll;
int main()
{
    //freopen("in.txt","r",stdin);
    //freopen("out.txt","w",stdout);
    int a,b,c;
    int d,e,f;
    scanf("%d%d%d",&a,&b,&c);
    scanf("%d%d%d",&d,&e,&f);
    bool flag=0;
    int amount=d+e+f;
    if(a>d) flag=1;
    if(a+b>d+e) flag=1;
    if(a+b+c>d+e+f) flag=1;
    
    if(flag) cout<<"NO\n";
    else cout<<"YES\n";
    return 0;
}