#include <cstdio>

char ch[20];

int m;int a[1001];

bool dfs(int x,int s){

    if (x>m){

        puts("YES"); printf("%d",a[1]);

        for (int i=2;i<=m;i++) printf(" %d",a[i]);puts("");

        return 1;

    }

    for(int i=10;i>s;i--)

    if (i!=a[x-1]&&ch[i-1]-'0') { a[x]=i;if (dfs(x+1,i-s)) return 1;}

    return 0;

}

int main(){

    //freopen("c.in","r",stdin);

    scanf("%s", ch);

    scanf("%d", &m);

    a[0]=0;if (!dfs(1,0)) puts("NO");

}