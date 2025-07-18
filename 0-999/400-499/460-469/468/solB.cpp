#include<bits/stdc++.h>



using namespace std;

map<int,int> m;

int n,a,b,sw,j,bi[100010];



void make_set(int x,int y,int s)

{

    bi[m[x]]=bi[m[y]]=s;

    m[x]=m[y]=0;

}



int main()

{

    scanf("%d%d%d",&n,&a,&b);

    if(b<a) swap(a,b),sw=1;

    for(int i=1;i<=n;i++) scanf("%d",&j),m[j]=i;



    for(auto x:m) {

        if(x.second) {

            if(m[b-x.first]) make_set(x.first,b-x.first,1-sw);

            else if(m[a-x.first]) make_set(x.first,a-x.first,sw);

            else {

                printf("NO\n");

                return 0;

            }

        }

    }

    printf("YES\n");

    for(int i=1;i<=n;i++) {

        printf("%d ",bi[i]);

    }

    return 0;

}