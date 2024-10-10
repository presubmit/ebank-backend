package db

import (
	"context"
	"database/sql"
	m "ebank/services/identity/models"
	"ebank/shared/db"
	"ebank/shared/utils"

	"github.com/google/uuid"
)

type IdentityDB interface {
	ExecTx(context.Context, func(db.Tx) error) error

	CreateUser(db.Tx, *m.User) (string, error)
	GetUser(db.Tx, string) (*m.User, error)
	FetchUser(db.Tx, string) (*m.User, error)
	GetUserByEmail(db.Tx, string) (*m.User, error)

	CreateCompany(db.Tx, *m.Company) (string, error)
	GetCompaniesByUserId(db.Tx, string) ([]*m.Company, error)

	GetEmployeeByUserId(db.Tx, string, string) (*m.Employee, error)
	CreateEmployee(db.Tx, *m.Employee) (string, error)
}

type DB struct {
	*db.DB
}

func (d *DB) CreateUser(tx db.Tx, u *m.User) (string, error) {
	userId := uuid.New().String()
	_, err := d.ExecStmt(tx, createUserStmt, userId, u.Email, u.Password, u.FirstName, u.LastName)
	if err != nil {
		return "", err
	}
	return userId, nil
}

func (d *DB) GetUser(tx db.Tx, id string) (*m.User, error) {
	row := d.QueryRowStmt(tx, getUserByIdStmt, id)
	user := &m.User{}
	if err := row.Scan(&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.CreatedAt); err != nil {
		return nil, err
	}
	return user, nil
}

func (d *DB) FetchUser(tx db.Tx, id string) (*m.User, error) {
	row := d.QueryRowStmt(tx, getUserByIdStmt, id)
	user := &m.User{}
	if err := row.Scan(&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.CreatedAt); err != nil {
		return nil, err
	}
	return user, nil
}

func (d *DB) GetUserByEmail(tx db.Tx, email string) (*m.User, error) {
	row := d.QueryRowStmt(tx, getUserByEmailStmt, email)
	user := &m.User{}
	if err := row.Scan(&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.CreatedAt); err != nil {
		return nil, err
	}
	return user, nil
}

func (d *DB) CreateCompany(tx db.Tx, c *m.Company) (string, error) {
	id := uuid.New().String()
	_, err := d.ExecStmt(tx, createCompanyStmt, id, c.Name, c.CreatedBy)
	if err != nil {
	}
	return id, err
}

func (d *DB) CreateEmployee(tx db.Tx, e *m.Employee) (string, error) {
	id := uuid.New().String()
	_, err := d.ExecStmt(tx, createEmployeeStmt,
		id,
		e.UserId,
		e.CompanyId,
		e.Email,
		e.Role,
		e.InvitationSent,
	)
	return id, err
}

func (d *DB) GetEmployeeByUserId(tx db.Tx, userId, companyId string) (*m.Employee, error) {
	row := d.QueryRowStmt(tx, getEmployeeByUserIdStmt, userId, companyId)
	var nullableUserId sql.NullString
	var nullableEmail sql.NullString
	e := &m.Employee{}
	if err := row.Scan(
		&e.ID,
		&nullableUserId,
		&e.CompanyId,
		&nullableEmail,
		&e.Role,
		&e.InvitationSent,
		&e.CreatedAt,
	); err != nil {
		return nil, err
	}
	e.UserId = utils.NullStringToString(nullableUserId)
	e.Email = utils.NullStringToString(nullableEmail)
	return e, nil
}

func (d *DB) GetCompaniesByUserId(tx db.Tx, userId string) ([]*m.Company, error) {
	rows, err := d.QueryStmt(tx, getCompaniesByUserIdStmt, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []*m.Company
	for rows.Next() {
		company := &m.Company{}
		err = rows.Scan(&company.ID, &company.Name, &company.CreatedAt)
		if err != nil {
			return nil, err
		}
		companies = append(companies, company)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return companies, nil
}

type statementKey string

var (
	createUserStmt           = statementKey("createUserStmt")
	getUserByIdStmt          = statementKey("getUserByIdStmt")
	getUserByEmailStmt       = statementKey("getUserByEmailStmt")
	createCompanyStmt        = statementKey("createCompanyStmt")
	getCompaniesByUserIdStmt = statementKey("getCompaniesByUserIdStmt")
	createEmployeeStmt       = statementKey("createEmployeeStmt")
	getEmployeeByUserIdStmt  = statementKey("getEmployeeByUserId")
)

func (d *DB) RegisterStmts() error {
	getUserBaseStmt := `SELECT id, email, password, first_name, last_name, created_at FROM users `
	return d.Register(map[interface{}]string{
		createUserStmt: db.InsertStmt("users", []string{
			"id", "email", "password", "first_name", "last_name",
		}),
		getUserByIdStmt:    getUserBaseStmt + `WHERE id = ?`,
		getUserByEmailStmt: getUserBaseStmt + `WHERE email = ?`,

		createCompanyStmt: db.InsertStmt("companies", []string{"id", "name", "created_by"}),
		getCompaniesByUserIdStmt: `
			SELECT c.id, c.name, c.created_at 
			FROM companies c JOIN employees e ON c.id=e.company_id 
			WHERE e.user_id = ?
		`,

		createEmployeeStmt: db.InsertStmt("employees", []string{
			"id", "user_id", "company_id", "email", "role", "invitation_sent",
		}),
		getEmployeeByUserIdStmt: `SELECT id, user_id, company_id, email, role, invitation_sent, created_at FROM employees WHERE user_id = ? AND company_id = ? AND is_active IS TRUE`,
	})
}

func New() *DB {
	db := &DB{
		DB: db.New(),
	}
	if err := db.RegisterStmts(); err != nil {
		panic(err)
	}
	return db
}
