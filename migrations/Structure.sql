/*M!999999\- enable the sandbox mode */ 
-- MariaDB dump 10.19-11.8.6-MariaDB, for Linux (x86_64)
--
-- Host: 127.0.0.1    Database: pos_kopitiam
-- ------------------------------------------------------
-- Server version	8.4.9

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*M!100616 SET @OLD_NOTE_VERBOSITY=@@NOTE_VERBOSITY, NOTE_VERBOSITY=0 */;

--
-- Table structure for table `audit_logs`
--

DROP TABLE IF EXISTS `audit_logs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `audit_logs` (
  `id` char(36) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` char(36) NOT NULL,
  `action` varchar(50) NOT NULL,
  `entity` varchar(50) NOT NULL,
  `entity_id` char(36) DEFAULT NULL,
  `old_value` json DEFAULT NULL,
  `new_value` json DEFAULT NULL,
  `timestamp` datetime(3) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_audit_logs_deleted_at` (`deleted_at`),
  KEY `idx_audit_user` (`user_id`),
  KEY `idx_audit_action` (`action`),
  KEY `idx_audit_entity` (`entity`),
  KEY `idx_audit_entity_id` (`entity_id`),
  KEY `idx_audit_timestamp` (`timestamp`),
  CONSTRAINT `fk_audit_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `branches`
--

DROP TABLE IF EXISTS `branches`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `branches` (
  `id` char(36) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(191) NOT NULL,
  `address` longtext,
  `phone` varchar(191) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_branches_deleted_at` (`deleted_at`),
  KEY `idx_branch_name` (`name`),
  KEY `idx_branch_phone` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `cash_movements`
--

DROP TABLE IF EXISTS `cash_movements`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `cash_movements` (
  `id` char(36) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `shift_id` char(36) NOT NULL,
  `type` varchar(10) NOT NULL,
  `amount` double NOT NULL,
  `reason` text,
  `recorded_by` char(36) NOT NULL,
  `recorded_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_cash_movements_deleted_at` (`deleted_at`),
  KEY `idx_cashmovement_shift` (`shift_id`),
  KEY `idx_cashmovement_type` (`type`),
  KEY `idx_cashmovement_user` (`recorded_by`),
  KEY `idx_cashmovement_date` (`recorded_at`),
  CONSTRAINT `fk_cashmovements_shift` FOREIGN KEY (`shift_id`) REFERENCES `shifts` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_cashmovements_user` FOREIGN KEY (`recorded_by`) REFERENCES `users` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `categories`
--

DROP TABLE IF EXISTS `categories`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `categories` (
  `id` char(36) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(191) NOT NULL,
  `sort_order` bigint DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_categories_name` (`name`),
  KEY `idx_categories_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `customers`
--

DROP TABLE IF EXISTS `customers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `customers` (
  `id` char(36) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(191) NOT NULL,
  `phone` varchar(191) NOT NULL,
  `email` varchar(191) DEFAULT NULL,
  `loyalty_points` bigint DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_customers_phone` (`phone`),
  UNIQUE KEY `idx_customers_email` (`email`),
  KEY `idx_customers_deleted_at` (`deleted_at`),
  KEY `idx_customer_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `discount_applications`
--

DROP TABLE IF EXISTS `discount_applications`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `discount_applications` (
  `id` char(36) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `order_id` char(36) NOT NULL,
  `discount_id` char(36) NOT NULL,
  `discount_amount` double DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_discount_applications_deleted_at` (`deleted_at`),
  KEY `idx_discountapp_order` (`order_id`),
  KEY `idx_discountapp_discount` (`discount_id`),
  CONSTRAINT `fk_discountapp_discount` FOREIGN KEY (`discount_id`) REFERENCES `discounts` (`id`) ON DELETE RESTRICT,
  CONSTRAINT `fk_discountapp_order` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `discounts`
--

DROP TABLE IF EXISTS `discounts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `discounts` (
  `id` char(36) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(191) NOT NULL,
  `type` varchar(20) NOT NULL,
  `value` double NOT NULL,
  `min_order_amount` double DEFAULT NULL,
  `max_discount` double DEFAULT NULL,
  `start_date` datetime(3) DEFAULT NULL,
  `end_date` datetime(3) DEFAULT NULL,
  `is_active` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`),
  KEY `idx_discounts_deleted_at` (`deleted_at`),
  KEY `idx_discount_name` (`name`),
  KEY `idx_discount_type` (`type`),
  KEY `idx_discount_active` (`is_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `menu_categories`
--

DROP TABLE IF EXISTS `menu_categories`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `menu_categories` (
  `menu_id` char(36) NOT NULL,
  `category_id` char(36) NOT NULL,
  PRIMARY KEY (`menu_id`,`category_id`),
  KEY `fk_menu_categories_category` (`category_id`),
  CONSTRAINT `fk_menu_categories_category` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`),
  CONSTRAINT `fk_menu_categories_menu` FOREIGN KEY (`menu_id`) REFERENCES `menus` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `menus`
--

DROP TABLE IF EXISTS `menus`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `menus` (
  `id` char(36) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(191) NOT NULL,
  `price` double NOT NULL,
  `is_active` tinyint(1) DEFAULT '1',
  `is_recipe_based` tinyint(1) DEFAULT '0',
  `daily_stock` bigint DEFAULT '0',
  `station` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_menus_deleted_at` (`deleted_at`),
  KEY `idx_menu_name` (`name`),
  KEY `idx_menu_active` (`is_active`),
  KEY `idx_menu_station` (`station`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `order_details`
--

DROP TABLE IF EXISTS `order_details`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `order_details` (
  `id` char(36) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `order_id` char(36) NOT NULL,
  `menu_id` char(36) NOT NULL,
  `menu_name` longtext NOT NULL,
  `quantity` bigint NOT NULL,
  `price` double NOT NULL,
  `cost_price` double DEFAULT NULL,
  `subtotal` double NOT NULL,
  `notes` text,
  `station` varchar(50) DEFAULT NULL,
  `is_voided` tinyint(1) DEFAULT '0',
  `void_reason` longtext,
  PRIMARY KEY (`id`),
  KEY `idx_order_details_deleted_at` (`deleted_at`),
  KEY `idx_details_order` (`order_id`),
  KEY `idx_details_menu` (`menu_id`),
  KEY `idx_details_voided` (`is_voided`),
  CONSTRAINT `fk_details_menu` FOREIGN KEY (`menu_id`) REFERENCES `menus` (`id`) ON DELETE RESTRICT,
  CONSTRAINT `fk_orders_order_details` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `order_returns`
--

DROP TABLE IF EXISTS `order_returns`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `order_returns` (
  `id` char(36) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `order_id` char(36) NOT NULL,
  `processed_by` char(36) NOT NULL,
  `reason` text NOT NULL,
  `original_amount` double DEFAULT NULL,
  `return_amount` double NOT NULL,
  `processed_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_order_returns_deleted_at` (`deleted_at`),
  KEY `idx_return_order` (`order_id`),
  KEY `idx_return_processor` (`processed_by`),
  KEY `idx_return_date` (`processed_at`),
  CONSTRAINT `fk_returns_order` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_returns_processor` FOREIGN KEY (`processed_by`) REFERENCES `users` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `orders`
--

DROP TABLE IF EXISTS `orders`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `orders` (
  `id` char(36) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `order_number` varchar(191) NOT NULL,
  `branch_id` char(36) DEFAULT NULL,
  `table_id` char(36) DEFAULT NULL,
  `employee_id` char(36) DEFAULT NULL,
  `customer_id` char(36) DEFAULT NULL,
  `shift_id` char(36) DEFAULT NULL,
  `order_type` varchar(20) NOT NULL,
  `status` varchar(20) NOT NULL,
  `subtotal` double DEFAULT NULL,
  `tax_amount` double DEFAULT NULL,
  `service_charge` double DEFAULT NULL,
  `discount_amount` double DEFAULT NULL,
  `total` double DEFAULT NULL,
  `notes` text,
  `has_returns` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_orders_order_number` (`order_number`),
  KEY `idx_orders_deleted_at` (`deleted_at`),
  KEY `idx_orders_branch_status` (`branch_id`,`status`),
  KEY `idx_orders_branch_created` (`branch_id`),
  KEY `idx_orders_table` (`table_id`),
  KEY `idx_orders_employee` (`employee_id`),
  KEY `idx_orders_customer` (`customer_id`),
  KEY `idx_orders_shift` (`shift_id`),
  KEY `idx_orders_type` (`order_type`),
  KEY `idx_orders_status` (`status`),
  KEY `idx_orders_created_at` (`created_at` DESC),
  CONSTRAINT `fk_orders_branch` FOREIGN KEY (`branch_id`) REFERENCES `branches` (`id`) ON DELETE RESTRICT,
  CONSTRAINT `fk_orders_customer` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`) ON DELETE SET NULL,
  CONSTRAINT `fk_orders_employee` FOREIGN KEY (`employee_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT,
  CONSTRAINT `fk_orders_shift` FOREIGN KEY (`shift_id`) REFERENCES `shifts` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `payment_splits`
--

DROP TABLE IF EXISTS `payment_splits`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `payment_splits` (
  `id` char(36) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `payment_id` char(36) NOT NULL,
  `payment_method` varchar(20) DEFAULT NULL,
  `amount` double NOT NULL,
  `reference_number` longtext,
  PRIMARY KEY (`id`),
  KEY `idx_payment_splits_deleted_at` (`deleted_at`),
  KEY `idx_split_payment` (`payment_id`),
  KEY `idx_split_method` (`payment_method`),
  CONSTRAINT `fk_payments_splits` FOREIGN KEY (`payment_id`) REFERENCES `payments` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `payments`
--

DROP TABLE IF EXISTS `payments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `payments` (
  `id` char(36) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `order_id` char(36) NOT NULL,
  `change_amount` double DEFAULT NULL,
  `status` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_payments_deleted_at` (`deleted_at`),
  KEY `idx_payment_order` (`order_id`),
  KEY `idx_payment_status` (`status`),
  CONSTRAINT `fk_orders_payments` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `raw_materials`
--

DROP TABLE IF EXISTS `raw_materials`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `raw_materials` (
  `id` char(36) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(191) NOT NULL,
  `unit` longtext NOT NULL,
  `cost_per_unit` double DEFAULT NULL,
  `current_stock` double DEFAULT '0',
  `min_stock_level` double DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_raw_materials_deleted_at` (`deleted_at`),
  KEY `idx_rawmaterial_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `recipe_ingredients`
--

DROP TABLE IF EXISTS `recipe_ingredients`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `recipe_ingredients` (
  `id` char(36) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `recipe_id` char(36) NOT NULL,
  `raw_material_id` char(36) NOT NULL,
  `quantity` double NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_recipe_ingredients_deleted_at` (`deleted_at`),
  KEY `idx_ingredient_recipe` (`recipe_id`),
  KEY `idx_ingredient_material` (`raw_material_id`),
  CONSTRAINT `fk_recipes_ingredients` FOREIGN KEY (`recipe_id`) REFERENCES `recipes` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `recipes`
--

DROP TABLE IF EXISTS `recipes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `recipes` (
  `id` char(36) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `menu_id` char(36) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_recipes_menu_id` (`menu_id`),
  KEY `idx_recipes_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_recipes_menu` FOREIGN KEY (`menu_id`) REFERENCES `menus` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `reservations`
--

DROP TABLE IF EXISTS `reservations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `reservations` (
  `id` char(36) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `branch_id` char(36) NOT NULL,
  `table_id` char(36) NOT NULL,
  `customer_id` char(36) DEFAULT NULL,
  `customer_name` longtext NOT NULL,
  `customer_phone` longtext NOT NULL,
  `reservation_date` datetime(3) DEFAULT NULL,
  `guest_count` bigint NOT NULL,
  `deposit_amount` double DEFAULT NULL,
  `status` varchar(20) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_reservations_deleted_at` (`deleted_at`),
  KEY `idx_reservation_branch` (`branch_id`),
  KEY `idx_reservation_table` (`table_id`),
  KEY `idx_reservation_customer` (`customer_id`),
  KEY `idx_reservation_date` (`reservation_date`),
  KEY `idx_reservation_status` (`status`),
  CONSTRAINT `fk_reservations_branch` FOREIGN KEY (`branch_id`) REFERENCES `branches` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_reservations_customer` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `shifts`
--

DROP TABLE IF EXISTS `shifts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `shifts` (
  `id` char(36) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `branch_id` char(36) NOT NULL,
  `opened_by` char(36) NOT NULL,
  `opening_time` datetime(3) DEFAULT NULL,
  `opening_cash` double DEFAULT NULL,
  `closed_by` char(36) DEFAULT NULL,
  `closing_time` datetime(3) DEFAULT NULL,
  `actual_closing_cash` double DEFAULT NULL,
  `status` varchar(20) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_shifts_deleted_at` (`deleted_at`),
  KEY `idx_shift_branch` (`branch_id`),
  KEY `idx_shift_opener` (`opened_by`),
  KEY `idx_shift_opening` (`opening_time`),
  KEY `idx_shift_closer` (`closed_by`),
  KEY `idx_shift_status` (`status`),
  CONSTRAINT `fk_shifts_branch` FOREIGN KEY (`branch_id`) REFERENCES `branches` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_shifts_closer` FOREIGN KEY (`closed_by`) REFERENCES `users` (`id`) ON DELETE SET NULL,
  CONSTRAINT `fk_shifts_opener` FOREIGN KEY (`opened_by`) REFERENCES `users` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `stock_adjustments`
--

DROP TABLE IF EXISTS `stock_adjustments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `stock_adjustments` (
  `id` char(36) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `raw_material_id` char(36) NOT NULL,
  `quantity_before` double DEFAULT NULL,
  `quantity_after` double NOT NULL,
  `reason` text NOT NULL,
  `adjusted_by` char(36) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_stock_adjustments_deleted_at` (`deleted_at`),
  KEY `idx_adjustment_material` (`raw_material_id`),
  KEY `idx_adjustment_user` (`adjusted_by`),
  CONSTRAINT `fk_adjustments_material` FOREIGN KEY (`raw_material_id`) REFERENCES `raw_materials` (`id`) ON DELETE RESTRICT,
  CONSTRAINT `fk_adjustments_user` FOREIGN KEY (`adjusted_by`) REFERENCES `users` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `tax_configs`
--

DROP TABLE IF EXISTS `tax_configs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `tax_configs` (
  `id` char(36) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `branch_id` char(36) DEFAULT NULL,
  `name` varchar(191) NOT NULL,
  `tax_type` varchar(20) NOT NULL,
  `percentage` double NOT NULL,
  `is_active` tinyint(1) DEFAULT '1',
  `effective_from` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_tax_configs_deleted_at` (`deleted_at`),
  KEY `idx_tax_branch` (`branch_id`),
  KEY `idx_tax_name` (`name`),
  KEY `idx_tax_type` (`tax_type`),
  KEY `idx_tax_active` (`is_active`),
  KEY `idx_tax_effective` (`effective_from`),
  CONSTRAINT `fk_tax_branch` FOREIGN KEY (`branch_id`) REFERENCES `branches` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` char(36) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(191) NOT NULL,
  `email` varchar(191) NOT NULL,
  `phone` varchar(191) DEFAULT NULL,
  `password_hash` longtext NOT NULL,
  `pin` longtext NOT NULL,
  `role` varchar(20) NOT NULL,
  `branch_id` char(36) DEFAULT NULL,
  `is_active` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_email` (`email`),
  KEY `idx_users_deleted_at` (`deleted_at`),
  KEY `idx_user_name` (`name`),
  KEY `idx_user_phone` (`phone`),
  KEY `idx_user_role` (`role`),
  KEY `idx_user_branch` (`branch_id`),
  KEY `idx_user_active` (`is_active`),
  CONSTRAINT `fk_users_branch` FOREIGN KEY (`branch_id`) REFERENCES `branches` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Temporary table structure for view `v_best_selling_items`
--

DROP TABLE IF EXISTS `v_best_selling_items`;
/*!50001 DROP VIEW IF EXISTS `v_best_selling_items`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8mb4;
/*!50001 CREATE VIEW `v_best_selling_items` AS SELECT
 1 AS `menu_id`,
  1 AS `menu_name`,
  1 AS `station`,
  1 AS `total_quantity`,
  1 AS `total_sales`,
  1 AS `order_count`,
  1 AS `avg_price`,
  1 AS `last_sold_at`,
  1 AS `qty_last_7d`,
  1 AS `qty_prev_7d` */;
SET character_set_client = @saved_cs_client;

--
-- Temporary table structure for view `v_cashier_performance`
--

DROP TABLE IF EXISTS `v_cashier_performance`;
/*!50001 DROP VIEW IF EXISTS `v_cashier_performance`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8mb4;
/*!50001 CREATE VIEW `v_cashier_performance` AS SELECT
 1 AS `cashier_id`,
  1 AS `cashier_name`,
  1 AS `branch_id`,
  1 AS `total_orders`,
  1 AS `total_sales`,
  1 AS `avg_order_value`,
  1 AS `active_days`,
  1 AS `last_order_at`,
  1 AS `sales_last_7d`,
  1 AS `sales_last_30d`,
  1 AS `void_rate_pct`,
  1 AS `unique_customers`,
  1 AS `dine_in_orders`,
  1 AS `take_away_orders` */;
SET character_set_client = @saved_cs_client;

--
-- Temporary table structure for view `v_daily_sales_summary`
--

DROP TABLE IF EXISTS `v_daily_sales_summary`;
/*!50001 DROP VIEW IF EXISTS `v_daily_sales_summary`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8mb4;
/*!50001 CREATE VIEW `v_daily_sales_summary` AS SELECT
 1 AS `sale_date`,
  1 AS `branch_id`,
  1 AS `total_orders`,
  1 AS `gross_revenue`,
  1 AS `total_cogs`,
  1 AS `net_profit`,
  1 AS `avg_order_value`,
  1 AS `cancelled_orders` */;
SET character_set_client = @saved_cs_client;

--
-- Temporary table structure for view `v_menu_profitability_live`
--

DROP TABLE IF EXISTS `v_menu_profitability_live`;
/*!50001 DROP VIEW IF EXISTS `v_menu_profitability_live`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8mb4;
/*!50001 CREATE VIEW `v_menu_profitability_live` AS SELECT
 1 AS `menu_id`,
  1 AS `menu_name`,
  1 AS `station`,
  1 AS `total_qty_sold`,
  1 AS `total_revenue`,
  1 AS `total_cogs`,
  1 AS `net_profit`,
  1 AS `profit_margin_pct`,
  1 AS `current_selling_price`,
  1 AS `current_recipe_cost`,
  1 AS `refreshed_at` */;
SET character_set_client = @saved_cs_client;

--
-- Temporary table structure for view `v_payment_method_summary`
--

DROP TABLE IF EXISTS `v_payment_method_summary`;
/*!50001 DROP VIEW IF EXISTS `v_payment_method_summary`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8mb4;
/*!50001 CREATE VIEW `v_payment_method_summary` AS SELECT
 1 AS `payment_date`,
  1 AS `branch_id`,
  1 AS `payment_method`,
  1 AS `total_payments`,
  1 AS `total_amount`,
  1 AS `order_count`,
  1 AS `avg_amount` */;
SET character_set_client = @saved_cs_client;

--
-- Temporary table structure for view `v_shift_reconciliation`
--

DROP TABLE IF EXISTS `v_shift_reconciliation`;
/*!50001 DROP VIEW IF EXISTS `v_shift_reconciliation`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8mb4;
/*!50001 CREATE VIEW `v_shift_reconciliation` AS SELECT
 1 AS `shift_id`,
  1 AS `branch_id`,
  1 AS `cashier_id`,
  1 AS `cashier_name`,
  1 AS `opening_time`,
  1 AS `closing_time`,
  1 AS `opening_cash`,
  1 AS `actual_closing_cash`,
  1 AS `total_cash_in`,
  1 AS `total_cash_out`,
  1 AS `expected_closing_cash`,
  1 AS `cash_difference`,
  1 AS `total_orders`,
  1 AS `total_sales`,
  1 AS `discrepancy`,
  1 AS `status` */;
SET character_set_client = @saved_cs_client;

--
-- Temporary table structure for view `v_void_and_return_logs`
--

DROP TABLE IF EXISTS `v_void_and_return_logs`;
/*!50001 DROP VIEW IF EXISTS `v_void_and_return_logs`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8mb4;
/*!50001 CREATE VIEW `v_void_and_return_logs` AS SELECT
 1 AS `event_date`,
  1 AS `branch_id`,
  1 AS `employee_id`,
  1 AS `employee_name`,
  1 AS `event_type`,
  1 AS `order_id`,
  1 AS `order_number`,
  1 AS `amount`,
  1 AS `reason`,
  1 AS `item_count`,
  1 AS `total_qty_voided` */;
SET character_set_client = @saved_cs_client;

--
-- Temporary table structure for view `v_weekly_sales_summary`
--

DROP TABLE IF EXISTS `v_weekly_sales_summary`;
/*!50001 DROP VIEW IF EXISTS `v_weekly_sales_summary`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8mb4;
/*!50001 CREATE VIEW `v_weekly_sales_summary` AS SELECT
 1 AS `week_start`,
  1 AS `branch_id`,
  1 AS `total_orders`,
  1 AS `gross_revenue`,
  1 AS `total_cogs`,
  1 AS `net_profit`,
  1 AS `active_days`,
  1 AS `last_order_at` */;
SET character_set_client = @saved_cs_client;

--
-- Final view structure for view `v_best_selling_items`
--

/*!50001 DROP VIEW IF EXISTS `v_best_selling_items`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_general_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`root`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `v_best_selling_items` AS select `od`.`menu_id` AS `menu_id`,`od`.`menu_name` AS `menu_name`,`m`.`station` AS `station`,sum((case when (`od`.`is_voided` = false) then `od`.`quantity` else 0 end)) AS `total_quantity`,sum((case when (`od`.`is_voided` = false) then `od`.`subtotal` else 0 end)) AS `total_sales`,count(distinct `o`.`id`) AS `order_count`,avg(`od`.`price`) AS `avg_price`,max(`o`.`created_at`) AS `last_sold_at`,sum((case when ((`od`.`is_voided` = false) and (`o`.`created_at` >= (now() - interval 7 day))) then `od`.`quantity` else 0 end)) AS `qty_last_7d`,sum((case when ((`od`.`is_voided` = false) and (`o`.`created_at` >= (now() - interval 14 day)) and (`o`.`created_at` < (now() - interval 7 day))) then `od`.`quantity` else 0 end)) AS `qty_prev_7d` from ((`order_details` `od` join `orders` `o` on(((`o`.`id` = `od`.`order_id`) and (`o`.`deleted_at` is null)))) left join `menus` `m` on((`m`.`id` = `od`.`menu_id`))) where ((`od`.`deleted_at` is null) and (`o`.`status` in ('confirmed','paid','served'))) group by `od`.`menu_id`,`od`.`menu_name`,`m`.`station` order by `total_quantity` desc */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `v_cashier_performance`
--

/*!50001 DROP VIEW IF EXISTS `v_cashier_performance`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_general_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`root`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `v_cashier_performance` AS select `o`.`employee_id` AS `cashier_id`,`u`.`name` AS `cashier_name`,`u`.`branch_id` AS `branch_id`,count(distinct `o`.`id`) AS `total_orders`,sum((case when (`od`.`is_voided` = false) then `od`.`subtotal` else 0 end)) AS `total_sales`,avg((case when (`od`.`is_voided` = false) then `od`.`subtotal` else NULL end)) AS `avg_order_value`,count(distinct cast(`o`.`created_at` as date)) AS `active_days`,max(`o`.`created_at`) AS `last_order_at`,sum((case when ((`od`.`is_voided` = false) and (`o`.`created_at` >= (now() - interval 7 day))) then `od`.`subtotal` else 0 end)) AS `sales_last_7d`,sum((case when ((`od`.`is_voided` = false) and (`o`.`created_at` >= (now() - interval 30 day))) then `od`.`subtotal` else 0 end)) AS `sales_last_30d`,(case when (count(`od`.`id`) > 0) then (((count((case when (`od`.`is_voided` = true) then 1 end)) * 1.0) / count(`od`.`id`)) * 100) else 0 end) AS `void_rate_pct`,count(distinct (case when (`o`.`customer_id` is not null) then `o`.`customer_id` end)) AS `unique_customers`,count(distinct (case when (`o`.`order_type` = 'dine_in') then `o`.`id` end)) AS `dine_in_orders`,count(distinct (case when (`o`.`order_type` = 'take_away') then `o`.`id` end)) AS `take_away_orders` from ((`orders` `o` join `users` `u` on((`u`.`id` = `o`.`employee_id`))) left join `order_details` `od` on(((`od`.`order_id` = `o`.`id`) and (`od`.`deleted_at` is null)))) where ((`o`.`deleted_at` is null) and (`o`.`status` in ('confirmed','paid','served')) and (`o`.`created_at` >= (now() - interval 90 day))) group by `o`.`employee_id`,`u`.`name`,`u`.`branch_id` order by `total_sales` desc */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `v_daily_sales_summary`
--

/*!50001 DROP VIEW IF EXISTS `v_daily_sales_summary`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_general_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`root`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `v_daily_sales_summary` AS select cast(`o`.`created_at` as date) AS `sale_date`,`o`.`branch_id` AS `branch_id`,count(distinct `o`.`id`) AS `total_orders`,sum((case when (`od`.`is_voided` = false) then `od`.`subtotal` else 0 end)) AS `gross_revenue`,sum((case when (`od`.`is_voided` = false) then (`od`.`cost_price` * `od`.`quantity`) else 0 end)) AS `total_cogs`,(sum((case when (`od`.`is_voided` = false) then `od`.`subtotal` else 0 end)) - sum((case when (`od`.`is_voided` = false) then (`od`.`cost_price` * `od`.`quantity`) else 0 end))) AS `net_profit`,avg((case when (`od`.`is_voided` = false) then `od`.`subtotal` else NULL end)) AS `avg_order_value`,count((case when (`o`.`status` = 'cancelled') then 1 end)) AS `cancelled_orders` from (`orders` `o` left join `order_details` `od` on(((`od`.`order_id` = `o`.`id`) and (`od`.`deleted_at` is null)))) where ((`o`.`deleted_at` is null) and (`o`.`status` in ('confirmed','paid','served'))) group by cast(`o`.`created_at` as date),`o`.`branch_id` order by `sale_date` desc */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `v_menu_profitability_live`
--

/*!50001 DROP VIEW IF EXISTS `v_menu_profitability_live`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_general_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`root`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `v_menu_profitability_live` AS select `od`.`menu_id` AS `menu_id`,`od`.`menu_name` AS `menu_name`,`m`.`station` AS `station`,sum((case when (`od`.`is_voided` = false) then `od`.`quantity` else 0 end)) AS `total_qty_sold`,sum((case when (`od`.`is_voided` = false) then `od`.`subtotal` else 0 end)) AS `total_revenue`,sum((case when (`od`.`is_voided` = false) then (`od`.`cost_price` * `od`.`quantity`) else 0 end)) AS `total_cogs`,(sum((case when (`od`.`is_voided` = false) then `od`.`subtotal` else 0 end)) - sum((case when (`od`.`is_voided` = false) then (`od`.`cost_price` * `od`.`quantity`) else 0 end))) AS `net_profit`,(case when (sum((case when (`od`.`is_voided` = false) then `od`.`subtotal` else 0 end)) > 0) then (((sum((case when (`od`.`is_voided` = false) then `od`.`subtotal` else 0 end)) - sum((case when (`od`.`is_voided` = false) then (`od`.`cost_price` * `od`.`quantity`) else 0 end))) / sum((case when (`od`.`is_voided` = false) then `od`.`subtotal` else 0 end))) * 100) else 0 end) AS `profit_margin_pct`,avg(`od`.`price`) AS `current_selling_price`,(select sum((`rm`.`cost_per_unit` * `ri`.`quantity`)) from ((`recipes` `r` join `recipe_ingredients` `ri` on((`ri`.`recipe_id` = `r`.`id`))) join `raw_materials` `rm` on((`rm`.`id` = `ri`.`raw_material_id`))) where (`r`.`menu_id` = `od`.`menu_id`) limit 1) AS `current_recipe_cost`,now() AS `refreshed_at` from ((`order_details` `od` join `orders` `o` on(((`o`.`id` = `od`.`order_id`) and (`o`.`deleted_at` is null)))) left join `menus` `m` on((`m`.`id` = `od`.`menu_id`))) where ((`od`.`deleted_at` is null) and (`o`.`status` in ('confirmed','paid','served'))) group by `od`.`menu_id`,`od`.`menu_name`,`m`.`station` order by `net_profit` desc */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `v_payment_method_summary`
--

/*!50001 DROP VIEW IF EXISTS `v_payment_method_summary`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_general_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`root`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `v_payment_method_summary` AS select cast(`p`.`created_at` as date) AS `payment_date`,`o`.`branch_id` AS `branch_id`,`ps`.`payment_method` AS `payment_method`,count(distinct `p`.`id`) AS `total_payments`,sum(`ps`.`amount`) AS `total_amount`,count(distinct `o`.`id`) AS `order_count`,avg(`ps`.`amount`) AS `avg_amount` from ((`payments` `p` join `payment_splits` `ps` on(((`ps`.`payment_id` = `p`.`id`) and (`ps`.`deleted_at` is null)))) join `orders` `o` on((`o`.`id` = `p`.`order_id`))) where ((`p`.`deleted_at` is null) and (`p`.`status` = 'success') and (`o`.`deleted_at` is null)) group by cast(`p`.`created_at` as date),`o`.`branch_id`,`ps`.`payment_method` order by `payment_date` desc,`total_amount` desc */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `v_shift_reconciliation`
--

/*!50001 DROP VIEW IF EXISTS `v_shift_reconciliation`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_general_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`root`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `v_shift_reconciliation` AS select `s`.`id` AS `shift_id`,`s`.`branch_id` AS `branch_id`,`s`.`opened_by` AS `cashier_id`,`u`.`name` AS `cashier_name`,`s`.`opening_time` AS `opening_time`,`s`.`closing_time` AS `closing_time`,`s`.`opening_cash` AS `opening_cash`,`s`.`actual_closing_cash` AS `actual_closing_cash`,coalesce(sum((case when (`cm`.`type` = 'in') then `cm`.`amount` else 0 end)),0) AS `total_cash_in`,coalesce(sum((case when (`cm`.`type` = 'out') then `cm`.`amount` else 0 end)),0) AS `total_cash_out`,(`s`.`opening_cash` + coalesce(sum((case when (`cm`.`type` = 'in') then `cm`.`amount` else -(`cm`.`amount`) end)),0)) AS `expected_closing_cash`,(`s`.`actual_closing_cash` - (`s`.`opening_cash` + coalesce(sum((case when (`cm`.`type` = 'in') then `cm`.`amount` else -(`cm`.`amount`) end)),0))) AS `cash_difference`,count(distinct `o`.`id`) AS `total_orders`,coalesce(sum(`p`.`total_paid`),0) AS `total_sales`,(case when (`s`.`actual_closing_cash` > 0) then (`s`.`actual_closing_cash` - (`s`.`opening_cash` + coalesce(sum((case when (`cm`.`type` = 'in') then `cm`.`amount` else -(`cm`.`amount`) end)),0))) else 0 end) AS `discrepancy`,`s`.`status` AS `status` from ((((`shifts` `s` join `users` `u` on((`u`.`id` = `s`.`opened_by`))) left join `cash_movements` `cm` on(((`cm`.`shift_id` = `s`.`id`) and (`cm`.`deleted_at` is null)))) left join `orders` `o` on(((`o`.`shift_id` = `s`.`id`) and (`o`.`deleted_at` is null) and (`o`.`status` in ('confirmed','paid','served'))))) left join (select `p`.`order_id` AS `order_id`,sum(`ps`.`amount`) AS `total_paid` from (`payments` `p` join `payment_splits` `ps` on(((`ps`.`payment_id` = `p`.`id`) and (`ps`.`deleted_at` is null)))) where ((`p`.`deleted_at` is null) and (`p`.`status` = 'success')) group by `p`.`order_id`) `p` on((`p`.`order_id` = `o`.`id`))) where (`s`.`deleted_at` is null) group by `s`.`id`,`s`.`branch_id`,`s`.`opened_by`,`u`.`name`,`s`.`opening_time`,`s`.`closing_time`,`s`.`opening_cash`,`s`.`actual_closing_cash`,`s`.`status` order by `s`.`opening_time` desc */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `v_void_and_return_logs`
--

/*!50001 DROP VIEW IF EXISTS `v_void_and_return_logs`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_general_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`root`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `v_void_and_return_logs` AS select cast(`o`.`updated_at` as date) AS `event_date`,`o`.`branch_id` AS `branch_id`,`o`.`employee_id` AS `employee_id`,`u`.`name` AS `employee_name`,'void_order' AS `event_type`,`o`.`id` AS `order_id`,`o`.`order_number` AS `order_number`,`o`.`total` AS `amount`,`o`.`notes` AS `reason`,count(`od`.`id`) AS `item_count`,sum(`od`.`quantity`) AS `total_qty_voided` from ((`orders` `o` join `users` `u` on((`u`.`id` = `o`.`employee_id`))) join `order_details` `od` on(((`od`.`order_id` = `o`.`id`) and (`od`.`deleted_at` is null) and (`od`.`is_voided` = true)))) where ((`o`.`deleted_at` is null) and (`o`.`status` = 'cancelled') and (`o`.`updated_at` >= (now() - interval 30 day))) group by cast(`o`.`updated_at` as date),`o`.`branch_id`,`o`.`employee_id`,`u`.`name`,`o`.`id`,`o`.`order_number`,`o`.`total`,`o`.`notes` union all select cast(`r`.`processed_at` as date) AS `event_date`,`o`.`branch_id` AS `branch_id`,`r`.`processed_by` AS `employee_id`,`u`.`name` AS `employee_name`,'return' AS `event_type`,`o`.`id` AS `order_id`,`o`.`order_number` AS `order_number`,`r`.`return_amount` AS `amount`,`r`.`reason` AS `reason`,0 AS `item_count`,0 AS `total_qty_voided` from ((`order_returns` `r` join `orders` `o` on((`o`.`id` = `r`.`order_id`))) join `users` `u` on((`u`.`id` = `r`.`processed_by`))) where ((`r`.`deleted_at` is null) and (`r`.`processed_at` >= (now() - interval 30 day))) order by `event_date` desc,`amount` desc */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `v_weekly_sales_summary`
--

/*!50001 DROP VIEW IF EXISTS `v_weekly_sales_summary`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_general_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`root`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `v_weekly_sales_summary` AS select (cast(`o`.`created_at` as date) - interval weekday(`o`.`created_at`) day) AS `week_start`,`o`.`branch_id` AS `branch_id`,count(distinct `o`.`id`) AS `total_orders`,sum((case when (`od`.`is_voided` = false) then `od`.`subtotal` else 0 end)) AS `gross_revenue`,sum((case when (`od`.`is_voided` = false) then (`od`.`cost_price` * `od`.`quantity`) else 0 end)) AS `total_cogs`,(sum((case when (`od`.`is_voided` = false) then `od`.`subtotal` else 0 end)) - sum((case when (`od`.`is_voided` = false) then (`od`.`cost_price` * `od`.`quantity`) else 0 end))) AS `net_profit`,count(distinct cast(`o`.`created_at` as date)) AS `active_days`,max(`o`.`created_at`) AS `last_order_at` from (`orders` `o` left join `order_details` `od` on(((`od`.`order_id` = `o`.`id`) and (`od`.`deleted_at` is null)))) where ((`o`.`deleted_at` is null) and (`o`.`status` in ('confirmed','paid','served'))) group by (cast(`o`.`created_at` as date) - interval weekday(`o`.`created_at`) day),`o`.`branch_id` order by `week_start` desc */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*M!100616 SET NOTE_VERBOSITY=@OLD_NOTE_VERBOSITY */;

-- Dump completed on 2026-05-20 16:15:15
