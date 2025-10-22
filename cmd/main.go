/**
 * FILE HEADER: cmd/main.go
 *
 * PURPOSE:
 * Titik masuk utama (entrypoint) untuk aplikasi SIMDOKPOL.
 * Versi final yang sadar-lokasi (location-aware) untuk build yang portabel.
 * Dengan perbaikan untuk kompatibilitas icon di Windows, Linux, dan macOS.
 * PERBAIKAN: Kompatibilitas path migrations dan template parsing untuk semua OS.
 * PERBAIKAN: Prioritas icon berbasis platform untuk systray yang optimal.
 * FITUR BARU: Auto-generate file .env dengan JWT secret yang aman.
 */
package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"simdokpol/internal/config"
	"simdokpol/internal/controllers"
	"simdokpol/internal/middleware"
	"simdokpol/internal/repositories"
	"simdokpol/internal/services"
	"strconv"
	"strings"
	"time"

	_ "simdokpol/docs"
	_ "time/tzdata"

	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	gormsqlite "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var version = "dev"

const (
	defaultPort = ":8080"
	localURL    = "http://localhost:8080"
)

func getExecutableDir() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("FATAL: Gagal mendapatkan path executable: %v", err)
	}
	return filepath.Dir(exePath)
}

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	exeDir := getExecutableDir()

	// Tentukan prioritas icon berdasarkan sistem operasi
	var iconPaths []string
	
	switch runtime.GOOS {
	case "windows":
		// Windows: prioritas .ico, fallback ke .png
		iconPaths = []string{
			filepath.Join(exeDir, "web", "static", "img", "icon.ico"),
			filepath.Join(exeDir, "web", "static", "img", "icon.png"),
		}
		log.Printf("INFO: Sistem operasi terdeteksi: Windows - mencoba format .ico terlebih dahulu")
	case "darwin":
		// macOS: prioritas .png, fallback ke .ico
		iconPaths = []string{
			filepath.Join(exeDir, "web", "static", "img", "icon.png"),
			filepath.Join(exeDir, "web", "static", "img", "icon.ico"),
		}
		log.Printf("INFO: Sistem operasi terdeteksi: macOS - mencoba format .png terlebih dahulu")
	default:
		// Linux dan OS lainnya: hanya .png
		iconPaths = []string{
			filepath.Join(exeDir, "web", "static", "img", "icon.png"),
		}
		log.Printf("INFO: Sistem operasi terdeteksi: Linux/%s - mencoba format .png", runtime.GOOS)
	}

	iconData, err := getIconFromPaths(iconPaths)
	if err != nil {
		log.Printf("PERINGATAN: Tidak dapat menemukan file icon untuk platform %s. Error: %v", runtime.GOOS, err)
		log.Printf("PERINGATAN: Systray akan menggunakan icon default sistem")
		log.Printf("INFO: Path yang dicoba: %v", iconPaths)
	} else {
		systray.SetIcon(iconData)
		log.Printf("INFO: Icon systray berhasil dimuat untuk platform: %s", runtime.GOOS)
	}

	systray.SetTitle("SIMDOKPOL")
	systray.SetTooltip("Sistem Informasi Manajemen Dokumen Kepolisian")

	mOpen := systray.AddMenuItem("Buka Aplikasi", "Buka SIMDOKPOL di browser")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Keluar", "Tutup aplikasi")

	go startWebServer()

	go func() {
		time.Sleep(1 * time.Second)
		notifyIconPath := filepath.Join(exeDir, "web", "static", "img", "icon.png")
		if err := beeep.Notify("SIMDOKPOL", "Aplikasi sedang berjalan di latar belakang.", notifyIconPath); err != nil {
			log.Printf("PERINGATAN: Gagal menampilkan notifikasi: %v", err)
		}
	}()

	go func() {
		time.Sleep(2 * time.Second)
		openBrowser(localURL)
	}()

	go func() {
		for {
			select {
			case <-mOpen.ClickedCh:
				openBrowser(localURL)
			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()
}

func onExit() {
	log.Println("INFO: Aplikasi SIMDOKPOL ditutup.")
}

func startWebServer() {
	exeDir := getExecutableDir()

	if err := ensureEnvFile(exeDir); err != nil {
		log.Printf("PERINGATAN: Gagal memastikan file .env: %v", err)
	}

	envPath := filepath.Join(exeDir, ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("PERINGATAN: Tidak dapat menemukan file .env di %s, menggunakan environment variables sistem.", envPath)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("FATAL: Gagal memuat konfigurasi: %v", err)
	}

	if !filepath.IsAbs(cfg.DBDSN) {
		dbName := strings.Split(cfg.DBDSN, "?")[0]
		queryParams := ""
		if strings.Contains(cfg.DBDSN, "?") {
			queryParams = "?" + strings.Split(cfg.DBDSN, "?")[1]
		}
		absDbPath := filepath.Join(exeDir, dbName)
		cfg.DBDSN = absDbPath + queryParams
		log.Printf("INFO: Menggunakan path database absolut: %s", cfg.DBDSN)
	}

	db, err := setupDatabase(cfg.DBDSN, exeDir)
	if err != nil {
		log.Fatalf("FATAL: Gagal setup database: %v", err)
	}

	repos, svcs, ctrls := setupDependencies(db, cfg)
	router := setupRouter(repos.UserRepo, svcs, ctrls, exeDir)

	log.Printf("INFO: Server web dimulai di %s", localURL)

	if err := router.Run(defaultPort); err != nil {
		log.Fatalf("FATAL: Gagal menjalankan server: %v", err)
	}
}

func ensureEnvFile(exeDir string) error {
	envPath := filepath.Join(exeDir, ".env")

	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		log.Println("INFO: File .env tidak ditemukan, membuat file baru dengan konfigurasi aman...")
		return createEnvFile(envPath)
	}

	if err := godotenv.Load(envPath); err != nil {
		log.Printf("PERINGATAN: Gagal membaca file .env: %v", err)
		return createEnvFile(envPath)
	}

	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "" || jwtSecret == "ganti-dengan-secret-key-yang-kuat" || jwtSecret == "will-be-auto-generated" {
		log.Println("INFO: JWT_SECRET_KEY tidak valid, melakukan regenerasi...")
		return updateEnvFile(envPath)
	}

	log.Println("INFO: File .env sudah dikonfigurasi dengan baik")
	return nil
}

func createEnvFile(envPath string) error {
	jwtSecret, err := generateSecureSecret(64)
	if err != nil {
		return fmt.Errorf("gagal generate JWT secret: %w", err)
	}

	content := fmt.Sprintf(`# SIMDOKPOL Configuration File
# File ini di-generate otomatis pada: %s
# JANGAN BAGIKAN FILE INI KEPADA SIAPAPUN!

# JWT Secret Key (RAHASIA - Jangan dibagikan!)
JWT_SECRET_KEY=%s

# Database Configuration
DB_DSN=simdokpol.db?_foreign_keys=on

# Server Port
PORT=8080
`, time.Now().Format("2006-01-02 15:04:05"), jwtSecret)

	if err := os.WriteFile(envPath, []byte(content), 0600); err != nil {
		return fmt.Errorf("gagal menulis file .env: %w", err)
	}

	log.Printf("INFO: File .env berhasil dibuat di: %s", envPath)
	log.Println("INFO: JWT_SECRET_KEY telah di-generate secara otomatis dengan algoritma cryptographically secure")
	return nil
}

func updateEnvFile(envPath string) error {
	content, err := os.ReadFile(envPath)
	if err != nil {
		return createEnvFile(envPath)
	}

	jwtSecret, err := generateSecureSecret(64)
	if err != nil {
		return fmt.Errorf("gagal generate JWT secret: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	updated := false

	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "JWT_SECRET_KEY=") {
			lines[i] = fmt.Sprintf("JWT_SECRET_KEY=%s", jwtSecret)
			updated = true
			break
		}
	}

	if !updated {
		lines = append(lines, fmt.Sprintf("JWT_SECRET_KEY=%s", jwtSecret))
	}

	newContent := strings.Join(lines, "\n")
	if err := os.WriteFile(envPath, []byte(newContent), 0600); err != nil {
		return fmt.Errorf("gagal update file .env: %w", err)
	}

	log.Println("INFO: JWT_SECRET_KEY berhasil di-regenerate dengan nilai yang aman")
	return nil
}

func generateSecureSecret(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}

	return string(b), nil
}

func getIconFromPaths(paths []string) ([]byte, error) {
	var lastErr error
	var attemptedPaths []string
	
	for _, path := range paths {
		attemptedPaths = append(attemptedPaths, path)
		
		// Validasi file terlebih dahulu
		if err := validateIconFile(path); err != nil {
			log.Printf("DEBUG: Path '%s' tidak valid: %v", path, err)
			lastErr = err
			continue
		}
		
		data, err := os.ReadFile(path)
		if err != nil {
			log.Printf("DEBUG: Gagal membaca file '%s': %v", path, err)
			lastErr = err
			continue
		}
		
		if len(data) == 0 {
			err := errors.New("file icon kosong")
			log.Printf("DEBUG: File '%s' tidak memiliki konten", path)
			lastErr = err
			continue
		}
		
		log.Printf("INFO: Berhasil memuat icon dari: %s (ukuran: %d bytes)", path, len(data))
		return data, nil
	}
	
	if lastErr != nil {
		return nil, fmt.Errorf("tidak ada file icon yang valid ditemukan di path: %v. Error terakhir: %w", attemptedPaths, lastErr)
	}
	return nil, fmt.Errorf("tidak ada file icon yang valid ditemukan di path: %v", attemptedPaths)
}

func validateIconFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file tidak ditemukan")
		}
		return fmt.Errorf("file tidak dapat diakses: %w", err)
	}
	
	if info.IsDir() {
		return errors.New("path adalah direktori, bukan file")
	}
	
	if info.Size() == 0 {
		return errors.New("file icon kosong")
	}
	
	// Validasi ekstensi
	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".png" && ext != ".ico" {
		return fmt.Errorf("format file tidak didukung: %s (hanya .png dan .ico yang diizinkan)", ext)
	}
	
	return nil
}

func openBrowser(url string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		log.Printf("PERINGATAN: Sistem operasi tidak didukung untuk membuka browser secara otomatis: %s", runtime.GOOS)
		return
	}

	if err := cmd.Start(); err != nil {
		log.Printf("PERINGATAN: Gagal membuka browser secara otomatis: %v", err)
		log.Printf("INFO: Silakan buka browser manual dan akses: %s", url)
	}
}

func setupDatabase(dsn string, exeDir string) (*gorm.DB, error) {
	db, err := gorm.Open(gormsqlite.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		),
	})
	if err != nil {
		return nil, fmt.Errorf("gagal membuka koneksi gorm: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan instance sql.DB: %w", err)
	}

	driver, err := sqlite.WithInstance(sqlDB, &sqlite.Config{})
	if err != nil {
		return nil, fmt.Errorf("gagal membuat instance driver sqlite: %w", err)
	}

	migrationsPath := filepath.Join(exeDir, "migrations")
	migrationsURL := fmt.Sprintf("file://%s", filepath.ToSlash(migrationsPath))
	log.Printf("INFO: Menggunakan path migrasi absolut: %s", migrationsURL)

	m, err := migrate.NewWithDatabaseInstance(migrationsURL, "sqlite", driver)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat instance migrasi dari '%s': %w", migrationsURL, err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("gagal menjalankan migrasi: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		log.Printf("PERINGATAN: Gagal memeriksa versi migrasi: %v", err)
	} else if err == nil {
		log.Printf("INFO: Migrasi database berhasil. Versi: %d, Dirty: %t", version, dirty)
	} else {
		log.Println("INFO: Migrasi database berhasil (tidak ada versi diterapkan).")
	}

	return db, nil
}

func setupRouter(userRepo repositories.UserRepository, svcs Services, ctrls Controllers, exeDir string) *gin.Engine {
	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard
	}
	router := gin.Default()

	templatePath := filepath.Join(exeDir, "web", "templates")
	templates := template.Must(
		template.Must(
			template.New("").
				Funcs(template.FuncMap{"ToUpper": strings.ToUpper}).
				ParseGlob(filepath.Join(templatePath, "*.html")),
		).ParseGlob(filepath.Join(templatePath, "partials", "*.html")),
	)

	router.SetHTMLTemplate(templates)

	router.Static("/static", filepath.Join(exeDir, "web", "static"))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/setup", ctrls.ConfigController.ShowSetupPage)
	router.POST("/api/setup", ctrls.ConfigController.SaveSetup)

	app := router.Group("")
	app.Use(middleware.SetupMiddleware(svcs.ConfigService))
	{
		app.GET("/login", func(c *gin.Context) {
			c.HTML(http.StatusOK, "login.html", gin.H{"Title": "Login"})
		})
		app.POST("/api/login", ctrls.AuthController.Login)
		app.POST("/api/logout", ctrls.AuthController.Logout)

		protected := app.Group("")
		protected.Use(middleware.AuthMiddleware(userRepo))
		{
			setupPageRoutes(protected, svcs)
			setupAPIRoutes(protected, ctrls)
		}
	}
	return router
}

func setupDependencies(db *gorm.DB, cfg *config.Config) (Repositories, Services, Controllers) {
	userRepo := repositories.NewUserRepository(db)
	residentRepo := repositories.NewResidentRepository(db)
	docRepo := repositories.NewLostDocumentRepository(db)
	configRepo := repositories.NewConfigRepository(db)
	auditRepo := repositories.NewAuditLogRepository(db)

	services.JWTSecretKey = []byte(cfg.JWTSecretKey)

	configService := services.NewConfigService(configRepo)
	auditService := services.NewAuditLogService(auditRepo)
	authService := services.NewAuthService(userRepo)
	dashboardService := services.NewDashboardService(docRepo, userRepo, configService)
	docService := services.NewLostDocumentService(db, docRepo, residentRepo, userRepo, auditService, configService)
	userService := services.NewUserService(userRepo, auditService, cfg)
	backupService := services.NewBackupService(cfg, configService, auditService)

	authController := controllers.NewAuthController(authService)
	dashboardController := controllers.NewDashboardController(dashboardService)
	docController := controllers.NewLostDocumentController(docService)
	userController := controllers.NewUserController(userService)
	configController := controllers.NewConfigController(configService, userService)
	auditController := controllers.NewAuditLogController(auditService)
	backupController := controllers.NewBackupController(backupService)
	settingsController := controllers.NewSettingsController(configService, auditService)

	return Repositories{UserRepo: userRepo},
		Services{ConfigService: configService, DocService: docService},
		Controllers{
			AuthController:      authController,
			DashboardController: dashboardController,
			DocController:       docController,
			UserController:      userController,
			ConfigController:    configController,
			AuditController:     auditController,
			BackupController:    backupController,
			SettingsController:  settingsController,
		}
}

func setupPageRoutes(router *gin.RouterGroup, svcs Services) {
	getUser := func(c *gin.Context) interface{} {
		user, _ := c.Get("currentUser")
		return user
	}

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard.html", gin.H{"Title": "Dasbor", "CurrentUser": getUser(c)})
	})

	router.GET("/documents", func(c *gin.Context) {
		c.HTML(http.StatusOK, "document_list.html", gin.H{"Title": "Daftar Dokumen Aktif", "CurrentUser": getUser(c), "PageType": "active"})
	})

	router.GET("/documents/archived", func(c *gin.Context) {
		c.HTML(http.StatusOK, "document_list.html", gin.H{"Title": "Arsip Dokumen", "CurrentUser": getUser(c), "PageType": "archived"})
	})

	router.GET("/documents/new", func(c *gin.Context) {
		c.HTML(http.StatusOK, "document_form.html", gin.H{"Title": "Buat Surat Baru", "CurrentUser": getUser(c), "IsEdit": false, "DocID": 0})
	})

	router.GET("/documents/:id/edit", func(c *gin.Context) {
		id := c.Param("id")
		c.HTML(http.StatusOK, "document_form.html", gin.H{"Title": "Edit Surat", "CurrentUser": getUser(c), "IsEdit": true, "DocID": id})
	})

	router.GET("/search", func(c *gin.Context) {
		query := c.Query("q")
		c.HTML(http.StatusOK, "search_results.html", gin.H{"Title": "Hasil Pencarian", "CurrentUser": getUser(c), "Query": query})
	})

	router.GET("/profile", func(c *gin.Context) {
		c.HTML(http.StatusOK, "profile.html", gin.H{"Title": "Profil Pengguna", "CurrentUser": getUser(c)})
	})

	router.GET("/panduan", func(c *gin.Context) {
		c.HTML(http.StatusOK, "panduan.html", gin.H{"Title": "Panduan Pengguna", "CurrentUser": getUser(c)})
	})

	router.GET("/tentang", func(c *gin.Context) {
		c.HTML(http.StatusOK, "tentang.html", gin.H{"Title": "Tentang Aplikasi", "CurrentUser": getUser(c), "AppVersion": version})
	})

	router.GET("/documents/:id/print", func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.HTML(http.StatusBadRequest, "error.html", gin.H{"Title": "Error", "CurrentUser": getUser(c), "ErrorMessage": "ID Dokumen tidak valid."})
			return
		}

		doc, err := svcs.DocService.FindByID(uint(id), c.GetUint("userID"))
		if err != nil {
			status := http.StatusNotFound
			message := "Dokumen tidak ditemukan."
			if errors.Is(err, services.ErrAccessDenied) {
				status = http.StatusForbidden
				message = "Anda tidak memiliki izin untuk mengakses dokumen ini."
			}
			c.HTML(status, "error.html", gin.H{"Title": "Error", "CurrentUser": getUser(c), "ErrorMessage": message})
			return
		}

		appConfig, err := svcs.ConfigService.GetConfig()
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "CurrentUser": getUser(c), "ErrorMessage": "Gagal memuat konfigurasi aplikasi."})
			return
		}

		c.HTML(http.StatusOK, "print_preview.html", gin.H{"Document": doc, "Now": time.Now(), "CurrentUser": getUser(c), "Config": appConfig})
	})

	adminRoutes := router.Group("")
	adminRoutes.Use(middleware.AdminAuthMiddleware())
	{
		adminRoutes.GET("/users", func(c *gin.Context) {
			c.HTML(http.StatusOK, "user_list.html", gin.H{"Title": "Manajemen Pengguna", "CurrentUser": getUser(c)})
		})

		adminRoutes.GET("/users/new", func(c *gin.Context) {
			c.HTML(http.StatusOK, "user_form.html", gin.H{"Title": "Tambah Pengguna", "CurrentUser": getUser(c), "IsEdit": false, "UserID": 0})
		})

		adminRoutes.GET("/users/:id/edit", func(c *gin.Context) {
			id, _ := strconv.Atoi(c.Param("id"))
			c.HTML(http.StatusOK, "user_form.html", gin.H{"Title": "Edit Pengguna", "CurrentUser": getUser(c), "IsEdit": true, "UserID": id})
		})

		adminRoutes.GET("/audit-logs", func(c *gin.Context) {
			c.HTML(http.StatusOK, "audit_log_list.html", gin.H{"Title": "Log Audit Sistem", "CurrentUser": getUser(c)})
		})

		adminRoutes.GET("/settings", func(c *gin.Context) {
			c.HTML(http.StatusOK, "settings.html", gin.H{"Title": "Pengaturan Sistem", "CurrentUser": getUser(c)})
		})
	}
}

func setupAPIRoutes(router *gin.RouterGroup, ctrls Controllers) {
	api := router.Group("/api")
	{
		api.GET("/stats", ctrls.DashboardController.GetStats)
		api.GET("/stats/monthly-issuance", ctrls.DashboardController.GetMonthlyChart)
		api.GET("/stats/item-composition", ctrls.DashboardController.GetItemCompositionChart)
		api.GET("/notifications/expiring-documents", ctrls.DashboardController.GetExpiringDocuments)
		api.PUT("/profile", ctrls.UserController.UpdateProfile)
		api.PUT("/profile/password", ctrls.UserController.ChangePassword)
		api.GET("/search", ctrls.DocController.SearchGlobal)
		api.POST("/documents", ctrls.DocController.Create)
		api.GET("/documents", ctrls.DocController.FindAll)
		api.GET("/documents/:id", ctrls.DocController.FindByID)
		api.PUT("/documents/:id", ctrls.DocController.Update)
		api.DELETE("/documents/:id", ctrls.DocController.Delete)

		adminAPI := router.Group("/api")
		adminAPI.Use(middleware.AdminAuthMiddleware())
		{
			adminAPI.GET("/users/operators", ctrls.UserController.FindOperators)
			adminAPI.POST("/users", ctrls.UserController.Create)
			adminAPI.GET("/users", ctrls.UserController.FindAll)
			adminAPI.GET("/users/:id", ctrls.UserController.FindByID)
			adminAPI.PUT("/users/:id", ctrls.UserController.Update)
			adminAPI.DELETE("/users/:id", ctrls.UserController.Delete)
			adminAPI.POST("/users/:id/activate", ctrls.UserController.Activate)
			adminAPI.GET("/audit-logs", ctrls.AuditController.FindAll)
			adminAPI.POST("/backups", ctrls.BackupController.CreateBackup)
			adminAPI.POST("/restore", ctrls.BackupController.RestoreBackup)
			adminAPI.GET("/settings", ctrls.SettingsController.GetSettings)
			adminAPI.PUT("/settings", ctrls.SettingsController.UpdateSettings)
		}
	}
}

type Repositories struct {
	UserRepo repositories.UserRepository
}

type Services struct {
	ConfigService services.ConfigService
	DocService    services.LostDocumentService
}

type Controllers struct {
	AuthController      *controllers.AuthController
	DashboardController *controllers.DashboardController
	DocController       *controllers.LostDocumentController
	UserController      *controllers.UserController
	ConfigController    *controllers.ConfigController
	AuditController     *controllers.AuditLogController
	BackupController    *controllers.BackupController
	SettingsController  *controllers.SettingsController
}