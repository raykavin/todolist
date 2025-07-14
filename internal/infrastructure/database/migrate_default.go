package database

/*
 * migrate_default.go
 *
 * This file defines default database migration.
 *
 * It should provide tools to apply, rollback, and version control schema changes
 * to keep the database structure in sync with your domain model.
 *
 */

import (
	"fmt"
	"todolist/internal/infrastructure/database/model"

	"gorm.io/gorm"
)

// MigrateDefault runs all database migrations for default database
func MigrateDefault(db *gorm.DB) error {
	models := []any{
		model.AuditLog{},
		model.LoginAttempt{},
		model.Person{},
		model.Tag{},
		model.Todo{},
		model.TodoDailyStatistics{},
		model.TodoTag{},
		model.TodoView{},
		model.User{},
		model.UserStatisticsView{},
	}

	// Auto migrate all models
	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	// Create indexes
	if err := createIndexes(db); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	// Create views
	if err := createViews(db); err != nil {
		return fmt.Errorf("failed to create views: %w", err)
	}

	// Create triggers
	if err := createTriggers(db); err != nil {
		return fmt.Errorf("failed to create triggers: %w", err)
	}

	return nil
}

// createIndexes creates additional indexes for better performance
func createIndexes(db *gorm.DB) error {
	indexes := []string{
		// Composite indexes for common queries
		`CREATE INDEX IF NOT EXISTS idx_todos_user_status ON todos(user_id, status) WHERE deleted_at IS NULL`,
		`CREATE INDEX IF NOT EXISTS idx_todos_user_priority ON todos(user_id, priority) WHERE deleted_at IS NULL`,
		`CREATE INDEX IF NOT EXISTS idx_todos_due_date_status ON todos(due_date, status) WHERE deleted_at IS NULL AND due_date IS NOT NULL`,

		// Partial indexes for active records
		`CREATE INDEX IF NOT EXISTS idx_users_active ON users(username) WHERE deleted_at IS NULL AND status = 'active'`,
		// Removi a condição com NOW()
		`CREATE INDEX IF NOT EXISTS idx_todos_overdue ON todos(user_id, due_date) WHERE deleted_at IS NULL AND status IN ('pending', 'in_progress')`,

		// Full-text search indexes
		`CREATE INDEX IF NOT EXISTS idx_todos_title_fts ON todos USING gin(to_tsvector('english', title))`,
		`CREATE INDEX IF NOT EXISTS idx_todos_description_fts ON todos USING gin(to_tsvector('english', description))`,
	}

	for _, index := range indexes {
		if err := db.Exec(index).Error; err != nil {
			return err
		}
	}

	return nil
}

// createViews creates database views for optimized queries
func createViews(db *gorm.DB) error {
	// Drop view
	for _, view := range []string{"todo_view", "user_statistics_view"} {
		if err := db.Exec(fmt.Sprintf(`DROP VIEW IF EXISTS %s CASCADE`, view)).Error; err != nil {
			return err
		}
	}

	views := []string{
		`CREATE VIEW todo_view AS
		SELECT
			t.id,
			t.user_id,
			u.username,
			p.name as person_name,
			t.title,
			t.description,
			t.status,
			t.priority,
			t.due_date,
			t.completed_at,
			CASE
				WHEN t.due_date < NOW() AND t.status IN ('pending', 'in_progress') THEN true
				ELSE false
			END as is_overdue,
			STRING_AGG(tag.name, ', ' ORDER BY tag.name) as tags,
			t.created_at,
			t.updated_at
		FROM todos t
		INNER JOIN users u ON u.id = t.user_id AND u.deleted_at IS NULL
		INNER JOIN people p ON p.id = u.person_id AND p.deleted_at IS NULL
		LEFT JOIN todo_tags tt ON tt.todo_id = t.id
		LEFT JOIN tags tag ON tag.id = tt.tag_id AND tag.deleted_at IS NULL
		WHERE t.deleted_at IS NULL
		GROUP BY t.id, u.username, p.name`,

		`CREATE VIEW user_statistics_view AS
		SELECT
			u.id as user_id,
			u.username,
			p.name as person_name,
			COUNT(t.id) as total_todos,
			COUNT(CASE WHEN t.status = 'completed' THEN 1 END) as completed_todos,
			COUNT(CASE WHEN t.status = 'pending' THEN 1 END) as pending_todos,
			COUNT(CASE WHEN t.status = 'in_progress' THEN 1 END) as in_progress_todos,
			COUNT(CASE WHEN t.status = 'cancelled' THEN 1 END) as cancelled_todos,
			COUNT(CASE WHEN t.due_date < NOW() AND t.status IN ('pending', 'in_progress') THEN 1 END) as overdue_todos,
			CASE
				WHEN COUNT(t.id) > 0 THEN ROUND(COUNT(CASE WHEN t.status = 'completed' THEN 1 END)::numeric / COUNT(t.id) * 100, 2)
				ELSE 0
			END as completion_rate,
			COALESCE(MAX(t.updated_at), u.created_at) as last_activity_at
		FROM users u
		INNER JOIN people p ON p.id = u.person_id AND p.deleted_at IS NULL
		LEFT JOIN todos t ON t.user_id = u.id AND t.deleted_at IS NULL
		WHERE u.deleted_at IS NULL
		GROUP BY u.id, u.username, p.name, u.created_at`,
	}

	for _, view := range views {
		if err := db.Exec(view).Error; err != nil {
			return err
		}
	}

	return nil
}

// createTriggers creates database triggers
func createTriggers(db *gorm.DB) error {
	triggers := []string{
		// Function to update todo daily statistics
		`CREATE OR REPLACE FUNCTION update_todo_daily_statistics()
		RETURNS TRIGGER AS $$
		BEGIN
			-- Update statistics for status changes
			IF TG_OP = 'UPDATE' AND OLD.status != NEW.status THEN
				-- Insert or update daily statistics
				INSERT INTO todo_daily_statistics (date, user_id, created, completed, cancelled, avg_time_to_complete)
				VALUES (
					CURRENT_DATE,
					NEW.user_id,
					0,
					CASE WHEN NEW.status = 'completed' THEN 1 ELSE 0 END,
					CASE WHEN NEW.status = 'cancelled' THEN 1 ELSE 0 END,
					CASE
						WHEN NEW.status = 'completed' AND NEW.completed_at IS NOT NULL
						THEN EXTRACT(EPOCH FROM (NEW.completed_at - NEW.created_at)) / 3600
						ELSE 0
					END
				)
				ON CONFLICT (date, user_id) DO UPDATE SET
					completed = todo_daily_statistics.completed + EXCLUDED.completed,
					cancelled = todo_daily_statistics.cancelled + EXCLUDED.cancelled,
					avg_time_to_complete = CASE
						WHEN todo_daily_statistics.completed + EXCLUDED.completed > 0
						THEN ((todo_daily_statistics.avg_time_to_complete * todo_daily_statistics.completed) + EXCLUDED.avg_time_to_complete) / (todo_daily_statistics.completed + EXCLUDED.completed)
						ELSE 0
					END;
			END IF;

			-- Update statistics for new todos
			IF TG_OP = 'INSERT' THEN
				INSERT INTO todo_daily_statistics (date, user_id, created, completed, cancelled, avg_time_to_complete)
				VALUES (CURRENT_DATE, NEW.user_id, 1, 0, 0, 0)
				ON CONFLICT (date, user_id) DO UPDATE SET
					created = todo_daily_statistics.created + 1;
			END IF;

			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql`,

		// Trigger for todo statistics
		`DROP TRIGGER IF EXISTS trigger_update_todo_statistics ON todos`,
		`CREATE TRIGGER trigger_update_todo_statistics
		AFTER INSERT OR UPDATE OF status ON todos
		FOR EACH ROW
		EXECUTE FUNCTION update_todo_daily_statistics()`,

		// Function to prevent deleting person with active user
		`CREATE OR REPLACE FUNCTION prevent_person_delete_with_user()
		RETURNS TRIGGER AS $$
		BEGIN
			IF EXISTS (SELECT 1 FROM users WHERE person_id = OLD.id AND deleted_at IS NULL) THEN
				RAISE EXCEPTION 'Cannot delete person with active user account';
			END IF;
			RETURN OLD;
		END;
		$$ LANGUAGE plpgsql`,

		// Trigger to prevent person deletion
		`DROP TRIGGER IF EXISTS trigger_prevent_person_delete ON people`,
		`CREATE TRIGGER trigger_prevent_person_delete
		BEFORE DELETE ON people
		FOR EACH ROW
		EXECUTE FUNCTION prevent_person_delete_with_user()`,
	}

	for _, trigger := range triggers {
		if err := db.Exec(trigger).Error; err != nil {
			return err
		}
	}

	return nil
}
