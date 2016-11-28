CREATE TABLE `object_access` (
  `access_id` int(11) NOT NULL AUTO_INCREMENT,
  `date_accessed` datetime NOT NULL,
  `ip_address` varchar(255) DEFAULT NULL,
  `host_name` varchar(255) DEFAULT NULL,
  `referer` text,
  `user_agent` varchar(255) DEFAULT NULL,
  `request_method` varchar(45) DEFAULT NULL,
  `path_info` varchar(255) DEFAULT NULL,
  `repo_object_id` int(11) NOT NULL,
  `purl_id` int(11) NOT NULL,
  PRIMARY KEY (`access_id`),
  KEY `ip` (`ip_address`) USING BTREE,
  KEY `date` (`date_accessed`),
  KEY `host_name` (`host_name`),
  KEY `user_agent` (`user_agent`),
  KEY `request_method` (`request_method`),
  KEY `purl_id` (`purl_id`) USING BTREE,
  KEY `repo_object_id` (`repo_object_id`)
) DEFAULT CHARSET=utf8;

CREATE TABLE `purl` (
  `purl_id` int(11) NOT NULL AUTO_INCREMENT,
  `repo_object_id` varchar(255) NOT NULL,
  `access_count` int(11) NOT NULL DEFAULT '0',
  `last_accessed` datetime DEFAULT NULL,
  `source_app` varchar(255) DEFAULT NULL,
  `date_created` datetime NOT NULL,
  PRIMARY KEY (`purl_id`),
  KEY `last_accessed` (`date_created`),
  KEY `access_count` (`access_count`),
  KEY `repo_object` (`repo_object_id`),
  KEY `date_created` (`date_created`)
) DEFAULT CHARSET=utf8;

CREATE TABLE `repo_object` (
  `repo_object_id` int(11) NOT NULL AUTO_INCREMENT,
  `filename` varchar(255) NOT NULL,
  `url` text NOT NULL,
  `date_added` datetime NOT NULL,
  `add_source_ip` varchar(255) NOT NULL DEFAULT '0.0.0.0',
  `date_modified` datetime DEFAULT NULL,
  `information` text,
  PRIMARY KEY (`repo_object_id`),
  KEY `filename` (`filename`),
  KEY `url` (`url`(50)),
  KEY `date_added` (`date_added`),
  KEY `info` (`information`(50))
) DEFAULT CHARSET=utf8;

CREATE TABLE `user` (
  `user_id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL,
  PRIMARY KEY (`user_id`),
  KEY `username` (`username`)
) DEFAULT CHARSET=utf8;

CREATE TABLE `user_roles` (
  `user_id` int(11) NOT NULL,
  `role_id` int(11) NOT NULL,
  PRIMARY KEY (`user_id`,`role_id`)
) DEFAULT CHARSET=utf8;

INSERT INTO purl (repo_object_id, access_count, last_accessed, date_created) VALUES
(1, 723, "2016-11-16 03:27:15", "2011-09-14 13:56:38"),
(2,  14, "2016-02-23 14:18:14", "2011-09-14 14:18:48"),
(3,   9, "2016-02-23 14:18:11", "2011-09-14 14:37:41"),
(4,   8, "2016-02-23 14:18:07", "2011-09-14 14:38:26"),
(5, 625, "2016-11-15 14:16:14", "2011-09-14 14:40:11"),
(6,  19, "2015-03-25 13:31:37", "2011-09-14 14:40:38"),
(7,   2, "2015-03-25 13:08:08", "2011-09-14 14:41:01"),
(8,   2, "2015-03-25 13:08:20", "2011-09-14 14:41:23"),
(9,   3, "2015-03-25 13:08:31", "2011-09-14 14:47:39"),
(10,  2, "2015-03-25 13:08:47", "2011-09-14 14:48:05");

INSERT INTO repo_object (filename, url, date_added, add_source_ip, date_modified, information) VALUES
("multilevel.zip",       "https://fedoraprod.library.nd.edu:8443/fedora/get/TEMP:4/DS4",                           "2011-09-14 14:47:38", "10.1.2.3",   NULL,                  "Published Paper"),
("834207.pdf",           "https://fedoraprod.library.nd.edu:8443/fedora/get/CATHOLLIC-PAMPHLET:834207/PDF1",       "2011-09-14 15:18:59", "10.1.2.3",   NULL,                  "Catholic Pamphlet 834207.pdf"),
("825854.pdf",           "https://fedoraprod.library.nd.edu:8443/fedora/get/CATHOLLIC-PAMPHLET:825854/PDF1",       "2011-10-17 14:53:11", "10.1.2.3",   NULL,                  "Catholic Pamphlet 825854.pdf"),
("cu31924026347082.pdf", "https://fedoraprod.library.nd.edu:8443/fedora/get/CYL:cu31924026347082/PDF1",            "2012-05-04 16:31:46", "10.1.2.3", NULL,                  "Catholic Youth Literature cu31924026347082.pdf"),
("000746007.pdf",        "https://fedoraprod.library.nd.edu:8443/fedora/get/CATHOLLIC-PAMPHLET:000746007/content", "2012-12-20 16:45:21", "10.1.2.3", NULL,                  "Catholic Pamphlet 000746007.pdf"),
("238501",               "http://archive.org/details/grammarscience02peargoog",                                    "2013-04-17 00:00:00", "0.0.0.0",        "2013-04-17 00:00:00", "Reformatting Unit: 238501"),
("1051413",              "http://catalog.hathitrust.org/Record/009783954",                                         "2013-04-17 00:00:00", "0.0.0.0",        "2013-04-17 00:00:00", "Reformatting Unit: 1051413"),
("9s161544m9s",          "http://curate.nd.edu/show/9s161544m9s",                                                  "2013-04-29 03:10:19", "0.0.0.0",        "2013-04-29 03:10:20", "CurateND - 9s161544m9s"),
("AIPE_for_Contras.pdf", "https://curate.nd.edu/downloads/hq37vm4368p",                                            "2011-09-14 14:56:12", "10.1.2.3",   "2016-02-19 14:31:14", "Published Paper"),
("210611",               "http://catalog.hathitrust.org/Record/009022229",                                         "2013-04-17 00:00:00", "0.0.0.0",        "2013-04-17 00:00:00", "Reformatting Unit: 210611");

