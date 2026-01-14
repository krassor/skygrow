-- Migration: Promo codes and test prices schema
-- Version: 002
-- Description: Create tables for test prices and promo codes

SET search_path TO profreport, public;

-- Create test_prices table - Стоимость каждого типа теста
CREATE TABLE IF NOT EXISTS profreport.test_prices (
    id SERIAL PRIMARY KEY,
    questionnaire_type VARCHAR(50) NOT NULL UNIQUE,
    price INTEGER NOT NULL DEFAULT 0,
    currency VARCHAR(10) NOT NULL DEFAULT 'RUB',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create promo_codes table - Промокоды
CREATE TABLE IF NOT EXISTS profreport.promo_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    code VARCHAR(100) NOT NULL UNIQUE,
    questionnaire_type VARCHAR(50) NOT NULL,
    final_price INTEGER NOT NULL DEFAULT 0,
    currency VARCHAR(10) NOT NULL DEFAULT 'RUB',
    expires_at TIMESTAMP NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_promo_codes_code ON profreport.promo_codes (code);

CREATE INDEX IF NOT EXISTS idx_promo_codes_questionnaire_type ON profreport.promo_codes (questionnaire_type);

CREATE INDEX IF NOT EXISTS idx_promo_codes_expires_at ON profreport.promo_codes (expires_at);

-- Add triggers for updated_at
CREATE TRIGGER update_test_prices_updated_at
    BEFORE UPDATE ON profreport.test_prices
    FOR EACH ROW
    EXECUTE FUNCTION profreport.update_updated_at_column();

CREATE TRIGGER update_promo_codes_updated_at
    BEFORE UPDATE ON profreport.promo_codes
    FOR EACH ROW
    EXECUTE FUNCTION profreport.update_updated_at_column();

-- Insert default test prices
INSERT INTO
    profreport.test_prices (
        questionnaire_type,
        price,
        currency
    )
VALUES ('ADULT', 890, 'RUB'),
    ('SCHOOLCHILD', 490, 'RUB')
ON CONFLICT (questionnaire_type) DO NOTHING;