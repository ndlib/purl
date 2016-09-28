---
category: projects
author:   Don Brower
title:    RepoPurl Spec
date:     2013-08-26
---

# Notes on the Repository PURL

These specs describe the repository purl service as it exists now (August 2013).

The current code is a perl script which is at
(https://svn.library.nd.edu/svn/svn_ndlibs_devel/applications/repopurl/trunk/PURL)[]

The pre-production service is at (http://repopurlpprd.library.nd.edu)[]

# Routes

## /admin

Return a summary page containing the total number of PURLs, the most used PURL, the number of PURLs used today, and the PURL with the highest usage count used today.
(Some of these fields may be "None Today" if no purls have been selected today).

There are also unclickable links to create a purl and to find a purl.

There is a search box to search PURLs.

## /admin/search

The result of searching. The search term is passed as form data.
If no search term is given, every purl is returned.

Results are returned as a table with the columns

* Purl ID
* File Name
* Repository URL
* PURL URL
* Access Count
* Date of last Access  (not time!)
* Description


## /view/:id

Return information about a specific purl.

* Purl ID
* Note
* File Name
* Date last accessed
* Repository URL
* PURL url
* Access count


## /view/:id/:filename

Proxy back or redirect the contents of the purl.
Not sure how it decides when to proxy back and when to redirect.
The code seems to have the following algorithm:

    does description start with "CurateND - "?
        YES: redirect (302) to repository location
        NO: proxy data out, setting the suggested file name

But the link [http://repopurlpprd.library.nd.edu/view/45]() is not a CurateND link, and it is redirected to the repository location.


## /purl/create

make a purl. Takes the following url parameters

    filename: the filename of the file
    url:      the repository url
    info:     the contents of the "note" field


## /query
## /purl/query

These seem to be broken. They always return the character '1'.

# Database schema dump

    mysql> show tables;
    +--------------------------+
    | Tables_in_repo_purl_pprd |
    +--------------------------+
    | object_access            |
    | purl                     |
    | repo_object              |
    | user                     |
    | user_roles               |
    +--------------------------+
    5 rows in set (0.00 sec)

    mysql> show columns from object_access;
    +----------------+--------------+------+-----+---------+----------------+
    | Field          | Type         | Null | Key | Default | Extra          |
    +----------------+--------------+------+-----+---------+----------------+
    | access_id      | int(11)      | NO   | PRI | NULL    | auto_increment |
    | date_accessed  | datetime     | NO   | MUL | NULL    |                |
    | ip_address     | varchar(255) | YES  | MUL | NULL    |                |
    | host_name      | varchar(255) | YES  | MUL | NULL    |                |
    | referer        | varchar(255) | YES  | MUL | NULL    |                |
    | user_agent     | varchar(255) | YES  | MUL | NULL    |                |
    | request_method | varchar(45)  | YES  | MUL | NULL    |                |
    | path_info      | varchar(255) | YES  |     | NULL    |                |
    | repo_object_id | int(11)      | NO   | MUL | NULL    |                |
    | purl_id        | int(11)      | NO   | MUL | NULL    |                |
    +----------------+--------------+------+-----+---------+----------------+
    10 rows in set (0.00 sec)

    mysql> show columns from purl;
    +----------------+--------------+------+-----+---------+----------------+
    | Field          | Type         | Null | Key | Default | Extra          |
    +----------------+--------------+------+-----+---------+----------------+
    | purl_id        | int(11)      | NO   | PRI | NULL    | auto_increment |
    | repo_object_id | varchar(255) | NO   | MUL | NULL    |                |
    | access_count   | int(11)      | NO   | MUL | 0       |                |
    | last_accessed  | datetime     | YES  |     | NULL    |                |
    | source_app     | varchar(255) | YES  |     | NULL    |                |
    | date_created   | datetime     | NO   | MUL | NULL    |                |
    +----------------+--------------+------+-----+---------+----------------+
    6 rows in set (0.00 sec)

    mysql> show columns from repo_object;
    +----------------+--------------+------+-----+---------+----------------+
    | Field          | Type         | Null | Key | Default | Extra          |
    +----------------+--------------+------+-----+---------+----------------+
    | repo_object_id | int(11)      | NO   | PRI | NULL    | auto_increment |
    | filename       | varchar(255) | NO   | MUL | NULL    |                |
    | url            | text         | NO   | MUL | NULL    |                |
    | date_added     | datetime     | NO   | MUL | NULL    |                |
    | add_source_ip  | varchar(255) | NO   |     | 0.0.0.0 |                |
    | date_modified  | datetime     | YES  |     | NULL    |                |
    | information    | text         | YES  | MUL | NULL    |                |
    +----------------+--------------+------+-----+---------+----------------+
    7 rows in set (0.01 sec)

    mysql> show columns from user;
    +----------+--------------+------+-----+---------+----------------+
    | Field    | Type         | Null | Key | Default | Extra          |
    +----------+--------------+------+-----+---------+----------------+
    | user_id  | int(11)      | NO   | PRI | NULL    | auto_increment |
    | username | varchar(255) | NO   | MUL | NULL    |                |
    +----------+--------------+------+-----+---------+----------------+
    2 rows in set (0.00 sec)

    mysql> show columns from user_roles;
    +---------+---------+------+-----+---------+-------+
    | Field   | Type    | Null | Key | Default | Extra |
    +---------+---------+------+-----+---------+-------+
    | user_id | int(11) | NO   | PRI | NULL    |       |
    | role_id | int(11) | NO   | PRI | NULL    |       |
    +---------+---------+------+-----+---------+-------+
    2 rows in set (0.00 sec)
