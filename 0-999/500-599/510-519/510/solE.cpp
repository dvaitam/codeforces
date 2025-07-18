#include <bits/stdc++.h>

#define pb push_back

using namespace std;



const int N = 205;

int a[N];

int n , odd , even , src , snk , vis_id;

vector<int> adj[N];

int flow[N][N] , vis[N];

vector<int> v;

vector<vector<int> > ans;





bool prime(int x , int y){

    if(x&1 && y&1)return 0;

    if(!(x&1) && !(y&1))return 0;

    int z = x + y;

    int sq = sqrt(z);

    for(int i = 2 ; i <= sq ; i++)

        if(z %i == 0)

            return 0;

    return 1;

}



void get(int u){

    v.pb(u);

    vis[u] = vis_id;

    for(int i = 0 ; i < adj[u].size();i++){

        int nxt = adj[u][i];

        if(vis[nxt] == vis_id || nxt == src || nxt == snk)continue;

        if((a[u]%2) == flow[nxt][u]){

            //cout << u << " " << nxt << " " << a[u] + a[nxt] << endl;

            assert(prime(a[u] , a[nxt]));

            get(nxt);

        }

    }

}





int dfs(int u , int f){

    if(!f || u == snk)return f;

    if(vis[u] == vis_id)return 0;

    vis[u] = vis_id;

    for(int i = 0 ; i < adj[u].size() ; i++){

        int nxt = adj[u][i];

        int x = dfs(nxt , min(f , flow[u][nxt]));

        if(x){

            flow[u][nxt] -= x;

            flow[nxt][u] += x;

            return x;

        }

    }

    return 0;

}



int max_flow(){



    int ret = 0 , x;

    vis_id++;

    while(x = dfs(src , 1e5)){

        ret += x;

        vis_id++;

    }

    return ret;

}



int main()

{

    scanf("%d",&n);

    for(int i = 1 ; i <= n ; i++){

        scanf("%d",&a[i]);

        if(a[i] & 1)odd++;

        else even++;

    }



    if(odd != even){

        puts("Impossible");

        return 0;

    }



    for(int i = 1; i <= n ; i++)

        for(int j = i+1 ; j <= n ; j++)

            if(prime(a[i] , a[j])){

                int o = i , e = j;

                if(a[e]%2 != 0)swap(o,e);

                adj[o].pb(e);

                adj[e].pb(o);

                flow[o][e] = 1;

            }



    snk = n+1;



    for(int i = 1 ; i <= n ;i++){

        if(a[i] % 2 != 0){

            adj[src].pb(i);

            flow[src][i] = 2;

        }else{

            adj[i].pb(snk);

            flow[i][snk] = 2;

        }

    }



    if(max_flow() != n){

        puts("Impossible");

        return 0;

    }





    vis_id++;

    for(int i = 1; i <= n ; i++)

        if(vis[i] != vis_id){

            v.clear();

            get(i);

            ans.pb(v);

            assert(prime(a[v[0]] , a[v[v.size() - 1]]));

        }





    cout << ans.size() << endl;

    for(int i = 0 ; i < ans.size() ; i++){

        printf("%d",ans[i].size());

        for(int j = 0 ; j < ans[i].size() ; j++)

            printf(" %d",ans[i][j]);

        puts("");

    }



    return 0;

}