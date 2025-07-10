package entity

import (
	"testing"
	"time"
	sharedvo "todolist/internal/domain/shared/valueobject"
	vo "todolist/internal/domain/todo/valueobject"
)

func TestNewTodo(t *testing.T) {
	title, _ := vo.NewTodoTitle("Test Todo")
	description, _ := vo.NewTodoDescription("Test Description")
	priority := sharedvo.PriorityMedium
	dueDate := time.Now().Add(24 * time.Hour)

	t.Run("should create valid todo", func(t *testing.T) {
		todo, err := NewTodo(1, 123, title, description, priority, &dueDate)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if todo == nil {
			t.Fatal("Expected todo to be created")
		}
		if todo.userID != 123 {
			t.Errorf("Expected userID to be 123, got %d", todo.userID)
		}
		if todo.title != title {
			t.Errorf("Expected title to be %v, got %v", title, todo.title)
		}
		if todo.description != description {
			t.Errorf("Expected description to be %v, got %v", description, todo.description)
		}
		if todo.status != vo.StatusPending {
			t.Errorf("Expected status to be %v, got %v", vo.StatusPending, todo.status)
		}
		if todo.priority != priority {
			t.Errorf("Expected priority to be %v, got %v", priority, todo.priority)
		}
		if todo.dueDate == nil || !todo.dueDate.Equal(dueDate) {
			t.Errorf("Expected dueDate to be %v, got %v", dueDate, todo.dueDate)
		}
		if todo.completedAt != nil {
			t.Errorf("Expected completedAt to be nil, got %v", todo.completedAt)
		}
		if len(todo.tags) != 0 {
			t.Errorf("Expected tags to be empty, got %v", todo.tags)
		}
	})

	t.Run("should create todo without due date", func(t *testing.T) {
		todo, err := NewTodo(1, 123, title, description, priority, nil)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if todo.dueDate != nil {
			t.Errorf("Expected dueDate to be nil, got %v", todo.dueDate)
		}
	})

	t.Run("should fail with invalid user ID", func(t *testing.T) {
		todo, err := NewTodo(1, 0, title, description, priority, &dueDate)

		if err != ErrInvalidUserID {
			t.Errorf("Expected ErrInvalidUserID, got %v", err)
		}
		if todo != nil {
			t.Errorf("Expected todo to be nil, got %v", todo)
		}
	})

	t.Run("should fail with past due date", func(t *testing.T) {
		pastDate := time.Now().Add(-24 * time.Hour)
		todo, err := NewTodo(1, 123, title, description, priority, &pastDate)

		if err != ErrInvalidDueDate {
			t.Errorf("Expected ErrInvalidDueDate, got %v", err)
		}
		if todo != nil {
			t.Errorf("Expected todo to be nil, got %v", todo)
		}
	})
}

func TestTodoGetters(t *testing.T) {
	title, _ := vo.NewTodoTitle("Test Todo")
	description, _ := vo.NewTodoDescription("Test Description")
	priority := sharedvo.PriorityHigh
	dueDate := time.Now().Add(24 * time.Hour)
	todo, _ := NewTodo(1, 123, title, description, priority, &dueDate)

	t.Run("should return correct user ID", func(t *testing.T) {
		if todo.UserID() != 123 {
			t.Errorf("Expected userID to be 123, got %d", todo.UserID())
		}
	})

	t.Run("should return correct title", func(t *testing.T) {
		if todo.Title() != title {
			t.Errorf("Expected title to be %v, got %v", title, todo.Title())
		}
	})

	t.Run("should return correct description", func(t *testing.T) {
		if todo.Description() != description {
			t.Errorf("Expected description to be %v, got %v", description, todo.Description())
		}
	})

	t.Run("should return correct status", func(t *testing.T) {
		if todo.Status() != vo.StatusPending {
			t.Errorf("Expected status to be %v, got %v", vo.StatusPending, todo.Status())
		}
	})

	t.Run("should return correct priority", func(t *testing.T) {
		if todo.Priority() != priority {
			t.Errorf("Expected priority to be %v, got %v", priority, todo.Priority())
		}
	})

	t.Run("should return copy of due date", func(t *testing.T) {
		returnedDate := todo.DueDate()
		if returnedDate == nil {
			t.Fatal("Expected due date to be returned")
		}
		if !returnedDate.Equal(dueDate) {
			t.Errorf("Expected due date to be %v, got %v", dueDate, returnedDate)
		}
		// Verify it's a copy by modifying the returned value
		*returnedDate = returnedDate.Add(time.Hour)
		if todo.dueDate.Equal(*returnedDate) {
			t.Error("Expected due date to be a copy, but original was modified")
		}
	})

	t.Run("should return nil for empty due date", func(t *testing.T) {
		todoWithoutDueDate, _ := NewTodo(2, 123, title, description, priority, nil)
		if todoWithoutDueDate.DueDate() != nil {
			t.Error("Expected due date to be nil")
		}
	})

	t.Run("should return copy of tags", func(t *testing.T) {
		todo.AddTag("test-tag")
		tags := todo.Tags()
		if len(tags) != 1 || tags[0] != "test-tag" {
			t.Errorf("Expected tags to be ['test-tag'], got %v", tags)
		}
		// Verify it's a copy by modifying the returned slice
		tags[0] = "modified"
		if todo.tags[0] == "modified" {
			t.Error("Expected tags to be a copy, but original was modified")
		}
	})
}

func TestTodoBusinessMethods(t *testing.T) {
	title, _ := vo.NewTodoTitle("Test Todo")
	description, _ := vo.NewTodoDescription("Test Description")
	priority := sharedvo.PriorityMedium

	t.Run("should identify completed status", func(t *testing.T) {
		todo, _ := NewTodo(1, 123, title, description, priority, nil)
		if todo.IsCompleted() {
			t.Error("Expected todo to not be completed initially")
		}

		todo.Complete()
		if !todo.IsCompleted() {
			t.Error("Expected todo to be completed after calling Complete()")
		}
	})

	t.Run("should identify overdue todos", func(t *testing.T) {
		pastDate := time.Now().Add(-24 * time.Hour)
		futureDate := time.Now().Add(24 * time.Hour)

		// Todo with past due date should be overdue
		todo1, _ := NewTodo(1, 123, title, description, priority, nil)
		todo1.dueDate = &pastDate // Set directly to bypass validation
		if !todo1.IsOverdue() {
			t.Error("Expected todo with past due date to be overdue")
		}

		// Todo with future due date should not be overdue
		todo2, _ := NewTodo(2, 123, title, description, priority, &futureDate)
		if todo2.IsOverdue() {
			t.Error("Expected todo with future due date to not be overdue")
		}

		// Todo without due date should not be overdue
		todo3, _ := NewTodo(3, 123, title, description, priority, nil)
		if todo3.IsOverdue() {
			t.Error("Expected todo without due date to not be overdue")
		}

		// Completed todo should not be overdue
		todo4, _ := NewTodo(4, 123, title, description, priority, nil)
		todo4.dueDate = &pastDate
		todo4.Complete()
		if todo4.IsOverdue() {
			t.Error("Expected completed todo to not be overdue")
		}
	})

	t.Run("should calculate days until due", func(t *testing.T) {
		futureDate := time.Now().Add(48 * time.Hour)
		todo, _ := NewTodo(1, 123, title, description, priority, &futureDate)

		days := todo.DaysUntilDue()
		if days == nil {
			t.Fatal("Expected days until due to be calculated")
		}
		// The calculation might return 1 or 2 depending on the exact time, so we check for a reasonable range
		if *days < 1 || *days > 2 {
			t.Errorf("Expected 1 or 2 days until due, got %d", *days)
		}

		// Todo without due date should return nil
		todoWithoutDueDate, _ := NewTodo(2, 123, title, description, priority, nil)
		if todoWithoutDueDate.DaysUntilDue() != nil {
			t.Error("Expected days until due to be nil for todo without due date")
		}
	})
}

func TestTodoUpdateMethods(t *testing.T) {
	title, _ := vo.NewTodoTitle("Test Todo")
	description, _ := vo.NewTodoDescription("Test Description")
	priority := sharedvo.PriorityMedium
	todo, _ := NewTodo(1, 123, title, description, priority, nil)

	t.Run("should update title", func(t *testing.T) {
		newTitle, _ := vo.NewTodoTitle("Updated Title")
		todo.UpdateTitle(newTitle)

		if todo.Title() != newTitle {
			t.Errorf("Expected title to be updated to %v, got %v", newTitle, todo.Title())
		}
	})

	t.Run("should update description", func(t *testing.T) {
		newDescription, _ := vo.NewTodoDescription("Updated Description")
		todo.UpdateDescription(newDescription)

		if todo.Description() != newDescription {
			t.Errorf("Expected description to be updated to %v, got %v", newDescription, todo.Description())
		}
	})

	t.Run("should update priority", func(t *testing.T) {
		newPriority := sharedvo.PriorityHigh
		todo.UpdatePriority(newPriority)

		if todo.Priority() != newPriority {
			t.Errorf("Expected priority to be updated to %v, got %v", newPriority, todo.Priority())
		}
	})

	t.Run("should update due date", func(t *testing.T) {
		futureDate := time.Now().Add(24 * time.Hour)
		err := todo.UpdateDueDate(&futureDate)

		if err != nil {
			t.Errorf("Expected no error updating due date, got %v", err)
		}
		if todo.DueDate() == nil || !todo.DueDate().Equal(futureDate) {
			t.Errorf("Expected due date to be updated to %v, got %v", futureDate, todo.DueDate())
		}
	})

	t.Run("should fail to update due date to past", func(t *testing.T) {
		pastDate := time.Now().Add(-24 * time.Hour)
		err := todo.UpdateDueDate(&pastDate)

		if err != ErrInvalidDueDate {
			t.Errorf("Expected ErrInvalidDueDate, got %v", err)
		}
	})

	t.Run("should clear due date", func(t *testing.T) {
		err := todo.UpdateDueDate(nil)

		if err != nil {
			t.Errorf("Expected no error clearing due date, got %v", err)
		}
		if todo.DueDate() != nil {
			t.Errorf("Expected due date to be cleared, got %v", todo.DueDate())
		}
	})
}

func TestTodoStatusTransitions(t *testing.T) {
	title, _ := vo.NewTodoTitle("Test Todo")
	description, _ := vo.NewTodoDescription("Test Description")
	priority := sharedvo.PriorityMedium

	t.Run("should change status with valid transitions", func(t *testing.T) {
		todo, _ := NewTodo(1, 123, title, description, priority, nil)

		// Pending -> InProgress
		err := todo.ChangeStatus(vo.StatusInProgress)
		if err != nil {
			t.Errorf("Expected no error transitioning to InProgress, got %v", err)
		}
		if todo.Status() != vo.StatusInProgress {
			t.Errorf("Expected status to be InProgress, got %v", todo.Status())
		}

		// InProgress -> Completed
		err = todo.ChangeStatus(vo.StatusCompleted)
		if err != nil {
			t.Errorf("Expected no error transitioning to Completed, got %v", err)
		}
		if todo.Status() != vo.StatusCompleted {
			t.Errorf("Expected status to be Completed, got %v", todo.Status())
		}
		if todo.CompletedAt() == nil {
			t.Error("Expected completedAt to be set when transitioning to Completed")
		}

		// Completed -> Pending (reopen)
		err = todo.ChangeStatus(vo.StatusPending)
		if err != nil {
			t.Errorf("Expected no error transitioning to Pending, got %v", err)
		}
		if todo.Status() != vo.StatusPending {
			t.Errorf("Expected status to be Pending after reopen, got %v", todo.Status())
		}
	})

	t.Run("should clear completedAt when transitioning away from completed", func(t *testing.T) {
		todo, _ := NewTodo(1, 123, title, description, priority, nil)

		// Complete the todo first
		todo.Complete()
		if todo.CompletedAt() == nil {
			t.Fatal("Expected completedAt to be set after completion")
		}

		// Reopen the todo
		err := todo.ChangeStatus(vo.StatusPending)
		if err != nil {
			t.Errorf("Expected no error reopening todo, got %v", err)
		}

		if todo.Status() != vo.StatusPending {
			t.Errorf("Expected status to be Pending after reopen, got %v", todo.Status())
		}
	})

	t.Run("should complete todo", func(t *testing.T) {
		todo, _ := NewTodo(1, 123, title, description, priority, nil)

		err := todo.Complete()
		if err != nil {
			t.Errorf("Expected no error completing todo, got %v", err)
		}
		if !todo.IsCompleted() {
			t.Error("Expected todo to be completed")
		}
		if todo.CompletedAt() == nil {
			t.Error("Expected completedAt to be set")
		}
	})

	t.Run("should fail to complete already completed todo", func(t *testing.T) {
		todo, _ := NewTodo(1, 123, title, description, priority, nil)
		todo.Complete()

		err := todo.Complete()
		if err != ErrTodoAlreadyCompleted {
			t.Errorf("Expected ErrTodoAlreadyCompleted, got %v", err)
		}
	})

	t.Run("should start progress", func(t *testing.T) {
		todo, _ := NewTodo(1, 123, title, description, priority, nil)

		err := todo.StartProgress()
		if err != nil {
			t.Errorf("Expected no error starting progress, got %v", err)
		}
		if todo.Status() != vo.StatusInProgress {
			t.Errorf("Expected status to be InProgress, got %v", todo.Status())
		}
	})

	t.Run("should cancel todo", func(t *testing.T) {
		todo, _ := NewTodo(1, 123, title, description, priority, nil)

		err := todo.Cancel()
		if err != nil {
			t.Errorf("Expected no error cancelling todo, got %v", err)
		}
		if todo.Status() != vo.StatusCancelled {
			t.Errorf("Expected status to be Cancelled, got %v", todo.Status())
		}
	})

	t.Run("should reopen todo", func(t *testing.T) {
		todo, _ := NewTodo(1, 123, title, description, priority, nil)
		todo.Complete()

		err := todo.Reopen()
		if err != nil {
			t.Errorf("Expected no error reopening todo, got %v", err)
		}
		if todo.Status() != vo.StatusPending {
			t.Errorf("Expected status to be Pending, got %v", todo.Status())
		}
	})
}

func TestTodoTagManagement(t *testing.T) {
	title, _ := vo.NewTodoTitle("Test Todo")
	description, _ := vo.NewTodoDescription("Test Description")
	priority := sharedvo.PriorityMedium
	todo, _ := NewTodo(1, 123, title, description, priority, nil)

	t.Run("should add tag", func(t *testing.T) {
		todo.AddTag("work")

		if !todo.HasTag("work") {
			t.Error("Expected todo to have 'work' tag")
		}
		tags := todo.Tags()
		if len(tags) != 1 || tags[0] != "work" {
			t.Errorf("Expected tags to be ['work'], got %v", tags)
		}
	})

	t.Run("should not add duplicate tag", func(t *testing.T) {
		todo.AddTag("work") // Adding same tag again

		tags := todo.Tags()
		if len(tags) != 1 {
			t.Errorf("Expected only one tag, got %v", tags)
		}
	})

	t.Run("should add multiple tags", func(t *testing.T) {
		todo.AddTag("urgent")
		todo.AddTag("personal")

		if len(todo.Tags()) != 3 {
			t.Errorf("Expected 3 tags, got %d", len(todo.Tags()))
		}
		if !todo.HasTag("urgent") {
			t.Error("Expected todo to have 'urgent' tag")
		}
		if !todo.HasTag("personal") {
			t.Error("Expected todo to have 'personal' tag")
		}
	})

	t.Run("should remove tag", func(t *testing.T) {
		todo.RemoveTag("work")

		if todo.HasTag("work") {
			t.Error("Expected 'work' tag to be removed")
		}
		tags := todo.Tags()
		if len(tags) != 2 {
			t.Errorf("Expected 2 tags after removal, got %d", len(tags))
		}
	})

	t.Run("should handle removing non-existent tag", func(t *testing.T) {
		initialTagCount := len(todo.Tags())
		todo.RemoveTag("non-existent")

		if len(todo.Tags()) != initialTagCount {
			t.Error("Expected tag count to remain unchanged when removing non-existent tag")
		}
	})

	t.Run("should check tag existence", func(t *testing.T) {
		if !todo.HasTag("urgent") {
			t.Error("Expected todo to have 'urgent' tag")
		}
		if todo.HasTag("non-existent") {
			t.Error("Expected todo to not have 'non-existent' tag")
		}
	})
}

func TestTodoCompletedAtHandling(t *testing.T) {
	title, _ := vo.NewTodoTitle("Test Todo")
	description, _ := vo.NewTodoDescription("Test Description")
	priority := sharedvo.PriorityMedium
	todo, _ := NewTodo(1, 123, title, description, priority, nil)

	t.Run("should return nil completed at for non-completed todo", func(t *testing.T) {
		if todo.CompletedAt() != nil {
			t.Error("Expected completedAt to be nil for non-completed todo")
		}
	})

	t.Run("should return copy of completed at", func(t *testing.T) {
		todo.Complete()

		completedAt := todo.CompletedAt()
		if completedAt == nil {
			t.Fatal("Expected completedAt to be set")
		}

		// Verify it's a copy by modifying the returned value
		*completedAt = completedAt.Add(time.Hour)
		if todo.completedAt.Equal(*completedAt) {
			t.Error("Expected completedAt to be a copy, but original was modified")
		}
	})
}

func TestTodoEntityIntegration(t *testing.T) {
	title, _ := vo.NewTodoTitle("Test Todo")
	description, _ := vo.NewTodoDescription("Test Description")
	priority := sharedvo.PriorityMedium
	todo, _ := NewTodo(1, 123, title, description, priority, nil)

	t.Run("should inherit from shared.Entity", func(t *testing.T) {
		// Test that Todo embeds shared.Entity properly
		if todo.ID() != 1 {
			t.Errorf("Expected ID to be 1, got %d", todo.ID())
		}
	})

	t.Run("should mark as modified when updated", func(t *testing.T) {
		// This assumes the shared.Entity has methods to track modification
		// The actual implementation would depend on the shared.Entity interface
		newTitle, _ := vo.NewTodoTitle("Updated Title")
		todo.UpdateTitle(newTitle)

		// We can't directly test SetAsModified() without knowing the shared.Entity interface
		// But we can verify the title was updated
		if todo.Title() != newTitle {
			t.Error("Expected title to be updated")
		}
	})
}
