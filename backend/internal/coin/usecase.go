package coin

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/memekita/backend/internal/models"
	"gorm.io/gorm"
)

type Usecase struct{ db *gorm.DB }

func NewUsecase(db *gorm.DB) *Usecase { return &Usecase{db: db} }

func (uc *Usecase) GetBalance(ctx context.Context, userID uuid.UUID) int {
	var ucModel models.UserCoin
	if err := uc.db.WithContext(ctx).Where("user_id = ?", userID).First(&ucModel).Error; err != nil {
		return 0
	}
	return ucModel.BalanceCoins
}

func (uc *Usecase) ListPackages(ctx context.Context) ([]models.CoinPackage, error) {
	var pkgs []models.CoinPackage
	err := uc.db.WithContext(ctx).Where("is_active = ?", true).Order("coins ASC").Find(&pkgs).Error
	return pkgs, err
}

func (uc *Usecase) Topup(ctx context.Context, userID, packageID uuid.UUID, purchaseToken string) (int, string, error) {
	if purchaseToken == "" {
		return 0, "", fmt.Errorf("token tidak valid")
	}
	var newBal int
	var txID string
	err := uc.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var pkg models.CoinPackage
		if err := tx.WithContext(ctx).Where("id = ?", packageID).First(&pkg).Error; err != nil {
			return err
		}
		var cur models.UserCoin
		if err := tx.WithContext(ctx).Where("user_id = ?", userID).First(&cur).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}
		newBal = cur.BalanceCoins + pkg.Coins
		if err := tx.WithContext(ctx).Exec("INSERT INTO user_coins (user_id, balance_coins) VALUES (?,?) ON CONFLICT (user_id) DO UPDATE SET balance_coins = ?", userID, newBal, newBal).Error; err != nil {
			return err
		}
		txID = uuid.NewString()
		if err := tx.WithContext(ctx).Create(&models.CoinTransaction{ID: uuid.MustParse(txID), UserID: userID, Type: "topup", AmountCoins: pkg.Coins, BalanceAfter: newBal, RefID: &packageID}).Error; err != nil {
			return err
		}
		return nil
	})
	return newBal, txID, err
}

func (uc *Usecase) SendTip(ctx context.Context, senderID, receiverID, postID uuid.UUID, amount int) (int, error) {
	if amount <= 0 {
		return 0, fmt.Errorf("jumlah tidak valid")
	}
	var newSender int
	err := uc.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var senderBal models.UserCoin
		if err := tx.WithContext(ctx).Where("user_id = ?", senderID).First(&senderBal).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}
		if senderBal.BalanceCoins < amount {
			return fmt.Errorf("Koin tidak cukup")
		}
		var receiverBal models.UserCoin
		if err := tx.WithContext(ctx).Where("user_id = ?", receiverID).First(&receiverBal).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}
		newSender = senderBal.BalanceCoins - amount
		newReceiver := receiverBal.BalanceCoins + amount
		if err := tx.WithContext(ctx).Exec("INSERT INTO user_coins (user_id, balance_coins) VALUES (?,?) ON CONFLICT (user_id) DO UPDATE SET balance_coins = ?", senderID, newSender, newSender).Error; err != nil {
			return err
		}
		if err := tx.WithContext(ctx).Exec("INSERT INTO user_coins (user_id, balance_coins) VALUES (?,?) ON CONFLICT (user_id) DO UPDATE SET balance_coins = ?", receiverID, newReceiver, newReceiver).Error; err != nil {
			return err
		}
		tipID := uuid.New()
		if err := tx.WithContext(ctx).Create(&models.Tip{ID: tipID, SenderID: senderID, ReceiverID: receiverID, PostID: &postID, AmountCoins: amount, Status: "completed"}).Error; err != nil {
			return err
		}
		if err := tx.WithContext(ctx).Create(&models.CoinTransaction{ID: uuid.New(), UserID: senderID, Type: "tip_sent", AmountCoins: -amount, BalanceAfter: newSender, RefID: &postID}).Error; err != nil {
			return err
		}
		if err := tx.WithContext(ctx).Create(&models.CoinTransaction{ID: uuid.New(), UserID: receiverID, Type: "tip_received", AmountCoins: amount, BalanceAfter: newReceiver, RefID: &postID}).Error; err != nil {
			return err
		}
		return nil
	})
	return newSender, err
}
