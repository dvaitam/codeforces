#include<bits/stdc++.h>

#define pii pair<int,int>

#define ll long long

#define endl '\n'

#define fastIO ios_base::sync_with_stdio(0);cin.tie(0);cout.tie(0)



using namespace std;



struct info

{

    int x,y,z;



    info(int _x,int _y,int _z)

    {

        x=_x;

        y=_y;

        z=_z;

    }

};



int main()

{

    fastIO;



    int t,n,m;

    bool f;



    cin>>t;



    while(t--)

    {

        cin>>n>>m;



        char s[n+1][m+1];

        int cnt[n+1];

        char ch;



        int num,total=0;

        for(int i=1; i<=n; i++)

        {

            num=0;

            for(int j=1; j<=m; j++)

            {

                cin>>ch;



                if(ch=='1')num++;

                s[i][j]=ch;

            }



            cnt[i]=num;

            total+=num;

        }



        if(total%n)

        {

            cout<<-1<<endl;

            continue;

        }



        int per=total/n;



        multiset<int>mi[m+1];

        vector<int>les;



        for(int i=1; i<=n; i++)

        {

            if(cnt[i]>per)

            {

                for(int j=1; j<=m; j++)

                    if(s[i][j]=='1')mi[j].insert(i);

            }



            else if(cnt[i]<per)les.push_back(i);

        }



        int sz=les.size();

        vector<info>ans;



        for(int i=0; i<sz; i++)

        {

           int ind=les[i];



           while(cnt[ind]<per)

           {

               for(int j=1;j<=m && cnt[ind]<per ;j++)

               {

                   if(s[ind][j]=='0')

                   {

                       vector<int>tr;

                       for(auto it:mi[j])

                       {

                           tr.push_back(it);

                           if(cnt[it]>per)

                           {

                               s[it][j]='0';





                               s[ind][j]='1';



                               cnt[ind]++;

                               cnt[it]--;



                               info tm(it,ind,j);

                               ans.push_back(tm);

                               break;

                           }

                       }



                       int tz=tr.size();

                       for(int l=0;l<tz;l++)mi[j].erase(tr[l]);

                   }

               }

           }

        }



        sz=ans.size();



        cout<<sz<<endl;



        for(int i=0;i<sz;i++)

        cout<<ans[i].x<<' '<<ans[i].y<<' '<<ans[i].z<<endl;

    }



    return 0;

}