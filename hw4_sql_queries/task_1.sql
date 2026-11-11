select name, age
  from users
 where age = (select max(age) from users)
   and city = 'Москва';