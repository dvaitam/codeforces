#include <iostream>

#include <vector>

#include <algorithm>



using namespace std;



int u, sz, idx[200001];



bool comp(const pair<int, int>& a, const pair<int, int>& b){

    return (idx[a.first] - idx[u] + sz)%sz < (idx[b.first] - idx[u] + sz)%sz;

}



vector<bool> pchk;

vector<vector<int>> edge;



void dfs(int ind){

    if (pchk[ind]) return;

    pchk[ind] = 1;



    for (auto e : edge[ind])

        dfs(e);



    cout<<ind<<"\n";

}



int main()

{

    ios_base::sync_with_stdio(false); cin.tie(NULL); cout.tie(NULL);

    int N, M; cin>>N>>M;

    vector<int> arr(N+1);

    for (int i=1; i<=N; i++)

        cin>>arr[i];



    int x, y;

    vector<vector<pair<int, int>>> f(N+1);

    edge = vector<vector<int>>(M+1);

    for (int i=1; i<=M; i++){

        cin>>x>>y;

        f[x].push_back({y, i});

        f[y].push_back({x, i});

    }



    vector<bool> chk(N+1, 0);

    for (int i=1; i<=N; i++){

        if (chk[i]) continue;

        

        int j = i;

        vector<int> cycle;

        while (!chk[j]){

            chk[j] = 1;

            idx[j] = cycle.size();

            cycle.push_back(j);

            j = arr[j];

        }



        sz = cycle.size();

        if (sz == 1) continue;



        for (int k=0; k<sz; k++){

            u = cycle[k];

            sort(f[u].begin(), f[u].end(), comp);

            x = u;

            for (int v=0; v+1<f[u].size(); v++)

                edge[f[u][v+1].second].push_back(f[u][v].second);

        }

    }



    pchk = vector<bool>(M+1, 0);

    for (int i=1; i<=M; i++){

        if (pchk[i]) continue;

        dfs(i);

    }



    return 0;

}