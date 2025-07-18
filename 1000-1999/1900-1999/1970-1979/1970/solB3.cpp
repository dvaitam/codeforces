#include<bits/stdc++.h>
#define ll long long
#define fs first
#define sn second
#define N 202400
using namespace std;
int n,x,v[N],anss[N];
pair<int,int> ans[N],a[N];
int main()
{
    scanf("%d",&n);
    for(int i=1;i<=n;++i){
        scanf("%d",&a[i].fs);
        a[i].sn=i;
    }
    sort(a+1,a+1+n);
    x=0;
    for(int i=1;i<n;++i)
        if(a[i].fs==a[i+1].fs){
            x=i;
            break;
        }
    if(x||a[1].fs==0){
        if(x){
            swap(a[x],a[1]);
            swap(a[x+1],a[2]);
        }
        ans[a[1].sn]={1,1};
        v[1]=1;
        if(a[1].fs==0)anss[a[1].sn]=a[1].sn;
        else anss[a[1].sn]=a[2].sn;
        for(int i=2;i<=n;++i){
            if(a[i].fs==0){
                ans[a[i].sn]={i,1};
                anss[a[i].sn]=a[i].sn;
                v[i]=1;
            }
            else if(a[i].fs<i){
                ans[a[i].sn]={i,v[i-a[i].fs]};
                v[i]=v[i-a[i].fs];
                anss[a[i].sn]=a[i-a[i].fs].sn;
            }
            else{
                ans[a[i].sn]={i,a[i].fs-i+2};
                v[i]=a[i].fs-i+2;
                anss[a[i].sn]=a[1].sn;
            }
        }
    }
    else{
        if(n==2){
            puts("NO");
            return 0;
        }
        ans[a[n].sn]={1,1};
        anss[a[n].sn]=a[n-1].sn;
        ans[a[n-1].sn]={n,2};
        anss[a[n-1].sn]=a[1].sn;
        for(int i=1;i<=n-2;++i){
            ans[a[i].sn]={i+1,1};
            anss[a[i].sn]=a[n].sn;
        }
    }
    puts("YES");
    for(int i=1;i<=n;++i)
        printf("%d %d\n",ans[i].fs,ans[i].sn);
    for(int i=1;i<=n;++i)
        printf("%d ",anss[i]);
    putchar(10);
    return 0;
}