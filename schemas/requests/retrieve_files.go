package requests

type FileResponse struct {
	ID          string `json:"file_id"`
	FileName    string `json:"file_name"`
	Path        string `json:"file_path"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at"`
	Size        int64  `json:"file_size"`
}
