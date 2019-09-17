-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

-- -----------------------------------------------------
-- Schema mydb
-- -----------------------------------------------------
-- -----------------------------------------------------
-- Schema transaction_table
-- -----------------------------------------------------
DROP SCHEMA IF EXISTS `transaction_table` ;

-- -----------------------------------------------------
-- Schema transaction_table
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `transaction_table` DEFAULT CHARACTER SET latin1 ;
USE `transaction_table` ;

-- -----------------------------------------------------
-- Table `transaction_table`.`points`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `transaction_table`.`points` ;

CREATE TABLE IF NOT EXISTS `transaction_table`.`points` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `point` BIGINT(20) NOT NULL,
  PRIMARY KEY (`id`))
ENGINE = InnoDB
AUTO_INCREMENT = 2
DEFAULT CHARACTER SET = latin1;


-- -----------------------------------------------------
-- Table `transaction_table`.`transactions`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `transaction_table`.`transactions` ;

CREATE TABLE IF NOT EXISTS `transaction_table`.`transactions` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `previous_date` DATETIME(6) NULL DEFAULT NULL,
  `date` DATETIME(6) NOT NULL,
  `previous` BIGINT(20) NULL DEFAULT NULL,
  `change` BIGINT(20) NOT NULL,
  `final` BIGINT(20) NOT NULL,
  PRIMARY KEY (`id`, `date`),
  UNIQUE INDEX `UNQ_transaction` (`date` ASC, `final` ASC),
  UNIQUE INDEX `UNQ_prev_transaction` (`previous_date` ASC, `previous` ASC),
  CONSTRAINT `REF_self_transaction`
    FOREIGN KEY (`previous_date` , `previous`)
    REFERENCES `transaction_table`.`transactions` (`date` , `final`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB
AUTO_INCREMENT = 4102
DEFAULT CHARACTER SET = latin1;

USE `transaction_table`;

DELIMITER $$

USE `transaction_table`$$
DROP TRIGGER IF EXISTS `transaction_table`.`transactions_BEFORE_INSERT` $$
USE `transaction_table`$$
CREATE
DEFINER=`root`@`%`
TRIGGER `transaction_table`.`transactions_BEFORE_INSERT`
BEFORE INSERT ON `transaction_table`.`transactions`
FOR EACH ROW
BEGIN
	IF (NEW.final < 0 OR (NEW.final <> COALESCE(NEW.`previous`, 0) + NEW.change)) THEN
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT='invalid data';
    END IF;
    
    IF (NEW.previous_date >= NEW.`date`) THEN
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT='invalid data';
	END IF;
    
    IF ((NEW.previous_date IS NULL OR NEW.`previous` IS NULL) AND (NEW.previous_date IS NOT NULL OR NEW.`previous` IS NOT NULL)) THEN
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT='invalid data';
    END IF;
END$$


DELIMITER ;

SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
