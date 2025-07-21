#include<bits/stdc++.h>
using namespace std;
inline int ask(int x){
  cout<<"? "<<x+1<<endl;
  int b; cin>>b; return b;
}
inline int solve(vector<set<int> > g){
  vector<int> f(g.size()),d(g.size()),m(g.size());
  function<void(int,int)> dfs=[&](int u,int p){
    f[u]=p,m[u]=d[u];
    for(int i:g[u])
      if(i!=p)d[i]=d[u]+1,dfs(i,u),m[u]=max(m[u],m[i]);
  };
  dfs(0,0);
  vector<bool> b(g.size());
  int l=max_element(d.begin(),d.end())-d.begin();
  if(ask(l))return l;
  for(int i=0;i<g.size();i++)
    if(m[i]-d[i]==51&&!b[i]){
      if(b[i]=true;ask(i)){
        if(i){
          while(ask(l),ask(i));
          return ask(f[i])?f[i]:f[f[f[i]]];
        }
        else break;
      }
      else g[f[i]].erase(i),g[i].erase(f[i]),dfs(0,0),i=-1;
    }
  for(int i=0;i<100;i++)ask(l);
  return 0;
}
int main(){
  ios::sync_with_stdio(false);
  int t; cin>>t;
  while(t--){
    int n; cin>>n;
    vector<set<int> > g(n);
    for(int i=1;i<n;i++){
      int u,v; cin>>u>>v;
      g[--u].emplace(--v),g[v].emplace(u);
    }
    int x=solve(g);
    cout<<"! "<<x+1<<endl;
  }
  return 0;
}