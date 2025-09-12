-- +goose Up
CREATE TABLE attendance_histories (
    id               BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT,
    employee_id      VARCHAR(50)      NOT NULL COLLATE utf8mb4_unicode_ci,
    attendance_id    VARCHAR(100)     NOT NULL COLLATE utf8mb4_unicode_ci,
    date_attendance  DATETIME         NOT NULL,
    attendance_type  TINYINT UNSIGNED NOT NULL COMMENT '1=In, 2=Out',
    description      TEXT,
    created_at       DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),

    KEY idx_histories_employee_id (employee_id),
    KEY idx_histories_attendance_id (attendance_id),
    KEY idx_histories_attendance_date (attendance_id, date_attendance),

    CONSTRAINT fk_history_employee
        FOREIGN KEY (employee_id)   REFERENCES employees(employee_id)
        ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_history_attendance
        FOREIGN KEY (attendance_id) REFERENCES attendances(attendance_id)
        ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +goose Down
DROP TABLE IF EXISTS attendance_histories;
