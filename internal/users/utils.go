package users


func cleanUp(m *Manager) error{
	query := `DELETE FROM users`
	_, err := m.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}