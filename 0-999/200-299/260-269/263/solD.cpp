#include<bits/stdc++.h>

using namespace std;

void read(int &x){

    register int c=getchar();

    x=0;

    for(;(c<48||c>57);c=getchar());

    for(;c>47&&c<58;c=getchar())

        x=(x<<1)+(x<<3)+c-48;

}

int a[200005],link[200005],head[100005],c[100005],n,i,j,k=1,u,v;

bool b[100005];

int main(){

    read(n);

    for(i=1;i<=n;i++)

        head[i]=b[i]=0;

    read(i);read(j);

    j=i<<1;

    while(i--){

        read(u);read(v);

        a[j]=v;

        link[j]=head[u];

        head[u]=j--;

        a[j]=u;

        link[j]=head[v];

        head[v]=j--;

    }

    j=c[1]=b[1]=1;

    do{

        i=head[j];

        while(i){

            if(!b[a[i]]){

                b[j=a[i]]=1;

                c[++k]=j;

                break;

            }

            i=link[i];

        }

    }while(i);

    for(i=1;i<=k;i++){

        j=head[c[i]];

        while(j){

            if(a[j]==c[k]){

                printf("%d\n",k-i+1);

                for(;i<=k;i++)

                    printf("%d ",c[i]);

                return 0;

            }

            j=link[j];

        }

    }

}