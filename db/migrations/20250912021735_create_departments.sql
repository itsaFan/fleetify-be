-- /20250912_create_departments.sql

-- +goose Up
CREATE TABLE departments (
  id                 BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  department_name    VARCHAR(255)    NOT NULL,
  max_clock_in_time  TIME            NOT NULL,
  max_clock_out_time TIME            NOT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY ux_departments_name (department_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +goose Down
DROP TABLE IF EXISTS departments;
