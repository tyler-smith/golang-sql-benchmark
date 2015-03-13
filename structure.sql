CREATE TABLE `tickets` (
	`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	`subdomain_id` int(11) NOT NULL,
	`subject` varchar(255) NOT NULL DEFAULT '',
	`state` varchar(255) NOT NULL DEFAULT 'open',
	PRIMARY KEY (`id`)
)
