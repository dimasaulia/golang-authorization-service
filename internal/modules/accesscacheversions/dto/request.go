package dto

type CreateAccessCacheVersionRequest struct {
	UserId  int64 `json:"user_id"`
	AppId   int64 `json:"app_id"`
	Version int64 `json:"version"`
}

type UpdateAccessCacheVersionRequest struct {
	UserId  *int64 `json:"user_id"`
	AppId   *int64 `json:"app_id"`
	Version *int64 `json:"version"`
}
