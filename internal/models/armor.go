package models

import (
	"errors"
)

type Armor struct {
	ID                  int    `json:"id"`                   // Уникальный идентификатор брони
	UserID              int    `json:"user_id"`              // ID пользователя-владельца
	Name                string `json:"name"`                 // Название предмета (например, "Кожаный доспех")
	Type                string `json:"type"`                 // Тип предмета (например, "лёгкий доспех")
	ArmorClass          int    `json:"armor_class"`          // Класс доспеха (например, 11)
	Modifier            string `json:"modifier"`             // Модификатор (например, "Ловк.")
	Cost                string `json:"cost"`                 // Стоимость (например, "30 золотых")
	Rarity              string `json:"rarity"`               // Редкость (например, "Обычная")
	StealthDisadvantage string `json:"stealth_disadvantage"` // Помеха для скрытности ("Да"/"Нет")
	StrengthRequirement string `json:"strength_requirement"` // Требование к силе ("Нет" или число)
	Weight              string `json:"weight"`               // Вес в кг (например, "3")
	UniqueStats         string `json:"unique_stats"`         // Уникальные показатели ("Нет" или текст)
	Charges             string `json:"charges"`              // Заряд ("Нет" или число)
	Modifiers           string `json:"modifiers"`            // Дополнительные модификаторы (JSON строка или текст)
	Photo               string `json:"photo"`                // URL фотографии брони
}

type ArmorRepository interface {
	Create(armor *Armor) (*Armor, error)         // Создание новой брони
	GetAll(userID int) ([]Armor, error)          // Получение списка всей брони пользователя
	GetByID(id, userID int) (*Armor, error)      // Получение брони по ID
	Update(id int, armor *Armor) (*Armor, error) // Обновление брони
	Delete(id, userID int) error                 // Удаление брони
	CheckBelonging(id, userID int) (bool, error) // Проверка принадлежности брони пользователю
}

type ArmorService interface {
	Create(armor *Armor) (*Armor, error)                     // Создание новой брони с валидацией
	GetAll(userID int) ([]Armor, error)                      // Получение списка всей брони пользователя
	GetByID(id, userID int) (*Armor, error)                  // Получение брони по ID
	Update(id int, userID int, armor *Armor) (*Armor, error) // Обновление брони с валидацией
	Delete(id, userID int) error                             // Удаление брони
}

var (
	ErrArmorNotFound = errors.New("armor not found")
)
