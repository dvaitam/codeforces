#include <iostream>

#include<vector>

#include<algorithm>

#include<math.h>

#include <bits/stdc++.h>



using namespace std;



bool isAnswer(double rad, vector< pair < pair<int,int>,int > > v,int s)

{

    long double r;

    long double diff=1000000-s;

    int temp=0;

    for(int i=0; i<v.size(); i++)

    {

        r=sqrt((v[i].first.first *v[i].first.first)+(v[i].first.second *v[i].first.second));

        // cout<<"current rad : "<<r<<endl;

        if(r<=rad)

        {

            temp+=v[i].second;

            //  cout<<temp<<endl;

        }

    }

    if(temp>=diff)

    {

       // cout <<fixed<<setprecision(10)<<  rad << " ---------" << endl;

        return true;

    }

    else

    {

        //cout <<fixed<<setprecision(10)<<  rad << endl;

        return false;

    }

}





int main()

{

    //magicity

    int n,p,x,y,z;

    vector<pair<pair<int,int>,int> >v;

    cin>>n>>p;

    for(int i=0; i<n; i++)

    {

        cin>>x>>y>>z;

        v.push_back(make_pair(make_pair(x,y),z));

    }

    long double st=0,ed=14143;

    while(ed-st >=0.0000001)

    {

        long double mid=(st+ed)/2;

        if(isAnswer(mid,v,p)==true)

            ed=mid;

        else

            st=mid;

    }

    //cout << ed << endl;

    if(isAnswer(ed,v,p)==false)

        cout<<"-1"<<endl;

    else

        cout<<fixed << setprecision(6) << ed<<endl;

    return 0;

}