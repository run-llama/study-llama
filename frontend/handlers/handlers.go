package handlers

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/run-llama/study-llama/frontend/agent"
	"github.com/run-llama/study-llama/frontend/auth"
	db "github.com/run-llama/study-llama/frontend/authdb"
	"github.com/run-llama/study-llama/frontend/files"
	"github.com/run-llama/study-llama/frontend/filesdb"
	"github.com/run-llama/study-llama/frontend/rules"
	"github.com/run-llama/study-llama/frontend/rulesdb"
	"github.com/run-llama/study-llama/frontend/templates"
)

func HandleSignUp(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	passwordR := c.FormValue("passwordRepeat")
	if password != passwordR {
		return c.SendStatus(400)
	}
	ctx := context.Background()
	sqlDb, err := auth.CreateNewDb()
	if err != nil {
		return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}
	queries := db.New(sqlDb)
	_, err = queries.GetUser(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			hashed_psw, err := auth.HashPassword(password)
			if err != nil {
				return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
			}
			_, err = queries.CreateUser(ctx, db.CreateUserParams{Username: username, HashedPassword: hashed_psw})
			if err != nil {
				return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
			} else {
				return templates.StatusBanner(nil).Render(c.Context(), c.Response().BodyWriter())
			}
		} else {
			return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
		}
	} else {
		return templates.StatusBanner(errors.New("user already exists")).Render(c.Context(), c.Response().BodyWriter())
	}
}

func HandleLogin(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	ctx := context.Background()

	sqlDb, err := auth.CreateNewDb()
	if err != nil {
		return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}
	queries := db.New(sqlDb)
	user, err := queries.GetUser(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return templates.StatusBanner(errors.New("there is no user with this username")).Render(c.Context(), c.Response().BodyWriter())
		} else {
			return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
		}
	}
	if !auth.CompareHashToPassword(password, user.HashedPassword) {
		return templates.StatusBanner(errors.New("wrong username or password")).Render(c.Context(), c.Response().BodyWriter())
	} else {
		sess_token, errSes := auth.GenerateToken(32)
		csrf_token, errCsrf := auth.GenerateToken(32)
		if errSes != nil || errCsrf != nil {
			return templates.StatusBanner(errors.New("an error occurred while generating your authentication credentials")).Render(c.Context(), c.Response().BodyWriter())
		}
		err = queries.UpdateUserTokensLogin(ctx, db.UpdateUserTokensLoginParams{SessionToken: pgtype.Text{String: sess_token, Valid: true}, CsrfToken: pgtype.Text{String: csrf_token, Valid: true}, Username: username})
		if err != nil {
			return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
		} else {
			c.Cookie(&fiber.Cookie{
				Name:     "session_token",
				Value:    sess_token,
				Expires:  time.Now().Add(24 * time.Hour),
				HTTPOnly: true,
			})
			c.Cookie(&fiber.Cookie{
				Name:     "csrf_token",
				Value:    csrf_token,
				Expires:  time.Now().Add(24 * time.Hour),
				HTTPOnly: false,
			})
			c.Set("HX-Redirect", "/categories")
			return c.SendStatus(fiber.StatusOK)
		}
	}
}

func HandleLogout(c *fiber.Ctx) error {
	_, err := auth.AuthorizePost(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "An error occurred: " + err.Error()})
	} else {
		ctx := context.Background()
		sqlDb, err := auth.CreateNewDb()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "An error occurred: " + err.Error()})
		}
		queries := db.New(sqlDb)
		st := c.Cookies("session_token", "")
		csrf := c.Cookies("csrf_token", "")
		err = queries.UpdateUserTokensLogout(ctx, db.UpdateUserTokensLogoutParams{SessionToken: pgtype.Text{String: st, Valid: true}, CsrfToken: pgtype.Text{String: csrf, Valid: true}})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "An error occurred: " + err.Error()})
		}
		c.Set("HX-Redirect", "/")
		return c.SendStatus(fiber.StatusOK)
	}
}

func HandleCreateRule(c *fiber.Ctx) error {
	user, err := auth.AuthorizePost(c)
	if err != nil {
		return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}
	ruleName := c.FormValue("rule_name")
	ruleType := c.FormValue("rule_type")
	ruleDes := c.FormValue("rule_description")
	db, err := rules.CreateNewDb()
	if err != nil {
		return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}
	queries := rulesdb.New(db)
	rules, err := queries.GetRules(context.Background(), user.Username)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
		}
	}
	rule, err := queries.CreateRule(context.Background(), rulesdb.CreateRuleParams{Username: user.Username, RuleName: ruleName, RuleType: ruleType, RuleDescription: ruleDes})
	if err != nil {
		return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}
	rules = append(rules, rule)
	return templates.RulesList(rules).Render(c.Context(), c.Response().BodyWriter())
}

func HandleUpdateRule(c *fiber.Ctx) error {
	user, err := auth.AuthorizePost(c)
	if err != nil {
		return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}
	ruleName := c.FormValue("rule_name")
	ruleType := c.FormValue("rule_type")
	ruleDes := c.FormValue("rule_description")
	db, err := rules.CreateNewDb()
	if err != nil {
		return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}
	queries := rulesdb.New(db)
	err = queries.UpdateRule(context.Background(), rulesdb.UpdateRuleParams{Username: user.Username, RuleName: ruleName, RuleType: ruleType, RuleDescription: ruleDes})
	if err != nil {
		return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}
	rules, err := queries.GetRules(context.Background(), user.Username)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
		}
	}
	return templates.RulesList(rules).Render(c.Context(), c.Response().BodyWriter())
}

func HandleDeleteRule(c *fiber.Ctx) error {
	user, err := auth.AuthorizePost(c)
	if err != nil {
		return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}
	ruleId := c.Params("id")
	ruleIdInt, err := strconv.Atoi(ruleId)
	if err != nil {
		return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}
	db, err := rules.CreateNewDb()
	if err != nil {
		return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}
	queries := rulesdb.New(db)
	err = queries.DeleteRule(context.Background(), int32(ruleIdInt))
	if err != nil {
		return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}
	rules, err := queries.GetRules(context.Background(), user.Username)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
		}
	}
	return templates.RulesList(rules).Render(c.Context(), c.Response().BodyWriter())
}

func HandleUploadFile(c *fiber.Ctx) error {
	user, err := auth.AuthorizePost(c)
	c.Set("Content-Type", "text/html")
	if err != nil {
		return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}
	file, err := c.FormFile("upload_file")
	if err != nil {
		return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}
	src, err := file.Open()
	if err != nil {
		return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}
	defer src.Close()
	fileId, err := files.UploadFile(src, file.Filename)
	if err != nil {
		return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}
	response, err := agent.ProcessFile(agent.InputFileEvent{FileId: fileId, FileName: file.Filename, Username: user.Username})
	if err != nil {
		return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}
	if response.Result.Error != nil {
		return templates.StatusBanner(errors.New(*response.Result.Error)).Render(c.Context(), c.Response().BodyWriter())
	}
	db, err := files.CreateNewDb()
	if err != nil {
		return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}
	queries := filesdb.New(db)
	files, err := queries.GetFiles(context.Background(), user.Username)
	if err != nil {
		return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}
	return templates.FilesList(files).Render(c.Context(), c.Response().BodyWriter())
}

func HandleDeleteFile(c *fiber.Ctx) error {
	user, err := auth.AuthorizePost(c)
	c.Set("Content-Type", "text/html")
	if err != nil {
		return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}
	fileId := c.Params("id")
	fileIdInt, err := strconv.Atoi(fileId)
	if err != nil {
		return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}
	db, err := files.CreateNewDb()
	if err != nil {
		return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}
	queries := filesdb.New(db)
	err = queries.DeleteFile(context.Background(), int32(fileIdInt))
	if err != nil {
		return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}
	files, err := queries.GetFiles(context.Background(), user.Username)
	if err != nil {
		return templates.StatusBanner(err).Render(c.Context(), c.Response().BodyWriter())
	}
	return templates.FilesList(files).Render(c.Context(), c.Response().BodyWriter())
}

func LoginRoute(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodGet {
		return c.SendStatus(fiber.StatusMethodNotAllowed)
	}
	c.Set("Content-Type", "text/html")
	return templates.SignIn().Render(c.Context(), c.Response().BodyWriter())
}

func SignUpRoute(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodGet {
		return c.SendStatus(fiber.StatusMethodNotAllowed)
	}
	c.Set("Content-Type", "text/html")
	return templates.SignUp().Render(c.Context(), c.Response().BodyWriter())
}

func PageDoesNotExistRoute(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodGet {
		return c.SendStatus(fiber.StatusMethodNotAllowed)
	}
	c.Set("Content-Type", "text/html")
	return templates.Page404().Render(c.Context(), c.Response().BodyWriter())
}

func HomeRoute(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodGet {
		return c.SendStatus(fiber.StatusMethodNotAllowed)
	}
	_, err := auth.AuthorizeGet(c)
	c.Set("Content-Type", "text/html")
	return templates.Home(err == nil).Render(c.Context(), c.Response().BodyWriter())
}

func CategoriesRoute(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodGet {
		return c.SendStatus(fiber.StatusMethodNotAllowed)
	}
	user, err := auth.AuthorizeGet(c)
	c.Set("Content-Type", "text/html")
	if err != nil {
		return templates.AuthFailedPage().Render(c.Context(), c.Response().BodyWriter())
	}
	db, err := rules.CreateNewDb()
	if err != nil {
		return templates.Page500(err).Render(c.Context(), c.Response().BodyWriter())
	}
	queries := rulesdb.New(db)
	rules, err := queries.GetRules(context.Background(), user.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return templates.RulesPage(user.Username, []rulesdb.Rule{}).Render(c.Context(), c.Response().BodyWriter())
		}
		return templates.Page500(err).Render(c.Context(), c.Response().BodyWriter())
	}

	return templates.RulesPage(user.Username, rules).Render(c.Context(), c.Response().BodyWriter())
}

func FilesRoute(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodGet {
		return c.SendStatus(fiber.StatusMethodNotAllowed)
	}
	user, err := auth.AuthorizeGet(c)
	c.Set("Content-Type", "text/html")
	if err != nil {
		return templates.AuthFailedPage().Render(c.Context(), c.Response().BodyWriter())
	}
	db, err := files.CreateNewDb()
	if err != nil {
		return templates.Page500(err).Render(c.Context(), c.Response().BodyWriter())
	}
	queries := filesdb.New(db)
	files, err := queries.GetFiles(context.Background(), user.Username)
	if err != nil {
		return templates.Page500(err).Render(c.Context(), c.Response().BodyWriter())
	}
	return templates.FilesPage(files).Render(c.Context(), c.Response().BodyWriter())
}
