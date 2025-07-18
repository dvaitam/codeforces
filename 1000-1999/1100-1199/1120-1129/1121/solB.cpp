#include<stdio.h>
#include<string.h>
#include<algorithm>
using namespace std;
const int N=1050;
const int MAXN=200050;
int n;
int a[N];
int cnt[MAXN];
int main(){
    int t,ans=0;
    scanf("%d",&n);
    for(int i=1;i<=n;i++)scanf("%d",&a[i]);
    for(int i=1;i<=n;i++)
        for(int j=i+1;j<=n;j++){
            ++cnt[a[i]+a[j]];
        }
    for(int i=1;i<=MAXN;i++)ans=max(ans,cnt[i]);
    printf("%d",ans);
    return 0;
}