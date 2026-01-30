package repository

import (
	"github.com/q67q67q67-commits/college-system-app/internal/db"
)

// GroupRow — группа с количеством студентов.
type GroupRow struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	AcademicYear string `json:"academic_year"`
	StudentCount int    `json:"student_count"`
}

// GroupStudentRow — студент в группе.
type GroupStudentRow struct {
	UserID   int64  `json:"user_id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	IsCurator bool  `json:"is_curator"`
}

// ListGroups возвращает все группы (для преподавателя/админа).
func ListGroups() ([]GroupRow, error) {
	rows, err := db.DB.Query(`
		SELECT g.id, g.name, COALESCE(g.description,''), COALESCE(g.academic_year,''),
		       (SELECT COUNT(*) FROM user_groups ug WHERE ug.group_id = g.id AND ug.is_curator = false)
		FROM groups g ORDER BY g.name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []GroupRow
	for rows.Next() {
		var gr GroupRow
		err := rows.Scan(&gr.ID, &gr.Name, &gr.Description, &gr.AcademicYear, &gr.StudentCount)
		if err != nil {
			return nil, err
		}
		list = append(list, gr)
	}
	return list, rows.Err()
}

// StudentsByGroupID возвращает студентов группы.
func StudentsByGroupID(groupID int64) ([]GroupStudentRow, error) {
	rows, err := db.DB.Query(`
		SELECT u.id, u.full_name, u.email, COALESCE(ug.is_curator, false)
		FROM user_groups ug
		JOIN users u ON u.id = ug.user_id
		WHERE ug.group_id = $1 ORDER BY u.full_name
	`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []GroupStudentRow
	for rows.Next() {
		var s GroupStudentRow
		err := rows.Scan(&s.UserID, &s.FullName, &s.Email, &s.IsCurator)
		if err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, rows.Err()
}
