package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Cakra17/JobTracker-Api/internal/config"
	"github.com/Cakra17/JobTracker-Api/internal/job"
	"github.com/Cakra17/JobTracker-Api/internal/ratelimiter"
	"github.com/Cakra17/JobTracker-Api/internal/user"
	"github.com/Cakra17/JobTracker-Api/pkg/jwt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	cfg := config.InitializeConfig()

	app := fiber.New(fiber.Config{
		Prefork: false,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))
	// jwt provider
	jwtProvider := jwt.NewJWTProvider(cfg.JWTSecret)

	// rate limiter
	rateLimiter := ratelimiter.NewRateLimiter(10, 2)

	// db connection
	db := connectToDB(cfg.Database)

	// repository
	userRepo := user.NewUserRepo(db)
	jobRepo := job.NewJobRepo(db)

	// handler
	userHandler := user.NewUserHandler(user.UserHandlerConfig{
		UserRepo: &userRepo,
		JWTProvider: &jwtProvider,
		SaltCost: cfg.BcryptSalt,
		RateLimiter: &rateLimiter,
	})

	jobHandler := job.NewJobHandler(job.JobHandlerConfig{
		JobRepo: &jobRepo,
		JWTProvider: &jwtProvider,
		RateLimiter: &rateLimiter,
	})

	// route registration
	userHandler.RegisterRoute(app)
	jobHandler.RegisterRoute(app)
	
	addr := fmt.Sprintf(":%s", cfg.AppPort)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	
	wg.Add(1)
	go rateLimiterCleanup(ctx, &wg, &rateLimiter, time.Minute * 10, time.Minute * 30)

	// gracefully shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := app.Listen(addr); err != nil {
			log.Fatalf("Failed to running server at port %s: %s", addr ,err)
		}
	}()

	receiveSignal := <-sig

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second * 15)
	defer shutdownCancel()

	log.Printf("received %v, app stopped!!...", receiveSignal)
	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Fatalf("Server is forced to shutdown: %v", err)
	} else {
		log.Println("Server is gracefully shutdown")
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
		close(sig)
	}()

	select {
	case <-done:
		log.Println("All background process stopped")
	case <-time.After(10 * time.Second):
		log.Println("Timeout waiting for background processes to stop")
	}

	log.Println("App shutdown complete")
}

func connectToDB(dbCfg config.DatabaseConfig) *sql.DB {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		dbCfg.Username, dbCfg.Password,
		dbCfg.Host, dbCfg.Port, dbCfg.Name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	db.SetMaxOpenConns(dbCfg.MaxOpenConnection)
	db.SetMaxIdleConns(dbCfg.MaxIdleConnection)
	db.SetConnMaxLifetime(time.Duration(dbCfg.MaxConnLifetime) * time.Minute)
	db.SetConnMaxIdleTime(time.Duration(dbCfg.MaxConnIdleTime) * time.Minute)

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	log.Println("Connected to MySQL successfully!")
	return db
}

func rateLimiterCleanup(ctx context.Context, wg *sync.WaitGroup, rl *ratelimiter.RateLimiter, interval, maxAge time.Duration) {
	defer wg.Done()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Printf("Start rate limiter cleanup with Interval: %v, MaxAge: %v", interval, maxAge)
	
	for {
		select {
		case <-ctx.Done():
			log.Println("Rate limiter cleanup stopping...")
			return
		case <-ticker.C:
			rl.Cleanup(maxAge)
			fmt.Println("Cleanup completed")
		}
	}
}