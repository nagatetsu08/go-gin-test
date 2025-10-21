package repositories

import (
	"errors"
	"gin-freemarket/models"

	"gorm.io/gorm"
)

// アイテムのリポジトリが持つべき基本機能（インターフェース）の実装
// メソッド名() (戻り値)
// メソッドの引数は基本的に値渡し。参照を渡すのはDBぐらい
type IItemRepository interface {

	// FindAllというメソッド名で、戻り値がmodels.Item型スライスへのポインタとerrorを返す（errorがない場合はnil）
	// 戻り値は参照を返すのが一般的（無駄なコピーを避けて省メモリ・省コストにしたい）
	FindAll() (*[]models.Item, error)

	// id検索は1件のみ返ってくるので、戻り値は*models.Itemとなる（FindAllは複数件返ってくる想定だから配列）
	FindById(itemId uint) (*models.Item, error)

	Create(newItem models.Item) (*models.Item, error)
	Update(newItem models.Item) (*models.Item, error)
	Delete(itemId uint) error
}

// アイテム情報をメモリ上に保存・取り扱うための「リポジトリ（倉庫）」となる構造体の定義
// データベースやファイルを直接使わず、一時的にメモリだけでデータ（アイテム一覧）を管理したい時に使う。（初っ端はDB使わないから）
// itemsというフィールドに全Itemを保持する。
type ItemMemoryRopository struct {
	items []models.Item
}

// ItemMemoryRopositoryのコンストラクタ
func NewItemMemoryRepository(items []models.Item) IItemRepository {
	// 作成した構造体のポインタを返す
	// &構造体{}とすると、その構造体のインスタンスをメモリ上に作り、そのポインタを取得する
	return &ItemMemoryRopository{items: items}
}

// ItemMemoryRopository型のポインタ（参照）を受け取る
// ポインタレシーバにすると、構造体のフィールド値を直接操作できたり、コピーせず効率よく扱える
// 元の構造体そのものを参照しているので、メソッド内から直接中身を変更ができるし、構造体が大きくてもパフォーマンスに影響がない

// Laravelとかでいうインスタンスをメソッドの頭にくっつけていると思ったらいい
func (r *ItemMemoryRopository) FindAll() (*[]models.Item, error) {
	// r.itemsのポインタを返す必要があるので&をつける
	return &r.items, nil
}

func (r *ItemMemoryRopository) FindById(itemId uint) (*models.Item, error) {
	for _, v := range r.items {
		if v.ID == itemId {
			return &v, nil
		}
	}
	return nil, errors.New("Item is not found")
}

func (r *ItemMemoryRopository) Create(newItem models.Item) (*models.Item, error) {
	newItem.ID = uint(len(r.items) + 1)
	r.items = append(r.items, newItem)
	return &newItem, nil
}

func (r *ItemMemoryRopository) Update(updateItem models.Item) (*models.Item, error) {
	for i, v := range r.items {
		if v.ID == updateItem.ID {
			r.items[i] = updateItem
			return &r.items[i], nil
		}
	}
	return nil, errors.New("Unexpected Error")
}

func (r *ItemMemoryRopository) Delete(itemId uint) error {
	for i, v := range r.items {
		if v.ID == itemId {
			// goにはスライス（配列）から特定のindexを削除するという処理がないので、以下のように実現している

			// 1.r.items[:i] は「削除対象より前の要素」の新しいスライス、r.items[i+1:] は「削除対象より後の要素」新しいスライス
			// 2.appendを使って「削除対象より前の要素」の新しいスライスに対し、r.items[i+1:] は「削除対象より後の要素」新しいスライスを合体させる
			// 3.ただし、スライス同士の結合はそのままではできないので、合体させるスライス（削除対象より後の要素）をスプレッド演算子を使って展開しながらappend
			// → これにより該当idのみを除外しつつ、新しいスライスを作成するということが可能になる

			r.items = append(r.items[:i], r.items[i+1:]...)
			return nil
		}
	}
	return errors.New("Item not Found")
}

type ItemRepository struct {
	db *gorm.DB
}

// Create implements IItemRepository.
func (r *ItemRepository) Create(newItem models.Item) (*models.Item, error) {
	// gormを介したDB登録では引数は参照を渡すこと
	result := r.db.Create(&newItem)

	if result.Error != nil {
		return nil, result.Error
	}
	return &newItem, nil
}

// Delete implements IItemRepository.
func (r *ItemRepository) Delete(itemId uint) error {
	deleteItem, err := r.FindById(itemId)
	if err != nil {
		return err
	}
	// 論理削除(deleted atに時刻が入るだけ)
	result := r.db.Delete(&deleteItem)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// FindAll implements IItemRepository.
func (r *ItemRepository) FindAll() (*[]models.Item, error) {

	// 検索結果を格納する変数
	var items []models.Item

	// 上記の変数の型はすでにmodels.Itemで定義されている。それに合わせた形で
	// データを取得&整形してくれる
	result := r.db.Find(&items)

	if result.Error != nil {
		return nil, result.Error
	}
	return &items, nil
}

// FindById implements IItemRepository.
func (r *ItemRepository) FindById(itemId uint) (*models.Item, error) {
	var item models.Item

	// 主キーがidであればカラムの指定はいらない
	// カラム指定の場合は次のような感じ
	// result := r.db.First(&item, "id = ?", itemId)
	result := r.db.First(&item, itemId)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("Item is not found")
		}
		return nil, result.Error
	}
	return &item, nil
}

// Update implements IItemRepository.
func (r *ItemRepository) Update(updateItem models.Item) (*models.Item, error) {

	// Saveメソッドは更新対象が存在すればupdate、存在しなければinsertといったアップサートを行う
	// updateItemには必要な部分のみを変更した1レコードが入っていて、それをそのまま上書きという感じ（=必要な部分のみ更新がはいる）
	result := r.db.Save(&updateItem)
	if result.Error != nil {
		return nil, result.Error
	}
	return &updateItem, nil
}

func NewItemRepository(db *gorm.DB) IItemRepository {
	return &ItemRepository{db: db}
}
