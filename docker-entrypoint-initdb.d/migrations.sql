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
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    race VARCHAR(50) NOT NULL,
    class VARCHAR(50) NOT NULL,
    level INTEGER NOT NULL CHECK (level >= 1 AND level <= 20),
    alignment VARCHAR(50),
    background VARCHAR(100),
    player_name VARCHAR(100),
    experience INTEGER DEFAULT 0 CHECK (experience >= 0),
    
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
    character_id VARCHAR(50) NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL,
    modifier INTEGER DEFAULT 0,
    proficient BOOLEAN DEFAULT FALSE,
    ability VARCHAR(20) NOT NULL
);

-- Создание таблицы снаряжения персонажа
CREATE TABLE IF NOT EXISTS character_equipment (
    id SERIAL PRIMARY KEY,
    character_id VARCHAR(50) NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT
);

-- Создание таблицы заклинаний персонажа
CREATE TABLE IF NOT EXISTS character_spells (
    id SERIAL PRIMARY KEY,
    character_id VARCHAR(50) NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT
);

-- Создание индексов для оптимизации запросов
CREATE INDEX IF NOT EXISTS idx_characters_name ON characters(name);
CREATE INDEX IF NOT EXISTS idx_characters_player_name ON characters(player_name);
CREATE INDEX IF NOT EXISTS idx_character_skills_character_id ON character_skills(character_id);
CREATE INDEX IF NOT EXISTS idx_character_equipment_character_id ON character_equipment(character_id);
CREATE INDEX IF NOT EXISTS idx_character_spells_character_id ON character_spells(character_id);