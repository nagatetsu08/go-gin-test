package services

import (
	"gin-freemarket/models"
	"gin-freemarket/repositories"
)

// サービスクラスにもinterfaceを作るのがお作法らしい
type IItemService interface {
	FindAll() (*[]models.Item, error)
}

// ItemServiceの本体（クラスに相当）
// repositories.IItemRepositoryはインタフェース。(newしたときの定義)
// インタフェースを定義することで差し替えが容易になる
type ItemService struct {
	repository repositories.IItemRepository
}

// コンストラクタ
func NewItemService(repository repositories.IItemRepository) IItemService {
	return &ItemService{repository: repository}
}

func (s *ItemService) FindAll() (*[]models.Item, error) {
	// ここを&s.repository.FindAll()とやらないのは、すでに利用しているrepository(IItemRepository)のFindAllの戻り値がポインタだから。
	// 同じメソッド名でわかりにくいが、リポジトリ経由で呼び出していて、リポジトリ側ですでに参照を返しているので、こちらでわざわざ参照を返す必要がない
	return s.repository.FindAll()
}
