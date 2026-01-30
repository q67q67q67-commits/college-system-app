package repository

import (
	"database/sql"
	"time"

	"github.com/q67q67q67-commits/college-system-app/internal/db"
)

// GradeRow — одна оценка.
type GradeRow struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	ScheduleID int64     `json:"schedule_id"`
	Grade      float64   `json:"grade"`
	GradeDate  time.Time `json:"grade_date"`
	Comment    string    `json:"comment"`
	Subject    string    `json:"subject,omitempty"`
}

// GradesByUserID возвращает оценки студента.
func GradesByUserID(userID int64) ([]GradeRow, error) {
	rows, err := db.DB.Query(`
		SELECT g.id, g.user_id, g.schedule_id, g.grade, g.grade_date, COALESCE(g.comment,''), s.subject
		FROM grades g
		JOIN schedules s ON s.id = g.schedule_id
		WHERE g.user_id = $1 ORDER BY g.grade_date DESC, s.subject
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []GradeRow
	for rows.Next() {
		var gr GradeRow
		err := rows.Scan(&gr.ID, &gr.UserID, &gr.ScheduleID, &gr.Grade, &gr.GradeDate, &gr.Comment, &gr.Subject)
		if err != nil {
			return nil, err
		}
		list = append(list, gr)
	}
	return list, rows.Err()
}

// GPAByUserID считает средний балл студента.
func GPAByUserID(userID int64) (float64, error) {
	var gpa sql.NullFloat64
	err := db.DB.QueryRow(`SELECT AVG(grade) FROM grades WHERE user_id = $1`, userID).Scan(&gpa)
	if err != nil {
		return 0, err
	}
	if gpa.Valid {
		return gpa.Float64, nil
	}
	return 0, nil
}

// CreateGrade добавляет оценку (преподаватель).
func CreateGrade(userID, scheduleID int64, grade float64, gradeDate string, comment string) (int64, error) {
	var id int64
	err := db.DB.QueryRow(`
		INSERT INTO grades (user_id, schedule_id, grade, grade_date, comment)
		VALUES ($1, $2, $3, $4::date, $5) RETURNING id
	`, userID, scheduleID, grade, gradeDate, comment).Scan(&id)
	return id, err
}

// UpdateGrade обновляет оценку.
func UpdateGrade(id int64, grade float64, comment string) error {
	_, err := db.DB.Exec(`UPDATE grades SET grade = $1, comment = $2 WHERE id = $3`, grade, comment, id)
	return err
}

// DeleteGrade удаляет оценку.
func DeleteGrade(id int64) error {
	_, err := db.DB.Exec(`DELETE FROM grades WHERE id = $1`, id)
	return err
}
