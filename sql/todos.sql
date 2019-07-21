DROP TABLE IF EXISTS todos;

CREATE TABLE `todos` (
    `id` varchar(10) NOT NULL,
    `title` varchar(100) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
