-- This is DB init code

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    tg_id BIGINT UNIQUE,
    name VARCHAR(255),
    username VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_premium BOOLEAN DEFAULT FALSE
);

-- Create an index on tg_id for faster lookups since it will be asked a lot
CREATE INDEX IF NOT EXISTS idx_users_tg_id ON users(tg_id);

-- Create a function to update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create a trigger to automatically update the updated_at column
CREATE OR REPLACE TRIGGER update_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Create conversations table
CREATE TABLE IF NOT EXISTS conversations (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    finished BOOLEAN DEFAULT FALSE
);

-- Create an index on user_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_conversations_user_id ON conversations(user_id);

-- Create a trigger to automatically update the updated_at column for conversations
CREATE OR REPLACE TRIGGER update_conversations_updated_at
BEFORE UPDATE ON conversations
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Create message_queue table
CREATE TABLE IF NOT EXISTS message_queue (
    id SERIAL PRIMARY KEY,
    conversation_id INTEGER REFERENCES conversations(id) ON DELETE CASCADE,
    question TEXT NOT NULL,
    answer TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) CHECK (status IN ('new', 'processing', 'ready', 'answered')) DEFAULT 'new',
    tg_question_id BIGINT UNIQUE,
    tg_answer_id BIGINT
);

-- Create an index on conversation_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_message_queue_conversation_id ON message_queue(conversation_id);

-- Create a trigger to automatically update the updated_at column for message_queue
CREATE OR REPLACE TRIGGER update_message_queue_updated_at
BEFORE UPDATE ON message_queue
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();