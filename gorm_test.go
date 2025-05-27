package belajargolanggorm

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

func OpenConnection() *gorm.DB {
	dialect := mysql.Open("root:@tcp(127.0.0.1:3306)/belajar_golang_gorm?charset=utf8mb4&parseTime=True&loc=Local")
	db, err := gorm.Open(dialect, &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Info),
		SkipDefaultTransaction: true, // untuk menghindari auto transaction
		PrepareStmt:            true, // untuk menggunakan prepared statement yang sudah pernah digunakan sehingga tidak perlu di prepare lagi
	})
	if err != nil {
		panic("failed to connect database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	return db
}

var db = OpenConnection()

func TestConnection(t *testing.T) {
	assert.NotNil(t, db)
}

func TestExecuteSql(t *testing.T) {
	err := db.Exec("INSERT INTO sample (name) VALUES (?)", "Salman").Error
	assert.Nil(t, err)

	err = db.Exec("INSERT INTO sample (name) VALUES (?)", "Seif").Error
	assert.Nil(t, err)

	err = db.Exec("INSERT INTO sample (name) VALUES (?)", "Man").Error
	assert.Nil(t, err)
}

type Sample struct {
	Id   int
	Name string
}

func TestRawQuery(t *testing.T) {
	var sample Sample
	err := db.Raw("SELECT * FROM sample WHERE id = ?", 1).Scan(&sample).Error
	assert.Nil(t, err)
	assert.Equal(t, "Salman", sample.Name)

	var samples []Sample
	err = db.Raw("SELECT * FROM sample").Scan(&samples).Error
	assert.Nil(t, err)
	assert.Equal(t, 3, len(samples))
}

func TestSqlRows(t *testing.T) {
	rows, err := db.Raw("SELECT * FROM sample").Rows()
	assert.Nil(t, err)
	defer rows.Close()

	var samples []Sample
	for rows.Next() {
		var id int
		var name string

		err := rows.Scan(&id, &name)
		assert.Nil(t, err)

		samples = append(samples, Sample{
			Id:   id,
			Name: name,
		})
	}

	assert.Equal(t, 3, len(samples))
}

func TestScanRows(t *testing.T) {
	rows, err := db.Raw("SELECT * FROM sample").Rows()
	assert.Nil(t, err)
	defer rows.Close()

	var samples []Sample
	for rows.Next() {
		err := db.ScanRows(rows, &samples)
		assert.Nil(t, err)
	}

	assert.Equal(t, 3, len(samples))
}

func TestCreateUser(t *testing.T) {
	user := User{
		ID:       1,
		Password: "",
		Name: Name{
			FirstName:  "Salman",
			LastName:   "Seif",
			MiddleName: "Man",
		},
		Information: "Belajar Golang Gorm",
	}

	response := db.Create(&user)
	assert.Nil(t, response.Error)
	assert.Equal(t, int64(1), response.RowsAffected)
}

func TestBatchInsert(t *testing.T) {
	var users []User

	for i := 2; i <= 10; i++ {
		users = append(users, User{
			ID:       i,
			Password: "123456",
			Name: Name{
				FirstName:  "User",
				LastName:   "Ke-" + strconv.Itoa(i),
				MiddleName: "Batch",
			},
			Information: "Belajar Golang Gorm",
		})
	}

	result := db.Create(&users)
	assert.Nil(t, result.Error)
	assert.Equal(t, int64(9), result.RowsAffected)
}

func TestTransactionSuccess(t *testing.T) {
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&User{
			ID:       11,
			Password: "123456",
			Name: Name{
				FirstName: "User 11",
			},
		}).Error
		if err != nil {
			return err
		}

		err = tx.Create(&User{
			ID:       12,
			Password: "123456",
			Name: Name{
				FirstName: "User 12",
			},
		}).Error
		if err != nil {
			return err
		}

		err = tx.Create(&User{
			ID:       13,
			Password: "123456",
			Name: Name{
				FirstName: "User 13",
			},
		}).Error
		if err != nil {
			return err
		}

		return nil
	})

	assert.Nil(t, err)
}

func TestTransactionRollback(t *testing.T) {
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&User{
			ID:       14,
			Password: "123456",
			Name: Name{
				FirstName: "User 14",
			},
		}).Error
		if err != nil {
			return err
		}

		err = tx.Create(&User{
			ID:       10,
			Password: "123456",
			Name: Name{
				FirstName: "User 10",
			},
		}).Error
		if err != nil {
			return err
		}

		return nil
	})

	assert.NotNil(t, err)
}

func TestManualTransactionSuccess(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	err := tx.Create(&User{ID: 15, Password: "123456", Name: Name{FirstName: "User 15"}}).Error
	assert.Nil(t, err)

	err = tx.Create(&User{ID: 16, Password: "123456", Name: Name{FirstName: "User 16"}}).Error
	assert.Nil(t, err)

	if err == nil {
		tx.Commit()
	}
}

func TestQuerySingleObject(t *testing.T) {
	user := User{}
	err := db.First(&user, 1).Error
	assert.Nil(t, err)
	assert.Equal(t, 1, user.ID)

	user = User{}
	err = db.Last(&user).Error
	assert.Nil(t, err)
	assert.Equal(t, 16, user.ID)
}

func TestQuerySingleObjectInlineCondition(t *testing.T) {
	user := User{}
	// kalau pake first maka datanya akan diurutkan
	// err := db.First(&user, "id = ?", 1).Error
	// kalau pake take maka datanya tidak diurutkan
	err := db.Take(&user, "id = ?", 1).Error
	assert.Nil(t, err)
	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "Salman", user.Name.FirstName)
}

func TestQueryAllObjects(t *testing.T) {
	var users []User
	err := db.Find(&users, "id in ?", []int{5, 6, 8}).Error
	assert.Nil(t, err)
	assert.Equal(t, 3, len(users))
}

func TestQueryCondition(t *testing.T) {
	var users []User
	err := db.Where("first_name like ?", "%User%").Where("password = ?", "123456").Find(&users).Error

	assert.Nil(t, err)
	assert.Equal(t, 15, len(users))
}

func TestOrCondition(t *testing.T) {
	var users []User
	err := db.Where("first_name like ?", "%User%").Or("password = ?", "123456").Find(&users).Error

	assert.Nil(t, err)
	assert.Equal(t, 16, len(users))
}

func TestNotOperator(t *testing.T) {
	var users []User
	err := db.Not("first_name like ?", "%User%").Where("password = ?", "123456").Find(&users).Error

	assert.Nil(t, err)
	assert.Equal(t, 1, len(users))
}

func TestSelectFields(t *testing.T) {
	var users []User
	err := db.Select("id", "first_name").Find(&users).Error

	assert.Nil(t, err)
	for _, user := range users {
		assert.NotNil(t, user.ID)
		assert.NotEqual(t, "", user.Name.FirstName)
	}

	assert.Equal(t, 16, len(users))
}

func TestStructCondition(t *testing.T) {
	userCondition := User{
		Name: Name{
			FirstName: "User 11",
			LastName:  "", // tidak bisa karena dianggap default value
		},
	}

	var users []User
	err := db.Where(&userCondition).Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 1, len(users))
}

func TestMapCondition(t *testing.T) {
	mapCondition := map[string]interface{}{
		"middle_name": "",
	}

	var users []User
	err := db.Where(mapCondition).Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 6, len(users))
}

func TestOrderLimitOffset(t *testing.T) {
	var users []User
	err := db.Order("id asc, first_name desc").Limit(5).Offset(5).Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 5, len(users))
}

type UserResponse struct {
	ID        int
	FirstName string
	LastName  string
}

func TestQueryNonModel(t *testing.T) {
	var users []UserResponse
	err := db.Model(&User{}).Select("id", "first_name", "last_name").Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 16, len(users))
}

func TestUpdate(t *testing.T) {
	user := User{}
	err := db.Take(&user, "id=?", 2).Error
	assert.Nil(t, err)

	user.Name.FirstName = "Adi"
	user.Name.MiddleName = ""
	user.Name.LastName = "Wijaya"
	user.Password = "rahasia"

	err = db.Save(&user).Error
	assert.Nil(t, err)
}

func TestUpdateSelectedColumn(t *testing.T) {
	user := User{}
	err := db.Model(&user).Where("id=?", 5).Updates(map[string]interface{}{
		"middle_name": "ada",
		"last_name":   "wong",
	}).Error
	assert.Nil(t, err)

	err = db.Model(&User{}).Where("id=?", 8).Update("first_name", "Ujang").Error
	assert.Nil(t, err)

	err = db.Where("id=?", 10).Updates(User{
		Name: Name{
			FirstName: "Steve",
		},
	}).Error
	assert.Nil(t, err)
}

func TestAutoIncrement(t *testing.T) {
	for i := 0; i < 10; i++ {
		userLog := UserLog{
			UserId: 1,
			Action: "Nothing",
		}

		err := db.Create(&userLog).Error
		assert.Nil(t, err)

		assert.NotEqual(t, 0, userLog.ID)
		fmt.Println(userLog.ID)
	}
}

// save bisa digunakan untuk insert dan update tapi lebih baik digunakan untuk auto increment
func TestSaveOrUpdate(t *testing.T) {
	userLog := UserLog{
		UserId: 1,
		Action: "Action baru",
	}
	err := db.Save(&userLog).Error // insert
	assert.Nil(t, err)

	userLog.UserId = 2
	err = db.Save(&userLog).Error // update
	assert.Nil(t, err)
}

func TestSaveOrUpdateNonAutoIncrement(t *testing.T) {
	user := User{
		ID: 99,
		Name: Name{
			FirstName: "Test",
		},
	}
	err := db.Save(&user).Error // insert
	assert.Nil(t, err)

	user.Name.FirstName = "Test 2 Updated"
	err = db.Save(&user).Error // update
	assert.Nil(t, err)
}

func TestConflict(t *testing.T) {
	user := User{
		ID: 88,
		Name: Name{
			FirstName: "Test",
		},
	}
	err := db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&user).Error
	assert.Nil(t, err)
}

func TestDelete(t *testing.T) {
	var user User
	err := db.Take(&user, "id = ?", 1).Error
	assert.Nil(t, err)

	err = db.Delete(&user).Error
	assert.Nil(t, err)

	err = db.Delete(&User{}, "id = ?", 2).Error
	assert.Nil(t, err)

	err = db.Where("id = ?", 4).Delete(&User{}).Error
	assert.Nil(t, err)

}

func TestSoftDeletes(t *testing.T) {
	todo := Todo{
		UserId:      1,
		Title:       "Belajar Golang Gorm",
		Description: "Belajar Golang Gorm",
	}
	err := db.Create(&todo).Error
	assert.Nil(t, err)

	err = db.Delete(&todo).Error
	assert.Nil(t, err)
	assert.NotNil(t, todo.DeletedAt)

	var todos []Todo
	err = db.Find(&todos).Error
	assert.Nil(t, err)
	assert.Equal(t, 0, len(todos))
}

func TestUnscoped(t *testing.T) {
	var todo Todo
	err := db.Unscoped().First(&todo, "id = ?", 2).Error
	assert.Nil(t, err)
	fmt.Println(todo)

	err = db.Unscoped().Delete(&todo).Error
	assert.Nil(t, err)

	var todos []Todo
	err = db.Unscoped().Find(&todos).Error
	assert.Nil(t, err)
}

func TestLock(t *testing.T) {
	err := db.Transaction(func(tx *gorm.DB) error {
		var user User
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Take(&user, "id = ?", 3).Error
		if err != nil {
			return err
		}

		user.Name.FirstName = "Salman Update"
		err = tx.Save(&user).Error

		return err
	})
	assert.Nil(t, err)
}

func TestCreateWallet(t *testing.T) {
	wallet := Wallet{
		ID:      "1",
		UserId:  3,
		Balance: 1000000,
	}

	err := db.Create(&wallet).Error
	assert.Nil(t, err)
}

// untuk preload kurang disarankan untuk one to one karena memanggil query lebih dari satu
func TestRetrieveRelation(t *testing.T) {
	var user User
	err := db.Model(&user).Preload("Wallet").Take(&user, "id=?", 3).Error
	assert.Nil(t, err)

	assert.Equal(t, 3, user.ID)
}

// untuk relasi one to one lebih disarankan untuk menggunakan join
func TestRetrieveRelationJoin(t *testing.T) {
	var user User
	err := db.Model(&user).Joins("Wallet").Take(&user, "users.id=?", 3).Error
	assert.Nil(t, err)

	assert.Equal(t, 3, user.ID)
}

func TestAutoCreateUpdate(t *testing.T) {
	user := User{
		ID:       20,
		Password: "Rahasia",
		Name: Name{
			FirstName: "Salman 2",
		},
		Wallet: Wallet{
			ID:      "20",
			UserId:  20,
			Balance: 1000000,
		},
	}

	err := db.Create(&user).Error
	assert.Nil(t, err)
}

func TestSkipAutoCreateUpdate(t *testing.T) {
	user := User{
		ID:       21,
		Password: "Rahasia",
		Name: Name{
			FirstName: "Salman 3",
		},
		Wallet: Wallet{
			ID:      "21",
			UserId:  21,
			Balance: 1000000,
		},
	}

	err := db.Omit(clause.Associations).Create(&user).Error
	assert.Nil(t, err)
}

func TestUserAndAddresses(t *testing.T) {
	user := User{
		ID:       51,
		Password: "Rahasia",
		Name: Name{
			FirstName: "Salman 51 test",
		},
		Wallet: Wallet{
			ID:      "51",
			UserId:  51,
			Balance: 1000000,
		},
		Addresses: []Address{
			{
				UserId:  51,
				Address: "Jl. Raya No 1",
			},
			{
				UserId:  51,
				Address: "Jl. Raya No 2",
			},
		},
	}

	err := db.Create(&user).Error
	assert.Nil(t, err)
}

func TestPreloadJoinOneToMany(t *testing.T) {
	var userPreload []User
	err := db.Model(&userPreload).Preload("Addresses").Joins("Wallet").Take(&userPreload, "users.id=?", 50).Error
	assert.Nil(t, err)
}

func TestBelongsTo(t *testing.T) {
	fmt.Println("Preload")
	var addresses []Address
	err := db.Model(&addresses).Preload("User").Find(&addresses).Error
	assert.Nil(t, err)

	fmt.Println("Join")
	addresses = []Address{}
	err = db.Model(&Address{}).Joins("User").Find(&addresses).Error
	assert.Nil(t, err)
}

func TestBelongsToWallet(t *testing.T) {
	fmt.Println("Preload")
	var wallets []Wallet
	err := db.Model(&wallets).Preload("User").Find(&wallets).Error
	assert.Nil(t, err)

	fmt.Println("Join")
	wallets = []Wallet{}
	err = db.Model(&Wallet{}).Joins("User").Find(&wallets).Error
	assert.Nil(t, err)
}

func TestCreateManyToMany(t *testing.T) {
	product := Product{
		ID:    1,
		Name:  "Product 1",
		Price: 10000,
	}

	err := db.Create(&product).Error
	assert.Nil(t, err)

	err = db.Table("user_like_product").Create(&map[string]interface{}{
		"user_id":    3,
		"product_id": 1,
	}).Error
	assert.Nil(t, err)

	err = db.Table("user_like_product").Create(&map[string]interface{}{
		"user_id":    5,
		"product_id": 1,
	}).Error
	assert.Nil(t, err)
}

func TestPreloadManyToMany(t *testing.T) {
	var product Product
	err := db.Preload("LikedByUsers").Take(&product, "id=?", 1).Error
	assert.Nil(t, err)
	assert.Equal(t, 2, len(product.LikedByUsers))
}

func TestPreloadManyToManyUser(t *testing.T) {
	var user User
	err := db.Preload("LikeProducts").Take(&user, "id=?", 5).Error
	assert.Nil(t, err)
	assert.Equal(t, 1, len(user.LikeProducts))
}

func TestAssociationFind(t *testing.T) {
	var product Product
	err := db.Take(&product, "id=?", 1).Error
	assert.Nil(t, err)

	var users []User
	err = db.Model(&product).Where("first_name LIKE ?", "User%").Association("LikedByUsers").Find(&users)
	assert.Nil(t, err)
}

func TestAssociationAppend(t *testing.T) {
	var user User
	err := db.Take(&user, "id=?", 5).Error
	assert.Nil(t, err)

	var product Product
	err = db.Take(&product, "id=?", 1).Error
	assert.Nil(t, err)

	err = db.Model(&product).Association("LikedByUsers").Append(&user)
	assert.Nil(t, err)
}

func TestAssociationReplace(t *testing.T) {
	err := db.Transaction(func(tx *gorm.DB) error {
		var user User
		err := tx.Take(&user, "id=?", 5).Error
		assert.Nil(t, err)

		wallet := Wallet{
			ID:      "1",
			UserId:  user.ID,
			Balance: 1000000,
		}

		err = tx.Model(&user).Association("Wallet").Replace(&wallet)
		return err
	})

	assert.Nil(t, err)
}

// delete hanya menghapus relasi di tabel pivot
func TestAssociationDelete(t *testing.T) {
	var user User
	err := db.Take(&user, "id=?", 5).Error
	assert.Nil(t, err)

	var product Product
	err = db.Take(&product, "id=?", 1).Error
	assert.Nil(t, err)

	err = db.Model(&product).Association("LikedByUsers").Delete(&user)
	assert.Nil(t, err)
}

// menghapus relasi dan table utama
func TestAssociationClear(t *testing.T) {
	var product Product
	err := db.Take(&product, "id=?", 1).Error
	assert.Nil(t, err)

	err = db.Model(&product).Association("LikedByUsers").Clear()
	assert.Nil(t, err)
}

func TestPreloadingWithCondition(t *testing.T) {
	var user User
	err := db.Preload("Wallet", "balance>?", 1000000).Take(&user, "id=?", 3).Error
	assert.Nil(t, err)
}

func TestNestedPreloading(t *testing.T) {
	var wallet Wallet
	err := db.Preload("User.Addresses").Take(&wallet, "id=?", 50).Error
	assert.Nil(t, err)
	fmt.Println("Wallet:", wallet)
}

func TestPreloadingAll(t *testing.T) {
	var user User
	err := db.Preload(clause.Associations).Take(&user, "id=?", 3).Error
	assert.Nil(t, err)
}

func TestJoinQuery(t *testing.T) {
	var users []User
	err := db.Joins("Join wallets w on w.user_id=users.id").Find(&users).Error
	assert.Nil(t, err)

	users = []User{}
	err = db.Joins("Wallet").Find(&users).Error // left join
	assert.Nil(t, err)
}

func TestJoinWithCondition(t *testing.T) {
	var users []User
	err := db.Joins("Join wallets w on w.user_id=users.id And w.balance > 500000").Find(&users).Error
	assert.Nil(t, err)

	users = []User{}
	err = db.Joins("Wallet").Where("Wallet.balance > ?", 500000).Find(&users).Error // left join
	assert.Nil(t, err)
}

func TestCount(t *testing.T) {
	var count int64
	err := db.Model(&User{}).Joins("Wallet").Where("Wallet.balance > ?", 500000).Count(&count).Error
	assert.Nil(t, err)
	fmt.Println("Count:", count)
}

type AggregationResult struct {
	TotalBalance int64
	MinBalance   int64
	MaxBalance   int64
	AvgBalance   float64
}

func TestAggregation(t *testing.T) {
	var result AggregationResult
	err := db.Model(&Wallet{}).Select("sum(balance) as total_balance, min(balance)as min_balance, max(balance) as max_balance, avg(balance) as avg_balance").Take(&result).Error
	assert.Nil(t, err)
	fmt.Println("Total Balance:", result.TotalBalance)
	fmt.Println("Min Balance:", result.MinBalance)
	fmt.Println("Max Balance:", result.MaxBalance)
	fmt.Println("Avg Balance:", result.AvgBalance)
}

func TestGroupByHaving(t *testing.T) {
	var result []AggregationResult
	err := db.Model(&Wallet{}).Select("sum(balance) as total_balance, min(balance)as min_balance, max(balance) as max_balance, avg(balance) as avg_balance").Joins("User").Group("User.id").Having("sum(balance) > ?", 1000).Take(&result).Error
	assert.Nil(t, err)
}

func TestContext(t *testing.T) {
	ctx := context.Background()
	var users []User
	err := db.WithContext(ctx).Find(&users).Error
	assert.Nil(t, err)
}

func BrokeWalletBalance(db *gorm.DB) *gorm.DB {
	return db.Where("balance = ?", 0)
}

func SultanWalletBalance(db *gorm.DB) *gorm.DB {
	return db.Where("balance > ?", 1000000)
}

func TestScopes(t *testing.T) {
	var wallets []Wallet
	err := db.Scopes(BrokeWalletBalance).Find(&wallets).Error
	assert.Nil(t, err)

	wallets = []Wallet{}
	err = db.Scopes(SultanWalletBalance).Find(&wallets).Error
	assert.Nil(t, err)
}

func TestMigrator(t *testing.T) {
	err := db.Migrator().AutoMigrate(&GuestBook{})
	assert.Nil(t, err)
}

func TestHookBeforeCreate(t *testing.T) {
	user := User{
		Password: "Rhs123",
		Name: Name{
			LastName: "Rand1",
		},
	}

	err := db.Create(&user).Error
	assert.Nil(t, err)

	fmt.Println("ID:", user.ID)
}
