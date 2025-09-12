-- +goose Up
CREATE TABLE employees (
  id             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  employee_id    VARCHAR(50)     NOT NULL,
  department_id  BIGINT UNSIGNED NOT NULL,
  name           VARCHAR(255)    NOT NULL,
  address        TEXT            NULL,
  created_at     DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at     DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY ux_employees_employee_id (employee_id),
  KEY ix_employees_department_id (department_id),
  CONSTRAINT fk_employees_department
    FOREIGN KEY (department_id) REFERENCES departments(id)
    ON UPDATE RESTRICT
    ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +goose Down
DROP TABLE IF EXISTS employees;
