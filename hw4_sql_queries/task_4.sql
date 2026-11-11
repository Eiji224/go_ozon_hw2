select distinct u.name
  from orders as o
  join users as u on o.user_id = u.id
 where o.item = 'наушники';