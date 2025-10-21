-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
    -- SCHEMA: ERP/Retail System
-- Table: companies


-- جدول HR Profiles
CREATE TABLE hr_profiles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    email VARCHAR(255) UNIQUE,
    password_hash TEXT NOT NULL,
    image VARCHAR(100) NOT NULL,
    company_name VARCHAR(100),
    job_position VARCHAR(100),
    rate REAL,
    total_rates_count INT DEFAULT 0,
    verified_profile BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- جدول HR Experience
CREATE TABLE hr_experience (
    id SERIAL PRIMARY KEY,
    hr_profile_id INT NOT NULL REFERENCES hr_profiles(id) ON DELETE CASCADE,
    name VARCHAR(255),
    start_date DATE,
    end_date DATE,
    job_position VARCHAR(255)
);
CREATE INDEX idx_hr_experience_hr_profile_id ON hr_experience(hr_profile_id);

-- جدول HR Job Roles
CREATE TABLE hr_job_roles (
    id SERIAL PRIMARY KEY,
    hr_profile_id INT NOT NULL REFERENCES hr_profiles(id) ON DELETE CASCADE,
    name VARCHAR(255),
    role_description TEXT,
    start_date DATE,
    done_rate INT,
    visible BOOLEAN DEFAULT TRUE
);
CREATE INDEX idx_hr_job_roles_hr_profile_id ON hr_job_roles(hr_profile_id);

-- جدول Employees
CREATE TABLE employees (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    job_field VARCHAR(255) NOT NULL,
    image VARCHAR(100) NOT NULL,
    address TEXT ,
    city VARCHAR(100) ,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    active_points BOOLEAN DEFAULT FALSE,
    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- جدول Badges
CREATE TABLE badges (
    id SERIAL PRIMARY KEY,
    hr_profile_id INT NOT NULL REFERENCES hr_profiles(id) ON DELETE CASCADE,
    created_date DATE NOT NULL,
    total_rates_number INT NOT NULL,
    rate REAL NOT NULL,
    job_position VARCHAR(255),
    current_job_roles TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
CREATE INDEX idx_badges_hr_profile_id ON badges(hr_profile_id);

-- جدول Rates
CREATE TABLE rates (
    id SERIAL PRIMARY KEY,
    hr_profile_id INT NOT NULL REFERENCES hr_profiles(id) ON DELETE CASCADE,
    employee_id INT NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    review_text TEXT NOT NULL CHECK (char_length(review_text) BETWEEN 5 AND 2000),
    rate_value REAL NOT NULL CHECK (rate_value >= 0 AND rate_value <= 5),
    rating_context TEXT,
    likes_count INT DEFAULT 0,
    is_verified BOOLEAN DEFAULT FALSE,
    hr_response TEXT,
    is_anonymous BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
CREATE INDEX idx_rates_hr_profile_id ON rates(hr_profile_id);
CREATE INDEX idx_rates_employee_id ON rates(employee_id);

-- جدول Rate Likes
CREATE TABLE rate_likes (
    id SERIAL PRIMARY KEY,
    rate_id INT NOT NULL REFERENCES rates(id) ON DELETE CASCADE,
    employee_id INT NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    is_like BOOLEAN NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
CREATE UNIQUE INDEX ux_rate_likes_rate_employee ON rate_likes(rate_id, employee_id);

-- جدول Badge Likes
CREATE TABLE badge_likes (
    id SERIAL PRIMARY KEY,
    badge_id INT NOT NULL REFERENCES badges(id) ON DELETE CASCADE,
    employee_id INT NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    is_like BOOLEAN NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
CREATE UNIQUE INDEX ux_badge_likes_badge_employee ON badge_likes(badge_id, employee_id);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop tables in reverse order of dependencies

-- DROP TABLE IF EXISTS applied_discounts CASCADE;
-- DROP TABLE IF EXISTS promotion_rules CASCADE;
-- DROP TABLE IF EXISTS promotions CASCADE;
-- DROP TABLE IF EXISTS customer_segments CASCADE;


-- DROP TABLE IF EXISTS credit_note_items;
-- DROP TABLE IF EXISTS credit_notes;
-- DROP TABLE IF EXISTS invoice_services;
-- DROP TABLE IF EXISTS invoice_items;
-- DROP TABLE IF EXISTS invoice_stages;
-- DROP TABLE IF EXISTS payments;
-- DROP TABLE IF EXISTS invoices;
-- DROP TABLE IF EXISTS stock_movements;
-- DROP TABLE IF EXISTS products;
-- DROP TABLE IF EXISTS categories;
-- DROP TABLE IF EXISTS employees;
-- DROP TABLE IF EXISTS users;
-- DROP TABLE IF EXISTS branches;
-- DROP TABLE IF EXISTS companies;
-- DROP TABLE IF EXISTS journal_entry_lines;
-- DROP TABLE IF EXISTS journal_entries;
-- DROP TABLE IF EXISTS accounts;
-- +goose StatementEnd
