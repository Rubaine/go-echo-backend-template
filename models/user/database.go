package user

import (
	"backend-template/models/postgresql"
	"database/sql"
	"time"

	"github.com/charmbracelet/log"

	"github.com/jackc/pgx/v4"
)

func ScanUser(row pgx.Row) (u User, err error) {

	var (
		id                                                 sql.NullInt64
		email, firstname, lastname, password, recoverToken sql.NullString
		admin, enable                                      sql.NullBool
	)

	err = row.Scan(
		&id,
		&email,
		&firstname,
		&lastname,
		&password,
		&recoverToken,
		&admin,
		&enable,
	)

	if err != nil {
		return
	}

	u = User{
		ID:           id.Int64,
		Email:        email.String,
		Firstname:    firstname.String,
		Lastname:     lastname.String,
		Password:     password.String,
		RecoverToken: recoverToken.String,
		Enable:       enable.Bool,
		Admin:        admin.Bool,
	}

	return
}

func GetSQLUserToken(email, password string) (token UserToken, err error) {

	query := "select * from account " +
		"where enable=true and email=$1 and password=crypt($2, password)"

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	row := sqlCo.QueryRow(postgresql.SQLCtx, query, email, password)
	u, err := ScanUser(row)

	if err == pgx.ErrNoRows {
		log.Info("Tentative de connexion infructueuse. Email inconnu ou compte désactivé.", "email", email)
		return
	} else if err != nil {
		log.Error("During GetSQLUserToken query", "error", err)
		return
	}

	token = UserToken{
		User:      u,
		CreatedAt: time.Now(),
	}

	return
}

func CreateAccount(email, firstname, lastname, password string) (id int64) {

	query := "insert into account (email, firstname, lastname, password) " +
		"VALUES ($1,$2,$3,crypt($4, gen_salt('bf'))) RETURNING id"

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return -1
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	id = time.Now().UnixNano()

	err = sqlCo.QueryRow(postgresql.SQLCtx, query, email, firstname, lastname, password).Scan(&id)
	if err != nil {
		log.Warn("During CreateAccount query", "error", err)
		return -1
	}
	return id
}

func DeleteAccount(id int64) (msg string) {

	query := "delete from account where id=$1"

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		msg = "Internal server error"
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	cmd, err := sqlCo.Exec(postgresql.SQLCtx, query, id)
	if err != nil {
		return "Internal server error"
	} else if cmd.RowsAffected() == 0 {
		return "Account not found"
	}
	return
}

func PasswordCheck(id int64, password string) (checked bool) {
	query := "select id from account where enable=true and id=$1 and password=crypt($2, password)"

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	var id_ int64
	err = sqlCo.QueryRow(postgresql.SQLCtx, query, id, password).Scan(&id_)
	if err == nil {
		return id_ == id
	}
	return
}

func CheckEmailAvailability(email string) (available bool) {
	query := "select id from account where email=$1"

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	var id int64
	err = sqlCo.QueryRow(postgresql.SQLCtx, query, email).Scan(&id)
	if err == pgx.ErrNoRows {
		return true
	}
	return
}

func GetUserById(UserId int64) (u User, err error) {

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	var query = "SELECT * FROM account WHERE enable=true and id=$1"

	row := sqlCo.QueryRow(postgresql.SQLCtx, query, UserId)
	u, err = ScanUser(row)

	return
}

func GetAllUsers() (users []User) {

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}

	defer sqlCo.Close(postgresql.SQLCtx)

	query := "SELECT * FROM account"

	rows, err := sqlCo.Query(postgresql.SQLCtx, query)
	if err != nil {
		log.Error("During GetAllUsers query", "error", err)
		return
	}
	defer rows.Close()

	for rows.Next() {

		u, err := ScanUser(rows)

		if err == pgx.ErrNoRows {
			return
		} else if err != nil {
			log.Error("During GetAllUsers scan", "error", err)
			return
		}

		users = append(users, u)
	}

	return
}

func GetUserByEmail(userEmail string) (u User, msg string) {

	if userEmail == "" {
		msg = "Empty email"
		return
	}

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		msg = "Internal server error"
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	query := "SELECT * FROM account WHERE email=$1"

	row := sqlCo.QueryRow(postgresql.SQLCtx, query, userEmail)
	u, err = ScanUser(row)

	if err == pgx.ErrNoRows {
		msg = "Username not found"
	} else if err != nil {
		log.Error("During GetUserByEmail query", "error", err)
		msg = "Internal server error"
	}

	return
}

func UpdateUser(user User, password bool) (ok bool) {

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	var query string
	var args []any
	if password {
		query = "UPDATE account set (email, firstname, lastname, password) = ($1,$2,$3,crypt($4, gen_salt('bf'))) " +
			"WHERE id=$5"
		args = []any{
			user.Email,
			user.Firstname,
			user.Lastname,
			user.Password,
			user.ID,
		}
	} else {
		query = "UPDATE account set (email, firstname, lastname) = ($1,$2,$3) " +
			"WHERE id=$4"
		args = []any{
			user.Email,
			user.Firstname,
			user.Lastname,
			user.ID,
		}
	}

	cmd, err := sqlCo.Exec(postgresql.SQLCtx, query, args...)
	ok = cmd.RowsAffected() == 1 && err == nil
	return
}

func CreateRecoverToken(email string) (string, error) {

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return "", err
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	var firstname, recoverToken sql.NullString

	query := "UPDATE account set recover_token=gen_random_uuid() where email=$1 RETURNING firstname, recover_token"
	err = sqlCo.QueryRow(postgresql.SQLCtx, query, email).Scan(&firstname, &recoverToken)
	if err == pgx.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return recoverToken.String, nil
}

func ResetPassword(token, password string) (ok bool) {

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return false
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	query := "UPDATE account set password=crypt($2, gen_salt('bf')), recover_token=null where recover_token=$1"
	cmd, err := sqlCo.Exec(postgresql.SQLCtx, query, token, password)
	return cmd.RowsAffected() == 1 && err == nil
}

func IsInOrganization(userId, orgId int64) (ok bool) {

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return false
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	query := "SELECT account FROM account_in_organization WHERE account=$1 and organization=$2"
	row := sqlCo.QueryRow(postgresql.SQLCtx, query, userId, orgId)

	var id int64
	err = row.Scan(&id)
	return err == nil
}

func ListOrganizationMembers(orgID int64) (users UserList, err error) {

	query := `SELECT a.* 
				FROM account a 
				JOIN account_in_organization aio ON a.id = aio.account 
				WHERE aio.organization = $1`

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	rows, err := sqlCo.Query(postgresql.SQLCtx, query, orgID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		user, err = ScanUser(rows)
		if err != nil {
			return
		}

		users = append(users, user)
	}

	return
}
