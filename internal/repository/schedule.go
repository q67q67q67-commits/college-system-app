package repository

import (
	"database/sql"
	"time"

	"github.com/narxoz-college/nc/internal/db"
)

// ScheduleRow — одна запись расписания.
type ScheduleRow struct {
	ID            int64     `json:"id"`
	GroupID       int64     `json:"group_id"`
	TeacherID     int64     `json:"teacher_id"`
	Subject       string    `json:"subject"`
	Room          string    `json:"room"`
	DayOfWeek     int       `json:"day_of_week"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	Note          string    `json:"note"`
	AcademicPeriod string   `json:"academic_period"`
}

// ScheduleAttachmentRow — ДЗ/уведомление к паре.
type ScheduleAttachmentRow struct {
	ID             int64     `json:"id"`
	ScheduleID     int64     `json:"schedule_id"`
	Title          string    `json:"title"`
	Body           string    `json:"body"`
	AttachmentType string    `json:"attachment_type"`
	CreatedAt      time.Time `json:"created_at"`
}

// SchedulesByGroupID возвращает расписание группы.
func SchedulesByGroupID(groupID int64) ([]ScheduleRow, error) {
	rows, err := db.DB.Query(`
		SELECT id, group_id, teacher_id, subject, COALESCE(room,''), COALESCE(note,''),
		       day_of_week, start_time, end_time, COALESCE(academic_period,'')
		FROM schedules WHERE group_id = $1 ORDER BY day_of_week, start_time
	`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []ScheduleRow
	for rows.Next() {
		var s ScheduleRow
		err := rows.Scan(&s.ID, &s.GroupID, &s.TeacherID, &s.Subject, &s.Room, &s.Note,
			&s.DayOfWeek, &s.StartTime, &s.EndTime, &s.AcademicPeriod)
		if err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, rows.Err()
}

// SchedulesByTeacherID возвращает расписание преподавателя.
func SchedulesByTeacherID(teacherID int64) ([]ScheduleRow, error) {
	rows, err := db.DB.Query(`
		SELECT id, group_id, teacher_id, subject, COALESCE(room,''), COALESCE(note,''),
		       day_of_week, start_time, end_time, COALESCE(academic_period,'')
		FROM schedules WHERE teacher_id = $1 ORDER BY day_of_week, start_time
	`, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []ScheduleRow
	for rows.Next() {
		var s ScheduleRow
		err := rows.Scan(&s.ID, &s.GroupID, &s.TeacherID, &s.Subject, &s.Room, &s.Note,
			&s.DayOfWeek, &s.StartTime, &s.EndTime, &s.AcademicPeriod)
		if err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, rows.Err()
}

// GroupIDByUserID возвращает group_id студента (первую группу).
func GroupIDByUserID(userID int64) (int64, error) {
	var id int64
	err := db.DB.QueryRow(`SELECT group_id FROM user_groups WHERE user_id = $1 LIMIT 1`, userID).Scan(&id)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return id, err
}

// AttachmentsByScheduleID возвращает дополнения к паре.
func AttachmentsByScheduleID(scheduleID int64) ([]ScheduleAttachmentRow, error) {
	rows, err := db.DB.Query(`
		SELECT id, schedule_id, title, COALESCE(body,''), attachment_type, created_at
		FROM schedule_attachments WHERE schedule_id = $1 ORDER BY created_at
	`, scheduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []ScheduleAttachmentRow
	for rows.Next() {
		var a ScheduleAttachmentRow
		err := rows.Scan(&a.ID, &a.ScheduleID, &a.Title, &a.Body, &a.AttachmentType, &a.CreatedAt)
		if err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, rows.Err()
}

// CreateAttachment создаёт ДЗ/уведомление к паре.
func CreateAttachment(scheduleID int64, title, body, attachmentType string) (int64, error) {
	var id int64
	err := db.DB.QueryRow(`
		INSERT INTO schedule_attachments (schedule_id, title, body, attachment_type)
		VALUES ($1, $2, $3, $4) RETURNING id
	`, scheduleID, title, body, attachmentType).Scan(&id)
	return id, err
}

// UpdateAttachment обновляет запись.
func UpdateAttachment(id int64, title, body string) error {
	_, err := db.DB.Exec(`UPDATE schedule_attachments SET title = $1, body = $2 WHERE id = $3`, title, body, id)
	return err
}

// DeleteAttachment удаляет запись.
func DeleteAttachment(id int64) error {
	_, err := db.DB.Exec(`DELETE FROM schedule_attachments WHERE id = $1`, id)
	return err
}
