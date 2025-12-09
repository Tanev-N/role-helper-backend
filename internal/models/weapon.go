package models

import (
	"errors"
)

type Weapon struct {
	ID          int    `json:"id"`           // Уникальный идентификатор оружия
	UserID      int    `json:"user_id"`      // ID пользователя-владельца
	Name        string `json:"name"`         // Название предмета (например, "Боевая коса")
	Type        string `json:"type"`         // Тип предмета (например, "простое рукопашное")
	Damage      string `json:"damage"`       // Урон (например, "1к8")
	Modifier    string `json:"modifier"`     // Модификатор (например, "Сила" или "Ловк.")
	Cost        string `json:"cost"`         // Стоимость (например, "30 золотых")
	Rarity      string `json:"rarity"`       // Редкость (например, "Обычная")
	Grip        string `json:"grip"`         // Хват ("Одноручное" или "Двуручное")
	RangeMeters string `json:"range_meters"` // Дальность в метрах (например, "2.5")
	Weight      string `json:"weight"`       // Вес в кг (например, "5")
	UniqueStats string `json:"unique_stats"` // Уникальные показатели ("Нет" или текст)
	Charges     string `json:"charges"`      // Заряд ("Нет" или число)
	Photo       string `json:"photo"`        // URL фотографии оружия
}

type WeaponRepository interface {
	Create(weapon *Weapon) (*Weapon, error)         // Создание нового оружия
	GetAll(userID int) ([]Weapon, error)            // Получение списка всего оружия пользователя
	GetByID(id, userID int) (*Weapon, error)        // Получение оружия по ID
	Update(id int, weapon *Weapon) (*Weapon, error) // Обновление оружия
	Delete(id, userID int) error                    // Удаление оружия
	CheckBelonging(id, userID int) (bool, error)    // Проверка принадлежности оружия пользователю
}

type WeaponService interface {
	Create(weapon *Weapon) (*Weapon, error)                     // Создание нового оружия с валидацией
	GetAll(userID int) ([]Weapon, error)                        // Получение списка всего оружия пользователя
	GetByID(id, userID int) (*Weapon, error)                    // Получение оружия по ID
	Update(id int, userID int, weapon *Weapon) (*Weapon, error) // Обновление оружия с валидацией
	Delete(id, userID int) error                                // Удаление оружия
}

var (
	ErrWeaponNotFound = errors.New("weapon not found")
)
