#include<bits/stdc++.h>

using namespace std;

const int maxn=2025;

int n,t,flg;

int p[maxn];

vector<int> vec;

int getpos(int x)

{

    for(int i=1;i<=n;i++) if(p[i]==x) return i;

    return assert(false),-1;

}

void solve(int x)

{

    vec.push_back(x),reverse(p+1,p+x+1);

}

int main()

{

    scanf("%d",&t);

    while(t--)

    {

        scanf("%d",&n),flg=1,vec.clear();

        for(int i=1;i<=n;i++) scanf("%d",&p[i]),flg&=(p[i]-i)%2==0;

        if(!flg)

        {

            printf("-1\n");

            continue;

        }

        for(int k=n/2;k>=1;k--)

        {

            solve(getpos(2*k+1));

            solve(getpos(2*k)-1);

            solve(getpos(2*k)+1);

            solve(3);

            solve(2*k+1);

        }

        printf("%d\n",vec.size());

        for(auto u:vec) printf("%d ",u);

        putchar('\n');

    }

    return 0;

}