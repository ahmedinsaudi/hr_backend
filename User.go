 package main

// import (
// 	"context"
// 	"database/sql"
// 	"strconv"
// 	"time"

// 	"github.com/gofiber/fiber/v3"
// 	"github.com/jmoiron/sqlx"
// 	"github.com/rs/zerolog"
// 	"golang.org/x/crypto/bcrypt"
// )


// type UserRepository interface {
// 	GetUseres(ctx context.Context, id int, pagination Pagination) ([]User, error)
// 	GetUser(ctx context.Context, id int) (User, error)
// 	GetUserByEmail(ctx context.Context, email string ) (*User, error)
// 	AddUser(ctx context.Context, User User) (int, error)
// 	UpdateUser(ctx context.Context, User User) (int, error)
// 	DeleteUser(ctx context.Context, id int) (int, error)
// }

// type PosUserRepository struct {
// 	DB *sqlx.DB
// }

// func NewPosUserRepository(db *sqlx.DB) UserRepository {
// 	return &PosUserRepository{DB: db}
// }

// func (r *PosUserRepository) GetUseres(ctx context.Context, id int, pagination Pagination) ([]User, error) {
// 	data := []User{}
// 	offset := (pagination.Page - 1) * pagination.Limit
// 	err := r.DB.SelectContext(ctx, &data, "SELECT * FROM Useres  LIMIT $1 OFFSET $2",pagination.Limit,offset)
// 	if err != nil {
// 		return []User{}, err
// 	}
// 	return data, err
// }

// func (r *PosUserRepository) GetUser(ctx context.Context, id int) (User, error) {
// 	data := new(User)
// 	err := r.DB.GetContext(ctx, &data, "SELECT * FROM Useres WHERE id =$1", id)
// 	if err != nil {
// 		return User{}, err
// 	}
// 	return *data, err
// }


// func (r *PosUserRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
// 	data := new(User)
// 	err := r.DB.GetContext(ctx, &data, "SELECT * FROM Useres WHERE email =$1", email)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return data, err
// }

// func (r *PosUserRepository) AddUser(ctx context.Context, user User) (int, error) {
//     // تحقق يدوي إن إما الإيميل موجود أو رقم الهاتف
//     if user.Email == "" && (user.Phone == nil || *user.Phone == "") {
//         return 0, nil
//     }

//     // تحقق من بعض الحقول الأخرى حسب الحاجة (مثلاً الاسم)
//     if len(user.Name) < 3 {
//         return 0, nil
//     }

//     // باقي عملية الإدخال مثل السابق
//     var insertedID int

//     query := `
//         INSERT INTO users (
//             name, email, password, user_type, created_at, role, 
//             loyalty_points, address, phone, user_level
//         ) VALUES (
//             :name, :email, :password, :user_type, :created_at, :role,
//             :loyalty_points, :address, :phone, :user_level
//         ) RETURNING id
//     `

//     if user.CreatedAt.IsZero() {
//         user.CreatedAt = time.Now()
//     }

//     rows, err := r.DB.NamedQueryContext(ctx, query, user)
//     if err != nil {
//         return 0, err
//     }
//     defer rows.Close()

//     if rows.Next() {
//         err = rows.Scan(&insertedID)
//         if err != nil {
//             return 0, err
//         }
//     } else if err = rows.Err(); err != nil {
//         return 0, err
//     } else {
//         return 0, sql.ErrNoRows
//     }

//     return insertedID, nil
// }


// func (r *PosUserRepository) UpdateUser(ctx context.Context, User User) (int, error) {

// 	result, err := r.DB.NamedExecContext(ctx, "Update Useres Set Name=:name where  id=:id", User)
// 	if err != nil {
// 		return 0, err
// 	}
// 	id, err := result.RowsAffected()
// 	if err != nil {
// 		return 0, err
// 	}
// 	return int(id), err
// }

// func (r *PosUserRepository) DeleteUser(ctx context.Context, UserId int) (int, error) {
// 	result, err := r.DB.ExecContext(ctx, "delete from Useres where id=$1", UserId)
// 	if err != nil {
// 		return 0, err
// 	}
// 	id, err := result.RowsAffected()
// 	if err != nil {
// 		return 0, err
// 	}
// 	return int(id), err
// }

// type UserService struct {
// 	Logger zerolog.Logger
// 	Repo   UserRepository
// }

// func NewUserService(logger zerolog.Logger, repo UserRepository) *UserService {
// 	return &UserService{Logger: logger.With().Str("layer", "service").Logger(), Repo: repo}
// }

// func (s *UserService) GetUseres(ctx context.Context, id int, pagination Pagination) ([]User, error) {
// 	data, err := s.Repo.GetUseres(ctx, id, pagination)
// 	if err != nil {
// 		serviceLogger := s.Logger.With().Caller().Logger()
// 		serviceLogger.Error().Err(err).Str("userID", "userID").Msg("Failed to fetch User from repository")
// 	}
// 	return data, err
// }

// func (s *UserService) GetUser(ctx context.Context, id int) (User, error) {
// 	data, err := s.Repo.GetUser(ctx, id)
// 	if err != nil {
// 		serviceLogger := s.Logger.With().Caller().Logger()
// 		serviceLogger.Error().Err(err).Str("userID", "userID").Msg("Failed to fetch User from repository where  id")
// 	}
// 	return data, err
// }

// func (s *UserService) AddUser(ctx context.Context, User User) (int, error) {
// 	data, err := s.Repo.AddUser(ctx, User)
// 	if err != nil {
// 		serviceLogger := s.Logger.With().Caller().Logger()
// 		serviceLogger.Error().Err(err).Str("userID", "userID").Msg("Failed to Add User from repository")
// 	}
// 	return data, err
// }

// func (s *UserService) UpdateUser(ctx context.Context, User User) (int, error) {
// 	data, err := s.Repo.UpdateUser(ctx, User)
// 	if err != nil {
// 		serviceLogger := s.Logger.With().Caller().Logger()
// 		serviceLogger.Error().Err(err).Str("userID", "userID").Msg("Failed to Update User from repository")
// 	}

// 	return data, err
// }

// func (s *UserService) DeleteUser(ctx context.Context, id int) (int, error) {
// 	data, err := s.Repo.DeleteUser(ctx, id)
// 	if err != nil {
// 		serviceLogger := s.Logger.With().Caller().Logger()
// 		serviceLogger.Error().Err(err).Str("userID", "userID").Msg("Failed to delete User from repository")
// 	}
// 	return data, err
// }


// type UserHandler struct {
// 	Logger  zerolog.Logger
// 	Service *UserService
// }

// func NewUserHandler(logger zerolog.Logger, serv *UserService) *UserHandler {
// 	return &UserHandler{Logger: logger.With().Str("layer", "service").Logger(), Service: serv}
// }

// func (h *UserHandler) GetUseres(ctx fiber.Ctx) error {
// 	companyId := ctx.Query("company_id")
// 	pagination := GetPagination(ctx)
// 	id, _ := strconv.Atoi(companyId)
// 	data, err := h.Service.GetUseres(ctx.Context(), id, pagination)

// 	if err != nil {
// 		handleLogging( h.Logger,err,"Failed to fetch Users from repository")
// 		return ctx.Status(401).JSON(fiber.Map{"error": "rrrr"})
// 	}

// 	return ctx.JSON(data)
// }

// func (h *UserHandler) GetUser(ctx fiber.Ctx) error {
// 	UserId := ctx.Query("User_id")
// 	id, _ := strconv.Atoi(UserId)
// 	data, err := h.Service.GetUser(ctx.Context(), id)

// 	if err != nil {
// 		handleLogging( h.Logger,err,"Failed to get User from repository")
// 		return ctx.Status(401).JSON(fiber.Map{"error": "rrrr"})
// 	}

// 	return ctx.JSON(data)
// }

// func (h *UserHandler) AddUser(ctx fiber.Ctx) error {

// 	User := ctx.Locals("bodyParse").(*User)

// 	datas, err := h.Service.AddUser(ctx.Context(), *User)

// 	if err != nil {
// 		handleLogging( h.Logger,err,"Failed to add User from repository")
// 		return ctx.Status(401).JSON(fiber.Map{"error": err.Error()})
// 	}

// 	return ctx.JSON(fiber.Map{"id": datas})
// }

// func (h *UserHandler) UpdateUser(ctx fiber.Ctx) error {
// 	User := ctx.Locals("bodyParse").(*User)

// 	data, err := h.Service.UpdateUser(ctx.Context(), *User)

// 	if err != nil {
// 		handleLogging( h.Logger,err,"Failed to update User from repository")
// 		return ctx.Status(401).JSON(fiber.Map{"error": "rrrr"})
// 	}

// 	return ctx.JSON(data)
// }

// func (h *UserHandler) DeleteUser(ctx fiber.Ctx) error {
// 	UserId := ctx.Query("User_id")
// 	id, _ := strconv.Atoi(UserId)
// 	data, err := h.Service.DeleteUser(ctx.Context(), id)

// 	if err != nil {
// 		handleLogging( h.Logger,err,"Failed to delete User from repository")
// 		return ctx.Status(401).JSON(fiber.Map{"error": "Failed to delete User from repository"})
// 	}

// 	return ctx.JSON(data)
// }



// func (h *UserHandler) Signup(c fiber.Ctx) error {


// 	req := c.Locals("bodyParse").(*SignupRequest)

// 	if req.Email == "" || req.Password == "" {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email and password are required"})
// 	}

// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not hash password"})
// 	}

// 	newUser := User{
// 		Email:    req.Email,
// 		Password: string(hashedPassword),
// 		Role:    req.Role,
// 	}

// 	newId,_ := h.Service.Repo.AddUser(c.Context(),newUser)

// 	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User registered successfully", "user_id": newId})
// }

// func (h *UserHandler) Login(c fiber.Ctx) error {

// 	creds := c.Locals("bodyParse").(*LoginCredentials)

// 	userFound,_ := h.Service.Repo.GetUserByEmail(c.Context(),creds.Email)

// 	if userFound == nil {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
// 	}

// 	if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(creds.Password)); err != nil {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
// 	}

// 	token, err := GenerateJWTToken(*userFound)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
// 	}

// 	return c.JSON(TokenResponse{Token: token})
// }
