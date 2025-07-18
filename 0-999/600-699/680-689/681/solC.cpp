#include<bits/stdc++.h>
using namespace std;
typedef long long ll;
typedef pair<int,int> ii;

#define X first
#define Y second
#define REP(i,a) for(i=0;i<a;++i)
#define REPP(i,a,b) for(i=a;i<b;++i)
#define FILL(a,x) memset(a,x,sizeof(a))
#define  mp make_pair
#define  pb push_back
#define sz(a) int((a).size())
#define  debug(ccc)  cout << #ccc << " = " << ccc << endl;
#define present(c,x) ((c).find(x) != (c).end())
priority_queue<int> pq;
vector<ii> a;
vector<ii>::iterator it, it1;
int main(){
  ios_base::sync_with_stdio(false);cin.tie(NULL);cout.tie(NULL);
  int n,i,temp,count;
  cin>>n;
  string insert ="insert";
  string min = "getMin";
  string rem = "removeMin";
  string str;
  count=0;
  for(i=0;i<n;i++)
  {
    cin>>str;
    if(!str.compare(insert))
      {
        cin>>temp;
        a.pb(mp(1,temp));
        pq.push(-(temp));
        count++;
      }
    else if(!str.compare(min))
      {
        cin>>temp;
        if(count==0 || (-1*pq.top()) > temp)
        {
            // std::cout << "hi2" << std::endl;
            int cd=temp;
            a.pb(mp(1,cd));
            pq.push(-(cd));
            count++;
        }
        else if((-1*pq.top()) <  temp)
        {
          while((-1*pq.top()) <  temp && count>0)
          {
            // std::cout << "hi1" << std::endl;
            a.pb(mp(3,-1*pq.top()));
            pq.pop();
            count--;
          }
          if(count==0 || (-1*pq.top()) > temp)
          {
              int cd=temp;
              a.pb(mp(1,cd));
              pq.push(-(cd));
              count++;
          }
        }
        a.pb(mp(2,temp));
      }
    else if(!str.compare(rem))
    {
      if(count==0)
      {
          a.pb(mp(1,1));
      }
      else
      {
        // std::cout << "remove" << std::endl;
        pq.pop();
        count--;
      }
      // std::cout << "hi" << std::endl;
      a.pb(mp(3,-1));
    }
  }
/*  for(it=a.begin();it!=a.end();it++)
  {

    if(it->X==1)
    {
      pq.push(-(it->Y));
      count++;
    }
    else if(it->X==2)
    {
      if(count==0 || (-1*pq.top()) > it->Y)
      {
          int cd=it->Y;
          a.insert(it,mp(1,cd));
          pq.push(-(cd));
          count++;
      }
      else if((-1*pq.top()) <  it -> Y)
      {
        while((-1*pq.top()) <  it -> Y)
        {
          it=a.insert(it,mp(3,-1*pq.top()));
          pq.pop();
          count--;
        }
      }
    }
    else if(it->X==3)
    {
      if(count==0)
      {
          it=a.insert(it,mp(1,1));
          pq.push(-1);
          count++;
      }
      else
      {
        pq.pop();
        count--;
      }
    }
  }*/
  cout<<a.size()<<endl;
  for(it=a.begin();it!=a.end();it++)
  {
    if(it->X==1)
      cout<<insert<<" "<<it->Y<<"\n";
    else if(it->X==2)
      cout<<min<<" "<<it->Y<<"\n";
    else if(it->X==3)
      cout<<rem<<"\n";
  }

  return 0;
}