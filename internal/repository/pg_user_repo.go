package repository

import (
    "database/sql"
    "errors"

    "sso/internal/model"
    // Using standard library `database/sql`; driver not locked here.
)

type pgUserRepo struct {
    db *sql.DB
}

func NewPostgresUserRepository(dsn string) (UserRepository, error) {
    db, err := sql.Open("pgx", dsn)
    if err != nil {
        return nil, err
    }

    if err := db.Ping(); err != nil {
        return nil, err
    }

    repo := &pgUserRepo{db: db}
    if err := repo.Migrate(); err != nil {
        return nil, err
    }

    return repo, nil
}

func (r *pgUserRepo) Migrate() error {
    _, err := r.db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id TEXT PRIMARY KEY,
            email TEXT UNIQUE NOT NULL,
            password TEXT NOT NULL,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
        );
    `)
    return err
}

func (r *pgUserRepo) Create(user model.User) error {
    _, err := r.db.Exec(`INSERT INTO users (id, email, password) VALUES ($1, $2, $3)`, user.ID, user.Email, user.Password)
    if err != nil {
        if err.Error() == "unique_violation" { // simplified check
            return errors.New("user exists")
        }
        return err
    }
    return nil
}

func (r *pgUserRepo) GetByEmail(email string) (model.User, error) {
    var u model.User
    row := r.db.QueryRow(`SELECT id, email, password FROM users WHERE email = $1`, email)
    if err := row.Scan(&u.ID, &u.Email, &u.Password); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return model.User{}, errors.New("user not found")
        }
        return model.User{}, err
    }
    return u, nil
}
