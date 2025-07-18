# include <iostream>
# include <algorithm>
# include <vector>
# define PII pair<int,int>
using namespace std;
const int MAX_N=100*1000+100;
int N;
PII p[MAX_N];
vector<int> team1;
vector<int> team2;
int main()
{
    ios::sync_with_stdio(false);
    cin >> N;
    for(int i=0;i<N;i++)
    {
        int temp;
        cin >> temp;
        p[i]=make_pair(temp,i+1);
    }
    sort(p,p+N);
    reverse(p,p+N);
    for(int i=0;i<N;i++)
    {
        int id=p[i].second;
        if(i%2==0)
            team1.push_back(id);
        else
            team2.push_back(id);
    }
    cout << (int)team1.size()<<endl;
    for(int i=0;i<(int)team1.size();i++)
        cout <<team1[i]<<" ";
    cout << endl;
    cout << (int)team2.size()<<endl;
    for(int i=0;i<(int)team2.size();i++)
        cout <<team2[i]<<" ";
    cout << endl;
    return 0;
}