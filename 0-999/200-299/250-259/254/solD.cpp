#include <bits/stdc++.h>
#include <iostream>

#include <cstdio>

#include <algorithm>

#include <set>

#include <vector>

#include <map>

#include <set>

#include <list>

#include <cmath>



using namespace std;



#define next roman_kaban



int iter1=0;

int iter2=0;

int k=0;

char a[1100];

int b[1100][1100];

int c[1100][1100];

int d[1100][1100];

int n,m,dist;

int xr[1000600];

int yr[1000600];



int r2,c2;

bool explode2(int x,int y);

bool explode1(int x,int y)

{

    if(x<0) return false;

    if(y<0) return false;

    if(x>n) return false;

    if(y>m) return false;

    if(b[x][y]==0) return false;



    c[x][y] = ++iter1;

    list<pair<pair<int,int>,int> > sp;

    sp.push_back(make_pair(make_pair(x,y),0));

    while(sp.begin()!=sp.end())

    {

        int dd = sp.front().second ;

        if(dd== dist) break;

        dd++;

        x = sp.front().first.first;

        y = sp.front().first.second;

        sp.pop_front();



        for(int xx = x-1;xx<=x+1;xx++)

        for(int yy = y-1;yy<=y+1;yy++)

            if(((xx!=x) ^ (yy!=y)) && b[xx][yy]==1 && c[xx][yy]!=iter1) {c[xx][yy] = iter1; sp.push_back(make_pair(make_pair(xx,yy),dd));}

    }

    int minx = 1;

    int maxx = n;

    int miny = 1;

    int maxy = m;

    bool all_dead = true;

    for(int i=0;i<k;i++)

    if(c[xr[i]][yr[i]]!=iter1)

    {

        all_dead = false;

        if(xr[i]-dist > minx) minx = xr[i] - dist;

        if(xr[i]+dist < maxx) maxx = xr[i] + dist;

        if(yr[i]-dist > miny) miny = yr[i] - dist;

        if(yr[i]+dist < maxy) maxy = yr[i] + dist;

    }

    r2 = c2 = -1;

    if(all_dead) return true;

    for(r2 = minx;r2<=maxx;r2++)

    for(c2 = miny;c2<=maxy;c2++)

    if(explode2(r2,c2)) return true;

    return false;

}



bool explode2(int x,int y)

{

    if(x<0) return false;

    if(y<0) return false;

    if(x>n) return false;

    if(y>m) return false;

    if(b[x][y]==0) return false;





    d[x][y] = ++iter2;

    list<pair<pair<int,int>,int> > sp;

    sp.push_back(make_pair(make_pair(x,y),0));

    while(sp.begin()!=sp.end())

    {

        int dd = sp.front().second ;

        if(dd== dist) break;

        dd++;

        x = sp.front().first.first;

        y = sp.front().first.second;

        sp.pop_front();



        for(int xx = x-1;xx<=x+1;xx++)

        for(int yy = y-1;yy<=y+1;yy++)

            if(((xx!=x) ^ (yy!=y)) && b[xx][yy]==1 && d[xx][yy]!=iter2) {d[xx][yy] = iter2; sp.push_back(make_pair(make_pair(xx,yy),dd));}

    }

    for(int i=0;i<k;i++)

    if(c[xr[i]][yr[i]]!=iter1 && d[xr[i]][yr[i]]!=iter2) return false;

    return true;

}



int main()

{

    freopen("input.txt","r",stdin);

    freopen("output.txt","w",stdout);

    fgets(a, sizeof(a), stdin);

    sscanf(a,"%d%d%d",&n,&m,&dist);

    k=0;

    for(int i=1;i<=n;i++)

    {

        fgets(a, sizeof(a), stdin);

        for(int j=1;j<=m;j++)

            if(a[j-1]=='R'){b[i][j] = 1; xr[k] = i; yr[k] = j; k++;}else

            if(a[j-1]=='.') b[i][j] = 1;

    }

    if(k> 2*(2*dist+1)*(2*dist+1) ){cout << -1<<endl; return 0;}

    for(int i=xr[0]-dist;i<=xr[0]+dist;i++)

    for(int j=yr[0]-dist;j<=yr[0]+dist;j++)

    if( explode1(i,j) )

    {

        if(i==r2 && j==c2) {r2=c2=-1;}

        if(r2==-1)

        {

            for(r2 = 1;r2<=n;r2++)

            for(c2 = 1;c2<=m;c2++)

            if(b[r2][c2]==1 && (r2!=i || c2!=j) )

            {

                cout << i<<' '<<j << ' '<<r2<<' '<<c2<<endl;

                return 0;

            }

        }

        cout << i<<' '<<j << ' '<<r2<<' '<<c2<<endl;

        return 0;

    }

    cout <<-1<<endl;

    return 0;

}

/*

30000000003312

0

*/