    package models

    import "time"

    type Submission struct {
        ID        uint       `gorm:"primaryKey"`
        Name      string     `json:"name" binding:"required"`
        Grid      [][]string `json:"grid" gorm:"-"`
        GridData  string     `json:"grid_data"`
        CreatedAt time.Time  `json:"created_at"`
    }

    // Structure utilitaire pour les réponses API
    type APIResponse struct {
        Success bool   `json:"success"`
        Message string `json:"message"`
    }
