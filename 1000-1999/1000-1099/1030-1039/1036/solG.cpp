#pragma once

#include <bits/stdc++.h>
using namespace std;

int n,m;

vector<int>graph[1000009];

int ind[1000009],outd[1000009];

map<int,int>in,out;

int hiss[25];

bool vis[1000009];

int val[1000009];

void call(int u)
{
    vis[u]=1;
    for(int i=0;i<graph[u].size();i++)
    {
        int v=graph[u][i];
        if(!vis[v])call(v);
        val[u]|=val[v];
    }
    if(outd[u]==0)val[u]|=(1<<out[u]);
    if(ind[u]==0)hiss[in[u]]=val[u];
    return;
}

int main()
{
    scanf("%d %d",&n,&m);
    int x,y;
    for(int i=1;i<=m;i++)
    {
        scanf("%d %d",&x,&y);
        graph[x].push_back(y);
        ind[y]++;
        outd[x]++;
    }
    int ii=0,oo=0;
    for(int i=1;i<=n;i++)
    {
        if(ind[i]==0)
        {
            in[i]=++ii;
        }
        if(outd[i]==0)
        {
            out[i]=++oo;
        }
    }
    for(int i=1;i<=n;i++)
    {
        if(vis[i])continue;
        call(i);
    }
    for(int i=0;i<ii;i++)
    {
        hiss[i]=(hiss[i+1]>>1);
    }
    //for(int i=0;i<2;i++)cout<<hiss[i]<<endl;
    for(int i=1;i<((1<<ii));i++)
    {
        if(i==((1<<ii)-1))
        {
            cout<<"YES"<<endl;
            return 0;
        }
        x=0;
        for(int j=0;j<ii;j++)
        {
            if((i&(1<<j)))x|=hiss[j];
        }
        if(__builtin_popcount(x)==__builtin_popcount(i))
        {
            cout<<"NO"<<endl;
            return 0;
        }
    }
}