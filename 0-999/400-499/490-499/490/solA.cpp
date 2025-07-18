#include<bits/stdc++.h>

#include<iostream>

#include<cstdio>

#include<queue>

using namespace std;

int n;

int m;

queue<int > q1;

queue<int > q2;

queue<int> q3;

int main()

{

    scanf("%d",&m);

    int minx=100000;

    while(!q1.empty())

        q1.pop();

    while(!q2.empty())

        q2.pop();

    while(!q3.empty())

        q3.pop();

    for(int i=1; i<=m; i++)

    {

        scanf("%d",&n);

        if(n==1)

        {

            q1.push(i);

        }

        if(n==2)

        {

            q2.push(i);

        }

        if(n==3)

        {

            q3.push(i);

        }

    }

    minx=min(q1.size(),min(q2.size(),q3.size()));

    printf("%d\n",minx);

    for(int i=1; i<=minx; i++)

    {

        printf("%d %d %d",q1.front(),q2.front(),q3.front());

        q1.pop();

        q2.pop();

        q3.pop();

        printf("\n");

    }

    return 0;

}