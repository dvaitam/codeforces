#include <cstdio>
#include <cmath>
#include <vector>
#include <utility>
#include <algorithm>

using namespace std;

const double eps = 1e-8;

int sgn(double x)
{
  return x > eps ? 1 : (x < -eps ? -1 : 0);
}

struct pt
{
  double x, y;
  pt(double _x = 0, double _y = 0): x(_x), y(_y) {}
  void input()
  {
    scanf("%lf%lf", &x, &y);
  }
  double len() const
  {
    return sqrt(x * x + y * y);
  }
};

pt operator+(const pt& p1, const pt& p2)
{
  return pt(p1.x + p2.x, p1.y + p2.y);
}

pt operator-(const pt& p1, const pt& p2)
{
  return pt(p1.x - p2.x, p1.y - p2.y);
}

double operator^(const pt& p1, const pt& p2)
{
  return p1.x * p2.x + p1.y * p2.y;
}

double operator*(const pt& p1, const pt& p2)
{
  return p1.x * p2.y - p1.y * p2.x;
}

const int MAXN = 110;

int n;
pt src, dst, pol[MAXN];

bool get_inter(const pt& p1, const pt& p2, const pt& p3, const pt& p4, pt& c)
{
  double d1 = (p2 - p1) * (p3 - p1), d2 = (p2 - p1) * (p4 - p1);
  double d3 = (p4 - p3) * (p1 - p3), d4 = (p4 - p3) * (p2 - p3);
  int s1 = sgn(d1), s2 = sgn(d2), s3 = sgn(d3), s4 = sgn(d4);
  if (s2 == 0)
    return false;
  c = pt((p3.x * d2 - p4.x * d1) / (d2 - d1), (p3.y * d2 - p4.y * d1) / (d2 - d1));
  return s1 * s2 <= 0 && s3 * s4 <= 0;
}

double get_cost(int id1, int id2, const pt& p1, const pt& p2)
{
  double cost = (pol[id1 + 1] - p1).len() + (pol[id2] - p2).len();
  for (int i = (id1 + 1) % n; i != id2; i = (i + 1) % n)
    cost += (pol[i + 1] - pol[i]).len();
  return cost;
}

double get_min_cost()
{
  vector<pair<int, pt> > vec;
  pol[n] = pol[0];
  for (int i = 0; i < n; ++i)
    {
      pt c;
      if (get_inter(src, dst, pol[i], pol[i + 1], c))
	vec.push_back(make_pair(i, c));
    }
  if (vec.size() != 2)
    return (dst - src).len();
  if (sgn((vec[0].second - src) ^ (vec[0].second - vec[1].second)) > 0)
    swap(vec[0], vec[1]);
  double tmp = (vec[1].second - vec[0].second).len() * 2.0;
  double tmp1 = get_cost(vec[0].first, vec[1].first, vec[0].second, vec[1].second);
  double tmp2 = get_cost(vec[1].first, vec[0].first, vec[1].second, vec[0].second);
  double cost = (vec[0].second - src).len() + (vec[1].second - dst).len() + min(tmp, min(tmp1, tmp2));
  return cost;
}

int main()
{
  src.input();
  dst.input();
  scanf("%d", &n);
  for (int i = 0; i < n; ++i)
    pol[i].input();
  printf("%.8lf\n", get_min_cost());
}