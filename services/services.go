package services

import (
	"fmt"

	"gorm.io/gorm"
)

type BaseService[T any] struct {
	Db *gorm.DB
}

type PaginationParams struct {
	Page     int
	PageSize int
}

func NewBaseService[T any](Db *gorm.DB) *BaseService[T] {
	return &BaseService[T]{Db: Db}
}

func (service *BaseService[T]) Create(item *T) error {
	return service.Db.Create(item).Error
}

func (service *BaseService[T]) GetAll() ([]T, error) {
	var items []T
	err := service.Db.Find(&items).Error
	return items, err
}

func (service *BaseService[T]) GetById(id uint) (*T, error) {
	var item T
	err := service.Db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (service *BaseService[T]) Find(conds ...interface{}) ([]T, error) {
	var items []T
	if len(conds) >= 1 {
		err := service.Db.Where(conds[0], conds[1:]...).Find(&items).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := service.Db.Find(&items).Error
		if err != nil {
			return nil, err
		}
	}
	return items, nil
}

func (service *BaseService[T]) First(dest *T, conds ...interface{}) error {
	if len(conds) >= 1 {
		err := service.Db.Where(conds[0], conds[1:]...).First(dest).Error
		if err != nil {
			return err
		}
	} else {
		err := service.Db.First(dest).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (service *BaseService[T]) Last(dest *T, conds ...interface{}) error {
	if len(conds) >= 1 {
		err := service.Db.Where(conds[0], conds[1:]...).Last(dest).Error
		if err != nil {
			return err
		}
	} else {
		err := service.Db.Last(dest).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (service *BaseService[T]) Update(getModelInstance func() T, conds interface{}, updates interface{}) error {
	modelInstance := getModelInstance()
	result := service.Db.Model(modelInstance).Where(conds).Updates(updates)
	return result.Error
}

type FilterOption func(*gorm.DB)

func OrderByTimestamp(ascending bool) FilterOption {
	return func(db *gorm.DB) {
		order := "created_at DESC"
		if ascending {
			order = "created_at ASC"
		}
		db.Order(order)
	}
}

func FilterBy(field string, value interface{}) FilterOption {
	return func(db *gorm.DB) {
		db.Where(fmt.Sprintf("%s = ?", field), value)
	}
}

func FilterWhere(condition string, args ...interface{}) FilterOption {
	return func(db *gorm.DB) {
		db.Where(condition, args...)
	}
}

func LimitResults(limit int) FilterOption {
	return func(db *gorm.DB) {
		db.Limit(limit)
	}
}

func GeneralFilter(field string, operator string, value interface{}) FilterOption {
	return func(db *gorm.DB) {
		switch operator {
		case "=", "!=", ">", ">=", "<", "<=":
			db.Where(fmt.Sprintf("%s %s ?", field, operator), value)
		case "IN":
			db.Where(fmt.Sprintf("%s IN (?)", field), value)
		case "NOT IN":
			db.Where(fmt.Sprintf("%s NOT IN (?)", field), value)
		case "LIKE":
			db.Where(fmt.Sprintf("%s LIKE ?", field), fmt.Sprintf("%%%v%%", value))
		case "IS NULL":
			db.Where(fmt.Sprintf("%s IS NULL", field))
		case "IS NOT NULL":
			db.Where(fmt.Sprintf("%s IS NOT NULL", field))
		case ">NOW()":
			db.Where(fmt.Sprintf("%s > NOW()", field))
		default:
			// If an unsupported operator is used, this will effectively make it a no-op filter
			fmt.Printf("Unsupported operator in GeneralFilter: %s", operator)
		}
	}
}

func (service *BaseService[T]) FindWithOptions(options ...FilterOption) ([]T, error) {
	var items []T
	query := service.Db

	for _, option := range options {
		option(query)
	}

	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func Paginate(query *gorm.DB, params PaginationParams) *gorm.DB {
	offset := (params.Page - 1) * params.PageSize
	return query.Offset(offset).Limit(params.PageSize)
}
