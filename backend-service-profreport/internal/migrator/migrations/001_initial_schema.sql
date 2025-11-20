-- Migration: Initial schema for users and questionnaires
-- Version: 001
-- Description: Create profreport schema and tables with auto-managed timestamps

-- Create schema if not exists
CREATE SCHEMA IF NOT EXISTS profreport;

-- Set search path to use profreport schema
SET search_path TO profreport, public;

-- Create users table
CREATE TABLE IF NOT EXISTS profreport.users (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create index on email for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_email ON profreport.users(email);

-- Create questionnaires table
CREATE TABLE IF NOT EXISTS profreport.questionnaires (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    payment_id UUID NOT NULL,
    payment_success BOOLEAN NOT NULL DEFAULT FALSE,
    questionnaire_type VARCHAR(50) NOT NULL DEFAULT 'ADULT',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_questionnaires_user FOREIGN KEY (user_id) REFERENCES profreport.users(id) ON DELETE CASCADE
);

-- Create indexes for faster lookups
CREATE INDEX IF NOT EXISTS idx_questionnaires_user_id ON profreport.questionnaires(user_id);
CREATE INDEX IF NOT EXISTS idx_questionnaires_payment_id ON profreport.questionnaires(payment_id);
CREATE INDEX IF NOT EXISTS idx_questionnaires_type ON profreport.questionnaires(questionnaire_type);

-- Create trigger function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION profreport.update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger for users table
DROP TRIGGER IF EXISTS update_users_updated_at ON profreport.users;
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON profreport.users
    FOR EACH ROW
    EXECUTE FUNCTION profreport.update_updated_at_column();

-- Create trigger for questionnaires table
DROP TRIGGER IF EXISTS update_questionnaires_updated_at ON profreport.questionnaires;
CREATE TRIGGER update_questionnaires_updated_at
    BEFORE UPDATE ON profreport.questionnaires
    FOR EACH ROW
    EXECUTE FUNCTION profreport.update_updated_at_column();
