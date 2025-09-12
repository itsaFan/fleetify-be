-- +goose Up
CREATE TABLE attendances (
    id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    employee_id     VARCHAR(50)     NOT NULL COLLATE utf8mb4_unicode_ci,
    attendance_id   VARCHAR(100)    NOT NULL COLLATE utf8mb4_unicode_ci,
    clock_in        DATETIME        NULL,
    clock_out       DATETIME        NULL,
    created_at      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uq_attendance_id (attendance_id),
    KEY idx_employee_id (employee_id),
    CONSTRAINT fk_attendance_employee
        FOREIGN KEY (employee_id) REFERENCES employees(employee_id)
        ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +goose Down
DROP TABLE IF EXISTS attendances;
