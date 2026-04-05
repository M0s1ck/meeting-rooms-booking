package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/internships-backend/test-backend-M0s1ck/internal/config"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/booking"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/schedule"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/slot"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
	"github.com/internships-backend/test-backend-M0s1ck/internal/infra/postgres"
	"github.com/internships-backend/test-backend-M0s1ck/internal/infra/postgres/repository"
)

var (
	adminUserID  = uuid.MustParse("48323e70-a002-447c-86c5-d640142aaef0")
	normalUserID = uuid.MustParse("230e467a-31b8-4a4a-88aa-d4757876b950")

	blueRoomID  = uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1")
	greenRoomID = uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa2")
)

func main() {
	ctx := context.Background()

	cfg := config.Load()
	fallBackCfg(cfg)

	db, err := postgres.Connect(ctx, cfg.Psg)
	if err != nil {
		log.Fatalf("connect postgres: %v", err)
	}
	defer db.Close()

	txGetter := trmpgx.DefaultCtxGetter

	scheduleRepo := repository.NewScheduleRepo(db, txGetter)
	slotRepo := repository.NewSlotRepo(db, txGetter)
	bookingRepo := repository.NewBookingRepo(db, txGetter)

	if err := seedUsers(ctx, db); err != nil {
		log.Fatalf("seed users: %v", err)
	}

	if err := seedRooms(ctx, db); err != nil {
		log.Fatalf("seed rooms: %v", err)
	}

	schedules, err := seedSchedules(ctx, scheduleRepo)
	if err != nil {
		log.Fatalf("seed schedules: %v", err)
	}

	if err := seedSlots(ctx, slotRepo, schedules); err != nil {
		log.Fatalf("seed slots: %v", err)
	}

	if err := seedBooking(ctx, db, bookingRepo); err != nil {
		log.Fatalf("seed booking: %v", err)
	}

	log.Println("seed completed")
}

func seedUsers(ctx context.Context, db *pgxpool.Pool) error {
	const query = `
		INSERT INTO users (id, email, role, password_hash, created_at, updated_at)
		VALUES
			($1, $2, $3, NULL, NOW(), NOW()),
			($4, $5, $6, NULL, NOW(), NOW())
		ON CONFLICT (id) DO NOTHING
	`

	_, err := db.Exec(
		ctx,
		query,
		adminUserID, "admin@example.com", user.RoleAdmin,
		normalUserID, "user@example.com", user.RoleUser,
	)
	if err != nil {
		return fmt.Errorf("insert users: %w", err)
	}

	return nil
}

func seedRooms(ctx context.Context, db *pgxpool.Pool) error {
	const query = `
		INSERT INTO rooms (id, name, description, capacity, created_at, updated_at)
		VALUES
			($1, $2, $3, $4, NOW(), NOW()),
			($5, $6, $7, $8, NOW(), NOW())
		ON CONFLICT (id) DO NOTHING
	`

	_, err := db.Exec(
		ctx,
		query,
		blueRoomID, "Blue Room", "Small meeting room", 4,
		greenRoomID, "Green Room", "Large meeting room", 8,
	)
	if err != nil {
		return fmt.Errorf("insert rooms: %w", err)
	}

	return nil
}

func seedSchedules(ctx context.Context, scheduleRepo interface {
	Create(context.Context, *schedule.Schedule) error
}) ([]schedule.Schedule, error) {
	now := time.Now().UTC()

	start0900, err := schedule.ParseTimeOfDay("09:00")
	if err != nil {
		return nil, err
	}
	end1200, err := schedule.ParseTimeOfDay("12:00")
	if err != nil {
		return nil, err
	}
	start1300, err := schedule.ParseTimeOfDay("13:00")
	if err != nil {
		return nil, err
	}
	end1800, err := schedule.ParseTimeOfDay("18:00")
	if err != nil {
		return nil, err
	}

	blueSchedule, err := schedule.New(
		blueRoomID,
		[]int{1, 2, 3, 4, 5},
		start0900,
		end1200,
		now,
	)
	if err != nil {
		return nil, fmt.Errorf("build blue schedule: %w", err)
	}

	greenSchedule, err := schedule.New(
		greenRoomID,
		[]int{1, 2, 3, 4, 5},
		start1300,
		end1800,
		now,
	)
	if err != nil {
		return nil, fmt.Errorf("build green schedule: %w", err)
	}

	for _, sch := range []*schedule.Schedule{blueSchedule, greenSchedule} {
		err := scheduleRepo.Create(ctx, sch)
		if err != nil && err != schedule.ErrAlreadyExists {
			return nil, fmt.Errorf("create schedule for room %s: %w", sch.RoomID, err)
		}
	}

	return []schedule.Schedule{*blueSchedule, *greenSchedule}, nil
}

func seedSlots(ctx context.Context, slotRepo interface {
	Add(context.Context, []slot.Slot) error
}, schedules []schedule.Schedule) error {
	gen := slot.NewGenerator()

	now := time.Now().UTC()
	from := startOfUTCDay(now.AddDate(0, 0, 1))
	to := from.AddDate(0, 0, 3)

	allSlots := make([]slot.Slot, 0)

	for _, sch := range schedules {
		allSlots = append(allSlots, gen.Generate(sch, from, to)...)
	}

	if err := slotRepo.Add(ctx, allSlots); err != nil {
		return fmt.Errorf("add slots: %w", err)
	}

	return nil
}

func seedBooking(ctx context.Context, db *pgxpool.Pool, bookingRepo *repository.BookingRepo) error {
	const query = `
		SELECT id
		FROM slots
		WHERE room_id = $1
		  AND start_at > $2
		ORDER BY start_at
		LIMIT 1
	`

	var slotID uuid.UUID
	err := db.QueryRow(ctx, query, blueRoomID, time.Now().UTC()).Scan(&slotID)
	if err != nil {
		return fmt.Errorf("select future slot: %w", err)
	}

	b, err := booking.NewActive(slotID, normalUserID, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("build booking: %w", err)
	}

	err = bookingRepo.Create(ctx, b)
	if err != nil && !errors.Is(err, booking.ErrAlreadyBooked) {
		return fmt.Errorf("create booking: %w", err)
	}

	return nil
}

func startOfUTCDay(t time.Time) time.Time {
	t = t.UTC()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func fallBackCfg(cfg *config.Config) {
	if cfg.Psg.Host == "" {
		cfg.Psg.Host = "localhost"
	}

	if cfg.Psg.User == "" {
		cfg.Psg.User = "psg_user"
	}

	if cfg.Psg.Password == "" {
		cfg.Psg.Password = "psg_pass"
	}

	if cfg.Psg.Name == "" {
		cfg.Psg.Name = "scheduler_db"
	}
}
