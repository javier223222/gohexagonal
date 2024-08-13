package admin

import (
	"database/sql"
	"fmt"
)

type MySQLAdminRepository struct {
	DB *sql.DB
}

// NewMySQLAdminRepository crea una nueva instancia de MySQLAdminRepository.
func NewMySQLAdminRepository(db *sql.DB) AdminRepository {
	return &MySQLAdminRepository{DB: db}
}

func (r *MySQLAdminRepository) Save(admin *Admin) error {
	query := "Insert into admin (username,name,lastName,fullname,email,number,idrol,createdBy,password) values (?,?,?,?,?,?,?,?,?)"
	_, err := r.DB.Exec(query, admin.Username, admin.Name, admin.LastName, admin.FullName, admin.Email, admin.Number, 1, admin.CreatedBy, admin.Password)
	return err

}

func (r *MySQLAdminRepository) GetByUsernameOrEmail(field string) (*Admin, error) {
	query := "Select id,username,name,lastName,fullname,email,number,idrol,createdAt,isDeleted,createdBy,password from admin where username = ? or email = ? and isDeleted = 0"
	row := r.DB.QueryRow(query, field, field)

	var admin Admin
	err := row.Scan(&admin.ID, &admin.Username, &admin.Name, &admin.LastName, &admin.FullName, &admin.Email, &admin.Number, &admin.IdRole, &admin.CreateAt, &admin.IsDeleted, &admin.CreatedBy, &admin.Password)
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *MySQLAdminRepository) Get(page int64, limit int64) ([]Admin, int64, int64, error) {
	offset := (page - 1) * limit

	// Consulta para contar el número total de registros
	var totalRecords int64
	countQuery := "SELECT COUNT(*) FROM admin WHERE isDeleted = 0"
	err := r.DB.QueryRow(countQuery).Scan(&totalRecords)
	if err != nil {
		return nil, 0, 0, err
	}

	// Cálculo del número total de páginas
	totalPages := (totalRecords + limit - 1) / limit

	// Consulta para obtener los registros paginados
	query := "SELECT id, username, name, lastName, fullname, email, number, idrol, createdAt, isDeleted, createdBy " +
		"FROM admin WHERE isDeleted = 0 LIMIT ? OFFSET ?"

	rows, err := r.DB.Query(query, limit, offset)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()

	var admins []Admin
	for rows.Next() {
		var admin Admin
		err := rows.Scan(&admin.ID, &admin.Username, &admin.Name, &admin.LastName, &admin.FullName, &admin.Email, &admin.Number, &admin.IdRole, &admin.CreateAt, &admin.IsDeleted, &admin.CreatedBy)
		if err != nil {
			return nil, 0, 0, err
		}
		admins = append(admins, admin)
	}

	return admins, totalRecords, totalPages, nil
}

func (r *MySQLAdminRepository) GetByID(id int64) (*Admin, error) {
	query := "Select id,username,name,lastName,fullname,email,number,idrol,createdAt,isDeleted,createdBy from admin where id = ? and isDeleted = 0"
	row := r.DB.QueryRow(query, id)

	var admin Admin
	err := row.Scan(&admin.ID, &admin.Username, &admin.Name, &admin.LastName, &admin.FullName, &admin.Email, &admin.Number, &admin.IdRole, &admin.CreateAt, &admin.IsDeleted, &admin.CreatedBy)
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *MySQLAdminRepository) Delete(id int64) error {
	query := "Update admin set isDeleted = 1 where id = ?"
	_, err := r.DB.Exec(query, id)
	return err
}

func (r *MySQLAdminRepository) UpdatePassword(id int64, password string) error {
	query := "Update admin set password = ? where id = ?"
	_, err := r.DB.Exec(query, password, id)
	return err
}

func (r *MySQLAdminRepository) GetByIDWithPassword(id int64) (*Admin, error) {
	query := "Select id,username,name,lastName,fullname,email,number,idrol,createdAt,isDeleted,createdBy,password from admin where id = ? and isDeleted = 0"
	row := r.DB.QueryRow(query, id)

	var admin Admin
	err := row.Scan(&admin.ID, &admin.Username, &admin.Name, &admin.LastName, &admin.FullName, &admin.Email, &admin.Number, &admin.IdRole, &admin.CreateAt, &admin.IsDeleted, &admin.CreatedBy, &admin.Password)
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *MySQLAdminRepository) UpdateUsername(id int64, username string) error {
	tx, err := r.DB.Begin()
    if err != nil {
        return fmt.Errorf("could not begin transaction: %v", err)
    }
	defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            panic(r) // Re-lanzar el pánico después de revertir
        } else if err != nil {
            tx.Rollback() // Revertir si hubo un error
        } else {
            err = tx.Commit() // Confirmar la transacción si todo salió bien
        }
    }()


	query := "Update admin set username = ? where id = ?"
    _, err =tx.Exec(query, username, id)
	if err != nil {
        return fmt.Errorf("could not update: %v", err)
    }


	return err
}

func (r *MySQLAdminRepository) UpdateName(id int64, name string) error {
	tx, err := r.DB.Begin()
    if err != nil {
        return fmt.Errorf("could not begin transaction: %v", err)
    }
	defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            panic(r) // Re-lanzar el pánico después de revertir
        } else if err != nil {
            tx.Rollback() // Revertir si hubo un error
        } else {
            err = tx.Commit() // Confirmar la transacción si todo salió bien
        }
    }()
	
	query := "Update admin set name = ? where id = ?"

	_, err = tx.Exec(query, name, id)
	if err != nil {
		return fmt.Errorf("could not update: %v", err)
	}
	query="Select name,lastName from admin where id = ?"
	row := tx.QueryRow(query,id)
	var admin Admin
	err = row.Scan(&admin.Name,&admin.LastName)
	if err != nil {

		return fmt.Errorf("could not update: %v", err)
	}
	query="Update admin set fullname = ? where id = ?"
	_, err = tx.Exec(query,admin.Name+" "+admin.LastName,id)
	if err != nil {
		return fmt.Errorf("could not update: %v", err)
	}
	
	


	return err
}	

func (r *MySQLAdminRepository) UpdateLastName(id int64, lastName string) error {
	tx, err := r.DB.Begin()
    if err != nil {
        return fmt.Errorf("could not begin transaction: %v", err)
    }
	defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            panic(r) // Re-lanzar el pánico después de revertir
        } else if err != nil {
            tx.Rollback() // Revertir si hubo un error
        } else {
            err = tx.Commit() // Confirmar la transacción si todo salió bien
        }
    }()
	
	query := "Update admin set LastName = ? where id = ?"

	_, err = tx.Exec(query, lastName, id)
	if err != nil {
		return fmt.Errorf("could not update: %v", err)
	}
	query="Select name,lastName from admin where id = ?"
	row := tx.QueryRow(query,id)
	var admin Admin
	err = row.Scan(&admin.Name,&admin.LastName)
	if err != nil {

		return fmt.Errorf("could not update: %v", err)
	}
	query="Update admin set fullname = ? where id = ?"
	_, err = tx.Exec(query,admin.Name+" "+admin.LastName,id)
	if err != nil {
		return fmt.Errorf("could not update: %v", err)
	}
	
	


	return err
}

func (r *MySQLAdminRepository) UpdateEmail(id int64, email string) error {
	query := "Update admin set email = ? where id = ?"
	_, err := r.DB.Exec(query, email, id)
	return err
}


func (r *MySQLAdminRepository) UpdateNumber(id int64, number string) error {
	query := "Update admin set number = ? where id = ?"
	_, err := r.DB.Exec(query, number, id)
	return err
}

func (r *MySQLAdminRepository) UpdateRole(id int64, role int64) error {
	query := "Update admin set idrol = ? where id = ?"
	_, err := r.DB.Exec(query, role, id)
	return err
}