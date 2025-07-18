#include<bits/stdc++.h>

using namespace std;

int main()

{

    int n,k;

    scanf("%d%d",&n,&k);

    int x=1;

    int ans;

    ans=(x+6*n-2)*k;

    printf("%d\n",ans);

    for(int i=1;i<=n;i++)

    {

        int a=x*k;

        int b=(x+2)*k;

        int c=(x+4)*k;

        int d;

        if((x+1)%3!=0)

            d=(x+1)*k;

        else

            d=(x+3)*k;

        printf("%d %d %d %d\n",a,b,c,d);

        x=x+6;

    }

}