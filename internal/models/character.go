package models

import (
	"errors"
)

type Character struct {
	ID           string `json:"id"`            // Уникальный идентификатор персонажа
	Name         string `json:"name"`          // Имя персонажа
	Race         string `json:"race"`          // Раса персонажа (любая)
	Class        string `json:"class"`        // Класс персонажа (любой)
	Level        int    `json:"level"`        // Уровень персонажа (1-20)
	Alignment    string `json:"alignment"`    // Мировоззрение персонажа
	Background   string `json:"background"`   // Предыстория персонажа
	PlayerName   string `json:"player_name"`  // Имя игрока
	Experience   int    `json:"experience"`   // Опыт персонажа
	
	// Основные характеристики (1-30)
	Strength     int `json:"strength"`     // Сила
	Dexterity    int `json:"dexterity"`   // Ловкость
	Constitution int `json:"constitution"` // Телосложение
	Intelligence int `json:"intelligence"` // Интеллект
	Wisdom       int `json:"wisdom"`      // Мудрость
	Charisma     int `json:"charisma"`    // Харизма
	
	// Модификаторы характеристик (рассчитываются автоматически)
	StrengthMod     int `json:"strength_mod"`     // Модификатор силы
	DexterityMod    int `json:"dexterity_mod"`    // Модификатор ловкости
	ConstitutionMod int `json:"constitution_mod"` // Модификатор телосложения
	IntelligenceMod int `json:"intelligence_mod"` // Модификатор интеллекта
	WisdomMod       int `json:"wisdom_mod"`       // Модификатор мудрости
	CharismaMod     int `json:"charisma_mod"`     // Модификатор харизмы
	
	// Боевые характеристики
	ProficiencyBonus int `json:"proficiency_bonus"` // Бонус мастерства (рассчитывается по уровню)
	Initiative       int `json:"initiative"`        // Инициатива (модификатор ловкости)
	ArmorClass       int `json:"armor_class"`       // Класс брони
	Speed            int `json:"speed"`             // Скорость передвижения
	HitPoints        int `json:"hit_points"`        // Текущие хиты
	MaxHitPoints     int `json:"max_hit_points"`    // Максимальные хиты
	TempHitPoints    int `json:"temp_hit_points"`   // Временные хиты
	HitDice          string `json:"hit_dice"`       // Кость хитов (например, "1d8")
	
	// Спасброски (булевы значения - владеет ли персонаж спасброском)
	StrengthSave     bool `json:"strength_save"`     // Спасбросок силы
	DexteritySave    bool `json:"dexterity_save"`    // Спасбросок ловкости
	ConstitutionSave bool `json:"constitution_save"` // Спасбросок телосложения
	IntelligenceSave bool `json:"intelligence_save"` // Спасбросок интеллекта
	WisdomSave       bool `json:"wisdom_save"`       // Спасбросок мудрости
	CharismaSave     bool `json:"charisma_save"`     // Спасбросок харизмы
	
	// Навыки персонажа
	Skills []CharacterSkill `json:"skills"`
	
	// Персональные черты персонажа
	PersonalityTraits string `json:"personality_traits"` // Черты характера
	Ideals           string `json:"ideals"`              // Идеалы
	Bonds            string `json:"bonds"`               // Привязанности
	Flaws            string `json:"flaws"`               // Слабости
	
	// Прочие характеристики
	Proficiencies string `json:"proficiencies"` // Владения (оружие, броня, инструменты)
	Languages    string `json:"languages"`     // Языки
	Senses       string `json:"senses"`        // Чувства (темное зрение и т.д.)
	
	// Снаряжение и заклинания
	Equipment []Equipment `json:"equipment"` // Снаряжение персонажа
	Spells    []Spell     `json:"spells"`    // Заклинания персонажа
	
	// Особенности и черты класса/расы
	Features string `json:"features"` // Особенности персонажа
	
	Photo string `json:"photo"` // URL фотографии персонажа
}

type CharacterSkill struct {
	Name       string `json:"name"`       // Название навыка
	Modifier   int    `json:"modifier"`   // Модификатор навыка (рассчитывается автоматически)
	Proficient bool   `json:"proficient"` // Владеет ли персонаж навыком
	Ability    string `json:"ability"`    // Характеристика, от которой зависит навык
}

type Equipment struct {
	Name        string `json:"name"`        // Название предмета
	Description string `json:"description"` // Описание предмета
}

type Spell struct {
	Name        string `json:"name"`        // Название заклинания
	Description string `json:"description"` // Описание заклинания
}

type CharacterShort struct {
	ID    string `json:"id"`    // Уникальный идентификатор персонажа
	Name  string `json:"name"`  // Имя персонажа
	Photo string `json:"photo"` // URL фотографии персонажа
}

type CharacterRepository interface {
	Create(character *Character) (*Character, error)                    // Создание нового персонажа
	GetAll() ([]CharacterShort, error)                                  // Получение списка всех персонажей
	FindByID(id string) (*Character, error)                             // Поиск персонажа по ID
	Update(id string, update *Character) (*Character, error)           // Обновление персонажа
	Delete(id string) error                                             // Удаление персонажа
}

type CharacterService interface {
	Create(create *Character) (*Character, error)                       // Создание нового персонажа с валидацией
	GetAll() ([]CharacterShort, error)                                  // Получение списка всех персонажей
	FindByID(id string) (*Character, error)                             // Поиск персонажа по ID
	Update(id string, update *Character) (*Character, error)           // Обновление персонажа с валидацией
	Delete(id string) error                                             // Удаление персонажа
}

var (
	ErrCharacterNotFound = errors.New("character not found")
)
