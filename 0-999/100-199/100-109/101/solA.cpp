#include <cstdio>
#include <cstring>

char s[100005];
int n,k,f[200];
int p[200];

int main()
{
    scanf("%s",s);
    scanf("%d",&k);
    n=strlen(s);
    memset(f,0,sizeof(f));
    for(int i=0;i<n;i++)
        f[s[i]]++;
    memset(p,0,sizeof(p));
    int m=n-k;
    int ans=0;
    for(int i=0;i<200 &&m>0;i++)
    {
        int x=0;
        for(int i=0;i<200;i++)
            if(f[i]>f[x])x=i;
        p[x]=f[x];
        m-=p[x];
        ans++;
        f[x]=0;
    }
    printf("%d\n",ans);
    int x=0;
    m=n-k;
    for(int i=0;i<n;i++)
        if(p[s[i]])
        {
            putchar(s[i]);
        }
    puts("");
    return 0;
}