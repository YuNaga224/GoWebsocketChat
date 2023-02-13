package main

import (
	"testing"
)

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	//値をセットしない場合
	client := new(client)
	url, err = authAvatar.GetAvatarURL(client)
	if err != ErrNoAvatarURL {
		t.Error("値が存在しない場合、AuthAvatar.GetAvatarURLは" + "ErrNoAvatarURLを返すべきです")

	}
	//値をセットした場合
	testUrl := "http://url-to-avatar/"
	client.UserData = map[string]interface{}{"avatar_url": testUrl}
	url, err = authAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("値が存在する場合authAvatar.GetAvatarURLは" + "エラーを返すべきではありません")

	} else {
		if url != testUrl {
			t.Error("AuthAvatart.GetAvatarURLは正しいURLを返すべきです")
		}
	}
}

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar gravatarAvatar
	client := new(client)
	client.userData = map[string]interface{}{"userid": "0bc83cb571cd1c50ba6f3e8a78ef1346"}
	url, err := gravatarAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("GravatarAvatar.GetAvatarURLはエラーを返すべきではありません")
	}
	if url != "//ww.gravatar.com/avatar/0bc83cb571cd1c50ba6f3e8a78ef1346" {
		t.Errorf("GravatarAvatar.GetAvatarURLが%sという誤った値を返しました", url)
	}
}
