#include<bits/stdc++.h>
using namespace std;
int i,i0,n,m,ans;
vector<int>va,vb;
int main()
{
    int a,b,c,d;
    scanf("%d %d %d %d",&a,&b,&c,&d);
    if(c==0&&d==0)
    {
        if(a==b+1)
        {
            printf("YES\n");
            for(int i=0;i<b;i++)
            {
                printf("0 1 ");
            }
            printf("0 ");
            return 0;
        }
    }
    if(a==0&&b==0)
    {
        if(c+1==d)
        {
            printf("YES\n");
            for(int i=0;i<c;i++)
            {
                printf("3 2 ");
            }
            printf("3 ");
            return 0;
        }
    }
    if(b-a+d==c&&b-a>=0)
    {
        printf("YES\n");
        int x=a,y=b-a,z=d;
        for(int i=0;i<x;i++)
        {
            printf("0 1 ");
        }
        for(int i=0;i<y;i++)
        {
            printf("2 1 ");
        }
        for(int i=0;i<z;i++)
        {
            printf("2 3 ");
        }
        return 0;
    }
    if(b-a+d-1==c&&b-a-1>=0)
    {
        printf("YES\n");
        int x=a,y=b-a-1,z=d;
        printf("1 ");
        for(int i=0;i<x;i++)
        {
            printf("0 1 ");
        }
        for(int i=0;i<y;i++)
        {
            printf("2 1 ");
        }
        for(int i=0;i<z;i++)
        {
            printf("2 3 ");
        }
        return 0;
    }
    if(b-a+d+1==c&&b-a>=0)
    {
        printf("YES\n");
        int x=a,y=b-a,z=d;
        for(int i=0;i<x;i++)
        {
            printf("0 1 ");
        }
        for(int i=0;i<y;i++)
        {
            printf("2 1 ");
        }
        for(int i=0;i<z;i++)
        {
            printf("2 3 ");
        }
        printf("2 ");
        return 0;
    }
    if(b-a+d==c&&b-a-1>=0)
    {
        printf("YES\n");
        int x=a,y=b-a-1,z=d;
        printf("1 ");
        for(int i=0;i<x;i++)
        {
            printf("0 1 ");
        }
        for(int i=0;i<y;i++)
        {
            printf("2 1 ");
        }
        for(int i=0;i<z;i++)
        {
            printf("2 3 ");
        }
        printf("2 ");
        return 0;
    }
    printf("NO\n");
    return 0;


    if(a==b)
    {
        if(c==d)
        {
            while(a&&b)
            {
                printf("0 1 0 1 ");
                a--,b--;
            }
            while(c&&d)
            {
                printf("2 3 2 3 ");
                c--,d--;
            }
        }
        else if(c==d+1)
        {
            while(a&&b)
            {
                printf("0 1 0 1 ");
                a--,b--;
            }
            while(c&&d)
            {
                printf("2 3 2 3 ");
                c--,d--;
            }
            printf("2");
        }
        else if(c==d-1)
        {
            if(a==0)
            {
                d--;
                printf("3 ");
                while(c&&d)
                {
                    printf("2 3 2 3 ");
                    c--,d--;
                }
            }
            else
            {
                printf("NO\n");
                return 0;
            }
        }
        else
        {
            printf("NO\n");
            return 0;
        }
    }
    else if(a==b+1)
    {
        if(c==d)
        {
            a--;
            printf("1 ");
            while(a&&b)
            {
                printf("0 1 0 1 ");
                a--,b--;
            }
            while(c&&d)
            {
                printf("2 3 2 3 ");
                c--,d--;
            }
        }
        else if(c==d+1)
        {
            a--;
            printf("1 ");
            while(a&&b)
            {
                printf("0 1 0 1 ");
                a--,b--;
            }
            while(c&&d)
            {
                printf("2 3 2 3 ");
                c--,d--;
            }
            printf("2");
        }
        else if(c==d-1)
        {
            if(a==0)
            {
                d--;
                printf("3 ");
                while(c&&d)
                {
                    printf("2 3 2 3 ");
                    c--,d--;
                }
            }
            else
            {
                printf("NO\n");
                return 0;
            }
        }
        else
        {
            printf("NO\n");
            return 0;
        }
    }
    else if(a==b-1)
    {
        if(c==0&&d==0)
        {
            printf("0 ");
            a--;
            while(a&&b)
            {
                printf("0 1 0 1 ");
                a--,b--;
            }
        }
        else
        {
            printf("NO\n");
            return 0;
        }
    }
    else
    {
        printf("NO\n");
        return 0;
    }
    return 0;
}