#include<bits/stdc++.h>
using namespace std;
int main()
{
    char table_card[2];
    int i,temp=0;
    cin>>table_card[0]>>table_card[1];
    char hand_card[5][2];
    for( i=0; i<5; i++)
    {
        for(int j=0; j<2; j++)
        {
            cin>>hand_card[i][j];
        }
    }
    for(i=0; i<5; i++)
    {
        if(table_card[0]==hand_card[i][0]||table_card[1]==hand_card[i][1])

            temp=1;
    }
    if(temp==1)
        cout<<"YES"<<endl;
    else
        cout<<"NO"<<endl;
    return 0;

}