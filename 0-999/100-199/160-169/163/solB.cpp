#include<bits/stdc++.h>  

using namespace std;  

const int maxn = 1e6;  

pair<pair<int,int>,int>p[maxn];  

int ans[maxn];  

int n,k,h;  

int check(double x)  

{  

    int t = 1;  

    for (int i = 1;i<=n;i++)  

    {  

       if (1.0*t*h<=p[i].first.second*x)  

        {  

            ans[t++]=p[i].second;  

            if (t>k)  

                return 1;  

        }  

    }  

    return 0;  

}  

int main()  

{  

    scanf("%d%d%d",&n,&k,&h);  

    for (int i = 1;i<=n;i++)  

        scanf("%d",&p[i].first.first);  

    for (int i = 1;i<=n;i++)  

        scanf("%d",&p[i].first.second);  

    for (int i = 1;i<=n;i++)  

        p[i].second=i;  

    sort(p+1,p+1+n);  

    double l = 0.0;  

    double r = 1000000000.0;  

    for (int i = 1;i<=100;i++)  

    {  

        double m = (l+r)/2.0;  

        if (check(m))  

            r=m;  

        else  

            l=m;  

    }  

    check(r);  

    for (int i = 1;i<=k;i++)  

        printf("%d ",ans[i]);  

    printf("\n");  

}