#include<bits/stdc++.h>
using namespace std;
int main()
{
    int g,t;
    string q;
    scanf("%d",&t);
    while(t--){
        int wq[211]={0},i;
        bool re,ck;
        re=ck=true;
        cin>>q;
        for(g=0;q[g];g++){
            wq[(int)q[g]]++;
        }
        for(g=97;g<125;g++){
            if(wq[g]>1){re=false;break;}
            if(wq[g]&&ck)i=g,ck=false;
            if(wq[g])
                if(g==i)i++;
                else {re=false;break;}
        }
        if(re)printf("Yes\n");
        else printf("No\n");
        q.clear();
    }
    return 0;
}