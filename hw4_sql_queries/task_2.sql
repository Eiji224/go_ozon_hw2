select    city, count(1) as users_qty
  from    users
 group by city
having    count(1) > 1
 order by users_qty desc;