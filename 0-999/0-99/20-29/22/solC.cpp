#include<cstdio>

#include<cmath>

#include<cstring>

#include<queue>

#include<stack>

#include<cstdlib>

#include<iomanip>

#include<string>

#include<vector>

#include<map>

#include<set>

#include<string>

#include<iostream>

#include<algorithm>

using namespace std;

typedef long long ll;

#define INF 0x3f3f3f3f



int main()

{

    int n,m,v;

    cin>>n>>m>>v;

    if(m<n-1||m>(n-1)*(n-2)/2+1)

    {

        puts("-1");

    }

    else

    {

        int hx=0,a,b;

        for(int i=1; i<=n; i++)

        {

            

            if(i!=v)

            {

                a=i,b=v;

                hx++;

                if(hx==n-1)

                    swap(a,b);

                printf("%d %d\n",a,b);

            }

        }

        for(int i=1; i<=n; i++)

        {

            if(hx==m)

                    break;

            if(i==v||i==b)

                continue;

            for(int j=i+1; j<=n; j++)

            {

                if(hx==m)

                    break;

                if(j!=v&&j!=b)

                {

                    printf("%d %d\n",i,j);

                    hx++;

                }

                if(hx==m)

                    break;

            }

            if(hx==m)

                break;

        }

    }

}