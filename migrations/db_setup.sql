CREATE TABLE IF NOT EXISTS users (id varchar(255)  NOT NULL, password varchar(255) NOT NULL, date_created DATETIME DEFAULT CURRENT_TIMESTAMP(), PRIMARY KEY (`id`));
CREATE TABLE IF NOT EXISTS loans (id VARCHAR(40), user_id varchar(250), amount int, term int, date_created DATE, status enum ('PENDING', 'APPROVED', 'PAID'), PRIMARY KEY  (`id`), FOREIGN KEY (`user_id`) REFERENCES users(id));
CREATE TABLE IF NOT EXISTS scheduled_repayments (id VARCHAR(40), loan_id VARCHAR(40), amount decimal(9,2), date DATE, status enum ('PENDING', 'PAID'), PRIMARY KEY  (`id`), FOREIGN KEY (`loan_id`) REFERENCES loans(id));
CREATE TABLE IF NOT EXISTS repayments (id VARCHAR(40), scheduled_repayment_id VARCHAR(40), amount decimal(9,2), date_created DATETIME DEFAULT CURRENT_TIMESTAMP(), PRIMARY KEY  (`id`), FOREIGN KEY (`scheduled_repayment_id`) REFERENCES scheduled_repayments(id));
INSERT IGNORE INTO users(id, password) VALUES ('aspire', 'aspire');
