#include<bits/stdc++.h>

using namespace std;

int main(){

    //freopen("5.in","r",stdin);

    int i,j,k,a,b,m,n,t,cp,can;

    scanf("%d",&t);

    while(t--){

        scanf("%d %d",&n,&m);

        int v[n][n], vv[n], vvv[n];

        memset(v, 0, sizeof v);

        memset(v, 0, sizeof vv);

        for(i=0;i<m;i++){

            scanf("%d %d",&a,&b);

            v[a-1][b-1]=1;

            v[b-1][a-1]=1;

        }

        cp=0;

        for(i=0;i<n;i++){

            can=0;

            for(j=0;j<n;j++)can+=v[i][j];

            vvv[i]=can;

            vv[i]=1-(can%2);

            cp+=vv[i];

        }

        /*for(i=0;i<n;i++){

            for(j=0;j<n;j++){

                cout<<v[i][j];

            }

            cout<<endl;

        }

        */

        printf("%d\n",cp);

        for(i=0;i<n;i++){

            j=0;

            if(vv[i]==0)continue;

            while(vvv[i]>0){

                //cout<<i<<endl;

                a=i;

                while(vvv[a]>0&&vv[a]>0){

                    for(k=0;k<n;k++){

                        if(v[a][k]==1){

                            if(j==0){

                                v[a][k]=2;

                                v[k][a]=0;

                            }

                            else{

                                v[a][k]=0;

                                v[k][a]=2;

                            }

                            vvv[a]--;

                            vvv[k]--;

                            a=k;

                            break;

                        }

                    }

                }

                if(vvv[i]%2==1)j=1;

                else j=0;

            }

        }

        for(i=1;i<=n;i++){

            for(j=1;j<=n;j++){

                if(v[i-1][j-1]==2)printf("%d %d\n",i,j);

                if(v[i-1][j-1]==1){

                    printf("%d %d\n",i,j);

                    v[j-1][i-1]=0;

                }

            }

        }

    }

}