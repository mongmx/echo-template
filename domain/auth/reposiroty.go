package auth

import (
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) CreateUser(u *User) error {
	tx := r.db.Begin()
	// err := u.BeforeCreate(tx)
	// if err != nil {
	// 	return err
	// }
	// _, err = tx.Exec("INSERT INTO users (id, email, password) VALUES (?, ?, ?)", u.ID, u.Email, u.Password)
	// if err != nil {
	// 	return err
	// }
	tx.Commit()
	return nil
}

func (r *Repo) GetUserByID(id string) (*User, error) {
	user := User{}
	// err := r.db.Get(&user, "SELECT * FROM users WHERE id = ?", id)
	// if err != nil {
	// 	return nil, err
	// }
	return &user, nil
}

func (r *Repo) GetUserByEmail(email string) (*User, error) {
	user := User{}
	// err := r.db.Get(&user, "SELECT * FROM users WHERE email = ?", email)
	// if err != nil {
	// 	return nil, err
	// }
	return &user, nil
}

func (r *Repo) CreateLog(h *History) error {
	// _, err := r.db.NamedExec("INSERT INTO log (user_id, name) VALUES (:user_id, :name)", h)
	// if err != nil {
	// 	return err
	// }
	return nil
}

func (r *Repo) StoreJWT(id string, a *Auth) error {
	// atExp := time.Duration(a.AccessTokenExpiresIn) * time.Minute
	// rtExp := time.Duration(a.RefreshTokenExpiresIn) * time.Minute
	tx := r.db.Begin()
	// sqlStr := "INSERT INTO tokens (token_id, user_id, expire) VALUES (?, ?, ?)"
	// _, err := tx.Exec(sqlStr, a.AccessToken, id, atExp)
	// if err != nil {
	// 	return err
	// }
	// _, err = tx.Exec(sqlStr, a.RefreshToken, id, rtExp)
	// if err != nil {
	// 	return err
	// }
	tx.Commit()
	return nil
}
