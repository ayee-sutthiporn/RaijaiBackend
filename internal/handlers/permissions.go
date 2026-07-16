package handlers

import "gorm.io/gorm"

// getBookRole returns the caller's role in a book ("OWNER"/"EDITOR"/"VIEWER"), or "" if not a member.
func getBookRole(db *gorm.DB, bookID, userID string) string {
	var role string
	db.Table("book_members").
		Select("role").
		Where("book_id = ? AND user_id = ?", bookID, userID).
		Scan(&role)
	return role
}

// isBookMember reports whether userID belongs to bookID in any role.
func isBookMember(db *gorm.DB, bookID, userID string) bool {
	return getBookRole(db, bookID, userID) != ""
}

// canEditInBook reports whether userID may create/modify resources inside bookID.
func canEditInBook(db *gorm.DB, bookID, userID string) bool {
	role := getBookRole(db, bookID, userID)
	return role == "OWNER" || role == "EDITOR"
}
