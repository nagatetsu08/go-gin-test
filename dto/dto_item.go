package dto

// CRUDでユーザーから渡ってくるパラメータ金型のようなもの(DTOという)
// パラメータを一旦この枠にはめることでバリデーションも行うことができる。
// PHPのように1個1個のパラメータに応じて処理するのではなく、メソッドごとにパラメータの型をきめておくと
// パラメータ周りの変更にも柔軟に対応できる

type CreateItemInput struct {
	// `json:〜`のことをタグという。JSON形式で渡ってきたパラメータnameをこのCreateItemInput.Nameに割り当てる（当てはめる）といういみになる
	// bindingはginバリデーションを示す。ちなみにバリデーション間にスペースとか入れるとエラーになる

	Name       string `json:"name" binding:"required,min=2"`
	Price      uint   `json:"price" binding:"required,min=1,max=99999999"`
	Desciption string `json:"description"`
}

type UpdateItemInput struct {
	// updateでは値が指定された部分のみを更新対象としたいので、型をポインタ型にする。
	// さらにbindingタグでomitnilをいれることで、ポインタ型のプロパティにnilがはいってきたときバリデーションがスキップされるという動きになる

	Name        *string `json:"name" binding:"omitnil,min=2"`
	Price       *uint   `json:"price" binding:"omitnil,min=1,max=99999999"`
	Description *string `json:"description"`
	SoldOut     *bool   `json:"soldout"`
}
