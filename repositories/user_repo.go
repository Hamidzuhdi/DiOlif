package repositories

import (
	"database/sql"
	"konveksi-app/models"
)

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (username, password, contact, address, created_at)
		VALUES (?, ?, ?, ?, NOW())`

	res, err := r.DB.Exec(query, user.Username, user.Password, user.Contact, user.Address)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int(id)
	return nil
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {
	query := `
		SELECT id, username, password, contact, address, created_at
		FROM users WHERE id = ?`

	var user models.User
	err := r.DB.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Password,
		&user.Contact, &user.Address, &user.CreatedAt,
	)

	return &user, err
}

func (r *UserRepository) GetAll() ([]models.User, error) {
	query := `
		SELECT id, username, password, contact, address, created_at
		FROM users ORDER BY username`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(
			&u.ID, &u.Username, &u.Password,
			&u.Contact, &u.Address, &u.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *UserRepository) Update(user *models.User) error {
	query := `
		UPDATE users 
		SET username = ?, password = ?, contact = ?, address = ? 
		WHERE id = ?`

	_, err := r.DB.Exec(query,
		user.Username, user.Password,
		user.Contact, user.Address,
		user.ID,
	)

	return err
}

func (r *UserRepository) Delete(id int) error {
	_, err := r.DB.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}
