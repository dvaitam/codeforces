//  ॐ

#include <bits/stdc++.h>

using namespace std;

#define PI 3.14159265358979323846

#define ll long long int





int main(){

   

    ios_base::sync_with_stdio(false);

    cin.tie(0);

    cout.tie(0);

 

    int test = 1;

    // cin>>test;





    while(test--){

                        

                        



                        int n,m;

                        cin>>n>>m;



                        vector<pair<int,int>> adj[n+1];



                        for(int i=0;i<m;i++){

                             int a,b,c;

                             cin>>a>>b>>c;

                             adj[a].push_back({b,c});

                             adj[b].push_back({a,c});

                        }



                        vector<double> ans(n+1,1e9);

                        vector<int> coeff(n+1,1e9),cons(n+1,1e9);





                        for(int i=1;i<=n;i++){

                                 if(coeff[i]!=1e9){

                                    continue;

                                 }

                                 

                                 queue<int> q;

                                 q.push(i);

                                 coeff[i]=1;

                                 cons[i]=0;



                                 long double x=1e9;

                                 vector<int> temp,vis;



                                 while(!q.empty()){

                                      int v=q.front();

                                      q.pop();



                                      // cout<<v<<": "<<coeff[v]<<' '<<cons[v]<<'\n';

                                      temp.push_back(-1*cons[v]*coeff[v]);

                                      vis.push_back(v);



                                      for(auto [u,c] : adj[v]){

                                            if(coeff[u]==1e9){

                                                 coeff[u]=-coeff[v];

                                                 cons[u]=c-cons[v];

                                                 q.push(u);

                                                 continue;

                                            }



                                            long double val1=x*coeff[u]+cons[u];

                                            long double val2=x*coeff[v]+cons[v];



                                            if(val1+val2==c){

                                                   continue;

                                            }



                                            if(x!=1e9){

                                                cout<<"NO";

                                                return 0;

                                            }



                                            if(coeff[u]+coeff[v]==0){

                                                cout<<"NO";

                                                return 0;

                                            }



                                            x=(long double)(c-cons[u]-cons[v])/(long double)(coeff[u]+coeff[v]);

                                      } 

                                 }



                                 if(x==1e9){

                                    sort(temp.begin(),temp.end());

                                    int sz=(int)temp.size();

                                    x=temp[sz/2];

                                 }



                                 for(auto u : vis){

                                      ans[u]=x*coeff[u]+cons[u];

                                 }

                        }



                        cout<<"YES\n";

                        cout<<fixed<<setprecision(5);



                        for(int i=1;i<=n;i++){

                            cout<<ans[i]<<' ';

                        }



                        



                        cout<<'\n';

        

    }

        return 0;

    

}