package services

import (
	"gin-freemarket/dto"
	"gin-freemarket/models"
	"gin-freemarket/repositories"
)

// サービスクラスにもinterfaceを作るのがお作法らしい
type IItemService interface {
	FindAll() (*[]models.Item, error)
	FindById(itemId uint, userId uint) (*models.Item, error)
	Create(createItemInput dto.CreateItemInput, userId uint) (*models.Item, error)
	Update(itemId uint, userId uint, updateItemInput dto.UpdateItemInput) (*models.Item, error)
	Delete(itemId uint, userId uint) error
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

func (s *ItemService) FindById(itemId uint, userId uint) (*models.Item, error) {
	return s.repository.FindById(itemId, userId)
}

func (s *ItemService) Create(createItemInput dto.CreateItemInput, userId uint) (*models.Item, error) {
	newItem := models.Item{
		Name:        createItemInput.Name,
		Price:       createItemInput.Price,
		Description: createItemInput.Desciption,
		SoldOut:     false,
		UserId:      userId,
	}

	return s.repository.Create(newItem)
}

func (s *ItemService) Update(itemId uint, userId uint, updateItemInput dto.UpdateItemInput) (*models.Item, error) {

	targetItem, err := s.FindById(itemId, userId)

	if err != nil {
		return nil, err
	}

	if updateItemInput.Name != nil {
		targetItem.Name = *updateItemInput.Name
	}
	if updateItemInput.Price != nil {
		targetItem.Price = *updateItemInput.Price
	}
	if updateItemInput.Description != nil {
		targetItem.Description = *updateItemInput.Description
	}
	if updateItemInput.SoldOut != nil {
		targetItem.SoldOut = *updateItemInput.SoldOut
	}

	// ここで*targetItemを渡しているのは、s.FindById(itemId)の結果がポインタで返ってくるから。
	// s.repository.Updateは普通の値を引数として要求しているので、ここでデシリアライズして値渡しをしている。
	// createは構造体をその時に作っていてそのまま渡しているので値渡しとなる。
	// よっぽど巨大なインスタンスを渡さないのであれば、参照渡しでOK
	return s.repository.Update(*targetItem)
}

func (s *ItemService) Delete(itemId uint, userId uint) error {
	return s.repository.Delete(itemId, userId)
}
