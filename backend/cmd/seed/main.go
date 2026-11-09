package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/memekita/backend/internal/config"
	"github.com/memekita/backend/internal/db"
	"github.com/memekita/backend/internal/models"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()
	gormDB, err := db.NewGORM(ctx, cfg)
	if err != nil { log.Fatal(err) }

	sections := []struct{slug,name,desc string; regional bool}{
		{"funny","Funny","Meme lucu sehari-hari",false},
		{"gaming","Gaming","Meme game",false},
		{"anime","Anime","Anime dan manga",false},
		{"indonesia","Indonesia","Humor nusantara",false},
		{"sunda","Sunda","Humor Sunda",true},
		{"jawa","Jawa","Humor Jawa",true},
		{"betawi","Betawi","Humor Betawi",true},
		{"otaku","Otaku","Wibu bersatu",false},
	}
	for _, s := range sections {
		if err := gormDB.WithContext(ctx).FirstOrCreate(&models.Section{}, models.Section{Slug: s.slug, Name: s.name}).Error; err != nil { log.Printf("seed error: %v", err) }
		// ensure description etc.
		if err := gormDB.WithContext(ctx).Model(&models.Section{}).Where("slug = ?", s.slug).Updates(map[string]interface{}{"description": s.desc, "is_regional": s.regional}).Error; err != nil { log.Printf("seed error: %v", err) }
	}
	fmt.Println("sections seeded")

	for i:=1;i<=10;i++ {
		username := fmt.Sprintf("user%d", i)
		email := fmt.Sprintf("user%d@memekita.id", i)
		hash := fmt.Sprintf("argon2$%d", i)
		id := uuid.New()
		if err := gormDB.WithContext(ctx).Create(&models.User{ID: id, Username: username, Email: email, PasswordHash: &hash, Role: "user", Status: "active"}).Error; err != nil { log.Printf("seed error: %v", err) }
		if err := gormDB.WithContext(ctx).Exec("INSERT INTO user_coins (user_id, balance_coins) VALUES (?,1000) ON CONFLICT DO NOTHING", id).Error; err != nil { log.Printf("seed error: %v", err) }
	}
	if err := gormDB.WithContext(ctx).Where("email = ?", "admin@memekita.id").FirstOrCreate(&models.User{ID: uuid.New(), Username: "admin", Email: "admin@memekita.id", PasswordHash: stringPtr("adminhash"), Role: "admin", Status: "active"}).Error; err != nil { log.Printf("seed error: %v", err) }
	fmt.Println("users seeded")

	for _, p := range []struct{coins, price int; sku string}{{100,15000,"coins_100"},{500,49000,"coins_500"},{1000,89000,"coins_1000"}} {
		if err := gormDB.WithContext(ctx).FirstOrCreate(&models.CoinPackage{}, models.CoinPackage{PlaySku: p.sku, Coins: p.coins, PriceIdr: p.price, IsActive: true}).Error; err != nil { log.Printf("seed error: %v", err) }
	}
	for _, t := range []struct{name, url string}{{"Drake","https://cdn.memekita.id/templates/drake.jpg"},{"Distracted Boyfriend","https://cdn.memekita.id/templates/boyfriend.jpg"}} {
		if err := gormDB.WithContext(ctx).FirstOrCreate(&models.MemeTemplate{}, models.MemeTemplate{Name: t.name, ImageURL: t.url, IsOfficial: true, Status: "approved"}).Error; err != nil { log.Printf("seed error: %v", err) }
	}

	var sectionIDs []uuid.UUID
	var secs []models.Section
	if err := gormDB.WithContext(ctx).Find(&secs).Error; err != nil { log.Printf("seed error: %v", err) }
	for _, s := range secs { sectionIDs = append(sectionIDs, s.ID) }
	var userIDs []uuid.UUID
	var users []models.User
	if err := gormDB.WithContext(ctx).Where("role = ?", "user").Find(&users).Error; err != nil { log.Printf("seed error: %v", err) }
	for _, u := range users { userIDs = append(userIDs, u.ID) }
	titles := []string{"Ketika dosen bilang tugas besok","Muka gue pas ditagih utang","Ngaku anak senja tapi takut gelap","Warteg solusi hidup","Macet Jakarta","Emak bilang terserah","Tidur jam 2 bangun jam 6","Niat diet Senin gagal Selasa","Chat dibaca doang overthinking","WiFi tetangga lebih kencang"}
	rand.Seed(time.Now().UnixNano())
	for i:=1;i<=50;i++ {
		uid := userIDs[rand.Intn(len(userIDs))]
		sid := sectionIDs[rand.Intn(len(sectionIDs))]
		title := fmt.Sprintf("%s #%d", titles[rand.Intn(len(titles))], i)
		mediaType := []string{"image","gif","video"}[rand.Intn(3)]
		if err := gormDB.WithContext(ctx).Create(&models.Post{ID: uuid.New(), UserID: uid, SectionID: sid, Title: title, MediaType: mediaType, MediaURL: stringPtr(fmt.Sprintf("https://cdn.memekita.id/sample/%d.jpg", i)), ThumbURL: stringPtr(fmt.Sprintf("https://cdn.memekita.id/sample/thumb-%d.jpg", i)), Status: "published", HotScore: rand.Float64()*10, Upvotes: rand.Intn(100), Downvotes: rand.Intn(10)}).Error; err != nil { log.Printf("seed error: %v", err) }
	}
	fmt.Println("posts seeded")
}
func stringPtr(s string) *string { return &s }
