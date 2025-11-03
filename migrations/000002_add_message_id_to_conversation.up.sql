-- Add message_id column to conversation_history table
ALTER TABLE conversation_history ADD COLUMN message_id VARCHAR(255);