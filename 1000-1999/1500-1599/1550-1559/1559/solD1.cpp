#define _CRT_SECURE_NO_WARNINGS
#include<iostream>
#include <bits/stdc++.h>
#include <iterator>
#include <unordered_map>
#define ll long long
#define all(a) (a).begin(),(a).end()
using namespace std;

void fast() {
    ios::sync_with_stdio(0);
    cin.tie(0);
    cout.tie(0);
}

void File() {
#ifndef ONLINE_JUDGE
    freopen("Input.txt", "r", stdin);
    freopen("Output.txt", "w", stdout);
#endif
}
int  n,m1,m2;
int findParent(int u,vector<int>& par)
{
    if(u==par[u])
        return u;
    return par[u]= findParent(par[u],par);
}
void join(int a,int b,vector<int>&par)
{
    a= findParent(a,par);
    b= findParent(b,par);
    if(a!=b)
        par[a]=b;
}
void solve()
{
   cin>>n>>m1>>m2;
   vector<int>par1(n);
   vector<int>par2(n);
    iota(all(par1),0);
    iota(all(par2),0);
   for(int i=0;i<m1;i++)
   {
       int u,v;
       cin>>u>>v;
       u--,v--;
       if(findParent(u,par1)!= findParent(v,par1))
           join(u,v,par1);
   }
    for(int i=0;i<m2;i++)
    {
        int u,v;
        cin>>u>>v;
        u--,v--;
        if(findParent(u,par2)!= findParent(v,par2))
            join(u,v,par2);
    }
    vector<pair<int,int>>ans;
    for(int i=0;i<n;i++)
    {
        for(int j=0;j<n;j++)
        {
            if((findParent(i,par1)!= findParent(j,par1))&&(findParent(i,par2)!= findParent(j,par2)))
            {
                ans.push_back({i+1,j+1});
                join(findParent(i,par1), findParent(j,par1),par1);
                join(findParent(i,par2), findParent(j,par2),par2);
            }
        }
    }
    cout<<ans.size()<<"\n";
    for(auto it:ans)
        cout<<it.first<<" "<<it.second<<"\n";
}
int main() {
    fast();
    File();
    int t = 1;
    //cin >> t;
    while (t--) { solve(); }

    return 0;
}