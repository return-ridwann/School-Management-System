package models

// User represents an account in the users table
type User struct {
	ID           int64  `json:"id" db:"id"`
	FullName     string `json:"full_name" db:"full_name"`
	Email        string `json:"email" db:"email"`
	PasswordHash string `json:"-" db:"password_hash"`
	Role         string `json:"role" db:"role"`
	CreatedAt    string `json:"created_at" db:"created_at"`
}

// Class represents a classroom and section
type Class struct {
	ID         int64  `json:"id" db:"id"`
	Name       string `json:"name" db:"name"`
	RoomNumber string `json:"room_number" db:"room_number"`
}

// Student represents a student profile linked to a user
type Student struct {
	ID            int64  `json:"id" db:"id"`
	UserID        int64  `json:"user_id" db:"user_id"`
	ClassID       int64  `json:"class_id" db:"class_id"`
	RollNumber    int    `json:"roll_number" db:"roll_number"`
	GuardianPhone string `json:"guardian_phone" db:"guardian_phone"`
}

// Subject represents a course subject
type Subject struct {
	ID        int64  `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Code      string `json:"code" db:"code"`
	ClassID   int64  `json:"class_id" db:"class_id"`
	TeacherID int64  `json:"teacher_id" db:"teacher_id"`
}

// Attendance represents a daily attendance entry
type Attendance struct {
	ID        int64  `json:"id" db:"id"`
	StudentID int64  `json:"student_id" db:"student_id"`
	Date      string `json:"date" db:"date"`
	Status    string `json:"status" db:"status"`
	MarkedBy  int64  `json:"marked_by" db:"marked_by"`
}

// Mark represents an exam score
type Mark struct {
	ID        int64   `json:"id" db:"id"`
	StudentID int64   `json:"student_id" db:"student_id"`
	SubjectID int64   `json:"subject_id" db:"subject_id"`
	ExamType  string  `json:"exam_type" db:"exam_type"`
	Score     float64 `json:"score" db:"score"`
}

// Notice represents a school announcement
type Notice struct {
	ID        int64  `json:"id" db:"id"`
	Title     string `json:"title" db:"title"`
	Content   string `json:"content" db:"content"`
	AuthorID  int64  `json:"author_id" db:"author_id"`
	CreatedAt string `json:"created_at" db:"created_at"`
}

// combines student record with user and class info
type StudentDetail struct {
	ID            int64  `json:"id" db:"id"`
	UserID        int64  `json:"user_id" db:"user_id"`
	FullName      string `json:"full_name" db:"full_name"`
	Email         string `json:"email" db:"email"`
	ClassID       int64  `json:"class_id" db:"class_id"`
	ClassName     string `json:"class_name" db:"class_name"`
	RollNumber    int    `json:"roll_number" db:"roll_number"`
	GuardianPhone string `json:"guardian_phone" db:"guardian_phone"`
}

// combines subject info with teacher and class details
type SubjectDetail struct {
	ID          int64  `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Code        string `json:"code" db:"code"`
	ClassID     int64  `json:"class_id" db:"class_id"`
	ClassName   string `json:"class_name" db:"class_name"`
	TeacherID   int64  `json:"teacher_id" db:"teacher_id"`
	TeacherName string `json:"teacher_name" db:"teacher_name"`
}

// combines exam mark with student and subject details
type MarkDetail struct {
	ID          int64   `json:"id" db:"id"`
	StudentID   int64   `json:"student_id" db:"student_id"`
	StudentName string  `json:"student_name" db:"student_name"`
	RollNumber  int     `json:"roll_number" db:"roll_number"`
	SubjectID   int64   `json:"subject_id" db:"subject_id"`
	SubjectName string  `json:"subject_name" db:"subject_name"`
	ExamType    string  `json:"exam_type" db:"exam_type"`
	Score       float64 `json:"score" db:"score"`
}

// combines attendance record with student and marker details
type AttendanceDetail struct {
	ID          int64  `json:"id" db:"id"`
	StudentID   int64  `json:"student_id" db:"student_id"`
	StudentName string `json:"student_name" db:"student_name"`
	RollNumber  int    `json:"roll_number" db:"roll_number"`
	Date        string `json:"date" db:"date"`
	Status      string `json:"status" db:"status"`
	MarkedBy    int64  `json:"marked_by" db:"marked_by"`
	MarkerName  string `json:"marker_name" db:"marker_name"`
}

// combines notice with author name
type NoticeDetail struct {
	ID         int64  `json:"id" db:"id"`
	Title      string `json:"title" db:"title"`
	Content    string `json:"content" db:"content"`
	AuthorID   int64  `json:"author_id" db:"author_id"`
	AuthorName string `json:"author_name" db:"author_name"`
	CreatedAt  string `json:"created_at" db:"created_at"`
}
