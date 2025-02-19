package internal

import "log"

const listStudents = `
SELECT student.*
FROM student
WHERE student.major = (?1)
  AND student.name IN (?2);`

func BuildINQueryExample() {
	major := "Math"
	names := []string{"Jacob", "Bob", "Jenny", "Kate"}
	query, args := buildINQuery(listStudents, []any{major, names})
	log.Printf("query: %s\n args: %v\n", query, args)
}
