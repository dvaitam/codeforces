// GOAT

#include<bits/stdc++.h>



using namespace std;



const int MAX = 5001;



vector<int> ls[MAX+1],dwa,niedwa;

int odw[MAX+1],czy_2[MAX+1];

bool flag = true;



void dfs(int ten){

    if(!flag){return;}

    odw[ten]=1;

    if(czy_2[ten]){ dwa.push_back(ten); }

    else{           niedwa.push_back(ten); } 

    for(int sas : ls[ten]){

        if(odw[sas]){

            if(czy_2[sas]==czy_2[ten]){

                flag = false;

                return;

            }

            continue;

        }

        czy_2[sas] = !czy_2[ten];

        dfs(sas);

    }

}



int main(){



    ios_base::sync_with_stdio(0);

    cin.tie(0);

    cout.tie(0);



    int n=0,m=0,n1=0,n2=0,n3=0;

    cin>>n>>m;

    cin>>n1>>n2>>n3;

    if(m>0&&n2==0){

        cout<<"NO";return 0;

    }

    vector<vector<bitset<MAX+1>>> dp(2);

    dp[0].assign(n+1,0);

    dp[1].assign(n+1,0);

    bitset<MAX+1> jest;

    for(int i=1;i<=m;++i){

        int u=0,v=0;

        cin>>u>>v;

        ls[u].push_back(v);

        ls[v].push_back(u);

    }

    jest[0]=1;

    int last=0;

    for(int i=1;i<=n;++i){

        if(odw[i]){continue;}

        niedwa.clear();dwa.clear();

        czy_2[i]=1;

        dfs(i);

        int a = dwa.size(), b = niedwa.size();

        if(a+b==1){continue;}

        if(!flag){

            cout<<"NO";

            return 0;

        }

        last++;

        bitset<MAX+1> jest2;

        for(int c=0;c<=n2;++c){

            if(c+a <= n2 && jest[c]){

                dp[last%2][c+a] = dp[(last+1)%2][c];

                for(int dod : dwa){

                    dp[last%2][c+a][dod] = 1;

                }

                jest2[c+a]=1;

            }

            if(c+b <= n2 && jest[c]){

                dp[last%2][c+b] = dp[(last+1)%2][c];

                for(int dod : niedwa){

                    dp[last%2][c+b][dod] = 1;

                }

                jest2[c+b]=1;

            }

        }

        jest = jest2;

    }

    for(int ile2=0;ile2<=n2;++ile2){

        if(!jest[ile2]){continue;}

        vector<int> odp(n+1);

        int N1=n1,N2=n2-ile2,N3=n3;

        for(int i = 1;i<=n;++i){

            if(dp[last%2][ile2][i]){

                odp[i]=2;

            }else{

                if(ls[i].empty() && N2){

                    odp[i]=2;--N2;

                }else{

                    if(N1){odp[i]=1;--N1;}

                    else if(N3){odp[i]=3;--N3;}

                }

            }

        }

        bool ok=true;

        odp[0]=1;

        for(int x : odp){

            if(!x){

                ok=false;

            }

        }

        odp[0]=0;

        if(ok){

            cout<<"YES\n";

            for(int x : odp){

                if(x==0){continue;}

                cout<<x;

            }

            return 0;

        }

    }

    cout<<"NO";

    return 0;

}