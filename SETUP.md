# Setting Up The Test Database

First install MySQL locally.
Then run the following inside the mysql command prompt.
In this example, I logged in as the root user so that I have permission to create the database.
Once I create the database I grant rights to change it to the annonyous user,
and then I load the seed data file, which will create the tables and load a handful of test records.

    $ mysql -u root
    mysql> create table repopurl;
    Query OK, 1 row affected (0.04 sec)

    mysql> grant all on repopurl to ''@'localhost';
    Query OK, 0 rows affected (0.04 sec)

    mysql> use repopurl;
    Reading table information for completion of table and column names
    You can turn off this feature to get a quicker startup with -A

    Database changed
    mysql> source test/seed_data.sql;
    Query OK, 0 rows affected (0.06 sec)

    Query OK, 0 rows affected, 1 warning (0.29 sec)

    Query OK, 0 rows affected (0.13 sec)

    Query OK, 0 rows affected (0.16 sec)

    Query OK, 0 rows affected (0.14 sec)

    Query OK, 10 rows affected (0.00 sec)
    Records: 10  Duplicates: 0  Warnings: 0

    Query OK, 10 rows affected (0.00 sec)
    Records: 10  Duplicates: 0  Warnings: 0

    mysql> exit;
    Bye

