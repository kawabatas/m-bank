-- +migrate Up
CREATE TABLE `users` (
  `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `users` (`name`) VALUES
('user1'),
('user2');

CREATE TABLE `balances` (
  `user_id` INT(11) UNSIGNED NOT NULL,
  `amount` INT(11) UNSIGNED NOT NULL DEFAULT '0',
  `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`),
  FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `payment_transactions` (
  `uuid` VARCHAR(255) NOT NULL,
  `user_id` INT(11) UNSIGNED NOT NULL,
  `amount` INT(11) NOT NULL,
  `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `try_time` DATETIME,
  `confirm_time` DATETIME,
  `cancel_time` DATETIME,
  PRIMARY KEY (`uuid`),
  FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +migrate Down
DROP TABLE IF EXISTS `payment_transactions`;
DROP TABLE IF EXISTS `balances`;
DROP TABLE IF EXISTS `users`;
