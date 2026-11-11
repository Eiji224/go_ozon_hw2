select u.name as name,
       coalesce(sum(o.price), 0) as spent
  from users as u
  left join orders as o on o.user_id = u.id
 group by u.name
 order by spent desc;