select item, price
  from orders
 where price = (select min(price) from orders where price > 10000);