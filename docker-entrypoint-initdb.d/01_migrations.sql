-- Создание таблицы пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    avatar_url VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы персонажей
CREATE TABLE IF NOT EXISTS characters (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    race VARCHAR(50) NOT NULL,
    class VARCHAR(50) NOT NULL,
    level INTEGER NOT NULL CHECK (level >= 1 AND level <= 20),
    alignment VARCHAR(50),
    background VARCHAR(100),
    player_name VARCHAR(100),
    experience INTEGER DEFAULT 0 CHECK (experience >= 0),
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Основные характеристики
    strength INTEGER NOT NULL CHECK (strength >= 1 AND strength <= 30),
    dexterity INTEGER NOT NULL CHECK (dexterity >= 1 AND dexterity <= 30),
    constitution INTEGER NOT NULL CHECK (constitution >= 1 AND constitution <= 30),
    intelligence INTEGER NOT NULL CHECK (intelligence >= 1 AND intelligence <= 30),
    wisdom INTEGER NOT NULL CHECK (wisdom >= 1 AND wisdom <= 30),
    charisma INTEGER NOT NULL CHECK (charisma >= 1 AND charisma <= 30),
    
    -- Модификаторы характеристик
    strength_mod INTEGER DEFAULT 0,
    dexterity_mod INTEGER DEFAULT 0,
    constitution_mod INTEGER DEFAULT 0,
    intelligence_mod INTEGER DEFAULT 0,
    wisdom_mod INTEGER DEFAULT 0,
    charisma_mod INTEGER DEFAULT 0,
    
    -- Боевые характеристики
    proficiency_bonus INTEGER DEFAULT 2,
    initiative INTEGER DEFAULT 0,
    armor_class INTEGER DEFAULT 10 CHECK (armor_class >= 0),
    speed INTEGER DEFAULT 30 CHECK (speed >= 0),
    hit_points INTEGER DEFAULT 0 CHECK (hit_points >= 0),
    max_hit_points INTEGER DEFAULT 0 CHECK (max_hit_points >= 0),
    temp_hit_points INTEGER DEFAULT 0 CHECK (temp_hit_points >= 0),
    hit_dice VARCHAR(20),
    
    -- Спасброски
    strength_save BOOLEAN DEFAULT FALSE,
    dexterity_save BOOLEAN DEFAULT FALSE,
    constitution_save BOOLEAN DEFAULT FALSE,
    intelligence_save BOOLEAN DEFAULT FALSE,
    wisdom_save BOOLEAN DEFAULT FALSE,
    charisma_save BOOLEAN DEFAULT FALSE,
    
    -- Персональные черты
    personality_traits TEXT,
    ideals TEXT,
    bonds TEXT,
    flaws TEXT,
    
    -- Прочее
    proficiencies TEXT,
    languages TEXT,
    senses TEXT,
    features TEXT,
    
    photo VARCHAR(255)
);

-- Создание таблицы навыков персонажа
CREATE TABLE IF NOT EXISTS character_skills (
    id SERIAL PRIMARY KEY,
    character_id INTEGER NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL,
    modifier INTEGER DEFAULT 0,
    proficient BOOLEAN DEFAULT FALSE,
    ability VARCHAR(20) NOT NULL
);

-- Создание таблицы снаряжения персонажа
CREATE TABLE IF NOT EXISTS character_equipment (
    id SERIAL PRIMARY KEY,
    character_id INTEGER NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS character_spells (
    id SERIAL PRIMARY KEY,
    character_id INTEGER NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT
);

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS games (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    master_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    description VARCHAR(1000)
);

CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    session_key VARCHAR(12) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    finished_at TIMESTAMPTZ,
    summary TEXT
);

CREATE TABLE session_players (
    id SERIAL PRIMARY KEY,
    session_id UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    character_id INTEGER NOT NULL REFERENCES characters(id),
    UNIQUE(session_id, user_id)
);

CREATE TABLE IF NOT EXISTS armor (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(100),
    armor_class INTEGER,
    modifier VARCHAR(50),
    cost VARCHAR(100),
    rarity VARCHAR(50),
    stealth_disadvantage VARCHAR(10),
    strength_requirement VARCHAR(50),
    weight VARCHAR(50),
    unique_stats TEXT,
    charges VARCHAR(50),
    modifiers TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE characters ADD COLUMN IF NOT EXISTS armor_id INTEGER REFERENCES armor(id) ON DELETE SET NULL;
