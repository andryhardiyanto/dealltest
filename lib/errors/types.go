package errors

import "net/http"

const (
	TypeUnauthorized         = "UNAUTHORIZED"
	TypeInternalServerError  = "INTERNAL_SERVER_ERROR"
	TypePanic                = "PANIC"
	TypeUnprocessableEntity  = "UNPROCESSABLE_ENTITY"
	TypeNotFound             = "NOT_FOUND"
	TypeUserNotFound         = "USER_NOT_FOUND"
	TypePasswordNotMatch     = "PASSWORD_NOT_MATCH"
	TypeBadRequest           = "BAD_REQUEST"
	TypeCannotDeleteSelf     = "CANNOT_DELETE_SELF"
	TypeInvalidRole          = "INVALID_ROLE"
	TypeInvalidCode          = "INVALID_CODE"
	TypeInvalidType          = "INVALID_TYPE"
	TypePageNotFound         = "PAGE_NOT_FOUND"
	TypePermissionDenied     = "PERMISSION_DENIED"
	TypeCronJobCurrencyError = "CRON_JOB_CURRENCY_ERROR"
	TypeRequestCannotEmpty   = "REQUEST_CANNOT_EMPTY"
	TypeEmailAlreadyExists   = "EMAIL_ALREADY_EXISTS"
)

var (
	ID = map[string]*Message{
		TypeCannotDeleteSelf: &Message{
			Language: "ID",
			Code:     http.StatusBadRequest,
			Message:  "tidak bisa delete account sendiri",
			Type:     TypeCannotDeleteSelf,
		},
		TypeInvalidRole: &Message{
			Language: "ID",
			Code:     http.StatusBadRequest,
			Message:  "Role tidak valid",
			Type:     TypeInvalidRole,
		},
		TypeEmailAlreadyExists: &Message{
			Language: "ID",
			Code:     http.StatusBadRequest,
			Message:  "Email sudah terdaftar",
			Type:     TypeEmailAlreadyExists,
		},
		TypeRequestCannotEmpty: &Message{
			Language: "ID",
			Code:     http.StatusBadRequest,
			Message:  "Request tidak boleh kosong",
			Type:     TypeRequestCannotEmpty,
		},
		TypePasswordNotMatch: &Message{
			Language: "ID",
			Code:     http.StatusBadRequest,
			Message:  "Password yang diinput salah",
			Type:     TypePasswordNotMatch,
		},
		TypeUserNotFound: &Message{
			Language: "ID",
			Code:     http.StatusNotFound,
			Message:  "User tidak ditemukan",
			Type:     TypeUserNotFound,
		},
		TypePermissionDenied: &Message{
			Code:     http.StatusForbidden,
			Language: "ID",
			Message:  "Anda tidak memiliki akses",
			Type:     TypePermissionDenied,
		},
		TypeUnauthorized: &Message{
			Code:     http.StatusUnauthorized,
			Language: "ID",
			Message:  "Anda tidak memiliki akses",
			Type:     TypeUnauthorized,
		},
		TypeInternalServerError: &Message{
			Language: "ID",
			Code:     http.StatusInternalServerError,
			Message:  "Server terjadi kesalahan, Mohon mencoba beberapa saat lagi",
			Type:     TypeInternalServerError,
		},
		TypePanic: &Message{
			Language: "ID",
			Code:     http.StatusUnprocessableEntity,
			Message:  "Request yang dikirim tidak valid, sehingga server tidak dapat memproses.",
			Type:     TypePanic,
		},
		TypeBadRequest: &Message{
			Language: "ID",
			Code:     http.StatusBadRequest,
			Message:  "Request yang dikirim tidak valid",
			Type:     TypeBadRequest,
		},
		TypeInvalidCode: &Message{
			Language: "ID",
			Code:     http.StatusBadRequest,
			Message:  "Code tidak valid",
			Type:     TypeInvalidCode,
		},
		TypeInvalidType: &Message{
			Language: "ID",
			Code:     http.StatusBadRequest,
			Message:  "Type tidak valid",
			Type:     TypeInvalidType,
		},
		TypeNotFound: &Message{
			Language: "ID",
			Code:     http.StatusNotFound,
			Message:  "Data tidak ditemukan",
			Type:     TypeNotFound,
		},
		TypeUnprocessableEntity: &Message{
			Language: "ID",
			Code:     http.StatusUnprocessableEntity,
			Message:  "Request yang dikirim tidak valid",
			Type:     TypeUnprocessableEntity,
		},
		TypePageNotFound: &Message{
			Language: "ID",
			Code:     http.StatusNotFound,
			Message:  "Halaman tidak ditemukan",
			Type:     TypePageNotFound,
		},
	}
)

var (
	Msg = map[string]map[string]*Message{
		"ID": ID,
	}
)
