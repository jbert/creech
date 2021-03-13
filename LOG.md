
Sat 13 Mar 11:18:12 GMT 2021 - "Mortality"
------------------------------------------

- Added a basic energy/food cost per tick and a 'Dead' state

- First bit of unexpected behaviour (not due to a bug :-) )
    - creeches will find a food source and sit and nibble at it
    - they only bite what they need
    - so once they find a food source, every tick they take a small nibble and
      do nothing else
    - this seems like an optimal survival strategy under current conditions!

- Write down some of the thoughts on genotype/phenotype in TODO

- Fixed turnHelper with the help of tests (yay), so the creeches can start to
  make some sensible choices


