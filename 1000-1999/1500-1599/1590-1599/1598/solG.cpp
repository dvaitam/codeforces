#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#define ll long long
#define MS 500010
#define md 2006091501
int nt[MS],le[MS];
void exkmp(char S[MS],int ls,char T[MS],int lt)
{
    nt[1]=lt;
    while(T[1+nt[2]]==T[2+nt[2]])nt[2]+=1;
    int ma=2+nt[2]-1,wz=2;
    for(int i=3,z=0;i<=lt;i++,z=0)
    {
        if(ma>=i)
        {
            z=ma-i+1;
            int t=nt[i-wz+1];
            if(t<z)z=t;
        }
        while(T[1+z]==T[i+z])z+=1;
        nt[i]=z;
        if(i+z-1>ma)ma=i+z-1,wz=i;
    }
    ma=wz=-1;
    for(int i=1,z=0;i<=ls;i++,z=0)
    {
        if(ma>=i)
        {
            z=ma-i+1;
            int t=nt[i-wz+1];
            if(t<z)z=t;
        }
        while(S[i+z]==T[1+z]&&z<ls&&z<lt)z+=1;
        le[i]=z;
        if(i+z-1>ma)ma=i+z-1,wz=i;
    }
}
char S[MS],T[MS];int ha[MS],mi[MS],ht,n;
int Hash(int l,int r)
{
    return (ha[r]-1ll*ha[l-1]*mi[r-l+1]%md+md)%md;
}
void check(int a,int b,int c)
{
    if(a<1||c>n)return;
    int h=(0ll+Hash(a,b)+Hash(b+1,c))%md;
    if(h!=ht)return;
    printf("%d %d\n%d %d",a,b,b+1,c);
    exit(0);
}
int main()
{
    scanf("%s%s",S+1,T+1);
    n=strlen(S+1);int m=strlen(T+1);
    for(int i=1;i<=m;i++)
        ht=(ht*10ll+T[i]-'0')%md;
    exkmp(S,n,T,m);mi[0]=1;
    for(int i=1;i<=n;i++)
    {
        mi[i]=10ll*mi[i-1]%md;
        ha[i]=(10ll*ha[i-1]+S[i]-'0')%md;
    }
    for(int l=1;l+m-1<=n;l++)
    {
        int r=l+m-1,c=le[l];
        if(c==m||S[l+c]>T[1+c])continue;
        int z=m-c;
        check(l,r,r+z);check(l-z,l-1,r);
        z-=1;
        if(z)check(l,r,r+z),check(l-z,l-1,r);
    }
    if(m>1)
    {
        for(int l=1;l+(m-1)*2-1<=n;l++)
            check(l,l+m-2,l+m*2-3);
    }
    return 0;
}