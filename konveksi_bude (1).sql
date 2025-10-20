-- phpMyAdmin SQL Dump
-- version 5.2.0
-- https://www.phpmyadmin.net/
--
-- Host: localhost:3306
-- Generation Time: Oct 20, 2025 at 08:12 AM
-- Server version: 8.0.30
-- PHP Version: 8.3.8

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `konveksi_bude`
--

-- --------------------------------------------------------

--
-- Table structure for table `customers`
--

CREATE TABLE `customers` (
  `id` int NOT NULL,
  `name` varchar(100) NOT NULL,
  `type` enum('TK','SD','SMP','Kelompok Tadarus','Lainnya') NOT NULL,
  `contact` varchar(100) DEFAULT NULL,
  `address` text,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `customers`
--

INSERT INTO `customers` (`id`, `name`, `type`, `contact`, `address`, `created_at`) VALUES
(1, 'sd muhammadiyah', 'SD', '0315754234', 'surabaya', '2025-05-30 03:30:30'),
(2, 'Aisyiyah', 'Lainnya', '086543446122', 'sidoarjo', '2025-05-30 03:35:11'),
(3, 'gahat', 'SD', '082387664533', 'yolo', '2025-06-05 08:02:42'),
(4, 'q', 'SD', '1234567890', 'ssqq', '2025-06-05 08:17:51');

-- --------------------------------------------------------

--
-- Table structure for table `customer_uniforms`
--

CREATE TABLE `customer_uniforms` (
  `id` int NOT NULL,
  `customer_id` int NOT NULL,
  `uniform_name` varchar(100) NOT NULL,
  `size` varchar(20) NOT NULL,
  `price` decimal(10,2) NOT NULL,
  `notes` text,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `customer_uniforms`
--

INSERT INTO `customer_uniforms` (`id`, `customer_id`, `uniform_name`, `size`, `price`, `notes`, `created_at`) VALUES
(1, 1, 'hw', 'L', '150000.00', '', '2025-05-30 03:31:09'),
(2, 1, 'olahraga', 'all size', '300000.00', 'cowo cewe', '2025-05-30 03:31:37'),
(3, 2, 'batik', 'L', '50000.00', 'nambah 100 harganya', '2025-05-30 03:35:11'),
(4, 2, 'batik', 'XL', '220000.00', '', '2025-05-30 03:35:11'),
(5, 2, 'PDH ', 'all size', '250000.00', '', '2025-05-30 03:35:11'),
(6, 2, 'taqwa', 'all size', '200000.00', '', '2025-05-30 03:35:11'),
(7, 3, 'hw', 'M', '450000.00', '', '2025-06-05 08:02:42'),
(8, 3, 'we', 'XL', '125000.00', '', '2025-06-05 08:02:42');

-- --------------------------------------------------------

--
-- Table structure for table `customer_uniform_price_history`
--

CREATE TABLE `customer_uniform_price_history` (
  `id` int NOT NULL,
  `customer_uniform_id` int NOT NULL,
  `old_price` decimal(10,2) NOT NULL,
  `changed_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `customer_uniform_price_history`
--

INSERT INTO `customer_uniform_price_history` (`id`, `customer_uniform_id`, `old_price`, `changed_at`) VALUES
(1, 1, '200000.00', '2025-05-30 03:32:03'),
(2, 1, '250000.00', '2025-05-30 03:32:23'),
(3, 3, '150000.00', '2025-06-01 13:49:51'),
(4, 3, '250000.00', '2025-06-01 13:57:08'),
(5, 8, '120000.00', '2025-06-05 17:51:41');

-- --------------------------------------------------------

--
-- Table structure for table `order_items`
--

CREATE TABLE `order_items` (
  `id` int NOT NULL,
  `transaction_id` int NOT NULL,
  `uniform_name` varchar(100) NOT NULL,
  `size` varchar(20) NOT NULL,
  `quantity` int NOT NULL,
  `unit_price` decimal(10,2) NOT NULL,
  `subtotal` decimal(12,2) GENERATED ALWAYS AS ((`quantity` * `unit_price`)) STORED,
  `notes` text
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `order_items`
--

INSERT INTO `order_items` (`id`, `transaction_id`, `uniform_name`, `size`, `quantity`, `unit_price`, `notes`) VALUES
(12, 1, 'batik', 'XL', 20, '220000.00', 'wdwd'),
(13, 1, 'batik', 'L', 15, '50000.00', ''),
(14, 1, 'taqwa', 'all size', 20, '200000.00', ''),
(15, 6, 'PDH ', 'all size', 13, '250000.00', ''),
(16, 6, 'taqwa', 'all size', 5, '200000.00', '');

-- --------------------------------------------------------

--
-- Table structure for table `student_order_items`
--

CREATE TABLE `student_order_items` (
  `id` int NOT NULL,
  `customer_id` int NOT NULL,
  `student_name` varchar(100) NOT NULL,
  `grade` varchar(10) DEFAULT NULL,
  `transaction_id` int NOT NULL,
  `uniform_name` varchar(100) NOT NULL,
  `size` varchar(20) NOT NULL,
  `quantity` int NOT NULL DEFAULT '1',
  `unit_price` decimal(10,2) NOT NULL,
  `subtotal` decimal(12,2) GENERATED ALWAYS AS ((`quantity` * `unit_price`)) STORED,
  `notes` text,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `student_order_items`
--

INSERT INTO `student_order_items` (`id`, `customer_id`, `student_name`, `grade`, `transaction_id`, `uniform_name`, `size`, `quantity`, `unit_price`, `notes`, `created_at`) VALUES
(1, 1, 'hamid', '3', 2, 'hw', 'L', 1, '150000.00', '', '2025-05-30 03:40:24'),
(2, 1, 'hamid', '3', 2, 'olahraga', 'all size', 1, '300000.00', '', '2025-05-30 03:40:24'),
(3, 1, 'zuhdi', '4', 2, 'hw', 'L', 1, '150000.00', '', '2025-05-30 03:40:24'),
(4, 1, 'wijanarko', '3', 2, 'olahraga', 'all size', 1, '300000.00', '', '2025-05-30 03:40:24'),
(5, 1, 'levina', '4', 2, 'olahraga', 'all size', 1, '300000.00', '', '2025-05-30 03:40:24'),
(6, 1, 'anjani', '5', 2, 'hw', 'L', 1, '150000.00', '', '2025-05-30 03:40:24'),
(7, 1, 'tobi', '4', 2, 'hw', 'L', 1, '150000.00', '', '2025-05-30 03:40:24'),
(8, 1, 'samson', '2', 2, 'hw', 'L', 1, '150000.00', '', '2025-05-30 03:40:24'),
(9, 1, 'pity', '3', 2, 'olahraga', 'all size', 1, '300000.00', '', '2025-05-30 03:40:24'),
(10, 1, 'boi', '4', 2, 'olahraga', 'all size', 1, '300000.00', '', '2025-05-30 03:40:24'),
(11, 1, 'thamuz', '3', 2, 'hw', 'L', 1, '150000.00', '', '2025-05-30 03:40:24'),
(12, 1, 'thamuz', '3', 2, 'olahraga', 'all size', 1, '300000.00', '', '2025-05-30 03:40:24'),
(13, 3, 'hamd', '3', 3, 'hw', 'M', 1, '450000.00', '', '2025-06-07 03:08:10'),
(14, 2, 'ss', '3', 7, 'batik', 'XL', 14, '220000.00', '', '2025-06-07 15:23:46'),
(15, 2, 'rf', '3', 7, 'taqwa', 'all size', 8, '200000.00', '', '2025-06-07 15:23:46');

-- --------------------------------------------------------

--
-- Table structure for table `transactions`
--

CREATE TABLE `transactions` (
  `id` int NOT NULL,
  `customer_id` int NOT NULL,
  `transaction_date` date NOT NULL,
  `payment_date` date DEFAULT NULL,
  `status` enum('pending','paid','cancelled') DEFAULT 'pending',
  `total_price` decimal(12,2) DEFAULT NULL,
  `notes` text,
  `nama_murid` varchar(30) DEFAULT NULL,
  `kelas` varchar(10) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `transactions`
--

INSERT INTO `transactions` (`id`, `customer_id`, `transaction_date`, `payment_date`, `status`, `total_price`, `notes`, `nama_murid`, `kelas`, `created_at`, `updated_at`) VALUES
(1, 2, '2025-06-07', '2025-06-01', 'paid', '9150000.00', '', NULL, NULL, '2025-05-30 03:36:53', '2025-06-06 08:17:22'),
(2, 1, '2025-05-31', '2025-06-04', 'pending', '2700000.00', '', NULL, NULL, '2025-05-30 03:40:24', '2025-05-30 03:40:24'),
(3, 3, '2025-07-05', '2025-07-05', 'pending', '450000.00', '', NULL, NULL, '2025-06-07 03:08:10', '2025-06-07 03:08:10'),
(6, 2, '2025-07-10', '2025-07-10', 'pending', '4250000.00', '', NULL, NULL, '2025-06-07 08:41:28', '2025-06-07 08:41:28'),
(7, 2, '2025-07-11', '2025-07-09', 'paid', '4680000.00', 'sdeeea', NULL, NULL, '2025-06-07 15:23:46', '2025-06-07 16:50:30');

-- --------------------------------------------------------

--
-- Table structure for table `users`
--

CREATE TABLE `users` (
  `id` int NOT NULL,
  `username` varchar(100) NOT NULL,
  `password` varchar(12) NOT NULL,
  `contact` varchar(100) DEFAULT NULL,
  `address` text,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `users`
--

INSERT INTO `users` (`id`, `username`, `password`, `contact`, `address`, `created_at`) VALUES
(1, 'admin', 'haha123', '083423667544', 'indonesia', '2025-05-30 03:29:49'),
(2, 'admin', 'admin123', '081234567890', 'Jl. Contoh No. 123', '2025-06-13 17:24:20'),
(3, 'operator', 'operator123', '081234567891', 'Jl. Operator No. 456', '2025-06-13 17:24:20'),
(4, 'manager', 'manager123', '081234567892', 'Jl. Manager No. 789', '2025-06-13 17:24:20');

-- --------------------------------------------------------

--
-- Stand-in structure for view `view_prices`
-- (See below for the actual view)
--
CREATE TABLE `view_prices` (
`id` int
,`customer_name` varchar(100)
,`uniform_name` varchar(100)
,`size` varchar(20)
,`price` decimal(10,2)
,`notes` text
);

-- --------------------------------------------------------

--
-- Structure for view `view_prices`
--
DROP TABLE IF EXISTS `view_prices`;

CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`localhost` SQL SECURITY DEFINER VIEW `view_prices`  AS SELECT `cu`.`id` AS `id`, `c`.`name` AS `customer_name`, `cu`.`uniform_name` AS `uniform_name`, `cu`.`size` AS `size`, `cu`.`price` AS `price`, `cu`.`notes` AS `notes` FROM (`customer_uniforms` `cu` join `customers` `c` on((`cu`.`customer_id` = `c`.`id`)))  ;

--
-- Indexes for dumped tables
--

--
-- Indexes for table `customers`
--
ALTER TABLE `customers`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `customer_uniforms`
--
ALTER TABLE `customer_uniforms`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `uniq_cust_uniform_size` (`customer_id`,`uniform_name`,`size`);

--
-- Indexes for table `customer_uniform_price_history`
--
ALTER TABLE `customer_uniform_price_history`
  ADD PRIMARY KEY (`id`),
  ADD KEY `customer_uniform_id` (`customer_uniform_id`);

--
-- Indexes for table `order_items`
--
ALTER TABLE `order_items`
  ADD PRIMARY KEY (`id`),
  ADD KEY `transaction_id` (`transaction_id`);

--
-- Indexes for table `student_order_items`
--
ALTER TABLE `student_order_items`
  ADD PRIMARY KEY (`id`),
  ADD KEY `customer_id` (`customer_id`),
  ADD KEY `transaction_id` (`transaction_id`);

--
-- Indexes for table `transactions`
--
ALTER TABLE `transactions`
  ADD PRIMARY KEY (`id`),
  ADD KEY `customer_id` (`customer_id`);

--
-- Indexes for table `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`id`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `customers`
--
ALTER TABLE `customers`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=5;

--
-- AUTO_INCREMENT for table `customer_uniforms`
--
ALTER TABLE `customer_uniforms`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=9;

--
-- AUTO_INCREMENT for table `customer_uniform_price_history`
--
ALTER TABLE `customer_uniform_price_history`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=6;

--
-- AUTO_INCREMENT for table `order_items`
--
ALTER TABLE `order_items`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=17;

--
-- AUTO_INCREMENT for table `student_order_items`
--
ALTER TABLE `student_order_items`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=16;

--
-- AUTO_INCREMENT for table `transactions`
--
ALTER TABLE `transactions`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=8;

--
-- AUTO_INCREMENT for table `users`
--
ALTER TABLE `users`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=5;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `customer_uniforms`
--
ALTER TABLE `customer_uniforms`
  ADD CONSTRAINT `customer_uniforms_ibfk_1` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`);

--
-- Constraints for table `customer_uniform_price_history`
--
ALTER TABLE `customer_uniform_price_history`
  ADD CONSTRAINT `customer_uniform_price_history_ibfk_1` FOREIGN KEY (`customer_uniform_id`) REFERENCES `customer_uniforms` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `order_items`
--
ALTER TABLE `order_items`
  ADD CONSTRAINT `order_items_ibfk_1` FOREIGN KEY (`transaction_id`) REFERENCES `transactions` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `student_order_items`
--
ALTER TABLE `student_order_items`
  ADD CONSTRAINT `student_order_items_ibfk_1` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `student_order_items_ibfk_2` FOREIGN KEY (`transaction_id`) REFERENCES `transactions` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `transactions`
--
ALTER TABLE `transactions`
  ADD CONSTRAINT `transactions_ibfk_1` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
