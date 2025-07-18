#include<bits/stdc++.h>



using namespace std;



typedef unsigned long long LL;

map<int,int> tag;

const int NN = 5000011;

const int P = 71;

int a[30],b[30];

char str[30];



struct HASH

{

    int head[NN],next[NN];

    LL state[NN];

    double dp[NN];

    int tol;



    void init()

    {

        memset(head,-1,sizeof(head));

        tol = 0;

    }

    double find(LL s)

    {

        for(int i=head[s%NN];i!=-1;i=next[i])

            if(state[i]==s) return dp[i];

        return -1;

    }

    bool insert(const LL& s,const double& f)

    {

        if(find(s) > -1) return false;

        dp[++tol] = f;

        state[tol] = s;

        next[tol] = head[s%NN];

        head[s%NN] = tol;

        return true;

    }

} hsh;

LL myhash(const int& n)

{

    for(int i = 1; i<=n; ++i) b[i] = a[i];

    sort( b+1, b+n);

    LL rtn = 0;

    for(int i = 1; i<=n; ++i)

    {

        rtn *= P;

        rtn += b[i];

    }

    return rtn;

}

inline void set_state(int* st,const int& x)

{ //reverse operation of get_state (decoding)

    st[0] = (x>>4) & 3;

    st[1] = (x>>2) & 3;

    st[2] = x & 3;

}

inline int get_state(const int& a,const int& b,const int& c)

{ //(encoding)

    if(a>c) return (((c<<2)+b)<<2)+a;; //reduce state num

    return (((a<<2)+b)<<2)+c;

}

double dfs(const int& n)

{

    LL s = myhash(n);

    double f = hsh.find(s), c[4] = {0,1e9,1e9,1e9};

    if(f>-1) return f; //memoization

    for(int i = 1; i<n; ++i)

    {

        int x[3], y[3],nn = n;

        int ai = a[i], an = a[n]; //n : top

        set_state(x,ai);

        set_state(y,an);

        if( (!!x[0]) + (!!x[1]) + (!!x[2]) == 1) continue; //only one block left

        if( y[0] && y[1] && y[2] ) //top layer complete

        {

            nn++;

            y[0] = y[1] = y[2] = 0;

        }

        for(int j=0; j<3; ++j) if(x[j])

        {

            //at least two blocks and one at middle, one at left ( or right )

            if((!!x[0])+(!!x[1])+(!!x[2])==2&&(j==1||!x[1]))continue;

            //-------------assume taking out block j ---------------//

            int t = x[j];                                           // 

            x[j] = 0;                                               //

            a[i] = get_state(x[0],x[1],x[2]);                       //

            //can't take any more from this layer                   //

            if( (!!x[0])+(!!x[1])+(!!x[2]) == 1 || !x[1] ) a[i]=0;  //

            //------------------------------------------------------//

            x[j] = t;

            for(int k = 0; k<3; ++k) if(!y[k])

            {//put the taken block j at any available slot at top layer

                y[k] = t; //assume put at k

                a[nn] = get_state(y[0],y[1],y[2]);

                y[k] = 0; //refreshing

                c[t] = min(c[t],dfs(nn));

            }

        }

        a[i] = ai;

        a[n] = an;

    }

    if(c[1]==1e9 && c[2]==1e9 && c[3]==1e9) //no available moves

    { 

        hsh.insert(s,0); 

        return 0.0; //game over so no more time

    }

    double p = 1./6.; //black

    if(c[1] == 1e9) p += 1.0/3.0, c[1] = 0; //green

    if(c[2] == 1e9) p += 1.0/3.0, c[2] = 0; //blue

    if(c[3] == 1e9) p += 1.0/6.0, c[3] = 0; //red

    //the formula on codeforce editorial

    f = (c[1]/3.0+c[2]/3.0+c[3]/6.0+1.0)/(1.0-p);

    hsh.insert(s,f);

    return f;

}

int main()

{

    tag['G'] = 1;

    tag['B'] = 2;

    tag['R'] = 3;

    int n;

    scanf("%d",&n);

    hsh.init();

    for(int i = 1; i<=n; ++i)

    {

        scanf("%s",str);

        a[i]=get_state(tag[str[0]],tag[str[1]],tag[str[2]]);

    }

    printf("%.15lf\n",dfs(n));

}