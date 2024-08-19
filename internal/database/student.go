package database

import (
	"GoAssignment/internal/student"
	"context"
	"database/sql"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
)

type StudentRow struct {
	ID          string
	Fname       sql.NullString
	Lname       sql.NullString
	Email       sql.NullString
	Gender      sql.NullString
	DateOfBirth sql.NullTime
	Address     sql.NullString
	CreatedBy   sql.NullString
	CreatedOn   sql.NullTime
	UpdatedBy   sql.NullString
	UpdatedOn   sql.NullTime
}

func convertStudentRowToStudent(s StudentRow) student.Student {
	dateOfBirth := student.CustomeTime{Time: s.DateOfBirth.Time}
	return student.Student{
		ID:          s.ID,
		Fname:       s.Fname.String,
		Lname:       s.Lname.String,
		Email:       s.Email.String,
		Gender:      s.Gender.String,
		DateOfBirth: dateOfBirth,
		Address:     s.Address.String,
		CreatedBy:   s.CreatedBy.String,
		CreatedOn:   s.CreatedOn.Time,
		UpdatedBy:   s.UpdatedBy.String,
		UpdatedOn:   s.UpdatedOn.Time,
	}
}

func (d *Database) CreateStudent(ctx context.Context, std student.Student) (student.Student, error) {
	std.ID = uuid.NewV4().String()
	dateOfBirthNullTime := sql.NullTime{Time: std.DateOfBirth.Time, Valid: !std.DateOfBirth.IsZero()}
	postRow := StudentRow{
		ID:          std.ID,
		Fname:       sql.NullString{String: std.Fname, Valid: true},
		Lname:       sql.NullString{String: std.Lname, Valid: true},
		Email:       sql.NullString{String: std.Email, Valid: true},
		Gender:      sql.NullString{String: std.Gender, Valid: true},
		DateOfBirth: dateOfBirthNullTime,
		Address:     sql.NullString{String: std.Address, Valid: true},
		CreatedBy:   sql.NullString{String: std.CreatedBy, Valid: true},
		CreatedOn:   sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedBy:   sql.NullString{String: std.UpdatedBy, Valid: true},
		UpdatedOn:   sql.NullTime{Time: time.Now(), Valid: false},
	}

	rows, err := d.Client.NamedExecContext(
		ctx,
		`INSERT INTO student 
        (id, fname, lname, email, gender, dateofbirth, address, createdby, createdon, updatedby, updatedon) 
     VALUES
        (:id, :fname, :lname, :email, :gender, :dateofbirth, :address, :createdby, :createdon, :updatedby, :updatedon)`,
		postRow,
	)
	if err != nil {
		return student.Student{}, fmt.Errorf("failed to insert Student: %w", err)
	}
	insertedID, err := rows.LastInsertId()
	if err != nil {
		return student.Student{}, fmt.Errorf("failed to close rows: %w", err)
	}
	fmt.Printf("New record ID: %d\n", insertedID)
	return std, nil
}

// GetStudent - retrieves a student from the database by ID
func (d *Database) GetStudent(ctx context.Context, uuid string) (student.Student, error) {
	// fetch stidentRow from the database and then convert to student.Student
	var stdRow StudentRow
	row := d.Client.QueryRowContext(
		ctx,
		`SELECT id, fname, lname, email, gender, dateofbirth, address, createdby, createdon, updatedby, updatedon
		FROM student
		WHERE id = ?`,
		uuid,
	)
	err := row.Scan(&stdRow.ID, &stdRow.Fname, &stdRow.Lname, &stdRow.Email, &stdRow.Gender, &stdRow.DateOfBirth, &stdRow.Address, &stdRow.CreatedBy, &stdRow.CreatedOn, &stdRow.UpdatedBy, &stdRow.UpdatedOn)
	if err != nil {
		return student.Student{}, fmt.Errorf("an error occurred fetching the student by uuid: %w", err)
	}
	// sqlx with context to ensure context cancelation is honoured
	return convertStudentRowToStudent(stdRow), nil
}

func (d *Database) GetStudents(ctx context.Context) ([]student.Student, error) {
	// Fetch all rows from the database
	var stdRows []StudentRow
	err := d.Client.Select(&stdRows, "SELECT * FROM student LIMIT 10")
	if err != nil {
		return nil, fmt.Errorf("fetchStudents %v", err)
	}

	// Convert each StudentRow to student.Student
	students := make([]student.Student, len(stdRows))
	for i, stdRow := range stdRows {
		students[i] = student.Student{
			ID:          stdRow.ID,
			Fname:       stdRow.Fname.String,
			Lname:       stdRow.Lname.String,
			Email:       stdRow.Email.String,
			Gender:      stdRow.Gender.String,
			DateOfBirth: student.CustomeTime{Time: stdRow.DateOfBirth.Time}, // Custom type handling
			Address:     stdRow.Address.String,
			CreatedBy:   stdRow.CreatedBy.String,
			CreatedOn:   stdRow.CreatedOn.Time, // Handling nulls directly
			UpdatedBy:   stdRow.UpdatedBy.String,
			UpdatedOn:   stdRow.UpdatedOn.Time, // Handling nulls directly
		}
	}

	return students, nil
}

// DeleteStudent - deletes a student from the database
func (d *Database) DeleteStudent(ctx context.Context, id string) error {
	_, err := d.Client.ExecContext(
		ctx,
		`DELETE FROM student WHERE id = ?`,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete comment from the database: %w", err)
	}
	return nil
}

// UpdateStudent - updates a student in the database
func (d *Database) UpdateStudent(ctx context.Context, id string, std student.Student) (student.Student, error) {
	dateOfBirthNullTime := sql.NullTime{Time: std.DateOfBirth.Time, Valid: !std.DateOfBirth.IsZero()}
	stdRow := StudentRow{
		ID:          id,
		Fname:       sql.NullString{String: std.Fname, Valid: true},
		Lname:       sql.NullString{String: std.Lname, Valid: true},
		Email:       sql.NullString{String: std.Email, Valid: true},
		Gender:      sql.NullString{String: std.Gender, Valid: true},
		DateOfBirth: dateOfBirthNullTime,
		Address:     sql.NullString{String: std.Address, Valid: true},
		CreatedBy:   sql.NullString{String: std.CreatedBy, Valid: true},
		CreatedOn:   sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedBy:   sql.NullString{String: std.UpdatedBy, Valid: true},
		UpdatedOn:   sql.NullTime{Time: time.Now(), Valid: true},
	}

	rows, err := d.Client.NamedExecContext(
		ctx,
		`UPDATE student 
		SET fname = :fname,
			lname = :lname,
			email = :email,
			gender = :gender,
			dateofbirth = :dateofbirth,
			address = :address,
			createdby = :createdby,
			createdon = :createdon,
			updatedby = :updatedby,
			updatedon = :updatedon
		WHERE id = :id`,
		stdRow,
	)
	if err != nil {
		return student.Student{}, fmt.Errorf("failed to insert Student: %w", err)
	}
	// if err := rows.Close(); err != nil {
	// 	return comment.Comment{}, fmt.Errorf("failed to close rows: %w", err)
	// }
	insertedID, err := rows.LastInsertId()
	if err != nil {
		return student.Student{}, fmt.Errorf("failed to close rows: %w", err)
	}
	fmt.Printf("New record ID: %d\n", insertedID)

	return convertStudentRowToStudent(stdRow), nil
}
