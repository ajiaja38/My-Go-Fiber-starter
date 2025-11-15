package repository

import (
	"learn/fiber/pkg/model/entity"
	"learn/fiber/pkg/model/res"
	"strings"

	"gorm.io/gorm"
)

type BlogRepository struct {
	db *gorm.DB
}

func NewBlogRepository(db *gorm.DB) *BlogRepository {
	return &BlogRepository{db: db}
}

func (r *BlogRepository) Create(blog *entity.Blog) error {
	return r.db.Create(blog).Error
}

func (r *BlogRepository) FindAllPagination(page, limit int, search string) (*[]res.FindBlogResponse, int64, error) {
	var blogs []res.FindBlogResponse = make([]res.FindBlogResponse, 0)
	var total int64

	search = "%" + strings.ToLower(search) + "%"

	queryCount := r.db.Raw(`
        SELECT COUNT(*) as total
        FROM blogs b
        JOIN users u ON b.user_id = u.id
        WHERE
            LOWER(b.title) LIKE ?
            OR LOWER(b.body) LIKE ?
            OR LOWER(u.username) LIKE ?
    `, search, search, search)

	if err := queryCount.Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	query := r.db.Raw(`
        SELECT
            b.id,
            b.title,
            b.body,
            b.image,
            b.user_id,
            u.username as owner,
            b.created_at,
            b.updated_at
        FROM blogs b
        JOIN users u ON b.user_id = u.id
        WHERE
            LOWER(b.title) LIKE ?
            OR LOWER(b.body) LIKE ?
            OR LOWER(u.username) LIKE ?
        ORDER BY b.created_at DESC
        LIMIT ? OFFSET ?
    `, search, search, search, limit, (page-1)*limit)

	if err := query.Scan(&blogs).Error; err != nil {
		return nil, 0, err
	}

	return &blogs, total, nil
}

func (r *BlogRepository) FindById(id string) (*res.FindBlogResponse, error) {
	var blog res.FindBlogResponse

	if row := r.db.Raw(`
        SELECT
            b.id,
            b.title,
            b.body,
            b.image,
            b.user_id,
            u.username as owner,
            b.created_at,
            b.updated_at
        FROM blogs b
        JOIN users u ON b.user_id = u.id
        WHERE b.id = ?
    `, id).Scan(&blog).RowsAffected; row == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &blog, nil
}
