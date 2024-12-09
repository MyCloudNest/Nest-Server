package requests

type UploadFileRequest struct {
	FileName    string `form:"file_name,omitempty" validate:"omitempty,min=1,max=255"`
	FilePath    string `form:"file_path,omitempty"`
	Description string `form:"description,omitempty"`
}
