#include <bits/stdc++.h>
using namespace std;
long long now,savenow,n,l,savel;
int main()
{
    scanf("%lld",&n);
    long long a;
    scanf("%lld",&a);
    savel=l=1; savenow=now=a;
    for(long long i=1;i<n;i++){
        scanf("%lld",&a);
        if(a>now/l){
            now=a;
            l=1;
            if(a>(savenow/savel)){
                savenow=a;
                savel=1;
            }
        }
        else if((now+a)/(l+1)==(now/l)){
            now+=a;
            l++;
            if((now/l)>=(savenow/savel)&&l>=savel){
                savel=l;
                savenow=now;
            }
        }
        else{
            now=a;
            l=1;
        }
    }
    printf("%lld",savel);
    return 0;
}