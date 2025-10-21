package repositories

import (
	"errors"
	"gin-freemarket/models"
)

// アイテムのリポジトリが持つべき基本機能（インターフェース）の実装
// メソッド名() (戻り値)
type IItemRepository interface {

	// FindAllというメソッド名で、戻り値がmodels.Item型スライスへのポインタとerrorを返す（errorがない場合はnil）
	// 戻り値は参照を返すのが一般的（無駄なコピーを避けて省メモリ・省コストにしたい）
	FindAll() (*[]models.Item, error)

	// id検索は1件のみ返ってくるので、戻り値は*models.Itemとなる（FindAllは複数件返ってくる想定だから配列）
	FindById(itemId uint) (*models.Item, error)

	Create(newItem models.Item) (*models.Item, error)
	Update(newItem models.Item) (*models.Item, error)
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
